package main

import (
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"net"
	"tiktok-backend/dal"
	publish "tiktok-backend/kitex_gen/publish/publishservice"
	"tiktok-backend/pkg/bound"
	"tiktok-backend/pkg/constants"
	"tiktok-backend/pkg/middleware"
	tracer2 "tiktok-backend/pkg/tracer"
)

func Init() {
	tracer2.InitJaeger(constants.PublishServiceName)
	dal.Init()
}

func main() {
	r, err := etcd.NewEtcdRegistry([]string{constants.EtcdAddress}) // r should not be reused.
	if err != nil {
		panic(err)
	}
	ip, err := constants.GetOutBoundIP()
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", ip+constants.PublishServicePort)
	if err != nil {
		panic(err)
	}
	Init()
	svr := publish.NewServer(new(PublishServiceImpl),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constants.PublishServiceName}), // server name
		server.WithMiddleware(middleware.CommonMiddleware),                                                // 中间件
		server.WithMiddleware(middleware.ServerMiddleware),                                                // 中间件
		server.WithServiceAddr(addr),                                                                      // address
		server.WithLimit(&limit.Option{MaxConnections: 1000, MaxQPS: 100}),                                // limit
		server.WithMuxTransport(),                                                                         // Multiplex
		server.WithSuite(trace.NewDefaultServerSuite()),                                                   // tracer
		server.WithBoundHandler(bound.NewCpuLimitHandler()),                                               // BoundHandler
		server.WithRegistry(r),                                                                            // registry
	)
	err = svr.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
