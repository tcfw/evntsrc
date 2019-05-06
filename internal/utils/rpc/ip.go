package utils

import (
	"context"
	"net"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

//RemoteIPFromContext returns the IP address of the RPC client, also taking proxy calls into account
func RemoteIPFromContext(ctx context.Context) net.IP {
	md, _ := metadata.FromIncomingContext(ctx)
	peer, _ := peer.FromContext(ctx)

	remoteIP := net.ParseIP(peer.Addr.String())
	if remoteIP == nil && len(md.Get("x-forwarded-for")) > 0 {
		remoteIP = net.ParseIP(md.Get("x-forwarded-for")[0])
	}

	return remoteIP
}
