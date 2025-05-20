package cip03

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

// DatasetEvent represents a dataset operation in modelgraph subspace
type DatasetEvent struct {
	*nostr.SubspaceOpEvent
	ProjectID    string
	TaskID       string
	Category     string
	Format       string
	Contributors []string
	Content      string
}

// SetDatasetInfo sets the dataset information
func (e *DatasetEvent) SetDatasetInfo(projectID, taskID, category, format string, contributors []string) {
	e.ProjectID = projectID
	e.TaskID = taskID
	e.Category = category
	e.Format = format
	e.Contributors = contributors

	e.Tags = append(e.Tags,
		nostr.Tag{"project_id", projectID},
		nostr.Tag{"task_id", taskID},
		nostr.Tag{"category", category},
		nostr.Tag{"format", format},
	)

	if len(contributors) > 0 {
		contributorsTag := nostr.Tag{"contributors"}
		contributorsTag = append(contributorsTag, contributors...)
		e.Tags = append(e.Tags, contributorsTag)
	}
}

// FinetuneEvent represents a finetune operation in modelgraph subspace
type FinetuneEvent struct {
	*nostr.SubspaceOpEvent
	ProjectID  string
	TaskID     string
	DatasetID  string
	ProviderID string
	ModelName  string
	Content    string
}

// SetFinetuneInfo sets the finetune information
func (e *FinetuneEvent) SetFinetuneInfo(projectID, taskID, datasetID, providerID, modelName string) {
	e.ProjectID = projectID
	e.TaskID = taskID
	e.DatasetID = datasetID
	e.ProviderID = providerID
	e.ModelName = modelName

	e.Tags = append(e.Tags,
		nostr.Tag{"project_id", projectID},
		nostr.Tag{"task_id", taskID},
		nostr.Tag{"dataset_id", datasetID},
		nostr.Tag{"provider_id", providerID},
		nostr.Tag{"model_name", modelName},
	)
}

// ConversationEvent represents a conversation operation in modelgraph subspace
type ConversationEvent struct {
	*nostr.SubspaceOpEvent
	SessionID       string
	UserID          string
	ModelID         string
	Timestamp       string
	InteractionHash string
	Content         string
}

// SetConversationInfo sets the conversation information
func (e *ConversationEvent) SetConversationInfo(sessionID, userID, modelID, timestamp, interactionHash string) {
	e.SessionID = sessionID
	e.UserID = userID
	e.ModelID = modelID
	e.Timestamp = timestamp
	e.InteractionHash = interactionHash

	e.Tags = append(e.Tags,
		nostr.Tag{"session_id", sessionID},
		nostr.Tag{"user_id", userID},
		nostr.Tag{"model_id", modelID},
		nostr.Tag{"timestamp", timestamp},
		nostr.Tag{"interaction_hash", interactionHash},
	)
}

// SessionEvent represents a session operation in modelgraph subspace
type SessionEvent struct {
	*nostr.SubspaceOpEvent
	SessionID string
	Action    string
	UserID    string
	StartTime string
	EndTime   string
	Content   string
}

// SetSessionInfo sets the session information
func (e *SessionEvent) SetSessionInfo(sessionID, action, userID, startTime, endTime string) {
	e.SessionID = sessionID
	e.Action = action
	e.UserID = userID
	e.StartTime = startTime
	e.EndTime = endTime

	e.Tags = append(e.Tags,
		nostr.Tag{"session_id", sessionID},
		nostr.Tag{"action", action},
		nostr.Tag{"user_id", userID},
		nostr.Tag{"start_time", startTime},
	)

	if endTime != "" {
		e.Tags = append(e.Tags, nostr.Tag{"end_time", endTime})
	}
}

// ParseModelGraphEvent parses a Nostr event into a modelgraph event
func ParseModelGraphEvent(evt nostr.Event) (nostr.SubspaceOpEventPtr, error) {
	// Extract common fields
	subspaceID := ""
	var authTag cip.AuthTag
	parents := []string{}

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
		return parseModelEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpDataset:
		return parseDatasetEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpCompute:
		return parseComputeEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpAlgo:
		return parseAlgoEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpValid:
		return parseValidEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpFinetune:
		return parseFinetuneEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpConversation:
		return parseConversationEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpSession:
		return parseSessionEvent(evt, subspaceID, operation, authTag, parents)
	default:
		return nil, fmt.Errorf("unknown operation type: %s", operation)
	}
}

func parseDatasetEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*DatasetEvent, error) {
	dataset := &DatasetEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Parents:    parents,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "project_id":
			dataset.ProjectID = tag[1]
		case "task_id":
			dataset.TaskID = tag[1]
		case "category":
			dataset.Category = tag[1]
		case "format":
			dataset.Format = tag[1]
		case "contributors":
			if len(tag) > 1 {
				dataset.Contributors = tag[1:]
			}
		}
	}

	return dataset, nil
}

func parseModelEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ModelEvent, error) {
	model := &ModelEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Parents:    parents,
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

func parseComputeEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ComputeEvent, error) {
	compute := &ComputeEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Parents:    parents,
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

func parseAlgoEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*AlgoEvent, error) {
	algo := &AlgoEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Parents:    parents,
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

func parseValidEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ValidEvent, error) {
	valid := &ValidEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Parents:    parents,
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

func parseFinetuneEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*FinetuneEvent, error) {
	finetune := &FinetuneEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Parents:    parents,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "project_id":
			finetune.ProjectID = tag[1]
		case "task_id":
			finetune.TaskID = tag[1]
		case "dataset_id":
			finetune.DatasetID = tag[1]
		case "provider_id":
			finetune.ProviderID = tag[1]
		case "model_name":
			finetune.ModelName = tag[1]
		}
	}

	return finetune, nil
}

func parseConversationEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ConversationEvent, error) {
	conversation := &ConversationEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Parents:    parents,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "session_id":
			conversation.SessionID = tag[1]
		case "user_id":
			conversation.UserID = tag[1]
		case "model_id":
			conversation.ModelID = tag[1]
		case "timestamp":
			conversation.Timestamp = tag[1]
		case "interaction_hash":
			conversation.InteractionHash = tag[1]
		}
	}

	return conversation, nil
}

func parseSessionEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*SessionEvent, error) {
	session := &SessionEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Parents:    parents,
		},
		Content: evt.Content,
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "session_id":
			session.SessionID = tag[1]
		case "action":
			session.Action = tag[1]
		case "user_id":
			session.UserID = tag[1]
		case "start_time":
			session.StartTime = tag[1]
		case "end_time":
			session.EndTime = tag[1]
		}
	}

	return session, nil
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

// NewDatasetEvent creates a new dataset event
func NewDatasetEvent(subspaceID string) (*DatasetEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindModelgraphDataset)
	if err != nil {
		return nil, err
	}
	return &DatasetEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewFinetuneEvent creates a new finetune event
func NewFinetuneEvent(subspaceID string) (*FinetuneEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindModelgraphFinetune)
	if err != nil {
		return nil, err
	}
	return &FinetuneEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewConversationEvent creates a new conversation event
func NewConversationEvent(subspaceID string) (*ConversationEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindModelgraphConversation)
	if err != nil {
		return nil, err
	}
	return &ConversationEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewSessionEvent creates a new session event
func NewSessionEvent(subspaceID string) (*SessionEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindModelgraphSession)
	if err != nil {
		return nil, err
	}
	return &SessionEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}
