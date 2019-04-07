package websocks

import (
	"context"

	streamauth "github.com/tcfw/evntsrc/pkg/streamauth/protos"
	"github.com/tcfw/evntsrc/pkg/tracing"
	"google.golang.org/grpc"
)

func (c *Client) validateAuth(ctx context.Context, auth *AuthCommand) error {
	//@TODO pass through passport instead
	conn, err := grpc.Dial("streamauth:443", tracing.GRPCClientOptions()...)
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := streamauth.NewStreamAuthServiceClient(conn)

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
