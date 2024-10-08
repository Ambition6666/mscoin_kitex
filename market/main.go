package main

import (
	cc "common/config"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/kitex-contrib/registry-etcd/retry"
	"grpc_common/kitex_gen/market/exchangerate"
	"grpc_common/kitex_gen/market/market"
	"market/config"
	"market/handler"
	"market/utils"
	"net"
	"os"
	"time"
)

func main() {

	// 配置初始化
	suite := cc.InitConfigClient(config.ServerName, config.ServerName, config.MID, config.EtcdAddr, config.GetConf())

	// 初始化工具
	utils.Init()
	defer utils.Close()

	// 日志注册
	klog.SetLogger(zap.NewLogger())
	klog.SetLevel(klog.LevelDebug)
	f, err := os.OpenFile("./log/output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	klog.SetOutput(f)

	// 服务注册
	addr, _ := net.ResolveTCPAddr("tcp", config.ServerAddr)

	retryConfig := retry.NewRetryConfig(
		retry.WithMaxAttemptTimes(10),
		retry.WithObserveDelay(20*time.Second),
		retry.WithRetryDelay(5*time.Second),
	)

	r, err := etcd.NewEtcdRegistryWithRetry([]string{config.EtcdAddr}, retryConfig) // r 不能重复使用.
	if err != nil {
		panic(err)
	}

	svr := server.NewServer(server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.ServerName}), server.WithRegistry(r), server.WithServiceAddr(addr), server.WithSuite(suite), server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	err = market.RegisterService(svr, handler.NewMarketImpl())
	if err != nil {
		panic(err)
	}
	err = exchangerate.RegisterService(svr, handler.NewExchangeRateImpl())
	if err != nil {
		panic(err)
	}

	err = svr.Run()
	if err != nil {
		panic(err)
	}
}
