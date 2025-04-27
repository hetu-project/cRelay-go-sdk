package cip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCausalityKey(t *testing.T) {
	key := NewCausalityKey(30300, 1)
	assert.Equal(t, uint32(30300), key.Key)
	assert.Equal(t, uint64(1), key.Counter)
}

func TestCalculateSubspaceID(t *testing.T) {
	sid := CalculateSubspaceID("test-subspace", "post=30300", "energy>1000")
	assert.True(t, isValidSubspaceID(sid))
	assert.Equal(t, 66, len(sid)) // 0x + 64 hex chars
}

func TestValidateSubspaceID(t *testing.T) {
	// Valid cases
	validSIDs := []string{
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"0x0000000000000000000000000000000000000000000000000000000000000000",
		"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
	}
	for _, sid := range validSIDs {
		assert.NoError(t, ValidateSubspaceID(sid))
	}

	// Invalid cases
	invalidSIDs := []string{
		"",      // empty
		"0x",    // too short
		"0x123", // too short
		"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",   // missing 0x
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeg", // invalid hex char
	}
	for _, sid := range invalidSIDs {
		assert.Error(t, ValidateSubspaceID(sid))
	}
}
