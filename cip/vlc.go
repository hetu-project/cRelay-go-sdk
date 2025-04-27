package cip

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// CausalityKey represents a causality key with its identifier and counter
type CausalityKey struct {
	Key     uint32 // causality key identifier
	Counter uint64 // Lamport clock
}

// Subspacekey represents a subspace with its identifier and operation clock
type Subspacekey struct {
	SubspaceID uint32                  // Subspace Identifier (32 bits)
	Keys       map[uint32]CausalityKey // Subspace operation clock
}

// NewCausalityKey creates a new causality key
func NewCausalityKey(key uint32, counter uint64) CausalityKey {
	return CausalityKey{
		Key:     key,
		Counter: counter,
	}
}

// NewSubspace creates a new subspace with the given ID
func NewSubspace(subspaceID uint32) *Subspacekey {
	return &Subspacekey{
		SubspaceID: subspaceID,
		Keys:       make(map[uint32]CausalityKey),
	}
}

// AddKey adds a new causality key to the subspace
func (s *Subspacekey) AddKey(key CausalityKey) {
	s.Keys[key.Key] = key
}

// GetKey returns the causality key for the given key ID
func (s *Subspacekey) GetKey(keyID uint32) (CausalityKey, bool) {
	key, exists := s.Keys[keyID]
	return key, exists
}

// UpdateCounter updates the counter for the given key ID
func (s *Subspacekey) UpdateCounter(keyID uint32, counter uint64) {
	if key, exists := s.Keys[keyID]; exists {
		key.Counter = counter
		s.Keys[keyID] = key
	}
}

// CalculateSubspaceID generates a unique subspace ID based on subspace_name, ops, and rules
func CalculateSubspaceID(subspaceName, ops, rules string) string {
	// Concatenate the components
	input := subspaceName + ops + rules
	// Calculate SHA256 hash
	hash := sha256.Sum256([]byte(input))
	// Convert to hex string with "0x" prefix
	return "0x" + hex.EncodeToString(hash[:])
}

// ValidateSubspaceID validates the format of a subspace ID
func ValidateSubspaceID(sid string) error {
	if !isValidSubspaceID(sid) {
		return fmt.Errorf("invalid subspace ID format: %s", sid)
	}
	return nil
}

// isValidSubspaceID checks if a subspace ID is valid
func isValidSubspaceID(sid string) bool {
	if len(sid) != 66 { // 0x + 64 hex chars
		return false
	}
	if !isValidHexString(sid[2:]) {
		return false
	}
	return true
}

// isValidHexString checks if a string is a valid hex string
func isValidHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
