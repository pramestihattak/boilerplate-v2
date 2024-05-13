package auth

import (
	"context"

	pb "boilerplate-v2/gen/auth"
	"boilerplate-v2/middleware"
	jwtPackage "boilerplate-v2/pkg/jwt"
	"boilerplate-v2/status"
	"boilerplate-v2/storage/postgres"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const grpcMaxMsgSize = 1024 * 1024 * 50 // 50 mb

type AuthService struct {
	pb.UnimplementedAuthServer
	logger  *logrus.Logger
	storage *postgres.Storage
	jwt     *jwtPackage.JWT
}

func NewService(logger *logrus.Logger, storage *postgres.Storage, jwt *jwtPackage.JWT) *AuthService {
	return &AuthService{
		logger:  logger,
		storage: storage,
		jwt:     jwt,
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
