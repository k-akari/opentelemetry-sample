package handler

import (
	"context"
	"fmt"

	pb "github.com/k-akari/opentelemetry-sample/go/proto/service_b/v1"
	cpb "github.com/k-akari/opentelemetry-sample/go/proto/service_c/v1"
)

type UserServiceHandler struct {
	usc cpb.UserServiceClient
	pb.UnimplementedUserServiceServer
}

func NewUserServiceHandler(usc cpb.UserServiceClient) *UserServiceHandler {
	return &UserServiceHandler{usc: usc}
}

func (s *UserServiceHandler) GetUser(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	user, err := s.usc.Get(ctx, &cpb.GetRequest{UserId: req.UserId})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &pb.GetResponse{
		User: &pb.User{
			UserId:   user.User.GetUserId(),
			UserName: user.User.GetUserName(),
		},
	}, nil
}
