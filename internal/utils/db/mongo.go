package db

import (
	"os"

	"github.com/globalsign/mgo"
)

//NewMongoDBSession returns a mongoDB session
func NewMongoDBSession() (*mgo.Session, error) {
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
