package main

import (
	"jobcenter/init"
	"jobcenter/task"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	init.Init()
	t := task.NewTask()
	t.Run()

	go func() {
		exit := make(chan os.Signal)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-exit:
			log.Println("监听到中断信号，终止程序")
			init.Close()
			t.Stop()
		}
	}()
	t.StartBlocking()
}
