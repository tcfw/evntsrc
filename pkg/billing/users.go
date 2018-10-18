package billing

import (
	"context"
	"os"

	userSvc "github.com/tcfw/evntsrc/pkg/users/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//ReqAuth TODO
type ReqAuth struct {
	Token string
}

//GetRequestMetadata TODO
func (a *ReqAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": a.Token,
	}, nil
}

//RequireTransportSecurity TODO
func (a *ReqAuth) RequireTransportSecurity() bool {
	return false
}

func newUserClient(ctx context.Context) (userSvc.UserServiceClient, error) {

	md, _ := metadata.FromIncomingContext(ctx)
	authReq := ReqAuth{}

	if auth := md.Get("authorization"); auth != nil {
		authReq.Token = auth[0]
	}

	userEndpoint, envExists := os.LookupEnv("USER_HOST")
	if envExists != true {
		userEndpoint = "users:443"
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}

	opts = append(opts, grpc.WithPerRPCCredentials(&authReq))

	conn, err := grpc.DialContext(ctx, userEndpoint, opts...)
	if err != nil {
		return nil, err
	}

	return userSvc.NewUserServiceClient(conn), nil
}
