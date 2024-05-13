package auth

import (
	"context"

	pb "boilerplate-v2/gen/auth"
)

func (s *AuthService) Me(ctx context.Context, req *pb.MeRequest) (*pb.MeResponse, error) {
	logger := s.logger.WithField("handler", "Me")
	logger.Info("Me called")

	userID, err := GetUserIDContext(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.MeResponse{Message: userID}, nil
}
