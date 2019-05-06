package main

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"

	emailSvc "github.com/tcfw/evntsrc/internal/emails/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
)

func main() {

	ctx := context.Background()

	opts := tracing.GRPCClientOptions()

	emailEndpoint, envExists := os.LookupEnv("EMAIL_HOST")
	if envExists != true {
		emailEndpoint = "emails:443"
	}

	conn, err := grpc.DialContext(ctx, emailEndpoint, opts...)
	if err != nil {
		panic(err)
	}

	svc := emailSvc.NewEmailServiceClient(conn)

	_, err = svc.Send(context.Background(), &emailSvc.Email{
		From:      "noreply@evntsrc.io",
		To:        []string{"tom@finao.com.au"},
		Subject:   "Welcome to your new EvntSrc.io account!",
		PlainText: "Welcome to Evntsrc.io. If you did not sign up, Please contact us at support@evntsrc.io",
		Html:      `<h3>Welcome to your new Account</h3><p>If you did not sign up, please contact us at <a href="mailto:support@evntsrc.io">support@evntsrc.io</a>`,
	})

	fmt.Printf("ERR: %v\n", err)

}
