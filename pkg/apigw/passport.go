package apigw

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	passport "github.com/tcfw/evntsrc/pkg/passport/protos"
	"google.golang.org/grpc"
)

const passportEndpoint = "passport:443"

func registerPassport(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) {
	if err := passport.RegisterAuthSeviceHandlerFromEndpoint(ctx, mux, passportEndpoint, opts); err != nil {
		panic(err)
	}
}
