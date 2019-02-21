package pkg

import (
	"os"
	"strings"
)

const (
	defaultBindAddr        = ":8080"
	defaultMetricsBindAddr = ":8081"

	defaultGrafanaURL         = "http://grafana:3000"
	defaultGrafanaCredentials = "admin:secret"

	defaultAMQPURI   = "amqp://guest:guest@rabbitmq:5672/"
	defaultAMQPQueue = "deployments"

	defaultKafkaBrokers = "kafka:9092"
	defaultKafkaTopic   = "events"
)

type Config struct {
	BindAddr        string
	MetricsBindAddr string

	AMQPURI   string
	AMQPQueue string

	KafkaBrokers []string
	KafkaTopic   string

	GrafanaURL         string
	GrafanaCredentials string

	SlackWebhookURL string
}

func LoadConfig() Config {
	return Config{
		MetricsBindAddr: getEnv("METRICS_BIND_ADDR", defaultMetricsBindAddr),
		BindAddr:        getEnv("BIND_ADDR", defaultBindAddr),

		AMQPURI:   getEnv("AMQP_URI", defaultAMQPURI),
		AMQPQueue: getEnv("AMQP_QUEUE", defaultAMQPQueue),

		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", defaultKafkaBrokers), ","),
		KafkaTopic:   getEnv("KAFKA_TOPIC", defaultKafkaTopic),

		GrafanaURL:         getEnv("GRAFANA_URL", defaultGrafanaURL),
		GrafanaCredentials: getEnv("GRAFANA_CREDENTIALS", defaultGrafanaCredentials),

		SlackWebhookURL: getEnv("SLACK_WEBHOOK_URL", ""),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return fallback
}
