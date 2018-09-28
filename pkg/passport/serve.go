package passport

import (
	fmt "fmt"
	"log"
	"net"
	"os"
	"reflect"
	"sync"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	pb "github.com/tcfw/evntsrc/pkg/passport/protos"
	userSvc "github.com/tcfw/evntsrc/pkg/users/protos"
	rpcUtils "github.com/tcfw/evntsrc/pkg/utils/rpc"
	events "github.com/tcfw/evntsrc/pkg/utils/sysevents"
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

//ReqAuth TODO
type ReqAuth struct {
	Token string
}

//GetRequestMetadata TODO
func (a *ReqAuth) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": a.Token,
	}, nil
}

//RequireTransportSecurity TODO
func (a *ReqAuth) RequireTransportSecurity() bool {
	return false
}

func newUserClient(ctx context.Context) (userSvc.UserServiceClient, error) {

	md, _ := metadata.FromIncomingContext(ctx)
	authReq := ReqAuth{}

	auth := md.Get("authorization")
	if auth != nil {
		authReq.Token = auth[0]
	}

	userEndpoint, envExists := os.LookupEnv("USER_HOST")
	if envExists != true {
		userEndpoint = "users:443"
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}

	opts = append(opts, grpc.WithPerRPCCredentials(&authReq))

	conn, err := grpc.DialContext(ctx, userEndpoint, opts...)
	if err != nil {
		return nil, err
	}

	return userSvc.NewUserServiceClient(conn), nil
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
		users, err := newUserClient(ctx)
		user, err := users.Find(ctx, &userSvc.UserRequest{Query: &userSvc.UserRequest_Email{Email: username}}, grpc.Header(&md))
		if err != nil {
			grpc.SendHeader(ctx, metadata.Pairs("Grpc-Metadata-X-RateLimit-Remaining", fmt.Sprintf("%d", remaining+1)))
			events.BroadcastNonStreamingEvent(ctx, events.AuthenticateEvent{Event: &events.Event{Type: "io.evntsrc.passport.limite_increased"}, Err: fmt.Sprintf("limit increased"), IP: remoteIP.String(), User: username})
			incRateLimit(username, remoteIP)
			return &pb.AuthResponse{Success: false}, status.Errorf(codes.Unauthenticated, "unknown username or password")
		}

		//Validate password hash
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.GetUserCreds().GetPassword()))
		if err != nil {
			grpc.SendHeader(ctx, metadata.Pairs("Grpc-Metadata-X-RateLimit-Remaining", fmt.Sprintf("%d", remaining+1)))
			events.BroadcastNonStreamingEvent(ctx, events.AuthenticateEvent{Event: &events.Event{Type: "io.evntsrc.passport.limite_increased"}, Err: fmt.Sprintf("limit increased"), IP: remoteIP.String(), User: user.Id})
			incRateLimit(username, remoteIP)
			return &pb.AuthResponse{Success: false}, status.Errorf(codes.Unauthenticated, "incorrect username or password")
		}

		clearRateLimit(username, remoteIP)
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

//Refresh TODO
func (s *Server) Refresh(context.Context, *pb.RefreshRequest) (*pb.AuthResponse, error) {
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
		return &pb.AuthResponse{Success: false}, fmt.Errorf("failed to login using provided tokens")
	}

	users, err := newUserClient(ctx)
	user, err := users.Find(ctx, &userSvc.UserRequest{Query: &userSvc.UserRequest_Email{Email: info.Email}}, grpc.Header(&md))
	if err != nil {
		incRateLimit(info.Email, remoteIP)
		return nil, status.Errorf(codes.Unauthenticated, "incorrect username or passport")
	}

	extraClaims := UserClaims(user)
	clearRateLimit(info.Email, remoteIP)

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

//RunGRPC starts the GRPC server
func RunGRPC(port int, tlsdir string) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tlsKeyDir = tlsdir

	grpcServer := grpc.NewServer()
	pb.RegisterAuthSeviceServer(grpcServer, NewServer())

	log.Println("Starting gRPC server")
	grpcServer.Serve(lis)
}