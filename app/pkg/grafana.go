package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	stdHttp "net/http"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-http/pkg/http"
	"github.com/ThreeDotsLabs/watermill/message"
)

type grafanaAnnotationsPayload struct {
	Text    string   `json:"text"`
	Tags    []string `json:"tags"`
	Time    int64    `json:"time"`
	TimeEnd int64    `json:"timeEnd"`
}

type grafanaParams struct {
	OccurredOn string
	Text       string
	Tags       []string
}

// * GrafanaMarshaller creates a new HTTP marshaller for Grafana
func GrafanaMarshaller(credentials string) http.MarshalMessageFunc {
	return func(url string, msg *message.Message) (*stdHttp.Request, error) {
		req, err := stdHttp.NewRequest(stdHttp.MethodPost, url, bytes.NewBuffer(msg.Payload))
		if err != nil {
			return nil, err
		}

		c := strings.Split(credentials, ":")
		req.SetBasicAuth(c[0], c[1])

		req.Header.Set("Content-Type", "application/json")

		return req, nil
	}
}

// * GrafanaHandler creates a new handler for Grafana
func GrafanaHandler(msg *message.Message) ([]*message.Message, error) {
	eventType := msg.Metadata.Get("event_type")
	params, err := grafanaParamsByType(eventType, msg.Payload)
	if err != nil {
		return nil, err
	}

	parsedTime, err := time.Parse(time.RFC3339, params.OccurredOn)
	if err != nil {
		return nil, err
	}

	timestamp := parsedTime.Unix() * 1000

	payload := grafanaAnnotationsPayload{
		Text:    params.Text,
		Tags:    params.Tags,
		Time:    timestamp,
		TimeEnd: timestamp,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	log.Println("Sending Grafana webhook:", string(payloadJSON))

	m := message.NewMessage(watermill.NewUUID(), payloadJSON)
	return []*message.Message{m}, nil
}

// * grafanaParamsByType returns Grafana parameters by event type
func grafanaParamsByType(eventType string, payload []byte) (grafanaParams, error) {
	switch eventType {
	case "commitPushed":
		event := CommitPushed{}
		err := json.Unmarshal(payload, &event)
		if err != nil {
			return grafanaParams{}, err
		}

		return grafanaParams{
			OccurredOn: event.OccurredOn,
			Text:       fmt.Sprintf("Commit %s pushed by %s", event.ID, event.Author),
			Tags:       []string{"pushed"},
		}, nil
	case "commitDeployed":
		event := CommitDeployed{}
		err := json.Unmarshal(payload, &event)
		if err != nil {
			return grafanaParams{}, err
		}

		return grafanaParams{
			OccurredOn: event.OccurredOn,
			Text:       fmt.Sprintf("Commit %s deployed to %s", event.ID, event.Env),
			Tags:       []string{fmt.Sprintf("deploy-%s", event.Env)},
		}, nil
	default:
		return grafanaParams{}, fmt.Errorf("unknown event type: %s", eventType)
	}
}
