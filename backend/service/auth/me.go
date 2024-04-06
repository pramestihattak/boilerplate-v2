package auth

import (
	"context"
	"fmt"
	"log"

	pb "boilerplate-v2/gen/auth"
	"boilerplate-v2/status"

	"google.golang.org/grpc/metadata"
)

func (s *AuthService) Me(ctx context.Context, req *pb.MeRequest) (*pb.MeResponse, error) {
	logger := s.logger.WithField("handler", "Me")
	logger.Info("Me called")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedReadMetadata)
	}

	if len(md["authorization"]) == 0 {
		logger.Errorf("fail to get user data: %v", "missing bearer token")
		return nil, status.ResponseFromCodeToErr(status.UserErrCode_MissingBearerToken)
	}

	auth, err := s.withAuth(md["authorization"][0])
	if err != nil {
		logger.Errorf("fail to login user: %v", err.Error())
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToGetAuthData)
	}

	log.Println("auth", auth.Email)
	log.Println("auth", auth.FullName)
	log.Println("auth", auth.UserID)

	return &pb.MeResponse{Message: "Hello "}, nil
}

func (s *AuthService) withAuth(token string) (*Auth, error) {
	if !s.jwt.IsValidToken(token) {
		return nil, fmt.Errorf("invalid token")
	}

	auth, err := s.jwt.GetClaims(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get jwt claims")
	}

	return &Auth{
		UserID:   auth.UserID,
		FullName: auth.FullName,
		Email:    auth.Email,
	}, nil
}
