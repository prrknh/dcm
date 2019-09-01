package handler

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/phayes/freeport"
	"github.com/prrknh/dcm/db"
	"github.com/prrknh/dcm/logger"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func CreateContainer(image string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			fmt.Println(w, "the request method is not supported")
			return
		}

		var mnt []mount.Mount

		initsql := r.PostFormValue("initsql")

		if len(initsql) > 0 {

			tmpDir, err := ioutil.TempDir("/tmp", "")
			if err != nil {
				log.Fatal(err)
			}
			f, err := os.Create(tmpDir + "/runtime.sql")
			if err != nil {
				log.Fatal(err)
			}
			if _, err := f.WriteString(initsql); err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := os.RemoveAll(tmpDir); err != nil {
					log.Fatalf(err.Error())
				}
			}()

			mnt = []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: tmpDir,
					Target: "/root/mount",
				},
			}
		}

		port, err := freeport.GetFreePort()
		if err != nil {
			log.Fatal(err)
		}
		strPort := strconv.Itoa(port)

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}

		con, er := cli.ContainerCreate(context.Background(),
			&container.Config{Image: image},
			&container.HostConfig{
				Mounts: mnt,
				PortBindings: nat.PortMap{
					"3306/tcp": []nat.PortBinding{
						{
							HostIP:   "",
							HostPort: strPort},
					},
				},
			},
			nil,
			image+time.Now().Format("20060102150405"))

		if er != nil {
			panic(er)
		}
		if err := cli.ContainerStart(context.Background(), con.ID, types.ContainerStartOptions{}); err != nil {
			fmt.Println(w, err.Error())
		}

		go func() {
			containerLog(con.ID)
		}()

		db.WaitInitialization(strPort)

		if _, err := fmt.Fprintf(w, "{\"containerId\": \"%s\", \"port\": %s}", con.ID, strPort); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func containerLog(containerId string) int {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return 1
	}
	r, err := c.ContainerLogs(context.Background(), containerId, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err != nil {
		log.Println(err)
		return 1
	}

	cl := logger.NewContainerLogger(containerId)

	_, err = stdcopy.StdCopy(cl, cl, r)
	if err != nil {
		log.Println(err)
		return 1
	}

	return 0
}
