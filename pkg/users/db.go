package users

import (
	"os"

	"github.com/globalsign/mgo"
)

//NewDBSession returns a mongoDB session
func NewDBSession() (*mgo.Session, error) {
	dbConnHost, exists := os.LookupEnv("DB_HOST")
	if exists == false {
		dbConnHost = "localhost:27017"
	}

	session, err := mgo.Dial(dbConnHost)
	if err != nil {
		return nil, err
	}

	session.SetMode(mgo.Monotonic, true)

	return session, nil
}
