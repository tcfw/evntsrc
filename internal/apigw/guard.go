package apigw

import (
	"net/http"
	"strings"

	"github.com/tcfw/evntsrc/internal/tracing"

	passport "github.com/tcfw/evntsrc/internal/passport/protos"
	"google.golang.org/grpc"
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

var (
	passportConn *grpc.ClientConn
)

//@TODO secure against session fixation
func authGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldValidate(r) {
			authToken := getAuthToken(r)

			if authToken == "" {
				http.Error(w, "Forbidden. No API Key provided", 403)
				return
			}

			if passportConn == nil {
				conn, err := grpc.Dial(passportEndpoint, tracing.GRPCClientOptions()...)
				if err != nil {
					http.Error(w, "Failed to validate token", 500)
					panic(err)
				}

				passportConn = conn
			}

			go func() {
				<-r.Context().Done()
			}()

			svc := passport.NewAuthSeviceClient(passportConn)

			response, err := svc.VerifyToken(r.Context(), &passport.VerifyTokenRequest{Token: authToken})
			if err != nil {
				http.Error(w, "Failed to validate token", 500)
				// panic(err)
				return
			}
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
