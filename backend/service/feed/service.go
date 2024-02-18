package feed

import (
	pb "boilerplate-v2/gen/feed"
	"boilerplate-v2/storage/postgres"
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const grpcMaxMsgSize = 1024 * 1024 * 50 // 50 mb

type FeedService struct {
	pb.UnimplementedFeedServer
	logger  *logrus.Logger
	storage *postgres.Storage
}

func NewFeedService(logger *logrus.Logger, storage *postgres.Storage) *FeedService {
	return &FeedService{
		logger:  logger,
		storage: storage,
	}
}

func RegisterService(s *FeedService) func(srv *grpc.Server) error {
	return func(srv *grpc.Server) error {
		pb.RegisterFeedServer(srv, s)
		return nil
	}
}

func RegisterGateway() func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
	return func(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) error {
		pb.RegisterFeedHandlerFromEndpoint(ctx, mux, addr, append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxMsgSize), grpc.MaxCallSendMsgSize(grpcMaxMsgSize))))
		return nil
	}
}
