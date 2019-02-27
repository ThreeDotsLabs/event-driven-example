package main

import (
	"github.com/ThreeDotsLabs/event-driven-example/pkg"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/metrics"
	"github.com/ThreeDotsLabs/watermill/message"
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

	err = pkg.SetupRouter(router, config, logger)
	if err != nil {
		panic(err)
	}

	err = router.Run()
	if err != nil {
		panic(err)
	}
}
