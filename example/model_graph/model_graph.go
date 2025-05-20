package main

import (
	"context"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
	"github.com/nbd-wtf/go-nostr/cip/cip03"
)

func main() {
	relays := []string{"ws://161.97.129.166:10547"}
	sk := nostr.GeneratePrivateKey()
	pub, _ := nostr.GetPublicKey(sk)

	// Create a subspace with modelgraph operations
	createEvent := nostr.NewSubspaceCreateEvent(
		"model_graph",
		cip.CommonPrjOps + cip.DefaultSubspaceOps + "," + cip.ModelGraphSubspaceOps, // Use modelgraph operations string
		"energy>10000",
		"Model Graph Example Subspace",
		"https://example.com/images/model_subspace.png",
	)
	createEvent.PubKey = pub
	createEvent.Sign(sk)

	// Join a subspace
	joinEvent := nostr.NewSubspaceJoinEvent(createEvent.SubspaceID)
	joinEvent.PubKey = pub
	joinEvent.Sign(sk)

	// Create a dataset event
	datasetEvent, err := cip03.NewDatasetEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	datasetEvent.PubKey = pub
	datasetEvent.SetDatasetInfo(
		"proj_001",
		"task_001",
		"training",
		"jsonl",
		[]string{"0xAlice", "0xBob"},
	)
	datasetEvent.Content = "{\"size\":1000,\"format\":\"jsonl\"}"
	datasetEvent.Sign(sk)

	// Create a finetune event
	finetuneEvent, err := cip03.NewFinetuneEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	finetuneEvent.PubKey = pub
	finetuneEvent.SetFinetuneInfo(
		"proj_001",
		"task_001",
		"dataset_001",
		"provider_001",
		"GPT-4",
	)
	finetuneEvent.Content = "{\"epochs\":3,\"learning_rate\":0.0001}"
	finetuneEvent.Sign(sk)

	// Create a conversation event
	conversationEvent, err := cip03.NewConversationEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	conversationEvent.PubKey = pub
	conversationEvent.SetConversationInfo(
		"session_001",
		"user_001",
		"model_001",
		"1712345678",
		"interaction_hash_001",
	)
	conversationEvent.Content = "{\"message\":\"Hello, how can I help you?\"}"
	conversationEvent.Sign(sk)

	// Create a session event
	sessionEvent, err := cip03.NewSessionEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	sessionEvent.PubKey = pub
	sessionEvent.SetSessionInfo(
		"session_001",
		"start",
		"user_001",
		"1712345678",
		"1712346000",
	)
	sessionEvent.Content = "{\"status\":\"active\"}"
	sessionEvent.Sign(sk)

	// Create a model event
	modelEvent, err := cip03.NewModelEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	modelEvent.PubKey = pub
	modelEvent.SetParent("parent_model_hash")
	modelEvent.SetContributions("0.8,0.2")
	modelEvent.Content = "{\"name\":\"GPT-4\",\"version\":\"1.0\"}"
	modelEvent.Sign(sk)

	// Create a compute event
	computeEvent, err := cip03.NewComputeEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	computeEvent.PubKey = pub
	computeEvent.ComputeType = "training"
	computeEvent.Content = "{\"gpu\":\"A100\",\"duration\":\"2h\"}"
	computeEvent.Sign(sk)

	// Create an algo event
	algoEvent, err := cip03.NewAlgoEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	algoEvent.PubKey = pub
	algoEvent.AlgoType = "transformer"
	algoEvent.Content = "{\"layers\":12,\"heads\":8}"
	algoEvent.Sign(sk)

	// Create a valid event
	validEvent, err := cip03.NewValidEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	validEvent.PubKey = pub
	validEvent.ValidResult = "passed"
	validEvent.Content = "{\"accuracy\":0.95,\"metrics\":{\"precision\":0.94,\"recall\":0.96}}"
	validEvent.Sign(sk)

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
			modelEvent.Event,
			datasetEvent.Event,
			computeEvent.Event,
			algoEvent.Event,
			validEvent.Event,
			finetuneEvent.Event,
			conversationEvent.Event,
			sessionEvent.Event,
		}

		for _, ev := range events {
			if err := relay.Publish(ctx, ev); err != nil {
				fmt.Printf("failed to publish event to %s: %v\n", url, err)
				continue
			}
			fmt.Printf("published event %s to %s\n", ev.ID, url)
		}
	}

	// Query for model events
	relay, err := nostr.RelayConnect(ctx, relays[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	var filter nostr.Filter
	filter = nostr.Filter{
		Kinds: []int{cip.KindModelgraphDataset},
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