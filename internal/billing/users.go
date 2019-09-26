package billing

import (
	"context"
	"os"

	"github.com/tcfw/evntsrc/internal/tracing"
	userSvc "github.com/tcfw/evntsrc/internal/users/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func withAuthContext(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	return metadata.NewOutgoingContext(ctx, md)
}

var (
	userConn *grpc.ClientConn
)

func newUserClient(ctx context.Context) (userSvc.UserServiceClient, error) {
	if userConn == nil {
		userEndpoint, envExists := os.LookupEnv("USER_HOST")
		if envExists != true {
			userEndpoint = "users:443"
		}
		conn, err := grpc.Dial(userEndpoint, tracing.GRPCClientOptions()...)
		if err != nil {
			return nil, err
		}

		userConn = conn
	}

	return userSvc.NewUserServiceClient(userConn), nil
}
