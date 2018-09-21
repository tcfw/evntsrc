package users

import (
	"testing"

	mgo "github.com/globalsign/mgo"
	assert "github.com/stretchr/testify/assert"
)

func TestNewDBSession(t *testing.T) {
	session, err := NewDBSession()

	assert.NoError(t, err)
	assert.IsType(t, &mgo.Session{}, session)

	session.Close()
}
