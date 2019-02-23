package pkg

import (
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type githubPayload struct {
	Commits []struct {
		ID        string `json:"id"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
		Author    struct {
			Name string `json:"name"`
		} `json:"author"`
	} `json:"commits"`
}

// GithubWebhookHandler receives GitHub webhooks and translates each commit details into a commitPushed event.
func GithubWebhookHandler(msg *message.Message) ([]*message.Message, error) {
	payload := githubPayload{}
	err := json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		return nil, err
	}

	log.Println("Received GitHub Webhook:", payload)

	var messages []*message.Message
	for _, commit := range payload.Commits {
		event := commitPushed{
			ID:         commit.ID,
			Message:    commit.Message,
			Author:     commit.Author.Name,
			OccurredOn: commit.Timestamp,
		}
		eventJSON, err := json.Marshal(event)
		if err != nil {
			return nil, err
		}

		log.Println("Sending push event", event)

		m := message.NewMessage(watermill.NewUUID(), eventJSON)
		m.Metadata.Set("event_type", "commitPushed")
		messages = append(messages, m)
	}

	return messages, nil
}
