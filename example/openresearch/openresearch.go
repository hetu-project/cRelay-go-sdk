package main

import (
	"context"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
	"github.com/nbd-wtf/go-nostr/cip/cip05"
)

func main() {
	relays := []string{"ws://161.97.129.166:10547"}
	sk := nostr.GeneratePrivateKey()
	pub, _ := nostr.GetPublicKey(sk)

	// Create a subspace with OpenResearch operations
	createEvent := nostr.NewSubspaceCreateEvent(
		"openresearch",
		cip.OpenResearchSubspaceOps,
		"energy>10000",
		"OpenResearch Example Subspace",
		"https://example.com/images/subspace.png",
	)
	createEvent.PubKey = pub
	createEvent.Sign(sk)

	// Join a subspace
	joinEvent := nostr.NewSubspaceJoinEvent(createEvent.SubspaceID)
	joinEvent.PubKey = pub
	joinEvent.Sign(sk)

	// Create a paper event
	paperEvent, err := cip04.NewPaperEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	paperEvent.PubKey = pub
	paperEvent.SetPaperInfo(
		"10.1234/example.2023",
		"pdf",
		[]string{"Alice Smith", "Bob Johnson", "Carol White"},
		[]string{"quantum computing", "machine learning", "optimization"},
		"2023",
		"Journal of Quantum Computing",
	)
	paperEvent.Content = `{"title":"Quantum Machine Learning: A Survey","abstract":"This paper provides a comprehensive survey of quantum machine learning...","url":"https://doi.org/10.1234/example.2023","file_hash":"ipfs://bafy..."}`
	paperEvent.Sign(sk)

	// Create an annotation event
	annotationEvent, err := cip04.NewAnnotationEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	annotationEvent.PubKey = pub
	annotationEvent.SetAnnotationInfo(
		paperEvent.Event.ID,
		"section:2,paragraph:3,offset:120,length:250",
		"comment",
		"",
	)
	annotationEvent.Content = "This finding contradicts the results from Smith et al. (2022), which might be due to different experimental conditions."
	annotationEvent.Sign(sk)

	// Create a review event
	reviewEvent, err := cip04.NewReviewEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	reviewEvent.PubKey = pub
	reviewEvent.SetReviewInfo(
		paperEvent.Event.ID,
		"4.5",
		map[string]string{
			"methodology":     "4",
			"novelty":         "5",
			"clarity":         "4",
			"reproducibility": "3",
		},
	)
	reviewEvent.Content = `{"summary":"This paper presents a novel approach to quantum machine learning...","strengths":"The methodology is robust and well-documented...","weaknesses":"The results section lacks detailed experimental parameters...","recommendations":"Authors should consider adding more implementation details..."}`
	reviewEvent.Sign(sk)

	// Create an AI analysis event
	aiAnalysisEvent, err := cip04.NewAIAnalysisEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	aiAnalysisEvent.PubKey = pub
	aiAnalysisEvent.SetAIAnalysisInfo(
		"literature_gap",
		[]string{paperEvent.Event.ID},
		"Identify research gaps and potential future directions in this paper",
	)
	aiAnalysisEvent.Content = `{"analysis_result":"Based on the provided paper, several research gaps emerge: 1) Limited exploration of quantum-classical hybrid approaches...","key_insights":["insight1","insight2"],"potential_directions":["direction1","direction2"]}`
	aiAnalysisEvent.Sign(sk)

	// Create a discussion event
	discussionEvent, err := cip04.NewDiscussionEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	discussionEvent.PubKey = pub
	discussionEvent.SetDiscussionInfo(
		"quantum_ml_ethics",
		"",
		[]string{paperEvent.Event.ID},
	)
	discussionEvent.Content = "The ethical implications of quantum machine learning should be addressed proactively. As discussed in the paper, quantum computers could potentially break current encryption standards."
	discussionEvent.Sign(sk)

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
			paperEvent.Event,
			annotationEvent.Event,
			reviewEvent.Event,
			aiAnalysisEvent.Event,
			discussionEvent.Event,
		}

		for _, ev := range events {
			if err := relay.Publish(ctx, ev); err != nil {
				fmt.Printf("failed to publish event to %s: %v\n", url, err)
				continue
			}
			fmt.Printf("published event %s to %s\n", ev.ID, url)
		}
	}

	// Query for paper events
	relay, err := nostr.RelayConnect(ctx, relays[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	var filter nostr.Filter
	filter = nostr.Filter{
		Kinds: []int{cip.KindOpenResearchPaper},
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
