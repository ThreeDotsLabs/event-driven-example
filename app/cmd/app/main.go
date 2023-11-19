package main

import (
	"context"

	"github.com/ThreeDotsLabs/event-driven-example/config"
	"github.com/ThreeDotsLabs/event-driven-example/internal"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func main() {
	// Create a new context and add a cancel function to it.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Setup the router.
	r := &internal.Router{
		Router: router,
		Logger: logger,
		Config: c,
	}
	err = r.SetupRouter()
	if err != nil {
		panic(err)
	}

	// Run the router.
	if err := router.Run(ctx); err != nil {
		panic(err)
	}

}
