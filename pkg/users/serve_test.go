package users

import (
	"context"
	"log"
	"testing"

	fake "github.com/icrowley/fake"
	assert "github.com/stretchr/testify/assert"
	"github.com/tcfw/evntsrc/pkg/passport"
	protos "github.com/tcfw/evntsrc/pkg/users/protos"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestCanRunServer(t *testing.T) {
	go RunGRPC(54765)

	_, err := grpc.Dial("localhost:54765", grpc.WithInsecure())
	assert.NoError(t, err)
}

func TestCreate(t *testing.T) {
	s := server{}

	mockUser := &protos.User{
		Name:  fake.FullName(),
		Email: fake.EmailAddress(),
	}

	user, err := s.Create(context.Background(), mockUser)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, mockUser.Name, user.Name)
	assert.Equal(t, mockUser.Email, user.Email)
	assert.NotEmpty(t, user.Id)

	//Clean up - manually delete created user
	session, err := NewDBSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	session.DB(dbName).C(dbCollection).RemoveId(user.Id)
}

func TestFind(t *testing.T) {
	s := server{}

	mockUser := &protos.User{
		Name:  fake.FullName(),
		Email: fake.EmailAddress(),
	}

	user, err := s.Create(context.Background(), mockUser)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Find by ID", func(t *testing.T) {
		fuser, err := s.Find(context.Background(), &protos.UserRequest{Query: &protos.UserRequest_Id{Id: user.Id}})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, user.Name, fuser.Name)
	})

	t.Run("Find by Email", func(t *testing.T) {
		fuser, err := s.Find(context.Background(), &protos.UserRequest{Query: &protos.UserRequest_Email{Email: user.Email}})
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, user.Email, fuser.Email)
	})

	//Clean up - manually delete created user
	session, err := NewDBSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	session.DB(dbName).C(dbCollection).RemoveId(user.Id)
}

func TestFindUsers(t *testing.T) {
	// s := server{}

}

func TestList(t *testing.T) {
	s := server{}

	mockUser := &protos.User{
		Name:  fake.FullName(),
		Email: fake.EmailAddress(),
	}

	session, err := NewDBSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	//Make sure collection is empty
	session.DB(dbName).C(dbCollection).DropCollection()

	user, err := s.Create(context.Background(), mockUser)
	if err != nil {
		t.Fatal(err)
	}

	list, err := s.List(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	if assert.NotNil(t, list) {
		assert.Len(t, list.Users, 1)
	}

	//Clean up - manually delete created user
	session.DB(dbName).C(dbCollection).RemoveId(user.Id)
}

func TestDelete(t *testing.T) {
	s := server{}

	mockUser := &protos.User{
		Name:  fake.FullName(),
		Email: fake.EmailAddress(),
	}

	session, err := NewDBSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	//Make sure collection is empty
	session.DB(dbName).C(dbCollection).DropCollection()

	t.Run("Delete by ID", func(t *testing.T) {

		user, err := s.Create(context.Background(), mockUser)
		if err != nil {
			t.Fatal(err)
		}

		s.Delete(context.Background(), &protos.UserRequest{Query: &protos.UserRequest_Id{Id: user.Id}})

		query := session.DB(dbName).C(dbCollection).FindId(user.Id)

		count, err := query.Count()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, count)
	})

	t.Run("Delete by email", func(t *testing.T) {
		user, err := s.Create(context.Background(), mockUser)
		if err != nil {
			t.Fatal(err)
		}

		s.Delete(context.Background(), &protos.UserRequest{Query: &protos.UserRequest_Email{Email: user.Email}})

		query := session.DB(dbName).C(dbCollection).FindId(user.Id)

		count, err := query.Count()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, 0, count)
	})

}

func TestCanSetPassword(t *testing.T) {
	s := server{}

	mockUser := &protos.User{
		Name:  fake.FullName(),
		Email: fake.EmailAddress(),
	}

	session, err := NewDBSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	//Create mock user
	user, err := s.Create(context.Background(), mockUser)
	assert.NoError(t, err)

	t.Run("as plain text", func(t *testing.T) {

		//Set password
		req := &protos.PasswordUpdateRequest{
			Id:       user.Id,
			Password: "test1234",
		}

		_, err = s.SetPassword(context.Background(), req)
		assert.NoError(t, err)

		//Verify
		vuser, err := s.Find(context.Background(), &protos.UserRequest{Query: &protos.UserRequest_Id{Id: user.Id}})
		assert.NoError(t, err)
		assert.NotEmpty(t, vuser.Password)
	})

	t.Run("as bcrypt", func(t *testing.T) {
		//Make hash:
		hash, err := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
		assert.NoError(t, err)

		//Set password
		req := &protos.PasswordUpdateRequest{
			Id:       user.Id,
			Password: string(hash),
		}

		_, err = s.SetPassword(context.Background(), req)
		assert.NoError(t, err)

		//Verify
		vuser, err := s.Find(context.Background(), &protos.UserRequest{Query: &protos.UserRequest_Id{Id: user.Id}})
		assert.NoError(t, err)
		assert.Equal(t, string(hash), vuser.Password)
	})

	//Cleanup
	session.DB(dbName).C(dbCollection).RemoveId(user.Id)
}

func TestUpdate(t *testing.T) {
	s := server{}

	mockUser := &protos.User{
		Name:  fake.FullName(),
		Email: fake.EmailAddress(),
	}

	session, err := NewDBSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	//Create mock user
	user, err := s.Create(context.Background(), mockUser)
	assert.NoError(t, err)

	//Modify user
	user.Name = fake.FullName()

	req := &protos.UserUpdateRequest{
		Id:   user.Id,
		User: user,
	}

	updatedUser, err := s.Update(context.Background(), req)
	assert.Equal(t, updatedUser.Name, user.Name)

	//Cleanup
	session.DB(dbName).C(dbCollection).RemoveId(updatedUser.Id)
}

func TestMe(t *testing.T) {
	s := server{}

	mockUser := &protos.User{
		Name:  "John Smith",
		Email: "johnsmith@example.com",
	}

	session, err := NewDBSession()
	if err != nil {
		t.Error(err)
	}
	defer session.Close()

	//Create mock user
	user, err := s.Create(context.Background(), mockUser)
	assert.NoError(t, err)

	testToken, err := passport.MakeTestToken(user)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Test token used: %s", testToken)

	md := metadata.New(map[string]string{"authorization": testToken})

	testContext := metadata.NewIncomingContext(context.Background(), md)

	me, err := s.Me(testContext, &protos.Empty{})
	if err != nil {
		t.Error(err)
	}
	if me == nil {
		t.Fatalf("Response is empty: %v", me)
	}
	assert.Contains(t, me.Email, mockUser.Email, "Returned user does not match expected user")

	//Cleanup
	session.DB(dbName).C(dbCollection).RemoveId(user.Id)
}
