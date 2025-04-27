package nostr

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/nbd-wtf/go-nostr/cip"
)

func TestSubspaceCreateEvent(t *testing.T) {
	// Test creating a subspace with basic operations
	createEvent := NewSubspaceCreateEvent(
		"test-subspace",
		cip.DefaultSubspaceOps,
		"energy>1000",
		"Test Subspace",
		"https://example.com/image.png",
	)

	// Verify event kind
	assert.Equal(t, cip.KindSubspaceCreate, createEvent.Kind)

	// Verify required tags
	requiredTags := map[string]bool{
		"d":             false,
		"sid":           false,
		"subspace_name": false,
		"ops":           false,
	}

	for _, tag := range createEvent.Tags {
		if len(tag) < 2 {
			continue
		}
		if _, exists := requiredTags[tag[0]]; exists {
			requiredTags[tag[0]] = true
		}
	}

	for tag, found := range requiredTags {
		assert.True(t, found, "missing required tag: %s", tag)
	}

	// Verify sid matches calculated hash
	calculatedSID := calculateSubspaceID(createEvent.SubspaceName, createEvent.Ops, createEvent.Rules)
	assert.Equal(t, calculatedSID, createEvent.SubspaceID)

	// Test validation
	err := ValidateSubspaceCreateEvent(createEvent)
	assert.NoError(t, err)

	// Test parsing
	parsedEvent, err := ParseSubspaceCreateEvent(createEvent.Event)
	assert.NoError(t, err)
	assert.Equal(t, createEvent.SubspaceID, parsedEvent.SubspaceID)
	assert.Equal(t, createEvent.SubspaceName, parsedEvent.SubspaceName)
	assert.Equal(t, createEvent.Ops, parsedEvent.Ops)
	assert.Equal(t, createEvent.Rules, parsedEvent.Rules)
}

func TestSubspaceJoinEvent(t *testing.T) {
	// Create a subspace first to get a valid sid
	createEvent := NewSubspaceCreateEvent(
		"test-subspace",
		cip.DefaultSubspaceOps,
		"energy>1000",
		"Test Subspace",
		"https://example.com/image.png",
	)

	// Test creating a join event
	joinEvent := NewSubspaceJoinEvent(createEvent.SubspaceID)

	// Verify event kind
	assert.Equal(t, cip.KindSubspaceJoin, joinEvent.Kind)

	// Verify required tags
	requiredTags := map[string]bool{
		"d":   false,
		"sid": false,
	}

	for _, tag := range joinEvent.Tags {
		if len(tag) < 2 {
			continue
		}
		if _, exists := requiredTags[tag[0]]; exists {
			requiredTags[tag[0]] = true
		}
	}

	for tag, found := range requiredTags {
		assert.True(t, found, "missing required tag: %s", tag)
	}

	// Test validation
	err := ValidateSubspaceJoinEvent(joinEvent)
	assert.NoError(t, err)

	// Test parsing
	parsedEvent, err := ParseSubspaceJoinEvent(joinEvent.Event)
	assert.NoError(t, err)
	assert.Equal(t, joinEvent.SubspaceID, parsedEvent.SubspaceID)
}