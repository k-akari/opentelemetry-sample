package handler

import (
	"context"

	pb "github.com/k-akari/opentelemetry-sample/go/proto/grpc_server_b/v1"
)

type UserServiceHandler struct {
	pb.UnimplementedUserServiceServer
}

func NewUserServiceHandler() *UserServiceHandler {
	return &UserServiceHandler{}
}

func (s *UserServiceHandler) GetUser(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	return &pb.GetResponse{
		User: &pb.User{
			UserId:   2,
			UserName: "Service B User",
		},
	}, nil
}
