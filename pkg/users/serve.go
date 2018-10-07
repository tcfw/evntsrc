package users

import (
	fmt "fmt"
	"log"
	"net"
	"regexp"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	protos "github.com/tcfw/evntsrc/pkg/users/protos"
	utils "github.com/tcfw/evntsrc/pkg/utils/authorization"
	"github.com/tcfw/evntsrc/pkg/utils/db"
	events "github.com/tcfw/evntsrc/pkg/utils/sysevents"
	"golang.org/x/crypto/bcrypt"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

const dbName = "users"
const dbCollection = "users"

type server struct {
	mu sync.Mutex // protects routeNotes
}

//server creates a ne struct to interface the auth server
func newServer() *server {
	return &server{}
}

//Create @TODO Validation
func (s *server) Create(ctx context.Context, request *protos.User) (*protos.User, error) {

	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)

	id := bson.NewObjectId().Hex()

	now := time.Now()

	request.CreatedAt = &now
	password, _ := validatePassword(request.Password)
	request.Password = *password
	request.Id = id

	err = collection.Insert(request)
	if err != nil {
		return nil, err
	}

	q := collection.FindId(id)
	if c, _ := q.Count(); c == 0 {
		return nil, fmt.Errorf("Failed to insert new user with id %s", id)
	}

	user := protos.User{}
	q.One(&user)

	events.BroadcastEvent(ctx, &events.UserEvent{Event: &events.Event{Type: "io.evntsrc.users.created"}, UserID: id})

	return &user, nil
}

//Delete TODO
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
		if err == nil {
			events.BroadcastEvent(ctx, &events.UserEvent{Event: &events.Event{Type: "io.evntsrc.users.deleted"}, UserID: user.Id})
		}
		return &protos.Empty{}, err
	case *protos.UserRequest_Email:
		err = collection.Remove(bson.M{"email": user.Email})
		if err == nil {
			events.BroadcastEvent(ctx, &events.UserEvent{Event: &events.Event{Type: "io.evntsrc.users.deleted"}, UserID: user.Id})
		}
		return &protos.Empty{}, err
	default:
		return &protos.Empty{}, fmt.Errorf("Unknown query type: %s", reqType)
	}
}

//Get TODO
func (s *server) Get(ctx context.Context, request *protos.UserRequest) (*protos.User, error) {
	user, err := s.Find(ctx, &protos.UserRequest{Query: &protos.UserRequest_Id{Id: request.GetId()}})
	if err != nil {
		return nil, err
	}

	user.Password = ""
	user.Id = bson.ObjectId(user.Id).Hex()
	return user, nil
}

//Find TODO
func (s *server) Find(ctx context.Context, request *protos.UserRequest) (*protos.User, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)

	var query *mgo.Query

	switch reqType := request.Query.(type) {
	case *protos.UserRequest_Id:
		query = collection.FindId(request.GetId())
	case *protos.UserRequest_Email:
		query = collection.Find(bson.M{"email": request.GetEmail()})
	default:
		return nil, fmt.Errorf("Unknown query type: %s", reqType)
	}

	if c, _ := query.Count(); c == 0 {
		log.Printf("No users found for query: %v", request)
		return nil, fmt.Errorf("Failed to find user")
	}

	user := protos.User{}
	query.One(&user)

	return &user, nil
}

//FindUsers TODO
func (s *server) FindUsers(request *protos.UserRequest, stream protos.UserService_FindUsersServer) error {
	return nil
}

//List TODO
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

//SetPassword TODO
func (s *server) SetPassword(ctx context.Context, request *protos.PasswordUpdateRequest) (*protos.Empty, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)

	query := collection.FindId(request.Id)

	if c, err := query.Count(); c == 0 || err != nil {
		return nil, fmt.Errorf("Failed to find user %s: %v", request.Id, err)
	}

	user := &protos.User{}

	err = query.One(user)
	if err != nil {
		return nil, err
	}

	password, _ := validatePassword(request.Password)
	user.Password = *password

	_, err = s.Update(ctx, &protos.UserUpdateRequest{Id: request.Id, User: user})
	return nil, err
}

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

//Update TODO
func (s *server) Update(ctx context.Context, request *protos.UserUpdateRequest) (*protos.User, error) {
	dbConn, err := db.NewMongoDBSession()
	if err != nil {
		return nil, err
	}
	defer dbConn.Close()

	collection := dbConn.DB(dbName).C(dbCollection)

	query := collection.FindId(request.Id)

	if c, err := query.Count(); c == 0 || err != nil {
		return nil, fmt.Errorf("Failed to find user %s: %v", request.Id, err)
	}

	err = collection.UpdateId(request.Id, request.User)
	if err != nil {
		return nil, err
	}

	query = collection.FindId(request.User.Id)
	if c, err := query.Count(); c == 0 || err != nil {
		return nil, fmt.Errorf("Failed to find user %s: %v", request.User.Id, err)
	}

	updatedUser := &protos.User{}
	err = query.One(updatedUser)
	if err != nil {
		return nil, err
	}

	events.BroadcastEvent(ctx, &events.UserEvent{Event: &events.Event{Type: "io.evntsrc.users.updated"}, UserID: request.User.Id})

	return updatedUser, nil
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

	user.Password = ""
	return user, nil
}

//RunGRPC starts the GRPC server
func RunGRPC(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	protos.RegisterUserServiceServer(grpcServer, newServer())

	log.Printf("Starting gRPC server (port %v)", port)
	grpcServer.Serve(lis)
}
