package apigw

import (
	"net/http"
	"strings"

	"github.com/tcfw/evntsrc/internal/tracing"

	passport "github.com/tcfw/evntsrc/internal/passport/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func getAuthToken(r *http.Request) string {
	if r.Header.Get("authorization") != "" {
		return r.Header.Get("authorization")
	}

	cookie, err := r.Cookie("session")
	if err != nil {
		return ""
	}
	return cookie.Value
}

//@TODO secure against session fixation
func authGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldValidate(r) {
			authToken := getAuthToken(r)

			if authToken == "" {
				http.Error(w, "Forbidden. No API Key provided", 403)
				return
			}

			conn, err := grpc.DialContext(r.Context(), passportEndpoint, tracing.GRPCClientOptions()...)
			if err != nil {
				panic(err)
			}

			go func() {
				<-r.Context().Done()
				if cerr := conn.Close(); cerr != nil {
					grpclog.Printf("Failed to close conn to %s: %v", passportEndpoint, cerr)
				}
			}()

			svc := passport.NewAuthSeviceClient(conn)

			response, err := svc.VerifyToken(r.Context(), &passport.VerifyTokenRequest{Token: authToken})
			if err != nil {
				panic(err)
			}
			conn.Close()
			if !response.Valid {
				http.Error(w, "Forbidden. Invalid API Key provided", 403)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func shouldValidate(r *http.Request) bool {
	return !strings.HasPrefix(r.URL.String(), "/v1/auth")
}
