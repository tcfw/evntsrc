package utils

import (
	"context"
	"fmt"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	passportSvc "github.com/tcfw/evntsrc/pkg/passport/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//ValidateAuthToken validates a token using the passport service
func ValidateAuthToken(ctx context.Context, token string) (bool, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, fmt.Errorf("failed to parse metadata")
	}

	auth := md.Get("authorization")
	if auth == nil || len(auth) == 0 {
		return false, fmt.Errorf("no authorization sent")
	}

	passportEndpoint, envExists := os.LookupEnv("PASSPORT_HOST")
	if envExists != true {
		passportEndpoint = "passport:443"
	}

	conn, err := grpc.DialContext(ctx, passportEndpoint, grpc.WithInsecure())
	if err != nil {
		return false, err
	}

	passport := passportSvc.NewAuthSeviceClient(conn)

	tokenResponse, err := passport.VerifyToken(ctx, &passportSvc.VerifyTokenRequest{Token: auth[0]})
	if err != nil {
		return false, err
	}
	if tokenResponse.Valid == false || tokenResponse.Revoked == true {
		return false, fmt.Errorf("Invalid auth token provided")
	}

	return true, nil
}

func getAuthToken(ctx context.Context) (string, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("failed to parse metadata")
	}

	auth := md.Get("grpcgateway-authorization")

	if auth == nil || len(auth) == 0 {
		auth = md.Get("authorization")
		if auth == nil || len(auth) == 0 {
			return "", fmt.Errorf("no authorization sent")
		}
	}

	return auth[0], nil
}

//ValidateAuthClaims validates a token in the Authorization header and returns claims map
func ValidateAuthClaims(ctx context.Context) (jwt.MapClaims, error) {

	token, err := getAuthToken(ctx)
	if err != nil {
		return nil, err
	}

	valid, err := ValidateAuthToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if valid == false {
		return nil, fmt.Errorf("Invalid token provided")
	}

	return TokenClaims(token)
}

//TokenClaimsFromContext fetches token claims via context
func TokenClaimsFromContext(ctx context.Context) (map[string]interface{}, error) {
	token, err := getAuthToken(ctx)
	if err != nil {
		return nil, err
	}

	return TokenClaims(token)
}

//TokenClaims parses a JWT token and returns it's body claims
func TokenClaims(token string) (map[string]interface{}, error) {
	jwtParser := &jwt.Parser{}
	claims := jwt.MapClaims{}
	_, _, err := jwtParser.ParseUnverified(token, &claims)

	return claims, err
}
