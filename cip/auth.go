package cip

import (
	"fmt"
	"strconv"
	"strings"
)

// Action represents the permission action type
type Action uint8

const (
	ActionRead    Action = 1 << iota // 1
	ActionWrite                      // 2
	ActionExecute                    // 4
)

// AuthTag represents the auth tag structure
type AuthTag struct {
	Action Action // Permission mask (1=read, 2=write, 4=execute)
	Key    uint32 // Causality key ID
	Exp    uint64 // Expiration clock value
}

// NewAuthTag creates a new auth tag
func NewAuthTag(action Action, key uint32, exp uint64) AuthTag {
	return AuthTag{
		Action: action,
		Key:    key,
		Exp:    exp,
	}
}

// ParseAuthTag parses an auth tag string into an AuthTag struct
func ParseAuthTag(authStr string) (AuthTag, error) {
	parts := strings.Split(authStr, ",")
	if len(parts) != 3 {
		return AuthTag{}, fmt.Errorf("invalid auth tag format: %s", authStr)
	}

	var auth AuthTag
	for _, part := range parts {
		kv := strings.Split(part, "=")
		if len(kv) != 2 {
			return AuthTag{}, fmt.Errorf("invalid auth tag part: %s", part)
		}

		switch kv[0] {
		case "action":
			action, err := strconv.ParseUint(kv[1], 10, 8)
			if err != nil {
				return AuthTag{}, fmt.Errorf("invalid action value: %s", kv[1])
			}
			auth.Action = Action(action)
		case "key":
			key, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return AuthTag{}, fmt.Errorf("invalid key value: %s", kv[1])
			}
			auth.Key = uint32(key)
		case "exp":
			exp, err := strconv.ParseUint(kv[1], 10, 64)
			if err != nil {
				return AuthTag{}, fmt.Errorf("invalid exp value: %s", kv[1])
			}
			auth.Exp = exp
		default:
			return AuthTag{}, fmt.Errorf("unknown auth tag field: %s", kv[0])
		}
	}

	return auth, nil
}

// String returns the string representation of an AuthTag
func (a AuthTag) String() string {
	return fmt.Sprintf("action=%d,key=%d,exp=%d", a.Action, a.Key, a.Exp)
}

// HasPermission checks if the auth tag has the specified permission
func (a AuthTag) HasPermission(action Action) bool {
	return a.Action&action != 0
}

// IsExpired checks if the auth tag is expired based on the current clock value
func (a AuthTag) IsExpired(currentClock uint64) bool {
	return a.Exp <= currentClock
}
