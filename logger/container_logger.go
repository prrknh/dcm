package logger

import (
	"github.com/fatih/color"
	"sync"
)

var counter int
var mu sync.Mutex

type ContainerLogger struct {
	ContainerId string
	ID          int
}

func NewContainerLogger(containerId string) *ContainerLogger {
	return &ContainerLogger{ContainerId: containerId, ID: getID()}
}

func (cl ContainerLogger) Write(data []byte) (n int, err error) {
	for i, b := range data {
		if b == '\n' {
			cl.colorLog(string(data[n:i+1]), cl.ID)
			n += i
		}
	}
	return len(data), nil
}

func (cl ContainerLogger) colorLog(item string, num int) {
	index := num % 6
	item = "[" + cl.ContainerId + "] " + item
	switch index {
	case 0:
		color.Red(item)
	case 1:
		color.Green(item)
	case 2:
		color.Yellow(item)
	case 3:
		color.Blue(item)
	case 4:
		color.Magenta(item)
	case 5:
		color.Cyan(item)
	}
}

func getID() (cnt int) {
	mu.Lock()
	defer mu.Unlock()
	cnt = counter
	counter++
	return cnt
}
