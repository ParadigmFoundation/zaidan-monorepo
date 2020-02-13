package logger

import (
	"context"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor is a gRPC UnaryServerInterceptor that logs requests
func UnaryServerInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	fn := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		now := time.Now()

		resp, err := handler(ctx, req)

		entry := logrus.NewEntry(logger).WithFields(Fields{
			"grpc.method": path.Base(info.FullMethod),
			"grpc.status": status.Code(err),
			"took":        time.Since(now),
		})

		if err != nil {
			entry.WithError(err).Error("Request finished with error")
		} else {
			entry.Info("Request finished successfully")
		}

		return resp, err
	}

	return fn
}
