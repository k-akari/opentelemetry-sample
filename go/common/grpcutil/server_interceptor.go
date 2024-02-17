package grpcconfig

import (
	"context"
	"errors"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logging "github.com/grpc-ecosystem/go-grpc-middleware/logging"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const healthCheckMethodName = "/grpc.health.v1.Health/Check"

func UnaryServerInterceptors(zapLogger *zap.Logger) grpc.ServerOption {
	if zapLogger == nil {
		panic("zap logger is required")
	}

	return grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(
				zapLogger,
				grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
					// will not log gRPC calls if it was a call to healthcheck and no error was raised
					if err == nil && fullMethodName == healthCheckMethodName {
						return false
					}

					// by default everything will be logged
					return true
				}),
				grpc_zap.WithCodes(grpc_logging.DefaultErrorToCode),
			),
			grpc_zap.PayloadUnaryServerInterceptor(
				zapLogger,
				func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
					// will not log gRPC calls if it was a call to healthcheck and no error was raised
					if fullMethodName == healthCheckMethodName {
						return false
					}

					// by default everything will be logged
					return true
				},
			),
			func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				resp, err := handler(ctx, req)
				if err == nil {
					return resp, nil
				}

				if errors.Is(err, context.Canceled) {
					zapLogger.Warn("context.Canceled", zap.Error(err))
					return resp, status.Error(codes.Canceled, err.Error())
				}
				if errors.Is(err, context.DeadlineExceeded) {
					zapLogger.Warn("context.DeadlineExceeded", zap.Error(err))
					return resp, status.Error(codes.DeadlineExceeded, err.Error())
				}

				zapLogger.Error("unary error", zap.Error(err))

				return resp, err
			},
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
				zapLogger.Error("panic recovered", zap.Any("panic", p))
				return status.Errorf(codes.Unknown, "panic triggered: %v", p)
			})),
		),
	)
}
