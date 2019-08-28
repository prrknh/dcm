package logger

import (
	"github.com/fatih/color"
	"sync"
)

type LoggerMan struct {
	LogChan chan string
	mu sync.Mutex
	cnt int
}

func (l LoggerMan) Write(p []byte)(n int, err error) {
	l.LogChan <- string(p)
	l.Log(l.LogChan)
	n = 0
	err = nil
	return
}

func (l LoggerMan)Log(logger chan string) {
	cnt := l.getCnt()
	for item := range logger {
		l.colorLog(item, cnt)
	}
}

func (l LoggerMan)colorLog(item string, num int) {
	index := num % 6
    switch index {
	case 0: color.Red(item)
	case 1: color.Green(item)
	case 2: color.Yellow(item)
	case 3: color.Blue(item)
	case 4: color.Magenta(item)
	case 5: color.Cyan(item)
	}
}

func (l *LoggerMan)getCnt() (cnt int){
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cnt++
	return l.cnt - 1
}
