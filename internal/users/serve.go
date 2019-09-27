package users

import (
	"errors"
	fmt "fmt"
	"log"
	"net"
	"regexp"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gogo/protobuf/types"

	emails "github.com/tcfw/evntsrc/internal/emails/protos"
	"github.com/tcfw/evntsrc/internal/tracing"
	"github.com/tcfw/evntsrc/internal/users/events"
	protos "github.com/tcfw/evntsrc/internal/users/protos"
	utils "github.com/tcfw/evntsrc/internal/utils/authorization"
	"github.com/tcfw/evntsrc/internal/utils/db"
	"github.com/tcfw/evntsrc/internal/utils/sysevents"
	"golang.org/x/crypto/bcrypt"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	dbName       = "users"
	dbCollection = "users"

	mdValidationtoken = "validation_token"
)

type server struct {
	mu         sync.Mutex // protects routeNotes
	emails     emails.EmailServiceClient
	emailsConn *grpc.ClientConn
}

//server creates a ne struct to interface the auth server
func newServer() *server {
	svr := &server{}

	emailsConn, err := grpc.Dial("dns:///emails:443", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		panic(err)
	}

	svr.emailsConn = emailsConn
	svr.emails = emails.NewEmailServiceClient(emailsConn)

	go svr.listenLoop()

	return svr
}

//Create takes in a request and creates a new user
func (s *server) Create(ctx context.Context, request *protos.User) (*protos.User, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)

	if err := s.validateUserRequest(request, true); err != nil {
		return nil, err
	}

	password, err := validatePassword(request.Password)
	if err != nil {
		return nil, err
	}

	id := bson.NewObjectId().Hex()

	request.CreatedAt = types.TimestampNow()
	request.Password = *password
	request.Id = id

	//Clear anything that was sent in reqest
	request.Metadata = map[string][]byte{
		mdValidationtoken: []byte(genValidationToken()),
	}
	request.Mfa = nil
	request.Status = protos.User_PENDING

	//check user already exists
	if existingUser, err := s.Find(ctx, &protos.UserRequest{Status: protos.UserRequest_ANY, Query: &protos.UserRequest_Email{Email: request.Email}}); err == nil {
		sysevents.BroadcastEvent(ctx, &events.Event{Event: &sysevents.Event{Type: events.BroadcastTypeRecreation}, UserID: existingUser.Id})

		if err := s.sendRecreationEmail(ctx, existingUser); err != nil {
			log.Printf("Failed to send recreation email: %s", err)
			return nil, err
		}

		//Return with request to avoid account enumeration
		request.Metadata = map[string][]byte{}
		return s.sanatiseUserForResponse(request), nil
	}

	if err = collection.Insert(request); err != nil {
		return nil, err
	}

	//Validate insert
	q := collection.FindId(id)
	if c, _ := q.Count(); c == 0 {
		return nil, fmt.Errorf("Failed to insert new user with id %s", id)
	}

	user := protos.User{}
	q.One(&user)

	sysevents.BroadcastEvent(ctx, &events.Event{Event: &sysevents.Event{Type: events.BroadcastTypeCreated}, UserID: id})
	if err := s.sendValidationEmail(ctx, &user); err != nil {
		log.Printf("Failed to send validation email: %s", err)
		return nil, err
	}

	user.Metadata = map[string][]byte{}
	return s.sanatiseUserForResponse(&user), nil
}

//sanatiseUserForResponse removes sensitive information from the user struct
//before sending back as response
func (s *server) sanatiseUserForResponse(user *protos.User) *protos.User {
	user.Password = ""
	return user
}

//validateUserRequest checks various fields in request
func (s *server) validateUserRequest(request *protos.User, checkPassword bool) error {
	if request.Name == "" {
		return fmt.Errorf("Name is required")
	}

	if request.Email == "" {
		return fmt.Errorf("Email is required")
	}

	if checkPassword && request.Password == "" {
		return fmt.Errorf("Password is required")
	}

	//picture should not be larger than 5MB
	if len(request.Picture) > 5*1<<20 {
		return fmt.Errorf("Picture too large - 5MB limit")
	}

	return nil
}

func (s *server) ValidateAccount(ctx context.Context, request *protos.ValidateRequest) (*protos.Empty, error) {
	notFound := fmt.Errorf("Unable to find matching token")

	user, err := s.Find(ctx, &protos.UserRequest{Status: protos.UserRequest_PENDING, Query: &protos.UserRequest_Email{Email: request.Email}})
	if err != nil {
		return nil, notFound
	}

	token, ok := user.Metadata[mdValidationtoken]
	if !ok {
		return nil, notFound
	}
	if string(token) == request.Token {
		user.Status = protos.User_ACTIVE
		delete(user.Metadata, mdValidationtoken)
	} else {
		return nil, notFound
	}

	err = s.updateUser(user)

	return &protos.Empty{}, err
}

//Delete deletes a user
func (s *server) Delete(ctx context.Context, request *protos.UserRequest) (*protos.Empty, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return &protos.Empty{}, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)
	user, err := s.Find(ctx, request)
	if err != nil {
		return &protos.Empty{}, err
	}

	switch reqType := request.Query.(type) {
	case *protos.UserRequest_Id:
		err = collection.RemoveId(user.Id)
	case *protos.UserRequest_Email:
		err = collection.Remove(bson.M{"email": user.Email})
	default:
		return &protos.Empty{}, fmt.Errorf("Unknown query type: %s", reqType)
	}

	if err == nil {
		sysevents.BroadcastEvent(ctx, &events.Event{Event: &sysevents.Event{Type: events.BroadcastTypeDeleted}, UserID: user.Id})
	}

	return &protos.Empty{}, err
}

//Get finds a user
func (s *server) Get(ctx context.Context, request *protos.UserRequest) (*protos.User, error) {
	user, err := s.Find(ctx, &protos.UserRequest{Query: &protos.UserRequest_Id{Id: request.GetId()}})
	if err != nil {
		return nil, err
	}

	return user, nil
}

//Find queries for a user based on email or id depending on the id
func (s *server) Find(ctx context.Context, request *protos.UserRequest) (*protos.User, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)

	var query *mgo.Query

	qBsonM := bson.M{}

	switch reqType := request.Query.(type) {
	case *protos.UserRequest_Id:
		qBsonM["_id"] = request.GetId()
	case *protos.UserRequest_Email:
		qBsonM["email"] = request.GetEmail()
	default:
		return nil, fmt.Errorf("Unknown query type: %s", reqType)
	}

	switch request.Status {
	default:
	case protos.UserRequest_ACTIVE:
		qBsonM["status"] = protos.User_ACTIVE
		break
	case protos.UserRequest_DELETED:
		qBsonM["status"] = protos.User_DELETED
		break
	case protos.UserRequest_PENDING:
		qBsonM["status"] = protos.User_PENDING
		break
	case protos.UserRequest_ANY:
		//don't filter
		break
	}

	query = collection.Find(qBsonM)

	if c, _ := query.Count(); c == 0 {
		log.Printf("No users found for query: %v", request)
		return nil, fmt.Errorf("Failed to find user")
	}

	user := protos.User{}
	query.One(&user)

	tracing.ActiveSpan(ctx).LogKV("user_id", user.Id)

	//Cannot sanatise output as password is required for passport svr
	return &user, nil
}

//FindUsers finds multiple users
func (s *server) FindUsers(request *protos.UserRequest, stream protos.UserService_FindUsersServer) error {
	//TODO(tcfw) is this different to List now?
	return nil
}

//List creates a list of all users in the DB
func (s *server) List(ctx context.Context, request *protos.Empty) (*protos.UserList, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)

	//Find All
	query := collection.Find(nil)

	list := &protos.UserList{}
	if err = query.All(&list.Users); err != nil {
		return nil, err
	}

	return list, nil
}

//SetPassword changes the users password
func (s *server) SetPassword(ctx context.Context, request *protos.PasswordUpdateRequest) (*protos.Empty, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)

	claims, err := utils.TokenClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if request.Id == "" {
		request.Id = claims["sub"].(string)
	} else if _, ok := claims["admin"]; request.Id != claims["sub"].(string) && !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized")
	}

	query := collection.FindId(request.Id)

	if c, err := query.Count(); c == 0 || err != nil {
		return nil, fmt.Errorf("Failed to find user %s: %v", request.Id, err)
	}

	user := &protos.User{}
	if err = query.One(user); err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.GetCurrentPassword())); err != nil {
		return nil, errors.New("Current password doesn't match")
	}

	password, _ := validatePassword(request.Password)

	collection.UpdateId(user.Id, bson.M{"$set": bson.M{"password": password}})

	return nil, err
}

//validatePassword checks for password complexity
func validatePassword(password string) (*string, error) {

	if password == "" {
		return &password, nil
	}

	regex, err := regexp.Compile("^\\$2[aby]?\\$\\d{1,2}\\$[.\\/A-Za-z0-9]{53}$")
	if err != nil {
		return nil, err
	}
	if regex.Match([]byte(password)) {
		return &password, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	stringHash := string(hash)
	return &stringHash, nil
}

//Update updates a users details given they exist
func (s *server) Update(ctx context.Context, request *protos.UserUpdateRequest) (*protos.User, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)

	claims, err := utils.TokenClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.validateUserRequest(request.User, false); err != nil {
		return nil, err
	}

	if request.Id == "" {
		request.Id = claims["sub"].(string)
	}
	if _, ok := claims["admin"]; request.Id != claims["sub"].(string) && !ok {
		return nil, status.Errorf(codes.PermissionDenied, "Unauthorized")
	}

	query := collection.FindId(request.Id)

	if c, err := query.Count(); c == 0 || err != nil {
		return nil, fmt.Errorf("Failed to find user %s: %v", request.Id, err)
	}

	if _, ok := claims["admin"]; request.Replace == true && ok {
		request.User.Id = request.Id
		err = collection.UpdateId(request.Id, request.User)
		if err != nil {
			return nil, err
		}
	} else {
		user := &protos.User{}
		if err = query.One(user); err != nil {
			return nil, err
		}

		user.Name = request.User.Name
		user.Email = request.User.Email
		user.Company = request.User.Company
		for k, v := range request.User.GetMetadata() {
			user.Metadata[k] = v
		}

		if err = collection.UpdateId(request.Id, user); err != nil {
			return nil, err
		}
	}

	sysevents.BroadcastEvent(ctx, &events.Event{Event: &sysevents.Event{Type: events.BroadcastTypeUpdated}, UserID: request.User.Id})

	return s.Get(ctx, &protos.UserRequest{Query: &protos.UserRequest_Id{Id: request.GetId()}})
}

func (s *server) updateUser(user *protos.User) error {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)
	query := collection.FindId(user.Id)

	if c, err := query.Count(); c == 0 || err != nil {
		return fmt.Errorf("Failed to find user %s: %v", user.Id, err)
	}

	//restore password
	origUser := &protos.User{}
	query.One(origUser)
	user.Password = origUser.Password

	return collection.UpdateId(user.Id, user)
}

//Me returns the user information based on the uid from a gwt token provided in the authorization metadata
func (s *server) Me(ctx context.Context, request *protos.Empty) (*protos.User, error) {
	claims, err := utils.TokenClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.Find(ctx, &protos.UserRequest{Query: &protos.UserRequest_Id{Id: claims["sub"].(string)}})
	if err != nil {
		return nil, err
	}

	return s.sanatiseUserForResponse(user), nil
}

//RunGRPC starts the GRPC server
func RunGRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	tracing.InitGlobalTracer("Users")

	grpcServer := grpc.NewServer(tracing.GRPCServerOptions()...)
	protos.RegisterUserServiceServer(grpcServer, newServer())

	log.Printf("Starting gRPC server (port %v)", port)
	grpcServer.Serve(lis)
}
