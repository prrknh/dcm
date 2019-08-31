package main

import (
	"context"
	"flag"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prrknh/dcm/handler"
	"net/http"

	"log"
)

func main() {

	image := initialize()

	http.HandleFunc("/create", handler.CreateContainer(image))
	http.ListenAndServe(":5000", nil)
}

func initialize() (image string) {
	flag.StringVar(&image, "image", "", "name of docker image")
	flag.Parse()
	if len(image) == 0 {
		log.Fatal("name must be set")
	}

	cli, _ := client.NewClientWithOpts(client.FromEnv)
	list, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		panic(err.Error())
	}
	for _, img := range list {
		for _, name := range img.RepoTags {
			if name == image+":latest" {
				return image
			}
		}
	}
	log.Fatalf("\"%s\" image not found", image)
	return
}
