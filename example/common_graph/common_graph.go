package main

import (
	"context"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
	"github.com/nbd-wtf/go-nostr/cip/cip02"
)

func main() {
	relays := []string{"ws://161.97.129.166:10547"}
	sk := nostr.GeneratePrivateKey()
	pub, _ := nostr.GetPublicKey(sk)

	// Create a subspace with all operations (basic + business)
	createEvent := nostr.NewSubspaceCreateEvent(
		"common_graph",
		cip.CommonPrjOps + "," + cip.CommonGraphOps, // Use default operations string
		"energy>10000",
		"Common Graph Example Subspace",
		"https://example.com/images/subspace.png",
	)
	createEvent.PubKey = pub
	createEvent.Sign(sk)

	// Join a subspace
	joinEvent := nostr.NewSubspaceJoinEvent(createEvent.SubspaceID)
	joinEvent.PubKey = pub
	joinEvent.Sign(sk)

	// Create a project
	projectEvent, err := cip02.NewProjectEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	projectEvent.PubKey = pub
	projectEvent.SetProjectInfo(
		"proj_001",
		"Quantum NLP",
		"Research on quantum natural language processing",
		[]string{"0xAlice", "0xBob"},
		"active",
	)
	projectEvent.Content = ""
	projectEvent.Sign(sk)

	// Create a task
	taskEvent, err := cip02.NewTaskEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	taskEvent.PubKey = pub
	taskEvent.SetTaskInfo(
		"proj_001",
		"task_001",
		"Literature Review",
		"0xBob",
		"in_progress",
		"1712345678",
	)
	taskEvent.Content = "Review recent papers on quantum NLP."
	taskEvent.Sign(sk)

	// Create an entity
	entityEvent, err := cip02.NewEntityEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	entityEvent.PubKey = pub
	entityEvent.SetEntityInfo("John_Smith", "person")
	entityEvent.Content = "{\"observations\":[\"Speaks fluent Spanish\"]}"
	entityEvent.Sign(sk)

	// Create a relation
	relationEvent, err := cip02.NewRelationEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	relationEvent.PubKey = pub
	relationEvent.SetRelationInfo("John_Smith", "Anthropic", "works_at", "")
	relationEvent.Content = ""
	relationEvent.Sign(sk)

	// Create an observation
	observationEvent, err := cip02.NewObservationEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	observationEvent.PubKey = pub
	observationEvent.SetObservationInfo("John_Smith", "Graduated in 2019")
	observationEvent.Content = ""
	observationEvent.Sign(sk)

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
			projectEvent.Event,
			taskEvent.Event,
			entityEvent.Event,
			relationEvent.Event,
			observationEvent.Event,
		}

		for _, ev := range events {
			if err := relay.Publish(ctx, ev); err != nil {
				fmt.Printf("failed to publish event to %s: %v\n", url, err)
				continue
			}
			fmt.Printf("published event %s to %s\n", ev.ID, url)
		}
	}

	// Query for project events
	relay, err := nostr.RelayConnect(ctx, relays[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	var filter nostr.Filter
	filter = nostr.Filter{
		Kinds: []int{cip.KindCommonGraphProject},
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