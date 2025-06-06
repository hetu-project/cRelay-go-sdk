package main

import (
	"context"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/cip"
	"github.com/nbd-wtf/go-nostr/cip/cip06"
)

func main() {
	relays := []string{"ws://161.97.129.166:10547"}
	sk := nostr.GeneratePrivateKey()
	pub, _ := nostr.GetPublicKey(sk)

	// Create a subspace with Social operations
	createEvent := nostr.NewSubspaceCreateEvent(
		"social",
		cip.SocialSubspaceOps,
		"energy>10000",
		"Social Example Subspace",
		"https://example.com/images/social.png",
	)
	createEvent.PubKey = pub
	createEvent.Sign(sk)

	// Join a subspace
	joinEvent := nostr.NewSubspaceJoinEvent(createEvent.SubspaceID)
	joinEvent.PubKey = pub
	joinEvent.Sign(sk)

	// Create a like event
	likeEvent, err := cip06.NewLikeEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	likeEvent.PubKey = pub
	likeEvent.SetLikeInfo(
		"paper_123",
		"user_456",
	)
	likeEvent.Content = "I really enjoyed reading this paper!"
	likeEvent.Sign(sk)

	// Create a collect event
	collectEvent, err := cip06.NewCollectEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	collectEvent.PubKey = pub
	collectEvent.SetCollectInfo(
		"paper_123",
		"user_456",
	)
	collectEvent.Content = "Saving this for later reference"
	collectEvent.Sign(sk)

	// Create a share event
	shareEvent, err := cip06.NewShareEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	shareEvent.PubKey = pub
	shareEvent.SetShareInfo(
		"paper_123",
		"user_456",
		"twitter",
		"42",
	)
	shareEvent.Content = "Check out this interesting paper!"
	shareEvent.Sign(sk)

	// Create a comment event
	commentEvent, err := cip06.NewCommentEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	commentEvent.PubKey = pub
	commentEvent.SetCommentInfo(
		"paper_123",
		"user_456",
		"comment_789",
	)
	commentEvent.Content = "This is a great point! I would add that..."
	commentEvent.Sign(sk)

	// Create a tag event
	tagEvent, err := cip06.NewTagEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	tagEvent.PubKey = pub
	tagEvent.SetTagInfo(
		"paper_123",
		"quantum",
	)
	tagEvent.Content = "Tagged as quantum computing related"
	tagEvent.Sign(sk)

	// Create a follow event
	followEvent, err := cip06.NewFollowEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	followEvent.PubKey = pub
	followEvent.SetFollowInfo(
		"user_456",
		"user_789",
	)
	followEvent.Content = "Following this researcher for their work in quantum computing"
	followEvent.Sign(sk)

	// Create a question event
	questionEvent, err := cip06.NewQuestionEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	questionEvent.PubKey = pub
	questionEvent.SetQuestionInfo(
		"paper_123",
		"user_456",
		"high",
	)
	questionEvent.Content = "How does this approach compare to the method proposed in Smith et al. (2022)?"
	questionEvent.Sign(sk)

	// Create a room event
	roomEvent, err := cip06.NewRoomEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	roomEvent.PubKey = pub
	roomEvent.SetRoomInfo(
		"Quantum Computing Discussion",
		"A room for discussing quantum computing topics",
		[]string{"user_456", "user_789", "user_101"},
	)
	roomEvent.Content = "Welcome to our quantum computing discussion room!"
	roomEvent.Sign(sk)

	// Create a message event
	messageEvent, err := cip06.NewMessageEvent(createEvent.SubspaceID)
	if err != nil {
		fmt.Printf("err: %s", err)
	}
	messageEvent.PubKey = pub
	messageEvent.SetMessageInfo(
		roomEvent.Event.ID,
		"message_123",
		[]string{"user_789"},
	)
	messageEvent.Content = "Hey @user_789, what do you think about the latest quantum supremacy results?"
	messageEvent.Sign(sk)

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
			likeEvent.Event,
			collectEvent.Event,
			shareEvent.Event,
			commentEvent.Event,
			tagEvent.Event,
			followEvent.Event,
			questionEvent.Event,
			roomEvent.Event,
			messageEvent.Event,
		}

		for _, ev := range events {
			if err := relay.Publish(ctx, ev); err != nil {
				fmt.Printf("failed to publish event to %s: %v\n", url, err)
				continue
			}
			fmt.Printf("published event %s to %s\n", ev.ID, url)
		}
	}

	// Query for like events
	relay, err := nostr.RelayConnect(ctx, relays[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	var filter nostr.Filter
	filter = nostr.Filter{
		Kinds: []int{cip.KindSocialLike},
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
