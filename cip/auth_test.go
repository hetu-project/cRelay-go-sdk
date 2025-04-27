package cip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthTag(t *testing.T) {
	// Test creating a new auth tag
	auth := NewAuthTag(ActionRead|ActionWrite, 30300, 1000)
	assert.Equal(t, ActionRead|ActionWrite, auth.Action)
	assert.Equal(t, uint32(30300), auth.Key)
	assert.Equal(t, uint64(1000), auth.Exp)

	// Test permission checks
	assert.True(t, auth.HasPermission(ActionRead))
	assert.True(t, auth.HasPermission(ActionWrite))
	assert.False(t, auth.HasPermission(ActionExecute))

	// Test expiration
	assert.False(t, auth.IsExpired(999))
	assert.True(t, auth.IsExpired(1000))
	assert.True(t, auth.IsExpired(1001))

	// Test string representation
	assert.Equal(t, "action=3,key=30300,exp=1000", auth.String())
}

func TestParseAuthTag(t *testing.T) {
	// Valid cases
	validTags := []string{
		"action=1,key=30300,exp=1000",
		"action=2,key=30301,exp=2000",
		"action=3,key=30302,exp=3000",
		"action=4,key=30303,exp=4000",
		"action=7,key=30304,exp=5000",
	}
	for _, tag := range validTags {
		auth, err := ParseAuthTag(tag)
		assert.NoError(t, err)
		assert.NotNil(t, auth)
	}

	// Invalid cases
	invalidTags := []string{
		"",                   // empty
		"action=1",           // missing fields
		"action=1,key=30300", // missing exp
		"action=1,key=30300,exp=1000,extra=value", // extra field
		"action=invalid,key=30300,exp=1000",       // invalid action
		"action=1,key=invalid,exp=1000",           // invalid key
		"action=1,key=30300,exp=invalid",          // invalid exp
		"action=1,key=30300,exp=1000,",            // trailing comma
		",action=1,key=30300,exp=1000",            // leading comma
	}
	for _, tag := range invalidTags {
		_, err := ParseAuthTag(tag)
		assert.Error(t, err)
	}
}
