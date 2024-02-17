package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/k-akari/opentelemetry-sample/go/grpc_server_a/internal/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/k-akari/opentelemetry-sample/go/common/grpcutil"
	pb "github.com/k-akari/opentelemetry-sample/go/proto/service_b/v1"
)

const (
	exitFail    = 1
	serviceName = "grpc_server_a"
)

func main() {
	l := grpcutil.NewZapLogger()

	if err := run(context.Background(), l); err != nil {
		l.Error(err.Error())
		os.Exit(exitFail)
	}
}

func run(ctx context.Context, l *zap.Logger) error {
	env, err := newEnv()
	if err != nil {
		return fmt.Errorf("failed to create env: %w", err)
	}

	s := grpc.NewServer(grpcutil.UnaryServerInterceptors(l))
	reflection.Register(s)

	hs := health.NewServer()
	hs.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)
	hs.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	healthpb.RegisterHealthServer(s, hs)

	h := handler.NewUserServiceHandler()
	pb.RegisterUserServiceServer(s, h)
	l.Info("Successfully registered handler.")

	errCh := make(chan error)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	go func() {
		defer close(errCh)
		if err := s.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("failed to serve: %v", err)
		}
	case <-ctx.Done():
		s.GracefulStop()
	}

	return nil
}
