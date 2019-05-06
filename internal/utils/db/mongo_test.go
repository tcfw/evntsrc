package db

import (
	"testing"

	mgo "github.com/globalsign/mgo"
	assert "github.com/stretchr/testify/assert"
)

func TestNewMongoDBSession(t *testing.T) {
	session, err := NewMongoDBSession()

	assert.NoError(t, err)
	assert.IsType(t, &mgo.Session{}, session)

	session.Close()
}
