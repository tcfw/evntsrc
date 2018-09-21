package apigw

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	users "github.com/tcfw/evntsrc/pkg/users/protos"
	"google.golang.org/grpc"
)

func registerUsers(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) {
	err := users.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "users:443", opts)
	if err != nil {
		panic(err)
	}
}
