package pkg

import (
	"encoding/json"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

// DeploySimulator simulates deployment messages for each commitPushed event.
type DeploySimulator struct {
	Env   string
	Delay time.Duration
}

// Handle receives a commitPushed event, waits for chosen duration and produces a commitDeployed event.
func (d DeploySimulator) Handle(msg *message.Message) ([]*message.Message, error) {
	if msg.Metadata.Get("event_type") != "commitPushed" {
		return nil, nil
	}

	event := commitPushed{}
	err := json.Unmarshal(msg.Payload, &event)
	if err != nil {
		return nil, err
	}

	baseTime, err := time.Parse(time.RFC3339, event.OccurredOn)
	if err != nil {
		return nil, err
	}

	log.Println("Received new event:", event)

	time.Sleep(d.Delay)

	payload := deploymentPayload{
		CommitID:  event.ID,
		Env:       d.Env,
		Timestamp: baseTime.Add(d.Delay).Format(time.RFC3339),
	}

	log.Println("Sending event:", payload)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	m := message.NewMessage(watermill.NewUUID(), payloadJSON)
	return []*message.Message{m}, nil
}
