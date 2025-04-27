package cip02

import (
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
)

// ModelEvent represents a model operation in modelgraph subspace
type ModelEvent struct {
	*nostr.SubspaceOpEvent
	ParentHash    string
	Contributions string
	Content       string
}

// SetContributions sets the contribution weights
func (e *ModelEvent) SetContributions(contributions string) {
	e.Contributions = contributions
	e.Tags = append(e.Tags, nostr.Tag{"contrib", contributions})
}

// SetParent sets the parent event hash
func (e *ModelEvent) SetParent(parentHash string) {
	e.ParentHash = parentHash
	e.Tags = append(e.Tags, nostr.Tag{"parent", parentHash})
}

// DataEvent represents a data operation in modelgraph subspace
type DataEvent struct {
	*nostr.SubspaceOpEvent
	Size    string
	Content string
}

// ComputeEvent represents a compute operation in modelgraph subspace
type ComputeEvent struct {
	*nostr.SubspaceOpEvent
	ComputeType string
	Content     string
}

// AlgoEvent represents an algo operation in modelgraph subspace
type AlgoEvent struct {
	*nostr.SubspaceOpEvent
	AlgoType string
	Content  string
}

// ValidEvent represents a valid operation in modelgraph subspace
type ValidEvent struct {
	*nostr.SubspaceOpEvent
	ValidResult string
	Content     string
}

// ParseModelGraphEvent parses a Nostr event into a modelgraph event
func ParseModelGraphEvent(evt nostr.Event) (nostr.SubspaceOpEventPtr, error) {
	// Extract common fields
	subspaceID := ""
	var authTag cip.AuthTag

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "sid":
			subspaceID = tag[1]
		case "auth":
			auth, err := cip.ParseAuthTag(tag[1])
			if err != nil {
				return nil, fmt.Errorf("failed to parse auth tag: %v", err)
			}
			authTag = auth
		}
	}

	// Get operation from kind
	operation, exists := cip.GetOpFromKind(evt.Kind)
	if !exists {
		return nil, fmt.Errorf("unknown kind value: %d", evt.Kind)
	}

	// Parse based on operation type
	switch operation {
	case cip.OpModel:
		return parseModelEvent(evt, subspaceID, operation, authTag)
	case cip.OpData:
		return parseDataEvent(evt, subspaceID, operation, authTag)
	case cip.OpCompute:
		return parseComputeEvent(evt, subspaceID, operation, authTag)
	case cip.OpAlgo:
		return parseAlgoEvent(evt, subspaceID, operation, authTag)
	case cip.OpValid:
		return parseValidEvent(evt, subspaceID, operation, authTag)
	default:
		return nil, fmt.Errorf("unknown operation type: %s", operation)
	}
}

func parseModelEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag) (*ModelEvent, error) {
	model := &ModelEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "parent":
			model.ParentHash = tag[1]
		case "contrib":
			model.Contributions = tag[1]
		}
	}

	return model, nil
}

func parseDataEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag) (*DataEvent, error) {
	data := &DataEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		if tag[0] == "size" {
			data.Size = tag[1]
		}
	}

	return data, nil
}

func parseComputeEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag) (*ComputeEvent, error) {
	compute := &ComputeEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		if tag[0] == "compute_type" {
			compute.ComputeType = tag[1]
		}
	}

	return compute, nil
}

func parseAlgoEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag) (*AlgoEvent, error) {
	algo := &AlgoEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		if tag[0] == "algo_type" {
			algo.AlgoType = tag[1]
		}
	}

	return algo, nil
}

func parseValidEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag) (*ValidEvent, error) {
	valid := &ValidEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		if tag[0] == "valid_result" {
			valid.ValidResult = tag[1]
		}
	}

	return valid, nil
}

// NewModelEvent creates a new model event
func NewModelEvent(subspaceID string) (*ModelEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindModelgraphModel)
	if err != nil {
		return nil, err
	}
	return &ModelEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewDataEvent creates a new data event
func NewDataEvent(subspaceID string) (*DataEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindModelgraphData)
	if err != nil {
		return nil, err
	}
	return &DataEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewComputeEvent creates a new compute event
func NewComputeEvent(subspaceID string) (*ComputeEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindModelgraphCompute)
	if err != nil {
		return nil, err
	}
	return &ComputeEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewAlgoEvent creates a new algo event
func NewAlgoEvent(subspaceID string) (*AlgoEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindModelgraphAlgo)
	if err != nil {
		return nil, err
	}
	return &AlgoEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewValidEvent creates a new valid event
func NewValidEvent(subspaceID string) (*ValidEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindModelgraphValid)
	if err != nil {
		return nil, err
	}
	return &ValidEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}
