package auth

import (
	"context"

	pb "boilerplate-v2/gen/auth"
	"boilerplate-v2/status"
	storageAuth "boilerplate-v2/storage/auth"
	"boilerplate-v2/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	logger := s.Logger.WithField("handler", "Login")
	logger.Info("Login called")
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedReadMetadata)
	}

	user, err := s.Storage.Login(ctx, &storageAuth.LoginInput{
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

	if !user.Verified {
		logger.Errorf("fail to login user: %v", "account hasn't been verified yet")
		return nil, status.ResponseFromCodeToErr(status.UserErrCode_AccountNotVerified)
	}

	if !util.ComparePasswords(user.Password, req.GetPassword()) {
		logger.Errorf("fail to login user: %v", "wrong password")
		return nil, status.ResponseFromCodeToErr(status.UserErrCode_LoginWrongPassword)
	}

	token, err := s.JWT.Sign(user)
	if err != nil {
		logger.Errorf("fail to login user: %v", err.Error())
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToLogin)
	}

	// Add token to header
	header := metadata.Pairs("jwt", token)
	if err := grpc.SetHeader(ctx, header); err != nil {
		logger.Errorf("failed to set header: %v", err)
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToSetHeader)
	}

	return &pb.LoginResponse{Token: token}, nil
}
