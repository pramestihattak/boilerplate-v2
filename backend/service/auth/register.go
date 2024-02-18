package auth

import (
	"context"
	"errors"

	pb "boilerplate-v2/gen/auth"

	"google.golang.org/grpc/metadata"
)

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.RegisterResponse{}, errors.New("fail to read metadata")
	}

	s.logger.Infof("Received: %v", req.GetName())
	return &pb.RegisterResponse{Message: "Hello " + req.GetName()}, nil
}
