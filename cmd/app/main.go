package main

import (
	"github.com/Furkan-Gulsen/Event-Driven-Architecture-with-Golang/config"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func main() {

	// Create a new logger.
	logger := watermill.NewStdLogger(false, false)
	c := config.LoadConfig()

	// Create a new router.
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	router.AddMiddleware(middleware.Recoverer)

	// Create a new Prometheus registry and serve it on the given address.
	prometheusRegistry, closeMetricsServer := metrics.CreateRegistryAndServeHTTP(c.MetricsBindAddr)
	defer closeMetricsServer()
	metricsBuilder := metrics.NewPrometheusMetricsBuilder(prometheusRegistry, "", "")
	metricsBuilder.AddPrometheusRouterMetrics(router)

}
