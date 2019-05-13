package websocks

import (
	"context"

	streamauth "github.com/tcfw/evntsrc/internal/streamauth/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
	"google.golang.org/grpc"
)

var (
	streamAuthConn *grpc.ClientConn
)

func (c *Client) validateAuth(ctx context.Context, auth *AuthCommand) error {
	if streamAuthConn == nil {
		conn, err := grpc.Dial("streamauth:443", tracing.GRPCClientOptions()...)
		if err != nil {
			return err
		}
		streamAuthConn = conn
	}

	cli := streamauth.NewStreamAuthServiceClient(streamAuthConn)

	sk, err := cli.ValidateKeySecret(ctx, &streamauth.KSRequest{Stream: auth.Stream, Key: auth.Key, Secret: auth.Secret})
	c.authKey = sk
	return err
}

func (c *Client) authFromHeader(ctx context.Context, apiKey string, apiSec string, stream int32) error {
	authCmd := &AuthCommand{Stream: stream, Key: apiKey, Secret: apiSec}

	err := c.validateAuth(ctx, authCmd)
	if err == nil {
		c.auth = authCmd
	}
	return err
}
