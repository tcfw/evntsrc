package apigw

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	billing "github.com/tcfw/evntsrc/pkg/billing/protos"
	"google.golang.org/grpc"
)

const billingEndpoint = "billing:443"

func registerBilling(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) {
	err := billing.RegisterBillingServiceHandlerFromEndpoint(ctx, mux, billingEndpoint, opts)
	if err != nil {
		panic(err)
	}
}