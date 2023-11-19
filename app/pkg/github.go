package pkg

import (
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

// * CommitPushed is the event that is sent when a commit is pushed to GitHub
type githubPushEvent struct {
	Commits []struct {
		ID        string `json:"id"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"`
		Author    struct {
			Name string `json:"name"`
		} `json:"author"`
	} `json:"commits"`
}

// * GithubWebhookHandler handles GitHub webhook events
func GithubWebhookHandler(msg *message.Message) ([]*message.Message, error) {
	pushEvent := githubPushEvent{}
	err := json.Unmarshal(msg.Payload, &pushEvent)
	if err != nil {
		return nil, err
	}

	log.Println("Received GitHub Webhook:", pushEvent)

	var messages []*message.Message
	for _, commit := range pushEvent.Commits {
		event := CommitPushed{
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
