package main

import (
	"fmt"
	"log"

	"github.com/nbd-wtf/go-nostr/cip/generator"
)

func main() {
	// Define a new CIP for a chat application
	chatCIP := generator.CIPDefinition{
		CIPName:     "chat",
		Package:     "cip100",
		Description: "Chat application CIP implementation",
		Events: []generator.EventDefinition{
			{
				EventName:   "MessageEvent",
				Operation:   "message",
				Kind:        30601,
				Description: "Represents a chat message",
				Fields: []generator.EventField{
					{
						FieldName: "Content",
						Type:      "string",
						Tag:       "content",
						Required:  true,
					},
					{
						FieldName: "RoomID",
						Type:      "string",
						Tag:       "room_id",
						Required:  true,
					},
					{
						FieldName: "ReplyTo",
						Type:      "string",
						Tag:       "reply_to",
						Required:  false,
					},
					{
						FieldName: "Mentions",
						Type:      "string",
						Tag:       "mentions",
						Required:  false,
						Multiple:  true,
					},
				},
			},
			{
				EventName:   "RoomEvent",
				Operation:   "room",
				Kind:        30602,
				Description: "Represents a chat room",
				Fields: []generator.EventField{
					{
						FieldName: "Name",
						Type:      "string",
						Tag:       "name",
						Required:  true,
					},
					{
						FieldName: "Description",
						Type:      "string",
						Tag:       "description",
						Required:  false,
					},
					{
						FieldName: "Members",
						Type:      "string",
						Tag:       "members",
						Required:  false,
						Multiple:  true,
					},
				},
			},
		},
	}

	// Generate the CIP implementation
	if err := generator.GenerateCIP(chatCIP); err != nil {
		log.Fatalf("Failed to generate CIP: %v", err)
	}

	fmt.Println("Successfully generated chat CIP implementation")
}
