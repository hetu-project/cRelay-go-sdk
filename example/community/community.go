package main

import (
	"context"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
	"github.com/nbd-wtf/go-nostr/cip/cip07"
)

func main() {
	relays := []string{"ws://161.97.129.166:10547"}
	sk := nostr.GeneratePrivateKey()
	pub, _ := nostr.GetPublicKey(sk)

	// Create a subspace with Community operations
	createEvent := nostr.NewSubspaceCreateEvent(
		"community",
		cip.CommunitySubspaceOps,
		"energy>10000",
		"Community Example Subspace",
		"https://example.com/images/community.png",
	)
	createEvent.PubKey = pub
	createEvent.Sign(sk)

	// Join a subspace
	joinEvent := nostr.NewSubspaceJoinEvent(createEvent.SubspaceID)
	joinEvent.PubKey = pub
	joinEvent.Sign(sk)

	// Create a community
	communityEvent, err := cip07.NewCommunityCreateEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	communityEvent.PubKey = pub
	communityEvent.SetCommunityCreateInfo(
		"quantum_computing_community",
		"Quantum Computing Research Community",
		"research",
	)
	communityEvent.Content = "A community for researchers working on quantum computing"
	communityEvent.Sign(sk)
	fmt.Printf("%s\n", communityEvent)

	// Create a community invite
	inviteEvent, err := cip07.NewCommunityInviteEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	inviteEvent.PubKey = pub
	inviteEvent.SetCommunityInviteInfo(
		communityEvent.Event.ID,
		"user_456",
		"user_789",
		"email",
	)
	inviteEvent.Content = "You're invited to join our quantum computing research community!"
	inviteEvent.Sign(sk)

	// Create a channel
	channelEvent, err := cip07.NewChannelCreateEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	channelEvent.PubKey = pub
	channelEvent.SetChannelCreateInfo(
		communityEvent.Event.ID,
		"quantum_algorithms",
		"Quantum Algorithms Discussion",
		"discussion",
	)
	channelEvent.Content = "A channel for discussing quantum algorithms and their implementations"
	channelEvent.Sign(sk)

	// Create a channel message
	messageEvent, err := cip07.NewChannelMessageEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	messageEvent.PubKey = pub
	messageEvent.SetChannelMessageInfo(
		channelEvent.Event.ID,
		"user_456",
		"",
	)
	messageEvent.Content = "Has anyone tried implementing the quantum Fourier transform using Qiskit?"
	messageEvent.Sign(sk)
	fmt.Printf("%s\n", messageEvent)

	// Create a reply message
	replyEvent, err := cip07.NewChannelMessageEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	replyEvent.PubKey = pub
	replyEvent.SetChannelMessageInfo(
		channelEvent.Event.ID,
		"user_789",
		messageEvent.Event.ID,
	)
	replyEvent.Content = "Yes, I have! Here's my implementation..."
	replyEvent.Sign(sk)
	fmt.Printf("%s\n", replyEvent)

	// publish the events to relays
	ctx := context.Background()
	for _, url := range relays {
		relay, err := nostr.RelayConnect(ctx, url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Publish all events
		events := []nostr.Event{
			createEvent.Event,
			joinEvent.Event,
			communityEvent.Event,
			inviteEvent.Event,
			channelEvent.Event,
			messageEvent.Event,
			replyEvent.Event,
		}

		for _, ev := range events {
			if err := relay.Publish(ctx, ev); err != nil {
				fmt.Printf("failed to publish event to %s: %v\n", url, err)
				continue
			}
			fmt.Printf("published event %s to %s\n", ev.ID, url)
		}
	}

	// Query for channel messages
	relay, err := nostr.RelayConnect(ctx, relays[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	var filter nostr.Filter
	filter = nostr.Filter{
		Kinds: []int{cip.KindCommunityChannelMessage},
		Limit: 3,
	}
	events, err := relay.QueryEvents(ctx, filter)
	if err != nil {
		fmt.Printf("failed to query events: %v\n", err)
		return
	}

	for event := range events {
		fmt.Println("------")
		fmt.Printf("Content: %s\n", event)
	}
}
