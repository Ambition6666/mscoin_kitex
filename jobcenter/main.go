package main

import (
	cc "common/config"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	"jobcenter/config"
	"jobcenter/task"
	"jobcenter/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 日志注册
	klog.SetLogger(zap.NewLogger())
	klog.SetLevel(klog.LevelDebug)
	f, err := os.OpenFile("./log/output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	klog.SetOutput(f)

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
