package pkg

import (
	"encoding/json"
	"log"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type deploymentPayload struct {
	CommitID  string `json:"commit_id"`
	Env       string `json:"env"`
	Timestamp string `json:"timestamp"`
}

// AMQPHandler receives a deployment message over AMQP and translates it into a commitDeployed event.
func AMQPHandler(msg *message.Message) ([]*message.Message, error) {
	payload := deploymentPayload{}
	err := json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		return nil, err
	}

	log.Println("Received AMQP message: ", payload)

	event := commitDeployed{
		ID:         payload.CommitID,
		Env:        payload.Env,
		OccurredOn: payload.Timestamp,
	}
	eventJSON, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	log.Println("Sending event:", event)

	m := message.NewMessage(watermill.NewUUID(), eventJSON)
	m.Metadata.Set("event_type", "commitDeployed")
	return []*message.Message{m}, nil
}
