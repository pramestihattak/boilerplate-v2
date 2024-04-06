package auth

import (
	"context"

	pb "boilerplate-v2/gen/auth"
	"boilerplate-v2/status"
	"boilerplate-v2/storage"
	"boilerplate-v2/util"

	"google.golang.org/grpc/metadata"
)

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	logger := s.logger.WithField("handler", "Register")
	logger.Info("Register called")
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedReadMetadata)
	}

	exist, err := s.storage.UserExist(ctx, req.GetEmail())
	if err != nil {
		logger.Errorf("fail to register user: %v", err.Error())
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToRegister)
	}

	if exist > 0 {
		logger.Errorf("fail to register user: %v", "account exist")
		return nil, status.ResponseFromCodeToErr(status.UserErrCode_AccountExist)
	}

	hashedPassword, err := util.HashAndSalt(req.GetPassword())
	if err != nil {
		logger.Errorf("fail to register user: %v", err.Error())
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToRegister)
	}

	userID, err := s.storage.Register(ctx, storage.Register{
		FullName:          req.GetFullName(),
		Email:             req.GetEmail(),
		Password:          hashedPassword,
		VerificationToken: util.RandomStringGenerator(10),
	})
	if err != nil {
		logger.Errorf("fail to register user: %v", err.Error())
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToRegister)
	}

	return &pb.RegisterResponse{Message: "Hello " + userID}, nil
}
