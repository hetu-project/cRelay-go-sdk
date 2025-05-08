package cip01

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
)

// PostEvent represents a post operation in governance subspace
type PostEvent struct {
	*nostr.SubspaceOpEvent
	ContentType string
}

// SetContentType sets the content type for the operation
func (e *PostEvent) SetContentType(contentType string) {
	e.ContentType = contentType
	e.Tags = append(e.Tags, nostr.Tag{"content_type", contentType})
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
	InviterAddr string
	Rules       string
}

// SetInviter sets the inviter eth address and rules
func (e *InviteEvent) SetInviter(inviterAddress, rules string) {
	e.InviterAddr = inviterAddress
	e.Tags = append(e.Tags, nostr.Tag{"inviter_addr", inviterAddress})
	if rules != "" {
		e.Tags = append(e.Tags, nostr.Tag{"rules", rules})
	}
}

// MintEvent represents a mint operation in governance subspace
type MintEvent struct {
	*nostr.SubspaceOpEvent
	TokenName     string
	TokenSymbol   string
	TokenDecimals string
	InitialSupply string
	DropRatio     string
}

// SetTokenInfo sets the token information for the mint operation
func (e *MintEvent) SetTokenInfo(name, symbol, decimals, initialSupply, dropRatio string) {
	e.TokenName = name
	e.TokenSymbol = symbol
	e.TokenDecimals = decimals
	e.InitialSupply = initialSupply
	e.DropRatio = dropRatio

	e.Tags = append(e.Tags,
		nostr.Tag{"token_name", name},
		nostr.Tag{"token_symbol", symbol},
		nostr.Tag{"token_decimals", decimals},
		nostr.Tag{"initial_supply", initialSupply},
		nostr.Tag{"drop_ratio", dropRatio},
	)
}

func (e *MintEvent) ParseRewardRules(rule string) map[int]int {
	rewardRules := make(map[int]int)
	segments := strings.Split(e.DropRatio, ",")
	for _, segment := range segments {
		parts := strings.Split(segment, ":")
		if len(parts) != 2 {
			continue
		}
		actionID, err1 := strconv.Atoi(parts[0])
		points, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil {
			rewardRules[actionID] = points
		}
	}
	return rewardRules
}

// ParseGovernanceEvent parses a Nostr event into a governance event
func ParseGovernanceEvent(evt nostr.Event) (nostr.SubspaceOpEventPtr, error) {
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
	case "post":
		return parsePostEvent(evt, subspaceID, operation, authTag, parents)
	case "propose":
		return parseProposeEvent(evt, subspaceID, operation, authTag, parents)
	case "vote":
		return parseVoteEvent(evt, subspaceID, operation, authTag, parents)
	case "invite":
		return parseInviteEvent(evt, subspaceID, operation, authTag, parents)
	case "mint":
		return parseMintEvent(evt, subspaceID, operation, authTag, parents)
	default:
		return nil, fmt.Errorf("unknown operation type: %s", operation)
	}

}

func parsePostEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*PostEvent, error) {
	post := &PostEvent{
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
		case "content_type":
			post.ContentType = tag[1]
		}
	}

	return post, nil
}

func parseProposeEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*ProposeEvent, error) {
	propose := &ProposeEvent{
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
		case "proposal_id":
			propose.ProposalID = tag[1]
		case "rules":
			propose.Rules = tag[1]
		}
	}

	return propose, nil
}

func parseVoteEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*VoteEvent, error) {
	vote := &VoteEvent{
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
		case "proposal_id":
			vote.ProposalID = tag[1]
		case "vote":
			vote.Vote = tag[1]
		}
	}

	return vote, nil
}

func parseInviteEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*InviteEvent, error) {
	invite := &InviteEvent{
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
		case "inviter_addr":
			invite.InviterAddr = tag[1]
		case "rules":
			invite.Rules = tag[1]
		}
	}

	return invite, nil
}

func parseMintEvent(evt nostr.Event, subspaceID, operation string, authTag cip.AuthTag, parents []string) (*MintEvent, error) {
	mint := &MintEvent{
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
		case "token_name":
			mint.TokenName = tag[1]
		case "token_symbol":
			mint.TokenSymbol = tag[1]
		case "token_decimals":
			mint.TokenDecimals = tag[1]
		case "initial_supply":
			mint.InitialSupply = tag[1]
		case "drop_ratio":
			mint.DropRatio = tag[1]
		}
	}

	return mint, nil
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

// NewMintEvent creates a new mint event
func NewMintEvent(subspaceID string) (*MintEvent, error) {
	baseEvent, err := nostr.NewSubspaceOpEvent(subspaceID, cip.KindGovernanceMint)
	if err != nil {
		return nil, err
	}
	return &MintEvent{
		SubspaceOpEvent: baseEvent,
	}, nil
}
