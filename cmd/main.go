package main

import (
	"math/rand"
	"time"

	"github.com/m110/webhooks/pkg"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/amqp"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/http"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/kafka"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func main() {
	logger := watermill.NewStdLogger(false, false)
	config := pkg.LoadConfig()

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	router.AddMiddleware(middleware.Recoverer)

	// Metrics setup
	prometheusRegistry, closeMetricsServer := metrics.CreateRegistryAndServeHTTP(config.MetricsBindAddr)
	defer closeMetricsServer()
	metricsBuilder := metrics.NewPrometheusMetricsBuilder(prometheusRegistry, "", "")
	metricsBuilder.AddPrometheusRouterMetrics(router)

	err = setupRouter(router, config, logger)
	if err != nil {
		panic(err)
	}

	err = router.Run()
	if err != nil {
		panic(err)
	}
}

func setupRouter(router *message.Router, c pkg.Config, logger watermill.LoggerAdapter) error {
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
		c.KafkaBrokers,
		kafka.DefaultMarshaler{},
		nil,
		logger)
	if err != nil {
		return err
	}

	kafkaConfig := kafka.SubscriberConfig{
		Brokers: c.KafkaBrokers,
	}
	kafkaSubscriber, err := kafka.NewSubscriber(
		kafkaConfig,
		nil,
		kafka.DefaultMarshaler{},
		logger)
	if err != nil {
		return err
	}

	grafanaPublisher, err := http.NewPublisher(
		http.PublisherConfig{
			MarshalMessageFunc: pkg.GrafanaMarshaller(c.GrafanaCredentials),
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
		pkg.GithubWebhookHandler,
	)

	router.AddHandler(
		"rabbitmq-to-kafka",
		c.AMQPQueue,
		amqpSubscriber,
		c.KafkaTopic,
		kafkaPublisher,
		pkg.AMQPHandler,
	)

	router.AddHandler(
		"kafka-to-grafana",
		c.KafkaTopic,
		kafkaSubscriber,
		c.GrafanaURL+"/api/annotations",
		grafanaPublisher,
		pkg.GrafanaHandler,
	)

	// Simulate commit deploys with delays
	stagingDelay := time.Second * time.Duration(rand.Intn(60)+30)
	productionDelay := stagingDelay + time.Second*time.Duration(rand.Intn(120)+60)

	router.AddHandler(
		"deploy-staging-simulator",
		c.KafkaTopic,
		kafkaSubscriber,
		c.AMQPQueue,
		amqpPublisher,
		pkg.DeploySimulator{"staging", stagingDelay}.Handle,
	)

	router.AddHandler(
		"deploy-production-simulator",
		c.KafkaTopic,
		kafkaSubscriber,
		c.AMQPQueue,
		amqpPublisher,
		pkg.DeploySimulator{"production", productionDelay}.Handle,
	)

	if c.SlackWebhookURL != "" {
		slackPublisher, err := http.NewPublisher(
			http.PublisherConfig{
				MarshalMessageFunc: pkg.SlackMarshaller,
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
			pkg.SlackHandler,
		)
	}

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
