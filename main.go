// Provides a Livepeer metrics exporter for Prometheus.
//
// It fetches various Livepeer metrics from different endpoints and exposes them via an HTTP server.
// The server provides a '8954/metrics' endpoint for Prometheus to scrape.
//
// The exporter has the following configuration environment variables.
//   - ORCHESTRATOR_ADDRESS - The address of the orchestrator to fetch data from.
//   - ORCHESTRATOR_ADDRESS_SECONDARY - The address of the secondary orchestrator to fetch data from. Used to
//     calculate the 'livepeer_orch_stake' metric. When set the LPT stake of this address is added to the LPT stake that is bonded by the orchestrator.
//   - FETCH_INTERVAL - How often to fetch data from the orchestrator.
//   - FETCH_TEST_STREAMS_INTERVAL - How often to fetch test streams data from the orchestrator. Implemented as a separate interval because the
//     test streams API takes a long time to respond.
//   - UPDATE_INTERVAL - How often to update metrics.
package main

import (
	"livepeer-exporter/exporters/orch_delegators_exporter"
	"livepeer-exporter/exporters/orch_info_exporter"
	"livepeer-exporter/exporters/orch_score_exporter"
	"livepeer-exporter/exporters/orch_test_streams_exporter"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Exporter default config values.
var (
	fetchIntervalDefault            = 1 * time.Minute
	testStreamsFetchIntervalDefault = 15 * time.Minute
	updateIntervalDefault           = 30 * time.Second
)

// Default config values.
func main() {
	log.Println("Starting Livepeer exporter...")

	// Retrieve orchestrator address.
	orchAddr := os.Getenv("ORCHESTRATOR_ADDRESS")
	if orchAddr == "" {
		log.Fatal("'ORCHESTRATOR_ADDRESS' environment variable should be set")
	}

	// Retrieve secondary orchestrator address.
	orchAddrSecondary := os.Getenv("ORCHESTRATOR_ADDRESS_SECONDARY")

	// Retrieve fetch interval.
	fetchIntervalStr := os.Getenv("FETCH_INTERVAL")
	var fetchInterval time.Duration
	if fetchIntervalStr == "" {
		fetchInterval = fetchIntervalDefault
	} else {
		var err error
		fetchInterval, err = time.ParseDuration(fetchIntervalStr)
		if err != nil {
			log.Fatalf("failed to parse 'FETCH_INTERVAL' environment variable: %v", err)
		}
	}

	// Retrieve test stream fetch interval.
	// NOTE: This is a separate interval because the test streams API takes a long time to respond.
	fetchTestStreamsIntervalStr := os.Getenv("FETCH_TEST_STREAMS_INTERVAL")
	var fetchTestStreamsInterval time.Duration
	if fetchTestStreamsIntervalStr == "" {
		fetchTestStreamsInterval = testStreamsFetchIntervalDefault
	} else {
		var err error
		fetchTestStreamsInterval, err = time.ParseDuration(fetchTestStreamsIntervalStr)
		if err != nil {
			log.Fatalf("failed to parse 'FETCH_TEST_STREAMS_INTERVAL' environment variable: %v", err)
		}
	}

	// Retrieve update interval.
	updateIntervalStr := os.Getenv("UPDATE_INTERVAL")
	var updateInterval time.Duration
	if updateIntervalStr == "" {
		updateInterval = updateIntervalDefault
	} else {
		var err error
		updateInterval, err = time.ParseDuration(updateIntervalStr)
		if err != nil {
			log.Fatalf("failed to parse 'UPDATE_INTERVAL' environment variable: %v", err)
		}
	}

	// Setup exporters.
	log.Println("Setting up exporters...")
	orchInfoExporter := orch_info_exporter.NewOrchInfoExporter(orchAddr, fetchInterval, updateInterval, orchAddrSecondary)
	orchScoreExporter := orch_score_exporter.NewOrchScoreExporter(orchAddr, fetchInterval, updateInterval)
	orchDelegatorsExporter := orch_delegators_exporter.NewOrchDelegatorsExporter(orchAddr, fetchInterval, updateInterval)
	orchTestStreamsExporter := orch_test_streams_exporter.NewOrchTestStreamsExporter(orchAddr, fetchTestStreamsInterval, updateInterval)

	// Start exporters.
	log.Println("Starting exporters...")
	orchInfoExporter.Start()
	orchScoreExporter.Start()
	orchDelegatorsExporter.Start()
	orchTestStreamsExporter.Start()

	// Expose the registered metrics via HTTP.
	log.Println("Exposing metrics via HTTP...")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9153", nil)
}
