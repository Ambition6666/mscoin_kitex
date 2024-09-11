package main

import (
	cc "common/config"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/logging/zap"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/kitex-contrib/registry-etcd/retry"
	"grpc_common/kitex_gen/ucenter/asset"
	"grpc_common/kitex_gen/ucenter/login"
	"grpc_common/kitex_gen/ucenter/member"
	"grpc_common/kitex_gen/ucenter/register"
	"net"
	"os"
	"time"
	"ucenter/config"
	"ucenter/handler"
	"ucenter/utils"
)

func main() {

	// 配置初始化
	suite := cc.InitConfigClient(config.ServerName, config.ServerName, config.MID, config.EtcdAddr, config.GetConf())

	// 工具初始化
	utils.Init()

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

	svr := server.NewServer(server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: config.ServerName}), server.WithRegistry(r), server.WithServiceAddr(addr), server.WithSuite(suite), server.WithRefuseTrafficWithoutServiceName())

	err = login.RegisterService(svr, handler.NewLoginImpl())
	if err != nil {
		panic(err)
	}

	err = register.RegisterService(svr, handler.NewRegisterImpl())
	if err != nil {
		panic(err)
	}

	err = asset.RegisterService(svr, handler.NewAssetImpl())
	if err != nil {
		panic(err)
	}

	err = member.RegisterService(svr, handler.NewMemberImpl())

	err = svr.Run()
	if err != nil {
		panic(err)
	}

}
