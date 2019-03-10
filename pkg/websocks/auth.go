package websocks

import (
	"context"

	streamauth "github.com/tcfw/evntsrc/pkg/streamauth/protos"
	"google.golang.org/grpc"
)

func (c *Client) validateAuth(auth *AuthCommand) error {
	//@TODO pass through passport instead
	conn, err := grpc.Dial("streamauth:443", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	cli := streamauth.NewStreamAuthServiceClient(conn)

	sk, err := cli.ValidateKeySecret(context.Background(), &streamauth.KSRequest{Stream: auth.Stream, Key: auth.Key, Secret: auth.Secret})
	c.authKey = sk
	return err
}

func (c *Client) authFromHeader(apiKey string, apiSec string, stream int32) error {
	authCmd := &AuthCommand{Stream: stream, Key: apiKey, Secret: apiSec}

	err := c.validateAuth(authCmd)
	if err == nil {
		c.auth = authCmd
	}
	return err
}
