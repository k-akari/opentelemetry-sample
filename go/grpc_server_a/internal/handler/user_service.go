package handler

import (
	"context"

	pb "github.com/k-akari/opentelemetry-sample/go/proto/grpc_server_a/v1"
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
			UserId:   1,
			UserName: "John Doe",
		},
	}, nil
}
