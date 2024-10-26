package auth

import (
	"context"

	pb "boilerplate-v2/gen/auth"
	"boilerplate-v2/status"

	"google.golang.org/grpc/metadata"
)

func (s *AuthService) Verification(ctx context.Context, req *pb.VerificationRequest) (*pb.VerificationResponse, error) {
	logger := s.Logger.WithField("handler", "Verification")
	logger.Info("Verification called")
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedReadMetadata)
	}

	if req.GetEmail() == "" || req.GetVerificationToken() == "" {
		logger.Errorf("fail to verify user: %v", "parameter required")
		return nil, status.ResponseFromCodeToErr(status.UserErrCode_VerificationMissingParameters)
	}

	exist, err := s.Storage.CheckedUserForVerification(ctx, req.GetEmail(), req.GetVerificationToken())
	if err != nil {
		logger.Errorf("fail to verify user: %v", err.Error())
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToVerify)
	}

	if exist == 0 {
		logger.Errorf("fail to verify user: %v", "account not found")
		return nil, status.ResponseFromCodeToErr(status.UserErrCode_AccountNotFound)
	}

	userID, err := s.Storage.VerifyUser(ctx, req.GetEmail())
	if err != nil {
		logger.Errorf("fail to verify user: %v", err.Error())
		return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToVerify)
	}

	return &pb.VerificationResponse{Message: "your account has been verified " + userID}, nil
}
