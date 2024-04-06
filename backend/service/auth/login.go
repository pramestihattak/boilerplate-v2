package auth

import (
	"context"

	pb "boilerplate-v2/gen/auth"
	"boilerplate-v2/status"
	"boilerplate-v2/storage"
	"boilerplate-v2/util"

	"google.golang.org/grpc/metadata"
)

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	logger := s.logger.WithField("handler", "Login")
	logger.Info("Login called")
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedReadMetadata)
	}

	user, err := s.storage.Login(ctx, &storage.LoginInput{
		Email: req.GetEmail(),
	})
	if err != nil {
		logger.Errorf("fail to login user: %v", err.Error())
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToLogin)
	}

	if user == nil {
		logger.Errorf("fail to login user: %v", "account not found")
		return nil, status.ResponseFromCodeToErr(status.UserErrCode_AccountNotFound)
	}

	if !util.ComparePasswords(user.Password, req.GetPassword()) {
		logger.Errorf("fail to login user: %v", "wrong password")
		return nil, status.ResponseFromCodeToErr(status.UserErrCode_LoginWrongPassword)
	}

	token, err := s.jwt.Sign(user)
	if err != nil {
		logger.Errorf("fail to login user: %v", err.Error())
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToLogin)
	}

	return &pb.LoginResponse{Token: token}, nil
}
