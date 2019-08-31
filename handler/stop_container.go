package handler

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"net/http"
)

func StopContainer() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		containerId := r.URL.Query().Get("containerId")
		if len(containerId) == 0 {
			io.WriteString(w, "container id must be set")
		}
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			panic(err)
		}

		go func() {
			if err := cli.ContainerStop(context.Background(), containerId, nil); err != nil {
				fmt.Println(err.Error())
				return
			}
			if err := cli.ContainerRemove(context.Background(), containerId, types.ContainerRemoveOptions{}); err != nil {
				fmt.Println(err.Error())
			}
		}()

		io.WriteString(w, "ok")
	}
}
