package rpc

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"tiktok-backend/kitex_gen/comment/commentservice"
	"tiktok-backend/pkg/constants"
	"tiktok-backend/pkg/middleware"
	"time"
)

var commentClient commentservice.Client

func initCommentRPC() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress}) // 服务注册发现中心
	if err != nil {
		panic(err)
	}

	c, err := commentservice.NewClient(
		constants.CommentServiceName,
		client.WithMiddleware(middleware.CommonMiddleware), // 通用中间件
		client.WithInstanceMW(middleware.ClientMiddleware), // 客户端中间件
		client.WithMuxConnection(1),                        // 多路复用
		client.WithRPCTimeout(3*time.Second),               // 设置 rpc 调用超时时间
		client.WithConnectTimeout(50*time.Millisecond),     // 设置 rpc 连接超时时间
		client.WithFailureRetry(retry.NewFailurePolicy()),  // 重试，默认2次，可以设置重试次数，熔断
		client.WithSuite(trace.NewDefaultClientSuite()),    // 链路追踪，默认使用 OpenTracing GlobalTracer
		client.WithResolver(r),                             // 服务发现
	)
	if err != nil {
		panic(err)
	}
	commentClient = c
}