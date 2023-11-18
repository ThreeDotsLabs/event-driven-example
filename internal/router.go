package internal

import (
	"fmt"

	"github.com/Furkan-Gulsen/Event-Driven-Architecture-with-Golang/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill-http/pkg/http"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Router struct {
	Router *message.Router
	Logger watermill.LoggerAdapter
	Config config.Config
}

// * SetupAmqp creates a new AMQP publisher and subscriber.
func (r *Router) SetupAmqp() (message.Publisher, message.Subscriber, error) {
	amqpConfig := amqp.NewDurableQueueConfig(r.Config.AMQPURI)

	amqpPublisher, err := amqp.NewPublisher(amqpConfig, r.Logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AMQP publisher: %w", err)
	}

	amqpSubscriber, err := amqp.NewSubscriber(amqpConfig, r.Logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create AMQP subscriber: %w", err)
	}

	return amqpPublisher, amqpSubscriber, nil
}

// * SetupHttpSubscriber creates a new HTTP subscriber.
func (r *Router) SetupHttpSubscriber() (message.Subscriber, error) {
	httpConfig := http.SubscriberConfig{
		UnmarshalMessageFunc: http.DefaultUnmarshalMessageFunc,
	}
	httpSubscriber, err := http.NewSubscriber(r.Config.BindAddr, httpConfig, r.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP subscriber: %w", err)
	}

	return httpSubscriber, nil
}

// * SetupKafka creates a new Kafka publisher and subscriber.
func (r *Router) SetupKafka() (message.Publisher, message.Subscriber, error) {
	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   r.Config.KafkaBrokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		r.Logger,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Kafka publisher: %w", err)
	}

	kafkaSubscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               r.Config.KafkaBrokers,
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: nil,
			ConsumerGroup:         "events",
		},
		r.Logger,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Kafka subscriber: %w", err)
	}

	return kafkaPublisher, kafkaSubscriber, nil
}

// * SetupGrafana creates a new Grafana publisher and subscriber.
func (r *Router) SetupGrafana() error {
	return nil
}

func (r *Router) SetupRouter() error {
	// amqpPublisher, amqpSubscriber, err := r.SetupAmqp()
	// if err != nil {
	// 	return err
	// }

	// httpSubscriber, err := r.SetupHttpSubscriber()
	// if err != nil {
	// 	return err
	// }

	// kafkaPublisher, kafkaSubscriber, err := r.SetupKafka()
	// if err != nil {
	// 	return err
	// }

	// grafanaPublisher, err := r.SetupGrafana()

	return nil
}
