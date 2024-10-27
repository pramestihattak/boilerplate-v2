package auth

import (
	"context"

	pb "boilerplate-v2/gen/auth"
	"boilerplate-v2/middleware"
	jwtPackage "boilerplate-v2/pkg/jwt"
	"boilerplate-v2/status"
	storageAuth "boilerplate-v2/storage/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const grpcMaxMsgSize = 1024 * 1024 * 50 // 50 mb

type AuthService struct {
	pb.UnimplementedAuthServer
	Logger  *logrus.Logger
	Storage storageAuth.PostgresStore
	JWT     jwtPackage.JWTInterface
}

func NewService(logger *logrus.Logger, storage storageAuth.PostgresStore, jwt *jwtPackage.JWT) *AuthService {
	return &AuthService{
		Logger:  logger,
		Storage: storage,
		JWT:     jwt,
	}
}

func RegisterService(s *AuthService) func(srv *grpc.Server) error {
	return func(srv *grpc.Server) error {
		pb.RegisterAuthServer(srv, s)
		return nil
	}
}

func RegisterGateway() func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
	return func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
		pb.RegisterAuthHandlerFromEndpoint(ctx, mux, addr, append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxMsgSize), grpc.MaxCallSendMsgSize(grpcMaxMsgSize))))
		return nil
	}
}

func GetUserIDContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(middleware.ContextKeyUserID).(string)
	if !ok {
		return "", status.ResponseFromCodeToErr(status.SystemErrCode_FailedToGetAuthData)
	}
	return userID, nil
}
