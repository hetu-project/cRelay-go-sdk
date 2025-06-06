package cip07

import (
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
)

// CommunityCreateEvent represents a community creation operation
type CommunityCreateEvent struct {
	*nostr.SubspaceOpEvent
	CommunityID string
	Name        string
	Type        string
}

// SetCommunityCreateInfo sets the community creation information
func (e *CommunityCreateEvent) SetCommunityCreateInfo(communityID, name, communityType string) {
	e.CommunityID = communityID
	e.Name = name
	e.Type = communityType

	e.Tags = append(e.Tags,
		nostr.Tag{"community_id", communityID},
		nostr.Tag{"name", name},
		nostr.Tag{"type", communityType},
	)
}

// CommunityInviteEvent represents a community invitation operation
type CommunityInviteEvent struct {
	*nostr.SubspaceOpEvent
	CommunityID string
	InviterID   string
	InviteeID   string
	Method      string
}

// SetCommunityInviteInfo sets the community invitation information
func (e *CommunityInviteEvent) SetCommunityInviteInfo(communityID, inviterID, inviteeID, method string) {
	e.CommunityID = communityID
	e.InviterID = inviterID
	e.InviteeID = inviteeID
	e.Method = method

	e.Tags = append(e.Tags,
		nostr.Tag{"community_id", communityID},
		nostr.Tag{"inviter_id", inviterID},
		nostr.Tag{"invitee_id", inviteeID},
		nostr.Tag{"method", method},
	)
}

// ChannelCreateEvent represents a channel creation operation
type ChannelCreateEvent struct {
	*nostr.SubspaceOpEvent
	CommunityID string
	ChannelID   string
	Name        string
	Type        string
}

// SetChannelCreateInfo sets the channel creation information
func (e *ChannelCreateEvent) SetChannelCreateInfo(communityID, channelID, name, channelType string) {
	e.CommunityID = communityID
	e.ChannelID = channelID
	e.Name = name
	e.Type = channelType

	e.Tags = append(e.Tags,
		nostr.Tag{"community_id", communityID},
		nostr.Tag{"channel_id", channelID},
		nostr.Tag{"name", name},
		nostr.Tag{"type", channelType},
	)
}

// ChannelMessageEvent represents a channel message operation
type ChannelMessageEvent struct {
	*nostr.SubspaceOpEvent
	ChannelID string
	UserID    string
	ReplyTo   string
}

// SetChannelMessageInfo sets the channel message information
func (e *ChannelMessageEvent) SetChannelMessageInfo(channelID, userID, replyTo string) {
	e.ChannelID = channelID
	e.UserID = userID
	e.ReplyTo = replyTo

	e.Tags = append(e.Tags,
		nostr.Tag{"channel_id", channelID},
		nostr.Tag{"user_id", userID},
	)

	if replyTo != "" {
		e.Tags = append(e.Tags, nostr.Tag{"reply_to", replyTo})
	}
}

// ParseCommunityEvent parses a Nostr event into a community event
func ParseCommunityEvent(evt nostr.Event) (nostr.SubspaceOpEventPtr, error) {
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
	case cip.OpCommunityCreate:
		return parseCommunityCreateEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpCommunityInvite:
		return parseCommunityInviteEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpChannelCreate:
		return parseChannelCreateEvent(evt, subspaceID, operation, authTag, parents)
	case cip.OpChannelMessage:
		return parseChannelMessageEvent(evt, subspaceID, operation, authTag, parents)
	default:
		return nil, fmt.Errorf("unknown operation type: %s", operation)
	}
}

func parseCommunityCreateEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*CommunityCreateEvent, error) {
	create := &CommunityCreateEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	create.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "community_id":
			create.CommunityID = tag[1]
		case "name":
			create.Name = tag[1]
		case "type":
			create.Type = tag[1]
		}
	}

	return create, nil
}

func parseCommunityInviteEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*CommunityInviteEvent, error) {
	invite := &CommunityInviteEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	invite.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "community_id":
			invite.CommunityID = tag[1]
		case "inviter_id":
			invite.InviterID = tag[1]
		case "invitee_id":
			invite.InviteeID = tag[1]
		case "method":
			invite.Method = tag[1]
		}
	}

	return invite, nil
}

func parseChannelCreateEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ChannelCreateEvent, error) {
	create := &ChannelCreateEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
			Parents:    parents,
		},
	}
	create.Event.Content = evt.Content

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "community_id":
			create.CommunityID = tag[1]
		case "channel_id":
			create.ChannelID = tag[1]
		case "name":
			create.Name = tag[1]
		case "type":
			create.Type = tag[1]
		}
	}

	return create, nil
}

func parseChannelMessageEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ChannelMessageEvent, error) {
	message := &ChannelMessageEvent{
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
		case "channel_id":
			message.ChannelID = tag[1]
		case "user_id":
			message.UserID = tag[1]
		case "reply_to":
			message.ReplyTo = tag[1]
		}
	}

	return message, nil
}

// NewCommunityCreateEvent creates a new community creation event
func NewCommunityCreateEvent(subspaceID string) (*CommunityCreateEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindCommunityCreate)
	if err != nil {
		return nil, err
	}
	return &CommunityCreateEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewCommunityInviteEvent creates a new community invitation event
func NewCommunityInviteEvent(subspaceID string) (*CommunityInviteEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindCommunityInvite)
	if err != nil {
		return nil, err
	}
	return &CommunityInviteEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewChannelCreateEvent creates a new channel creation event
func NewChannelCreateEvent(subspaceID string) (*ChannelCreateEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindCommunityChannelCreate)
	if err != nil {
		return nil, err
	}
	return &ChannelCreateEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewChannelMessageEvent creates a new channel message event
func NewChannelMessageEvent(subspaceID string) (*ChannelMessageEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindCommunityChannelMessage)
	if err != nil {
		return nil, err
	}
	return &ChannelMessageEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}
