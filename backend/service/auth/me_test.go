package auth_test

import (
	"context"
	"testing"

	pb "boilerplate-v2/gen/auth"
	"boilerplate-v2/middleware"
	jwtMock "boilerplate-v2/pkg/jwt/mock"
	"boilerplate-v2/service/auth"
	"boilerplate-v2/status"
	storageAuthMock "boilerplate-v2/storage/auth/mock"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

func TestMe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageAuthMock.NewMockPostgresStore(ctrl)
	mockJWT := jwtMock.NewMockJWTInterface(ctrl)

	authService := &auth.AuthService{
		Storage: mockStorage,
		JWT:     mockJWT,
		Logger:  logrus.New(),
	}

	tests := []struct {
		name           string
		setupContext   func() context.Context
		expectedResult *pb.MeResponse
		expectedError  error
	}{
		{
			name: "Successful Me request",
			setupContext: func() context.Context {
				ctx := context.WithValue(context.Background(), middleware.ContextKeyUserID, "user123")
				return metadata.NewIncomingContext(ctx, metadata.New(map[string]string{"key": "value"}))
			},
			expectedResult: &pb.MeResponse{Message: "user123"},
			expectedError:  nil,
		},
		{
			name: "Failed to get user ID from context",
			setupContext: func() context.Context {
				return metadata.NewIncomingContext(context.Background(), metadata.New(map[string]string{"key": "value"}))
			},
			expectedResult: nil,
			expectedError:  status.ResponseFromCodeToErr(status.SystemErrCode_FailedToGetAuthData),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.setupContext()
			result, err := authService.Me(ctx, &pb.MeRequest{})

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
