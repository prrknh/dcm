package main

import (
	"github.com/prrknh/mikasa_container_api/handler"
	"github.com/prrknh/mikasa_container_api/logger"
	"net/http"
)

func main() {
	logs := make(chan string)
	logger := logger.LoggerMan{LogChan: logs}
	go logger.Log(logs)

	http.HandleFunc("/create", handler.CreateContainer(logger))
	http.ListenAndServe(":5000", nil)

}

