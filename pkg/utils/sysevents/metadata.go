package sysevents

import (
	"context"
	"log"

	utils "github.com/tcfw/evntsrc/pkg/utils/authorization"
)

func appendContextUserInfo(ctx context.Context, event EventInterface) EventInterface {
	claims, err := utils.TokenClaimsFromContext(ctx)
	if err == nil {
		md := event.GetMetadata()

		if md == nil {
			md = map[string]interface{}{}
		}

		md["user"] = claims["sub"]
		event.SetMetadata(md)
	} else {
		log.Printf("failed to fetch claims: %s", err)
	}

	return event
}
