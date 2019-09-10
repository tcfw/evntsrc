package passport

import (
	"encoding/json"
	fmt "fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/tcfw/evntsrc/internal/passport/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
	userSvc "github.com/tcfw/evntsrc/internal/users/protos"
	utils "github.com/tcfw/evntsrc/internal/utils/authorization"
	rpcUtils "github.com/tcfw/evntsrc/internal/utils/rpc"
	events "github.com/tcfw/evntsrc/internal/utils/sysevents"
	"github.com/tstranex/u2f"
	"golang.org/x/crypto/bcrypt"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

//Server basic Server construct
type Server struct {
	mu sync.Mutex // protects routeNotes
}

//NewServer creates a ne struct to interface the auth server
func NewServer() *Server {
	s := &Server{}
	return s
}

func withAuthContext(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	return metadata.NewOutgoingContext(ctx, md)
}

var (
	userConn *grpc.ClientConn
)

func newUserClient() (userSvc.UserServiceClient, error) {
	if userConn == nil {
		userEndpoint, envExists := os.LookupEnv("USER_HOST")
		if envExists != true {
			userEndpoint = "users:443"
		}

		conn, err := grpc.Dial(userEndpoint, tracing.GRPCClientOptions()...)
		if err != nil {
			return nil, err
		}
		userConn = conn
	}

	return userSvc.NewUserServiceClient(userConn), nil
}

func validateMFA(user *userSvc.User) (*pb.AuthResponse, error) {
	mfa := &pb.MFAResponse{}
	mfa.Type = reflect.TypeOf(user.GetMfa()).String()

	switch mfaType := user.Mfa.MFA.(type) {
	case *userSvc.MFA_FIDO:
		U2FChallenge, err := u2f.NewChallenge("evntsrc.io", []string{"evntsrc.io"})
		if err != nil {
			return nil, err
		}
		userReg := &u2f.Registration{}
		if err = json.Unmarshal(user.Mfa.GetFIDO().Registration, userReg); err != nil {
			return nil, err
		}
		r := U2FChallenge.SignRequest([]u2f.Registration{*userReg})
		mfa.Challenge = r.Challenge
	case *userSvc.MFA_TOTP:
	case *userSvc.MFA_SMS:
		//...
	default:
		return nil, fmt.Errorf("Unknown MFA type %v", mfaType)
	}

	return &pb.AuthResponse{Success: false, MFAResponse: mfa}, nil
}

//Authenticate takes in oneof a authentication types and tries to generate tokens
func (s *Server) Authenticate(ctx context.Context, request *pb.AuthRequest) (*pb.AuthResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	remoteIP := rpcUtils.RemoteIPFromContext(ctx)

	ok, ttl, _ := checkIPRateLimit(remoteIP)
	if ok == false {
		grpc.SendHeader(ctx, metadata.Pairs("Grpc-Metadata-X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix())))
		events.BroadcastNonStreamingEvent(ctx, events.AuthenticateEvent{Event: &events.Event{Type: "io.evntsrc.passport.limite_reached"}, Err: fmt.Sprintf("limit reached: ip"), IP: remoteIP.String()})
		return nil, status.Errorf(codes.ResourceExhausted, "IP rate limit exceeded. Wait %s before making another request", ttl)
	}

	extraClaims := make(map[string]interface{})

	switch authType := request.Creds.(type) {
	case *pb.AuthRequest_UserCreds:
		username := request.GetUserCreds().GetUsername()

		//Rate limit
		ok, ttl, remaining := checkUserRateLimit(username, remoteIP)
		if ok == false {
			events.BroadcastNonStreamingEvent(ctx, events.AuthenticateEvent{Event: &events.Event{Type: "io.evntsrc.passport.limite_reached"}, Err: fmt.Sprintf("limit reached: user"), IP: remoteIP.String(), User: username})
			return nil, status.Errorf(codes.ResourceExhausted, "User rate limit exceeded. Wait %s before making another request", ttl)
		}

		//Validate user creds against the user service
		users, err := newUserClient()
		if err != nil {
			return nil, err
		}

		user, err := users.Find(withAuthContext(ctx), &userSvc.UserRequest{Query: &userSvc.UserRequest_Email{Email: username}}, grpc.Header(&md))
		if err != nil {
			grpc.SendHeader(ctx, metadata.Pairs("Grpc-Metadata-X-RateLimit-Remaining", fmt.Sprintf("%d", remaining+1)))
			events.BroadcastNonStreamingEvent(ctx, events.AuthenticateEvent{Event: &events.Event{Type: "io.evntsrc.passport.limite_increased"}, Err: fmt.Sprintf("limit increased"), IP: remoteIP.String(), User: username})
			incRateLimit(username, remoteIP)
			log.Printf("Unknown user %s @ %s", username, remoteIP)
			return &pb.AuthResponse{Success: false}, status.Errorf(codes.Unauthenticated, "unknown username or password")
		}

		if user.Mfa != nil && request.GetUserCreds().GetMFA() == "" {
			return validateMFA(user)
		}

		//Validate password hash
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.GetUserCreds().GetPassword())); err != nil {
			grpc.SendHeader(ctx, metadata.Pairs("Grpc-Metadata-X-RateLimit-Remaining", fmt.Sprintf("%d", remaining+1)))
			events.BroadcastNonStreamingEvent(ctx, events.AuthenticateEvent{Event: &events.Event{Type: "io.evntsrc.passport.limite_increased"}, Err: fmt.Sprintf("limit increased"), IP: remoteIP.String(), User: user.Id})
			incRateLimit(username, remoteIP)
			log.Printf("Password mismatch for user %s @ %s", username, remoteIP)
			return &pb.AuthResponse{Success: false}, status.Errorf(codes.Unauthenticated, "incorrect username or password")
		}

		// clearRateLimit(username, remoteIP)
		extraClaims = UserClaims(user)

	case *pb.AuthRequest_OauthClientSecretCreds:
	case *pb.AuthRequest_OAuthCodeCreds:
	default:
		events.BroadcastEvent(ctx, events.AuthenticateEvent{
			Event:    &events.Event{Type: "io.evntsrc.passport.authenticate"},
			AuthType: fmt.Sprintf("%v", authType),
			Success:  false,
			Err:      "unknown_type",
		})
		return nil, fmt.Errorf("Unknown auth type: %v", authType)
	}

	tokenString, token, _ := MakeNewToken(extraClaims)
	refreshToken := MakeNewRefresh()

	events.BroadcastEvent(ctx, events.AuthenticateEvent{
		Event:    &events.Event{Type: "io.evntsrc.passport.authenticate"},
		AuthType: fmt.Sprintf("%v", reflect.TypeOf(request.Creds)),
		Success:  true,
	})

	nounce, err := genNounce()
	if err != nil {
		return nil, fmt.Errorf("Failed to generate nounce")
	}

	cookie := http.Cookie{Name: "snounce", Value: nounce, Domain: "evntsrc.io", Path: "/", Expires: time.Now().Add(1 * time.Hour), Secure: true}
	grpc.SendHeader(ctx, metadata.Pairs("Grpc-Metadata-Set-Cookie", cookie.String()))

	log.Printf("Successful login for %s @ %s", token.Claims.(jwt.MapClaims)["sub"].(string), remoteIP)

	return &pb.AuthResponse{
		Success: true,
		Tokens: &pb.Tokens{
			Token:         *tokenString,
			TokenExpire:   &pb.Timestamp{Seconds: token.Claims.(jwt.MapClaims)["exp"].(int64)},
			RefreshToken:  *refreshToken,
			RefreshExpire: &pb.Timestamp{Seconds: time.Now().Add(time.Hour * 8).Unix()},
		},
	}, nil
}

//Refresh allows the creation of a auth token given a refresh token validating the refresh
//token has not been used before and has not expired
func (s *Server) Refresh(context.Context, *pb.RefreshRequest) (*pb.AuthResponse, error) {
	//TODO(tcfw)
	return &pb.AuthResponse{}, nil
}

//VerifyToken takes in a VerifyTokenRequest, validates the token in that request
func (s *Server) VerifyToken(ctx context.Context, request *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	if request.Token == "" {
		return nil, fmt.Errorf("Invalid token format")
	}

	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return GetKeyPublic()
	})

	if err != nil || !token.Valid {
		peer, _ := peer.FromContext(ctx)
		log.Printf("Invalid Token Used by '%s' due to '%s'", peer.Addr, err)
		return &pb.VerifyTokenResponse{
			Valid: false,
		}, nil
	}

	claims, _ := token.Claims.(jwt.MapClaims)

	isRevoked, err := isTokenRevoked(claims["jti"].(string))
	if err != nil {
		return nil, fmt.Errorf("Failed to validate token")
	}
	if isRevoked {
		return &pb.VerifyTokenResponse{
			Valid:   false,
			Revoked: true,
		}, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Failed to get metadata from request context")
	}
	snounces := md.Get("grpcgateway-snounce")
	if len(snounces) != 0 {
		snounce := snounces[0]
		newSnounce, err := validateSNounce(claims["jti"].(string), snounce)
		if err != nil {
			return nil, err
		}

		cookie := http.Cookie{Name: "snounce", Value: newSnounce, Domain: "evntsrc.io", Path: "/", Expires: time.Now().Add(1 * time.Hour), Secure: true}
		grpc.SendHeader(ctx, metadata.Pairs("Grpc-Metadata-Set-Cookie", cookie.String()))
	}

	return &pb.VerifyTokenResponse{
		Valid:       token.Valid,
		TokenExpire: &pb.Timestamp{Seconds: int64(claims["exp"].(float64))},
	}, nil
}

//SocialLogin validates remote idP tokens and creates users and passes back auth tokens
func (s *Server) SocialLogin(ctx context.Context, request *pb.SocialRequest) (*pb.AuthResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	remoteIP := rpcUtils.RemoteIPFromContext(ctx)

	info, err := getSocialInfo(request)

	if err != nil || info.Email == "" {
		return &pb.AuthResponse{Success: false}, fmt.Errorf("failed to login using provided tokens: %s", err)
	}

	users, err := newUserClient()
	if err != nil {
		return nil, err
	}

	user, err := users.Find(withAuthContext(ctx), &userSvc.UserRequest{Query: &userSvc.UserRequest_Email{Email: info.Email}}, grpc.Header(&md))
	if err != nil {
		incRateLimit(info.Email, remoteIP)
		return nil, status.Errorf(codes.Unauthenticated, "incorrect username or passport")
	}

	extraClaims := UserClaims(user)
	// clearRateLimit(info.Email, remoteIP)

	tokenString, token, _ := MakeNewToken(extraClaims)
	refreshToken := MakeNewRefresh()

	events.BroadcastEvent(ctx, events.AuthenticateEvent{
		Event:    &events.Event{Type: "io.evntsrc.passport.authenticate"},
		AuthType: request.GetProvider(),
		Success:  true,
		User:     user.Id,
		IP:       remoteIP.String(),
	})

	return &pb.AuthResponse{
		Success: true,
		Tokens: &pb.Tokens{
			Token:         *tokenString,
			TokenExpire:   &pb.Timestamp{Seconds: token.Claims.(jwt.MapClaims)["exp"].(int64)},
			RefreshToken:  *refreshToken,
			RefreshExpire: &pb.Timestamp{Seconds: time.Now().Add(time.Hour * 8).Unix()},
		},
	}, nil
}

//RevokeToken adds a token to the revoked tokens list
func (s *Server) RevokeToken(ctx context.Context, request *pb.Revoke) (*pb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to parse metadata")
	}

	auth := md.Get("grpcgateway-authorization")

	if auth == nil || len(auth) == 0 {
		auth = md.Get("authorization")
		if auth == nil || len(auth) == 0 {
			return nil, fmt.Errorf("no authorization sent")
		}
	}

	token := auth[0]

	resp, err := s.VerifyToken(ctx, &pb.VerifyTokenRequest{Token: token})
	if err != nil || resp.Revoked || !resp.Valid {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized")
	}

	claims, err := utils.TokenClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	request, err = s.populateRevokeRequestDefaults(request, claims)
	if err != nil {
		return nil, err
	}

	err = revokeToken(claims, request.Reason)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *Server) populateRevokeRequestDefaults(req *pb.Revoke, claims map[string]interface{}) (*pb.Revoke, error) {
	if req.Id == "" {
		req.Id = claims["sub"].(string)
	}

	if req.Jti == "" {
		req.Jti = claims["jti"].(string)
	}

	if _, ok := claims["admin"]; req.Id != claims["sub"].(string) && !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized")
	}

	if req.Reason == "" {
		req.Reason = "LOGOUT"
	}

	return req, nil
}

//RunGRPC starts the GRPC server
func RunGRPC(port int, tlsdir string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tracing.InitGlobalTracer("Passport")

	tlsKeyDir = tlsdir

	grpcServer := grpc.NewServer(tracing.GRPCServerOptions()...)
	pb.RegisterAuthSeviceServer(grpcServer, NewServer())

	log.Println("Starting gRPC server")
	grpcServer.Serve(lis)
}
