package rpc

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"
	trace "github.com/kitex-contrib/tracer-opentracing"
	"tiktok-backend/kitex_gen/feed"
	"tiktok-backend/kitex_gen/feed/feedservice"
	"tiktok-backend/pkg/constants"
	"tiktok-backend/pkg/errno"
	"tiktok-backend/pkg/middleware"
	"time"
)

var feedClient feedservice.Client

func initFeedRPC() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddress}) // 服务注册发现中心
	if err != nil {
		panic(err)
	}

	c, err := feedservice.NewClient(
		constants.FeedServiceName,
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
	feedClient = c
}

func Feed(ctx context.Context, req *feed.DouyinFeedRequest) ([]*feed.Video, int64, error) {
	resp, err := feedClient.Feed(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode != 0 {
		return nil, 0, errno.NewErrNo(resp.StatusCode, resp.StatusMsg)
	}
	return resp.VideoList, resp.NextTime, nil
}
