package cip02

import (
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
)

// ProjectEvent represents a project operation in common graph
type ProjectEvent struct {
	*nostr.SubspaceOpEvent
	ProjectID string
	Name      string
	Desc      string
	Members   []string
	Status    string
}

// SetProjectInfo sets the project information
func (e *ProjectEvent) SetProjectInfo(projectID, name, desc string, members []string, status string) {
	e.ProjectID = projectID
	e.Name = name
	e.Desc = desc
	e.Members = members
	e.Status = status

	e.Tags = append(e.Tags,
		nostr.Tag{"project_id", projectID},
		nostr.Tag{"name", name},
		nostr.Tag{"desc", desc},
		nostr.Tag{"status", status},
	)

	if len(members) > 0 {
		membersTag := nostr.Tag{"members"}
		membersTag = append(membersTag, members...)
		e.Tags = append(e.Tags, membersTag)
	}
}

// TaskEvent represents a task operation in common graph
type TaskEvent struct {
	*nostr.SubspaceOpEvent
	ProjectID string
	TaskID    string
	Title     string
	Assignee  string
	Status    string
	Deadline  string
}

// SetTaskInfo sets the task information
func (e *TaskEvent) SetTaskInfo(projectID, taskID, title, assignee, status, deadline string) {
	e.ProjectID = projectID
	e.TaskID = taskID
	e.Title = title
	e.Assignee = assignee
	e.Status = status
	e.Deadline = deadline

	e.Tags = append(e.Tags,
		nostr.Tag{"project_id", projectID},
		nostr.Tag{"task_id", taskID},
		nostr.Tag{"title", title},
		nostr.Tag{"assignee", assignee},
		nostr.Tag{"status", status},
		nostr.Tag{"deadline", deadline},
	)
}

// EntityEvent represents an entity operation in common graph
type EntityEvent struct {
	*nostr.SubspaceOpEvent
	EntityName string
	EntityType string
}

// SetEntityInfo sets the entity information
func (e *EntityEvent) SetEntityInfo(entityName, entityType string) {
	e.EntityName = entityName
	e.EntityType = entityType

	e.Tags = append(e.Tags,
		nostr.Tag{"entity_name", entityName},
		nostr.Tag{"entity_type", entityType},
	)
}

// RelationEvent represents a relation operation in common graph
type RelationEvent struct {
	*nostr.SubspaceOpEvent
	From         string
	To           string
	RelationType string
	Context      string
}

// SetRelationInfo sets the relation information
func (e *RelationEvent) SetRelationInfo(from, to, relationType, context string) {
	e.From = from
	e.To = to
	e.RelationType = relationType
	e.Context = context

	e.Tags = append(e.Tags,
		nostr.Tag{"from", from},
		nostr.Tag{"to", to},
		nostr.Tag{"relation_type", relationType},
	)

	if context != "" {
		e.Tags = append(e.Tags, nostr.Tag{"context", context})
	}
}

// ObservationEvent represents an observation operation in common graph
type ObservationEvent struct {
	*nostr.SubspaceOpEvent
	EntityName  string
	Observation string
}

// SetObservationInfo sets the observation information
func (e *ObservationEvent) SetObservationInfo(entityName, observation string) {
	e.EntityName = entityName
	e.Observation = observation

	e.Tags = append(e.Tags,
		nostr.Tag{"entity_name", entityName},
		nostr.Tag{"observation", observation},
	)
}

// ParseCommonGraphEvent parses a Nostr event into a common graph event
func ParseCommonGraphEvent(evt nostr.Event) (nostr.SubspaceOpEventPtr, error) {
	// Extract common fields
	subspaceID := ""
	parents := []string{}
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
		case "parent":
			parents = append(parents, tag[1:]...)
		}
	}

	// Get operation from kind
	operation, exists := cip.GetOpFromKind(evt.Kind)
	if !exists {
		return nil, fmt.Errorf("unknown kind value: %d", evt.Kind)
	}

	// Parse based on operation type
	switch operation {
	case cip.OpProject:
		return parseProjectEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpTask:
		return parseTaskEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpEntity:
		return parseEntityEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpRelation:
		return parseRelationEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpObservation:
		return parseObservationEvent(evt, subspaceID, operation, authTag, parents)
	default:
		return nil, fmt.Errorf("unknown operation type: %s", operation)
	}
}

func parseProjectEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ProjectEvent, error) {
	project := &ProjectEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "project_id":
			project.ProjectID = tag[1]
		case "name":
			project.Name = tag[1]
		case "desc":
			project.Desc = tag[1]
		case "members":
			project.Members = tag[1:]
		case "status":
			project.Status = tag[1]
		}
	}

	return project, nil
}

func parseTaskEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*TaskEvent, error) {
	task := &TaskEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "project_id":
			task.ProjectID = tag[1]
		case "task_id":
			task.TaskID = tag[1]
		case "title":
			task.Title = tag[1]
		case "assignee":
			task.Assignee = tag[1]
		case "status":
			task.Status = tag[1]
		case "deadline":
			task.Deadline = tag[1]
		}
	}

	return task, nil
}

func parseEntityEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*EntityEvent, error) {
	entity := &EntityEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "entity_name":
			entity.EntityName = tag[1]
		case "entity_type":
			entity.EntityType = tag[1]
		}
	}

	return entity, nil
}

func parseRelationEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*RelationEvent, error) {
	relation := &RelationEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "from":
			relation.From = tag[1]
		case "to":
			relation.To = tag[1]
		case "relation_type":
			relation.RelationType = tag[1]
		case "context":
			relation.Context = tag[1]
		}
	}

	return relation, nil
}

func parseObservationEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ObservationEvent, error) {
	observation := &ObservationEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "entity_name":
			observation.EntityName = tag[1]
		case "observation":
			observation.Observation = tag[1]
		}
	}

	return observation, nil
}

// NewProjectEvent creates a new project event
func NewProjectEvent(subspaceID string) (*ProjectEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindCommonGraphProject)
	if err != nil {
		return nil, err
	}
	return &ProjectEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewTaskEvent creates a new task event
func NewTaskEvent(subspaceID string) (*TaskEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindCommonGraphTask)
	if err != nil {
		return nil, err
	}
	return &TaskEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewEntityEvent creates a new entity event
func NewEntityEvent(subspaceID string) (*EntityEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindCommonGraphEntity)
	if err != nil {
		return nil, err
	}
	return &EntityEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewRelationEvent creates a new relation event
func NewRelationEvent(subspaceID string) (*RelationEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindCommonGraphRelation)
	if err != nil {
		return nil, err
	}
	return &RelationEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewObservationEvent creates a new observation event
func NewObservationEvent(subspaceID string) (*ObservationEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindCommonGraphObservation)
	if err != nil {
		return nil, err
	}
	return &ObservationEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}
