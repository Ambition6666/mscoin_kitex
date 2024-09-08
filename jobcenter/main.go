package main

import (
	cc "common/config"
	"jobcenter/config"
	"jobcenter/task"
	"jobcenter/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cc.InitConfigClient(config.ServerName, config.ServerName, config.MID, config.EtcdAddr, config.GetConf())
	utils.Init()
	defer utils.Close()
	t := task.NewTask()
	t.Run()

	go func() {
		exit := make(chan os.Signal)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-exit:
			log.Println("监听到中断信号，终止程序")
			utils.Close()
			t.Stop()
		}
	}()
	t.StartBlocking()
}
