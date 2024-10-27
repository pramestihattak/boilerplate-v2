package auth_test

import (
	"context"
	"testing"

	pb "boilerplate-v2/gen/auth"
	jwtMock "boilerplate-v2/pkg/jwt/mock"
	"boilerplate-v2/service/auth"
	"boilerplate-v2/status"
	storageAuth "boilerplate-v2/storage/auth"
	storageAuthMock "boilerplate-v2/storage/auth/mock"
	"boilerplate-v2/util"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storageAuthMock.NewMockPostgresStore(ctrl)
	mockJWT := jwtMock.NewMockJWTInterface(ctrl)

	authService := &auth.AuthService{
		Storage: mockStorage,
		JWT:     mockJWT,
		Logger:  logrus.New(),
	}

	hashedPassword, _ := util.HashAndSalt("password123")
	user := &storageAuth.LoginOutput{
		Email:    "test@example.com",
		Password: hashedPassword,
		Verified: true,
	}

	userNotValidated := &storageAuth.LoginOutput{
		Email:    "test@example.com",
		Password: hashedPassword,
		Verified: false,
	}

	tests := []struct {
		name          string
		setupMocks    func(input *pb.LoginRequest)
		input         *pb.LoginRequest
		expectedToken string
		expectedError error
	}{
		{
			name: "Successful login",
			input: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(input *pb.LoginRequest) {
				mockStorage.EXPECT().Login(gomock.Any(), &storageAuth.LoginInput{
					Email: input.GetEmail(),
				}).Return(user, nil)
				mockJWT.EXPECT().Sign(user).Return("valid_token", nil)
			},
			expectedToken: "valid_token",
			expectedError: nil,
		},
		{
			name: "not verified",
			input: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(input *pb.LoginRequest) {
				mockStorage.EXPECT().Login(gomock.Any(), &storageAuth.LoginInput{
					Email: input.GetEmail(),
				}).Return(userNotValidated, nil)
			},
			expectedError: status.ResponseFromCodeToErr(status.UserErrCode_AccountNotVerified),
		},
		{
			name: "not found",
			input: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMocks: func(input *pb.LoginRequest) {
				mockStorage.EXPECT().Login(gomock.Any(), &storageAuth.LoginInput{
					Email: input.GetEmail(),
				}).Return(nil, nil)
			},
			expectedError: status.ResponseFromCodeToErr(status.UserErrCode_AccountNotFound),
		},
		{
			name: "wrong password",
			input: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			setupMocks: func(input *pb.LoginRequest) {
				mockStorage.EXPECT().Login(gomock.Any(), &storageAuth.LoginInput{
					Email: input.GetEmail(),
				}).Return(user, nil)
			},
			expectedError: status.ResponseFromCodeToErr(status.UserErrCode_LoginWrongPassword),
		},
		{
			name: "system error",
			input: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "password",
			},
			setupMocks: func(input *pb.LoginRequest) {
				mockStorage.EXPECT().Login(gomock.Any(), &storageAuth.LoginInput{
					Email: input.GetEmail(),
				}).Return(nil, assert.AnError)
			},
			expectedError: status.ResponseFromCodeToErr(status.SystemErrCode_FailedToLogin),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks(tt.input)

			mdVals := map[string]string{
				"key": "value",
			}
			ctx := metadata.NewIncomingContext(context.Background(), metadata.New(mdVals))
			ctx = grpc.NewContextWithServerTransportStream(ctx, &mockServerTransportStream{})
			resp, err := authService.Login(ctx, tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedToken, resp.Token)
			}
		})
	}
}

// mockServerTransportStream is a mock implementation of the grpc.ServerTransportStream interface
type mockServerTransportStream struct{}

func (m *mockServerTransportStream) Method() string {
	return "foo"
}

func (m *mockServerTransportStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m *mockServerTransportStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m *mockServerTransportStream) SetTrailer(md metadata.MD) error {
	return nil
}
