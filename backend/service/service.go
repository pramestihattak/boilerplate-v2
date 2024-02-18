package service

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func RegisterServices(srv *grpc.Server, services ...func(srv *grpc.Server) error) error {
	for _, svc := range services {
		if err := svc(srv); err != nil {
			return err
		}
	}
	return nil
}

func RegisterGateways(mux *runtime.ServeMux, gwAddr string, opts []grpc.DialOption,
	services ...func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error) error {
	ctx := context.Background()
	for _, svc := range services {
		if err := svc(ctx, mux, gwAddr, opts); err != nil {
			return err
		}
	}
	return nil
}
