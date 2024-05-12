package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "boilerplate-v2/gen/auth"
	"boilerplate-v2/status"
)

func AuthUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	switch req.(type) {
	case *pb.MeRequest:
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.ResponseFromCodeToErr(status.SystemErrCode_FailedReadMetadata)
		}

		if len(md["authorization"]) == 0 {
			return nil, status.ResponseFromCodeToErr(status.UserErrCode_MissingBearerToken)
		}
	}

	return handler(ctx, req)
}
