package cip06

import (
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
)

// LikeEvent represents a like operation in social subspace
type LikeEvent struct {
	*nostr.SubspaceOpEvent
	ObjectID string
	UserID   string
}

// SetLikeInfo sets the like information
func (e *LikeEvent) SetLikeInfo(objectID, userID string) {
	e.ObjectID = objectID
	e.UserID = userID

	e.Tags = append(e.Tags,
		nostr.Tag{"object_id", objectID},
		nostr.Tag{"user_id", userID},
	)
}

// CollectEvent represents a collect operation in social subspace
type CollectEvent struct {
	*nostr.SubspaceOpEvent
	ObjectID string
	UserID   string
}

// SetCollectInfo sets the collect information
func (e *CollectEvent) SetCollectInfo(objectID, userID string) {
	e.ObjectID = objectID
	e.UserID = userID

	e.Tags = append(e.Tags,
		nostr.Tag{"object_id", objectID},
		nostr.Tag{"user_id", userID},
	)
}

// ShareEvent represents a share operation in social subspace
type ShareEvent struct {
	*nostr.SubspaceOpEvent
	ObjectID string
	UserID   string
	Platform string
	Clicks   string
}

// SetShareInfo sets the share information
func (e *ShareEvent) SetShareInfo(objectID, userID, platform, clicks string) {
	e.ObjectID = objectID
	e.UserID = userID
	e.Platform = platform
	e.Clicks = clicks

	e.Tags = append(e.Tags,
		nostr.Tag{"object_id", objectID},
		nostr.Tag{"user_id", userID},
		nostr.Tag{"platform", platform},
		nostr.Tag{"clicks", clicks},
	)
}

// CommentEvent represents a comment operation in social subspace
type CommentEvent struct {
	*nostr.SubspaceOpEvent
	ObjectID string
	UserID   string
	Parent   string
}

// SetCommentInfo sets the comment information
func (e *CommentEvent) SetCommentInfo(objectID, userID, parent string) {
	e.ObjectID = objectID
	e.UserID = userID
	e.Parent = parent

	e.Tags = append(e.Tags,
		nostr.Tag{"object_id", objectID},
		nostr.Tag{"user_id", userID},
	)

	if parent != "" {
		e.Tags = append(e.Tags, nostr.Tag{"parent", parent})
	}
}

// TagEvent represents a tag operation in social subspace
type TagEvent struct {
	*nostr.SubspaceOpEvent
	ObjectID string
	Tag      string
}

// SetTagInfo sets the tag information
func (e *TagEvent) SetTagInfo(objectID, tag string) {
	e.ObjectID = objectID
	e.Tag = tag

	e.Tags = append(e.Tags,
		nostr.Tag{"object_id", objectID},
		nostr.Tag{"tag", tag},
	)
}

// FollowEvent represents a follow operation in social subspace
type FollowEvent struct {
	*nostr.SubspaceOpEvent
	UserID   string
	TargetID string
}

// SetFollowInfo sets the follow information
func (e *FollowEvent) SetFollowInfo(userID, targetID string) {
	e.UserID = userID
	e.TargetID = targetID

	e.Tags = append(e.Tags,
		nostr.Tag{"user_id", userID},
		nostr.Tag{"target_id", targetID},
	)
}

// UnfollowEvent represents an unfollow operation in social subspace
type UnfollowEvent struct {
	*nostr.SubspaceOpEvent
	UserID   string
	TargetID string
}

// SetUnfollowInfo sets the unfollow information
func (e *UnfollowEvent) SetUnfollowInfo(userID, targetID string) {
	e.UserID = userID
	e.TargetID = targetID

	e.Tags = append(e.Tags,
		nostr.Tag{"user_id", userID},
		nostr.Tag{"target_id", targetID},
	)
}

// QuestionEvent represents a question operation in social subspace
type QuestionEvent struct {
	*nostr.SubspaceOpEvent
	ObjectID string
	UserID   string
	Quality  string
}

// SetQuestionInfo sets the question information
func (e *QuestionEvent) SetQuestionInfo(objectID, userID, quality string) {
	e.ObjectID = objectID
	e.UserID = userID
	e.Quality = quality

	e.Tags = append(e.Tags,
		nostr.Tag{"object_id", objectID},
		nostr.Tag{"user_id", userID},
		nostr.Tag{"quality", quality},
	)
}

// RoomEvent represents a room operation in social subspace
type RoomEvent struct {
	*nostr.SubspaceOpEvent
	Name        string
	Description string
	Members     []string
}

// SetRoomInfo sets the room information
func (e *RoomEvent) SetRoomInfo(name, description string, members []string) {
	e.Name = name
	e.Description = description
	e.Members = members

	e.Tags = append(e.Tags,
		nostr.Tag{"name", name},
		nostr.Tag{"description", description},
	)

	if len(members) > 0 {
		membersTag := nostr.Tag{"members"}
		membersTag = append(membersTag, members...)
		e.Tags = append(e.Tags, membersTag)
	}
}

// MessageEvent represents a message operation in social subspace
type MessageEvent struct {
	*nostr.SubspaceOpEvent
	RoomID   string
	ReplyTo  string
	Mentions []string
}

// SetMessageInfo sets the message information
func (e *MessageEvent) SetMessageInfo(roomID, replyTo string, mentions []string) {
	e.RoomID = roomID
	e.ReplyTo = replyTo
	e.Mentions = mentions

	e.Tags = append(e.Tags,
		nostr.Tag{"room_id", roomID},
	)

	if replyTo != "" {
		e.Tags = append(e.Tags, nostr.Tag{"reply_to", replyTo})
	}

	if len(mentions) > 0 {
		mentionsTag := nostr.Tag{"mentions"}
		mentionsTag = append(mentionsTag, mentions...)
		e.Tags = append(e.Tags, mentionsTag)
	}
}

// ParseSocialEvent parses a Nostr event into a social event
func ParseSocialEvent(evt nostr.Event) (nostr.SubspaceOpEventPtr, error) {
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
	case cip.OpLike:
		return parseLikeEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpCollect:
		return parseCollectEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpShare:
		return parseShareEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpComment:
		return parseCommentEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpTag:
		return parseTagEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpFollow:
		return parseFollowEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpUnfollow:
		return parseUnfollowEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpQuestion:
		return parseQuestionEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpRoom:
		return parseRoomEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpMessage:
		return parseMessageEvent(evt, subspaceID, operation, authTag, parents)
	default:
		return nil, fmt.Errorf("unknown operation type: %s", operation)
	}
}

func parseLikeEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*LikeEvent, error) {
	like := &LikeEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	like.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "object_id":
			like.ObjectID = tag[1]
		case "user_id":
			like.UserID = tag[1]
		}
	}

	return like, nil
}

func parseCollectEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*CollectEvent, error) {
	collect := &CollectEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	collect.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "object_id":
			collect.ObjectID = tag[1]
		case "user_id":
			collect.UserID = tag[1]
		}
	}

	return collect, nil
}

func parseShareEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ShareEvent, error) {
	share := &ShareEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	share.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "object_id":
			share.ObjectID = tag[1]
		case "user_id":
			share.UserID = tag[1]
		case "platform":
			share.Platform = tag[1]
		case "clicks":
			share.Clicks = tag[1]
		}
	}

	return share, nil
}

func parseCommentEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*CommentEvent, error) {
	comment := &CommentEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	comment.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "object_id":
			comment.ObjectID = tag[1]
		case "user_id":
			comment.UserID = tag[1]
		case "parent":
			comment.Parent = tag[1]
		}
	}

	return comment, nil
}

func parseTagEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*TagEvent, error) {
	tag := &TagEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	tag.Event.Content = evt.Content

	for _, t := range evt.Tags {
		if len(t) < 2 {
			continue
		}
		switch t[0] {
		case "object_id":
			tag.ObjectID = t[1]
		case "tag":
			tag.Tag = t[1]
		}
	}

	return tag, nil
}

func parseFollowEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*FollowEvent, error) {
	follow := &FollowEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	follow.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "user_id":
			follow.UserID = tag[1]
		case "target_id":
			follow.TargetID = tag[1]
		}
	}

	return follow, nil
}

func parseUnfollowEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*UnfollowEvent, error) {
	unfollow := &UnfollowEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	unfollow.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "user_id":
			unfollow.UserID = tag[1]
		case "target_id":
			unfollow.TargetID = tag[1]
		}
	}

	return unfollow, nil
}

func parseQuestionEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*QuestionEvent, error) {
	question := &QuestionEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	question.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "object_id":
			question.ObjectID = tag[1]
		case "user_id":
			question.UserID = tag[1]
		case "quality":
			question.Quality = tag[1]
		}
	}

	return question, nil
}

func parseRoomEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*RoomEvent, error) {
	room := &RoomEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	room.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "name":
			room.Name = tag[1]
		case "description":
			room.Description = tag[1]
		case "members":
			room.Members = tag[1:]
		}
	}

	return room, nil
}

func parseMessageEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*MessageEvent, error) {
	message := &MessageEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	message.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "room_id":
			message.RoomID = tag[1]
		case "reply_to":
			message.ReplyTo = tag[1]
		case "mentions":
			message.Mentions = tag[1:]
		}
	}

	return message, nil
}

// NewLikeEvent creates a new like event
func NewLikeEvent(subspaceID string) (*LikeEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindSocialLike)
	if err != nil {
		return nil, err
	}
	return &LikeEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewCollectEvent creates a new collect event
func NewCollectEvent(subspaceID string) (*CollectEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindSocialCollect)
	if err != nil {
		return nil, err
	}
	return &CollectEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewShareEvent creates a new share event
func NewShareEvent(subspaceID string) (*ShareEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindSocialShare)
	if err != nil {
		return nil, err
	}
	return &ShareEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewCommentEvent creates a new comment event
func NewCommentEvent(subspaceID string) (*CommentEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindSocialComment)
	if err != nil {
		return nil, err
	}
	return &CommentEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewTagEvent creates a new tag event
func NewTagEvent(subspaceID string) (*TagEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindSocialTag)
	if err != nil {
		return nil, err
	}
	return &TagEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewFollowEvent creates a new follow event
func NewFollowEvent(subspaceID string) (*FollowEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindSocialFollow)
	if err != nil {
		return nil, err
	}
	return &FollowEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewUnfollowEvent creates a new unfollow event
func NewUnfollowEvent(subspaceID string) (*UnfollowEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindSocialUnfollow)
	if err != nil {
		return nil, err
	}
	return &UnfollowEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewQuestionEvent creates a new question event
func NewQuestionEvent(subspaceID string) (*QuestionEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindSocialQuestion)
	if err != nil {
		return nil, err
	}
	return &QuestionEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewRoomEvent creates a new room event
func NewRoomEvent(subspaceID string) (*RoomEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindSocialRoom)
	if err != nil {
		return nil, err
	}
	return &RoomEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewMessageEvent creates a new message event
func NewMessageEvent(subspaceID string) (*MessageEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindSocialMessage)
	if err != nil {
		return nil, err
	}
	return &MessageEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}
