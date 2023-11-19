package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

// * CommitPushed represents a commit pushed event.
type slackMessagePayload struct {
	Text string `json:"text"`
}

// * SlackMarshaller creates a new HTTP request with a JSON payload.
func SlackMarshaller(url string, msg *message.Message) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(msg.Payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// * SlackHandler receives an event and translates it into a message payload compatible with Slack REST API.
func SlackHandler(msg *message.Message) ([]*message.Message, error) {
	eventType := msg.Metadata.Get("event_type")
	text, err := slackTextByType(eventType, msg.Payload)
	if err != nil {
		return nil, err
	}

	payload := slackMessagePayload{
		Text: text,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	log.Println("Sending Slack webhook:", string(payloadJSON))

	m := message.NewMessage(watermill.NewUUID(), payloadJSON)
	return []*message.Message{m}, nil
}

// * slackTextByType returns a text message for a given event type.
func slackTextByType(eventType string, payload []byte) (string, error) {
	switch eventType {
	case "commitPushed":
		event := CommitPushed{}
		err := json.Unmarshal(payload, &event)
		if err != nil {
			return "", err
		}

		text := fmt.Sprintf(":rocket: Commit `%s` *pushed* by %s: _%s_",
			event.ID, event.Author, event.Message)
		return text, nil
	case "commitDeployed":
		event := CommitDeployed{}
		err := json.Unmarshal(payload, &event)
		if err != nil {
			return "", err
		}

		text := fmt.Sprintf(":heavy_check_mark: Commit `%s` deployed to *%s*", event.ID, event.Env)
		return text, nil
	default:
		return "", fmt.Errorf("unknown event type: %s", eventType)
	}
}
