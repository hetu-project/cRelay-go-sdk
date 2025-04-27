package main

import (
	"context"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
	"github.com/nbd-wtf/go-nostr/cip/cip01"
	"github.com/nbd-wtf/go-nostr/cip/cip02"
)

// AllOps returns the combined operations string including both basic and business operations
func AllOps() string {
	return cip.DefaultSubspaceOps + "," + cip.ModelGraphSubspaceOps
}

func main() {
	relays := []string{"ws://161.97.129.166:10547"}
	sk := nostr.GeneratePrivateKey()
	pub, _ := nostr.GetPublicKey(sk)

	// Create a subspace with all operations (basic + business)
	createEvent := nostr.NewSubspaceCreateEvent(
		"modelgraph",
		AllOps(), // Use combined operations string
		"energy>1000",
		"Desci AI Model collaboration subspace",
		"https://causality-graph.com/images/subspace.png",
	)
	createEvent.PubKey = pub
	createEvent.Sign(sk)

	// Join a subspace
	joinEvent := nostr.NewSubspaceJoinEvent(createEvent.SubspaceID)
	joinEvent.PubKey = pub
	joinEvent.Sign(sk)

	// Create a post operation (basic operation)
	postEvent, err := cip01.NewPostEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	postEvent.PubKey = pub
	postEvent.SetContentType("markdown")
	postEvent.SetParent("parent-hash")
	postEvent.Content = "# Subspace update \n We have completed the model optimization!"
	postEvent.Sign(sk)

	// Create a proposal (basic operation)
	proposeEvent, err := cip01.NewProposeEvent(createEvent.SubspaceID)
	proposeEvent.PubKey = pub
	proposeEvent.SetProposal("prop_001", "energy>2000")
	proposeEvent.Content = "Increase the energy requirement for subspace addition to 2000"
	proposeEvent.Sign(sk)

	// Create a vote (basic operation)
	voteEvent, err := cip01.NewVoteEvent(createEvent.SubspaceID)
	voteEvent.PubKey = pub
	voteEvent.SetVote("prop_001", "yes")
	voteEvent.Content = "Agree to increase the energy requirements"
	voteEvent.Sign(sk)

	// Create an invite (basic operation)
	inviteEvent, err := cip01.NewInviteEvent(createEvent.SubspaceID)
	inviteEvent.PubKey = pub
	inviteEvent.SetInvite("<Charlie's ETH Address>", "energy>1000")
	inviteEvent.Content = "Invite Charlie join into Desci AI subspace"
	inviteEvent.Sign(sk)

	// Create a model operation (business operation)
	modelEvent, err := cip02.NewModelEvent(createEvent.SubspaceID)
	modelEvent.PubKey = pub
	modelEvent.SetParent("parent-hash")
	modelEvent.SetContributions("base:0.1,data:0.6,algo:0.3")
	modelEvent.Content = "ipfs://bafy..."
	modelEvent.Sign(sk)

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
			postEvent.Event,
			proposeEvent.Event,
			voteEvent.Event,
			inviteEvent.Event,
			modelEvent.Event,
		}

		for _, ev := range events {
			if err := relay.Publish(ctx, ev); err != nil {
				fmt.Printf("failed to publish event to %s: %v\n", url, err)
				continue
			}
			fmt.Printf("published event %s to %s\n", ev.ID, url)
		}
	}
	relay, err := nostr.RelayConnect(ctx, relays[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	var filter nostr.Filter
	filter = nostr.Filter{
		Kinds: []int{30300},
		// limit = 3, get the three most recent notes
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
