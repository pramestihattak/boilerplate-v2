package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "boilerplate-v2/gen/auth"
	jwtPackage "boilerplate-v2/pkg/jwt"

	"boilerplate-v2/status"
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	ContextKeyUserID = contextKey("userID")
)

type Middleware struct {
	jwt *jwtPackage.JWT
}

func NewMiddleware(jwt *jwtPackage.JWT) *Middleware {
	return &Middleware{
		jwt: jwt,
	}
}

func (m *Middleware) AuthUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	var err error
	switch req.(type) {
	case *pb.MeRequest:
		ctx, err = m.WithAuth(ctx)
		if err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}

func (m *Middleware) WithAuth(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.ResponseFromCodeToErr(status.SystemErrCode_FailedReadMetadata)
	}

	if len(md["authorization"]) == 0 {
		return ctx, status.ResponseFromCodeToErr(status.UserErrCode_MissingBearerToken)
	}

	token := md["authorization"][0]
	if !m.jwt.IsValidToken(token) {
		return ctx, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToGetAuthData)
	}

	auth, err := m.jwt.GetClaims(token)
	if err != nil {
		return ctx, status.ResponseFromCodeToErr(status.SystemErrCode_FailedToGetAuthData)
	}

	ctx = context.WithValue(ctx, ContextKeyUserID, auth.UserID)

	return ctx, nil
}
