package handler

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/phayes/freeport"
	"github.com/prrknh/dcm/db"
	"github.com/prrknh/dcm/logger"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func CreateContainer() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}

		var mnt []mount.Mount

		if len(r.URL.Query().Get("initsql")) > 0 {

			tmpDir, errr := ioutil.TempDir("", "mikasa_unittest")
			if errr != nil {
				log.Fatal(err)
			}
			defer os.Remove(tmpDir)

			fp, _ := ioutil.TempFile(tmpDir, "xxx")
			defer fp.Close()

			hoge := mount.Mount{Type: mount.TypeVolume, Source: tmpDir, Target: "/root/runtime.sql"}

			mnt = append(mnt, hoge)
		}

		port, err := freeport.GetFreePort()
		if err != nil {
			log.Fatal(err)
		}
		strPort := strconv.Itoa(port)

		container, er := cli.ContainerCreate(context.Background(),
			&container.Config{Image: "mikasa_unittest"},
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
			"mikasa_unittest_"+time.Now().Format("20060102150405"))

		if er != nil {
			panic(er)
		}
		cli.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})

		go func() {
			containerLog(container.ID)
		}()

		db.WaitInitialization(strPort)

		io.WriteString(w, container.ID)
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
