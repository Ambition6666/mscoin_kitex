package main

import (
	cc "common/config"
	"exchange/config"
	"exchange/consumer"
	"exchange/handler"
	"exchange/processor"
	"exchange/utils"
	"exchange/utils/rpc"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/kitex-contrib/registry-etcd/retry"
	exchange "grpc_common/kitex_gen/exchange/order"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// 配置初始化
	suite := cc.InitConfigClient(config.ServerName, config.ServerName, config.MID, config.EtcdAddr, config.GetConf())

	// 日志注册
	klog.SetLogger(zap.NewLogger())
	klog.SetLevel(klog.LevelDebug)
	f, err := os.OpenFile("./log/output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// 工具初始化
	utils.Init()
	defer utils.Close()

	// rpc服务注册
	rpc.Init()

	defer f.Close()
	klog.SetOutput(f)

	// 启动消费者
	factory := processor.InitCoinTradeFactory()
	factory.Init()
	con := consumer.NewRocketmqConsumer(factory)
	con.Run()

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

	svr := exchange.NewServer(handler.NewOrderImpl(), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.ServerName}), server.WithRegistry(r), server.WithServiceAddr(addr), server.WithSuite(suite), server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
