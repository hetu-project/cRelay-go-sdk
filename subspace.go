package nostr

import (
	jsonutils "encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nbd-wtf/go-nostr/cip"
)

// SubspaceCreateEvent represents a subspace creation event
type SubspaceCreateEvent struct {
	Event
	SubspaceID   string
	SubspaceName string
	Ops          string
	Rules        string
	Description  string
	ImageURL     string
}

// calculateSubspaceID generates a unique subspace ID based on subspace_name, ops, and rules
func calculateSubspaceID(subspaceName, ops, rules string) string {
	return cip.CalculateSubspaceID(subspaceName, ops, rules)
}

// NewSubspaceCreateEvent creates a new subspace creation event
func NewSubspaceCreateEvent(subspaceName, ops, rules, description, imageURL string) *SubspaceCreateEvent {
	// Calculate subspace ID
	sid := calculateSubspaceID(subspaceName, ops, rules)

	evt := &SubspaceCreateEvent{
		Event: Event{
			Kind:      cip.KindSubspaceCreate,
			CreatedAt: Timestamp(time.Now().Unix()),
		},
		SubspaceID:   sid,
		SubspaceName: subspaceName,
		Ops:          ops,
		Rules:        rules,
		Description:  description,
		ImageURL:     imageURL,
	}

	// Set tags
	evt.Tags = Tags{
		Tag{"d", cip.OpSubspaceCreate},
		Tag{"sid", sid},
		Tag{"subspace_name", subspaceName},
		Tag{"ops", ops},
	}
	if rules != "" {
		evt.Tags = append(evt.Tags, Tag{"rules", rules})
	}

	// Set content
	content := map[string]string{
		"desc":    description,
		"img_url": imageURL,
	}
	contentBytes, _ := jsonutils.Marshal(content)
	evt.Content = string(contentBytes)

	return evt
}

// ValidateSubspaceCreateEvent validates a SubspaceCreateEvent
func ValidateSubspaceCreateEvent(evt *SubspaceCreateEvent) error {
	// 1. Verify event kind
	if evt.Kind != cip.KindSubspaceCreate {
		return fmt.Errorf("invalid event kind: expected %d, got %d", cip.KindSubspaceCreate, evt.Kind)
	}

	// 2. Verify required tags
	requiredTags := map[string]bool{
		"d":             false,
		"sid":           false,
		"subspace_name": false,
		"ops":           false,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		if _, exists := requiredTags[tag[0]]; exists {
			requiredTags[tag[0]] = true
		}
	}

	// Check if all required tags are present
	for tag, found := range requiredTags {
		if !found {
			return fmt.Errorf("missing required tag: %s", tag)
		}
	}

	// 3. Verify sid matches the calculated hash
	calculatedSID := calculateSubspaceID(evt.SubspaceName, evt.Ops, evt.Rules)
	if evt.SubspaceID != calculatedSID {
		return fmt.Errorf("invalid subspace ID: expected %s, got %s", calculatedSID, evt.SubspaceID)
	}

	// 4. Verify content is valid JSON with required fields
	var content struct {
		Desc   string `json:"desc"`
		ImgURL string `json:"img_url"`
	}
	if err := jsonutils.Unmarshal([]byte(evt.Content), &content); err != nil {
		return fmt.Errorf("invalid content format: %v", err)
	}
	if content.Desc == "" {
		return fmt.Errorf("missing description in content")
	}

	// 5. Verify ops format
	// ops should be in format "key1=value1,key2=value2,..."
	opsParts := strings.Split(evt.Ops, ",")
	for _, part := range opsParts {
		if !strings.Contains(part, "=") {
			return fmt.Errorf("invalid ops format: %s", evt.Ops)
		}
	}

	return nil
}

// ParseSubspaceCreateEvent parses a raw Event into a SubspaceCreateEvent
func ParseSubspaceCreateEvent(evt Event) (*SubspaceCreateEvent, error) {
	// Create new SubspaceCreateEvent
	subspaceEvt := &SubspaceCreateEvent{
		Event: evt,
	}

	// Extract fields from tags
	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "sid":
			subspaceEvt.SubspaceID = tag[1]
		case "subspace_name":
			subspaceEvt.SubspaceName = tag[1]
		case "ops":
			subspaceEvt.Ops = tag[1]
		case "rules":
			subspaceEvt.Rules = tag[1]
		}
	}

	// Parse content
	var content struct {
		Desc   string `json:"desc"`
		ImgURL string `json:"img_url"`
	}
	if err := jsonutils.Unmarshal([]byte(evt.Content), &content); err != nil {
		return nil, fmt.Errorf("failed to parse content: %v", err)
	}
	subspaceEvt.Description = content.Desc
	subspaceEvt.ImageURL = content.ImgURL

	// Validate the parsed event
	if err := ValidateSubspaceCreateEvent(subspaceEvt); err != nil {
		return nil, fmt.Errorf("invalid subspace create event: %v", err)
	}

	return subspaceEvt, nil
}

// SubspaceJoinEvent represents a subspace join event
type SubspaceJoinEvent struct {
	Event
	SubspaceID string
}

// NewSubspaceJoinEvent creates a new subspace join event
func NewSubspaceJoinEvent(subspaceID string) *SubspaceJoinEvent {
	evt := &SubspaceJoinEvent{
		Event: Event{
			Kind:      cip.KindSubspaceJoin,
			CreatedAt: Timestamp(time.Now().Unix()),
		},
		SubspaceID: subspaceID,
	}

	evt.Tags = Tags{
		Tag{"d", cip.OpSubspaceJoin},
		Tag{"sid", subspaceID},
	}

	return evt
}

// ValidateSubspaceJoinEvent validates a SubspaceJoinEvent
func ValidateSubspaceJoinEvent(evt *SubspaceJoinEvent) error {
	// 1. Verify event kind
	if evt.Kind != cip.KindSubspaceJoin {
		return fmt.Errorf("invalid event kind: expected %d, got %d", cip.KindSubspaceJoin, evt.Kind)
	}

	// 2. Verify required tags
	requiredTags := map[string]bool{
		"d":   false,
		"sid": false,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		if _, exists := requiredTags[tag[0]]; exists {
			requiredTags[tag[0]] = true
		}
	}

	// Check if all required tags are present
	for tag, found := range requiredTags {
		if !found {
			return fmt.Errorf("missing required tag: %s", tag)
		}
	}

	// 3. Verify sid format (should be a valid hex string with 0x prefix)
	if err := cip.ValidateSubspaceID(evt.SubspaceID); err != nil {
		return err
	}

	return nil
}

// ParseSubspaceJoinEvent parses a raw Event into a SubspaceJoinEvent
func ParseSubspaceJoinEvent(evt Event) (*SubspaceJoinEvent, error) {
	joinEvt := &SubspaceJoinEvent{
		Event: evt,
	}

	// Extract fields from tags
	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		if tag[0] == "sid" {
			joinEvt.SubspaceID = tag[1]
		}
	}

	// Validate the parsed event
	if err := ValidateSubspaceJoinEvent(joinEvt); err != nil {
		return nil, fmt.Errorf("invalid subspace join event: %v", err)
	}

	return joinEvt, nil
}

// SubspaceOpEventPtr represents a governance subspace event
type SubspaceOpEventPtr interface {
	GetSubspaceID() string
	GetOperation() string
	GetAuthTag() cip.AuthTag
}

// SubspaceOpEvent represents a subspace operation event
type SubspaceOpEvent struct {
	Event
	SubspaceID string
	Operation  string
	AuthTag    cip.AuthTag
	Parents    []string
}

func (e *SubspaceOpEvent) GetSubspaceID() string   { return e.SubspaceID }
func (e *SubspaceOpEvent) GetOperation() string    { return e.Operation }
func (e *SubspaceOpEvent) GetAuthTag() cip.AuthTag { return e.AuthTag }

// NewSubspaceOpEvent creates a new subspace operation event
func NewSubspaceOpEvent(subspaceID string, kind int) (*SubspaceOpEvent, error) {
	operation, exist := cip.GetOpFromKind(kind)
	if !exist {
		return nil, errors.New("Not existed operation!")
	}
	evt := &SubspaceOpEvent{
		Event: Event{
			Kind:      kind,
			CreatedAt: Timestamp(time.Now().Unix()),
		},
		SubspaceID: subspaceID,
		Operation:  operation,
	}

	evt.Tags = Tags{
		Tag{"d", "subspace_op"},
		Tag{"sid", subspaceID},
		Tag{"op", operation},
	}

	return evt, nil
}

// SetAuth sets the auth tag for the operation
func (e *SubspaceOpEvent) SetAuth(action cip.Action, key uint32, exp uint64) {
	e.AuthTag = cip.NewAuthTag(action, key, exp)
	e.Tags = append(e.Tags, Tag{"auth", e.AuthTag.String()})
}

// SetParents sets the parent event hash
func (e *SubspaceOpEvent) SetParents(parentHashSet []string) {
	e.Parents = parentHashSet
	parents := Tag{"parent"}
	parents = append(parents, parentHashSet...)
	e.Tags = append(e.Tags, parents)
}