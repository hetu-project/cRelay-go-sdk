package cip01

import (
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
)

// PostEvent represents a post operation in governance subspace
type PostEvent struct {
	*nostr.SubspaceOpEvent
	ContentType string
	ParentHash  string
}

// SetContentType sets the content type for the operation
func (e *PostEvent) SetContentType(contentType string) {
	e.ContentType = contentType
	e.Tags = append(e.Tags, nostr.Tag{"content_type", contentType})
}

// SetParent sets the parent event hash
func (e *PostEvent) SetParent(parentHash string) {
	e.ParentHash = parentHash
	e.Tags = append(e.Tags, nostr.Tag{"parent", parentHash})
}

// ProposeEvent represents a propose operation in governance subspace
type ProposeEvent struct {
	*nostr.SubspaceOpEvent
	ProposalID string
	Rules      string
}

// SetProposal sets the proposal ID and rules
func (e *ProposeEvent) SetProposal(proposalID, rules string) {
	e.ProposalID = proposalID
	e.Rules = rules
	e.Tags = append(e.Tags, nostr.Tag{"proposal_id", proposalID})
	if rules != "" {
		e.Tags = append(e.Tags, nostr.Tag{"rules", rules})
	}
}

// VoteEvent represents a vote operation in governance subspace
type VoteEvent struct {
	*nostr.SubspaceOpEvent
	ProposalID string
	Vote       string
}

// SetVote sets the vote for a proposal
func (e *VoteEvent) SetVote(proposalID, vote string) {
	e.ProposalID = proposalID
	e.Vote = vote
	e.Tags = append(e.Tags, nostr.Tag{"proposal_id", proposalID}, nostr.Tag{"vote", vote})
}

// InviteEvent represents an invite operation in governance subspace
type InviteEvent struct {
	*nostr.SubspaceOpEvent
	InviteePubkey string
	Rules         string
}

// SetInvite sets the invitee pubkey and rules
func (e *InviteEvent) SetInvite(inviteePubkey, rules string) {
	e.InviteePubkey = inviteePubkey
	e.Tags = append(e.Tags, nostr.Tag{"invitee_pubkey", inviteePubkey})
	if rules != "" {
		e.Tags = append(e.Tags, nostr.Tag{"rules", rules})
	}
}

// ParseGovernanceEvent parses a Nostr event into a governance event
func ParseGovernanceEvent(evt nostr.Event) (nostr.SubspaceOpEventPtr, error) {
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
	case "post":
		return parsePostEvent(evt, subspaceID, operation, authTag)
	case "propose":
		return parseProposeEvent(evt, subspaceID, operation, authTag)
	case "vote":
		return parseVoteEvent(evt, subspaceID, operation, authTag)
	case "invite":
		return parseInviteEvent(evt, subspaceID, operation, authTag)
	default:
		return nil, fmt.Errorf("unknown operation type: %s", operation)
	}

}

func parsePostEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag) (*PostEvent, error) {
	post := &PostEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
		},
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "content_type":
			post.ContentType = tag[1]
		case "parent":
			post.ParentHash = tag[1]
		}
	}

	return post, nil
}

func parseProposeEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag) (*ProposeEvent, error) {
	propose := &ProposeEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
		},
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "proposal_id":
			propose.ProposalID = tag[1]
		case "rules":
			propose.Rules = tag[1]
		}
	}

	return propose, nil
}

func parseVoteEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag) (*VoteEvent, error) {
	vote := &VoteEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
		},
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "proposal_id":
			vote.ProposalID = tag[1]
		case "vote":
			vote.Vote = tag[1]
		}
	}

	return vote, nil
}

func parseInviteEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag) (*InviteEvent, error) {
	invite := &InviteEvent{
		SubspaceOpEvent: &nostr.SubspaceOpEvent{
			SubspaceID: subspaceID,
			Operation:  operation,
			AuthTag:    authTag,
			Event:      evt,
		},
	}

	for _, tag := range evt.Tags {
		if len(tag) < 2 {
			continue
		}
		switch tag[0] {
		case "invitee_pubkey":
			invite.InviteePubkey = tag[1]
		case "rules":
			invite.Rules = tag[1]
		}
	}

	return invite, nil
}

// NewPostEvent creates a new post event
func NewPostEvent(subspaceID string) (*PostEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindGovernancePost)
	if err != nil {
		return nil, err
	}
	return &PostEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewProposeEvent creates a new propose event
func NewProposeEvent(subspaceID string) (*ProposeEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindGovernancePropose)
	if err != nil {
		return nil, err
	}
	return &ProposeEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewVoteEvent creates a new vote event
func NewVoteEvent(subspaceID string) (*VoteEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindGovernanceVote)
	if err != nil {
		return nil, err
	}
	return &VoteEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}

// NewInviteEvent creates a new invite event
func NewInviteEvent(subspaceID string) (*InviteEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindGovernanceInvite)
	if err != nil {
		return nil, err
	}
	return &InviteEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}
