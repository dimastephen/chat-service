package interceptor

import (
	"context"
	"github.com/dimastephen/chatServer/internal/metrics"
	"google.golang.org/grpc"
)

func MetricInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metrics.IncRequestCounter()

	res, err := handler(ctx, req)
	if err != nil {
		metrics.IncResponseCounter("error", info.FullMethod)
	} else {
		metrics.IncResponseCounter("success", info.FullMethod)
	}
	return res, err
}
