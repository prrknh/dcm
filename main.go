package main

import (
	"context"
	"flag"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/prrknh/dcm/handler"
	"net/http"
	"strings"

	"log"
)

func main() {

	image := initialize()

	log.Println("============================================")
	log.Println("============= start server =================")
	log.Println("============================================")

	http.HandleFunc("/create", handler.CreateContainer(image))
	http.HandleFunc("/stop", handler.StopContainer())

	if err := http.ListenAndServe(":5000", logRequest(http.DefaultServeMux)); err != nil {
		log.Fatalf(err.Error())
	}
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
			if strings.HasPrefix(name, image) {
				return image
			}
		}
	}
	log.Fatalf("\"%s\" image not found", image)
	return
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
