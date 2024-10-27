package auth_test

import (
	"context"
	"testing"

	pb "boilerplate-v2/gen/auth"
	"boilerplate-v2/service/auth"
	"boilerplate-v2/status"
	storageAuthMock "boilerplate-v2/storage/auth/mock"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

func TestVerification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageAuthMock.NewMockPostgresStore(ctrl)
	logger := logrus.New()

	service := &auth.AuthService{
		Logger:  logger,
		Storage: mockStorage,
	}

	tests := []struct {
		name           string
		req            *pb.VerificationRequest
		setupMock      func(req *pb.VerificationRequest)
		expectedResult *pb.VerificationResponse
		expectedError  error
	}{
		{
			name: "Successful verification",
			req: &pb.VerificationRequest{
				Email:             "test@example.com",
				VerificationToken: "valid_token",
			},
			setupMock: func(req *pb.VerificationRequest) {
				mockStorage.EXPECT().CheckedUserForVerification(gomock.Any(), req.GetEmail(), req.GetVerificationToken()).Return(1, nil)
				mockStorage.EXPECT().VerifyUser(gomock.Any(), req.GetEmail()).Return("user_id", nil)
			},
			expectedResult: &pb.VerificationResponse{Message: "your account has been verified user_id"},
			expectedError:  nil,
		},
		{
			name: "Missing parameters",
			req: &pb.VerificationRequest{
				Email: "test@example.com",
			},
			setupMock:      func(req *pb.VerificationRequest) {},
			expectedResult: nil,
			expectedError:  status.ResponseFromCodeToErr(status.UserErrCode_VerificationMissingParameters),
		},
		{
			name: "User not found",
			req: &pb.VerificationRequest{
				Email:             "nonexistent@example.com",
				VerificationToken: "invalid_token",
			},
			setupMock: func(req *pb.VerificationRequest) {
				mockStorage.EXPECT().CheckedUserForVerification(gomock.Any(), req.GetEmail(), req.GetVerificationToken()).Return(0, nil)
			},
			expectedResult: nil,
			expectedError:  status.ResponseFromCodeToErr(status.UserErrCode_AccountNotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(tt.req)

			mdVals := map[string]string{
				"key": "value",
			}
			ctx := metadata.NewIncomingContext(context.Background(), metadata.New(mdVals))
			result, err := service.Verification(ctx, tt.req)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
