package pkg

import (
	"math/rand"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill-http/v2/pkg/http"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

func SetupRouter(router *message.Router, c Config, logger watermill.LoggerAdapter) error {
	amqpPublisher, err := amqp.NewPublisher(amqp.NewDurableQueueConfig(c.AMQPURI), logger)
	if err != nil {
		return err
	}

	amqpConfig := amqp.NewDurableQueueConfig(c.AMQPURI)
	amqpSubscriber, err := amqp.NewSubscriber(amqpConfig, logger)
	if err != nil {
		return err
	}

	httpConfig := http.SubscriberConfig{
		UnmarshalMessageFunc: http.DefaultUnmarshalMessageFunc,
	}
	httpSubscriber, err := http.NewSubscriber(c.BindAddr, httpConfig, logger)
	if err != nil {
		return err
	}

	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   c.KafkaBrokers,
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		return err
	}

	kafkaSubscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:     c.KafkaBrokers,
			Unmarshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		return err
	}

	grafanaPublisher, err := http.NewPublisher(
		http.PublisherConfig{
			MarshalMessageFunc: GrafanaMarshaller(c.GrafanaCredentials),
		}, logger)
	if err != nil {
		return err
	}

	router.AddHandler(
		"http-to-kafka",
		"/",
		httpSubscriber,
		c.KafkaTopic,
		kafkaPublisher,
		GithubWebhookHandler,
	)

	router.AddHandler(
		"rabbitmq-to-kafka",
		c.AMQPQueue,
		amqpSubscriber,
		c.KafkaTopic,
		kafkaPublisher,
		AMQPHandler,
	)

	router.AddHandler(
		"kafka-to-grafana",
		c.KafkaTopic,
		kafkaSubscriber,
		c.GrafanaURL+"/api/annotations",
		grafanaPublisher,
		GrafanaHandler,
	)

	if c.SlackWebhookURL != "" {
		slackPublisher, err := http.NewPublisher(
			http.PublisherConfig{
				MarshalMessageFunc: SlackMarshaller,
			}, logger)
		if err != nil {
			return err
		}

		router.AddHandler(
			"kafka-to-slack",
			c.KafkaTopic,
			kafkaSubscriber,
			c.SlackWebhookURL,
			slackPublisher,
			SlackHandler,
		)
	}

	// Simulate commit deploys with delays
	stagingDelay := time.Second * time.Duration(rand.Intn(60)+30)
	productionDelay := stagingDelay + time.Second*time.Duration(rand.Intn(120)+60)

	router.AddHandler(
		"deploy-staging-simulator",
		c.KafkaTopic,
		kafkaSubscriber,
		c.AMQPQueue,
		amqpPublisher,
		DeploySimulator{"staging", stagingDelay}.Handle,
	)

	router.AddHandler(
		"deploy-production-simulator",
		c.KafkaTopic,
		kafkaSubscriber,
		c.AMQPQueue,
		amqpPublisher,
		DeploySimulator{"production", productionDelay}.Handle,
	)

	go func() {
		// Start HTTP server only after the router is running
		<-router.Running()
		err = httpSubscriber.StartHTTPServer()
		if err != nil {
			panic(err)
		}
	}()

	return nil
}
