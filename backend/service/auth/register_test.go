package auth_test

import (
	"context"
	"testing"

	pb "boilerplate-v2/gen/auth"
	jwtMock "boilerplate-v2/pkg/jwt/mock"
	"boilerplate-v2/service/auth"
	"boilerplate-v2/status"
	storageAuthMock "boilerplate-v2/storage/auth/mock"

	// storageAuth "boilerplate-v2/storage/auth"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/metadata"
)

func TestRegister(t *testing.T) {
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
		setupMocks     func(input *pb.RegisterRequest)
		input          *pb.RegisterRequest
		expectedResult *pb.RegisterResponse
		expectedError  error
	}{
		{
			name: "Successful registration",
			input: &pb.RegisterRequest{
				FullName: "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			setupMocks: func(input *pb.RegisterRequest) {
				mockStorage.EXPECT().UserExist(gomock.Any(), input.GetEmail()).Return(0, nil)
				mockStorage.EXPECT().Register(gomock.Any(), gomock.Any()).Return("user123", nil)
			},
			expectedResult: &pb.RegisterResponse{Message: "Hello user123"},
			expectedError:  nil,
		},
		{
			name: "User already exists",
			input: &pb.RegisterRequest{
				FullName: "Jane Doe",
				Email:    "jane@example.com",
				Password: "password456",
			},
			setupMocks: func(input *pb.RegisterRequest) {
				mockStorage.EXPECT().UserExist(gomock.Any(), input.GetEmail()).Return(1, nil)
			},
			expectedResult: nil,
			expectedError:  status.ResponseFromCodeToErr(status.UserErrCode_AccountExist),
		},
		{
			name: "Registration failed",
			input: &pb.RegisterRequest{
				FullName: "Bob Smith",
				Email:    "bob@example.com",
				Password: "password789",
			},
			setupMocks: func(input *pb.RegisterRequest) {
				mockStorage.EXPECT().UserExist(gomock.Any(), input.GetEmail()).Return(0, nil)
				mockStorage.EXPECT().Register(gomock.Any(), gomock.Any()).Return("", assert.AnError)
			},
			expectedResult: nil,
			expectedError:  status.ResponseFromCodeToErr(status.SystemErrCode_FailedToRegister),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(tt.input)

			mdVals := map[string]string{
				"key": "value",
			}
			ctx := metadata.NewIncomingContext(context.Background(), metadata.New(mdVals))
			result, err := authService.Register(ctx, tt.input)

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
