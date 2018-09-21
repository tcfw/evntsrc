package passport

import (
	"testing"

	assert "github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	ten := RandomString(10)
	assert.Equal(t, 10, len(ten), "String length mismatch")
}
