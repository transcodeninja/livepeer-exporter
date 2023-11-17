// Provides a Livepeer metrics exporter for Prometheus.
//
// It fetches various Livepeer metrics from different endpoints and exposes them via an HTTP server.
// The server provides a '8954/metrics' endpoint for Prometheus to scrape.
//
// The exporter has the following configuration environment variables:
//   - LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS - The address of the orchestrator to fetch data from.
//   - LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS_SECONDARY - The address of the secondary orchestrator to fetch data from. Used to
//     calculate the 'livepeer_orch_stake' metric. When set the LPT stake of this address is added to the LPT stake that is bonded by the orchestrator.
//   - LIVEPEER_EXPORTER_INFO_FETCH_INTERVAL - How often to fetch general orchestrator information.
//   - LIVEPEER_EXPORTER_SCORE_FETCH_INTERVAL - How often to fetch score data for the orchestrator.
//   - LIVEPEER_EXPORTER_DELEGATORS_FETCH_INTERVAL - How often to fetch delegators data for the orchestrator.
//   - LIVEPEER_EXPORTER_TEST_STREAMS_FETCH_INTERVAL - How often to fetch the test streams data for the orchestrator.
//   - LIVEPEER_EXPORTER_TICKETS_FETCH_INTERVAL - How often to fetch tickets data for the orchestrator.
//   - LIVEPEER_EXPORTER_REWARDS_FETCH_INTERVAL - How often to fetch rewards data for the orchestrator.
//   - LIVEPEER_EXPORTER_CRYPTO_PRICES_FETCH_INTERVAL - How often to fetch crypto prices.
//   - LIVEPEER_EXPORTER_INFO_UPDATE_INTERVAL - How often to update the orchestrator info metrics.
//   - LIVEPEER_EXPORTER_SCORE_UPDATE_INTERVAL - How often to update the orchestrator score metrics.
//   - LIVEPEER_EXPORTER_DELEGATORS_UPDATE_INTERVAL - How often to update the orchestrator delegators metrics.
//   - LIVEPEER_EXPORTER_TEST_STREAMS_UPDATE_INTERVAL - How often to update the orchestrator test streams metrics.
//   - LIVEPEER_EXPORTER_TICKETS_UPDATE_INTERVAL - How often to update the orchestrator tickets metrics.
//   - LIVEPEER_EXPORTER_REWARDS_UPDATE_INTERVAL - How often to update the orchestrator rewards metrics.
//   - LIVEPEER_EXPORTER_CRYPTO_PRICES_UPDATE_INTERVAL - How often to update the crypto prices metrics.
package main

import (
	"livepeer-exporter/exporters/crypto_prices_exporter"
	"livepeer-exporter/exporters/orch_delegators_exporter"
	"livepeer-exporter/exporters/orch_info_exporter"
	"livepeer-exporter/exporters/orch_rewards_exporter"
	"livepeer-exporter/exporters/orch_score_exporter"
	"livepeer-exporter/exporters/orch_test_streams_exporter"
	"livepeer-exporter/exporters/orch_tickets_exporter"
	"livepeer-exporter/util"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Exporter default config values.
var (
	// Fetch intervals.
	infoFetchIntevalDefault         = 1 * time.Minute
	scoreFetchIntevalDefault        = 1 * time.Minute
	delegatorsFetchIntevalDefault   = 5 * time.Minute
	testStreamsFetchIntervalDefault = 1 * time.Hour
	ticketsFetchIntervalDefault     = 1 * time.Hour
	rewardsFetchIntervalDefault     = 12 * time.Hour
	cryptoPricesFetchInterval       = 1 * time.Minute

	// Update intervals.
	infoUpdateIntervalDefault         = 1 * time.Minute
	scoreUpdateIntervalDefault        = 1 * time.Minute
	delegatorsUpdateIntervalDefault   = 5 * time.Minute
	testStreamsUpdateIntervalDefault  = 1 * time.Hour
	ticketsUpdateIntervalDefault      = 1 * time.Hour
	rewardsUpdateIntervalDefault      = 12 * time.Hour
	cryptoPricesUpdateIntervalDefault = 1 * time.Minute
)

// Default config values.
func main() {
	log.Println("Starting Livepeer exporter...")

	// Retrieve orchestrator address.
	orchAddr := os.Getenv("LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS")
	if orchAddr == "" {
		log.Fatal("'LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS' environment variable should be set")
	}

	// Retrieve secondary orchestrator address.
	orchAddrSecondary := os.Getenv("LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS_SECONDARY")

	// Retrieve fetch intervals.
	infoFetchInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_INFO_FETCH_INTERVAL", infoFetchIntevalDefault)
	scoreFetchInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_SCORE_FETCH_INTERVAL", scoreFetchIntevalDefault)
	delegatorsFetchInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_DELEGATORS_FETCH_INTERVAL", delegatorsFetchIntevalDefault)
	testStreamFetchInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_TEST_STREAMS_FETCH_INTERVAL", testStreamsFetchIntervalDefault)
	ticketsFetchInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_TICKETS_FETCH_INTERVAL", ticketsFetchIntervalDefault)
	rewardsFetchInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_REWARDS_FETCH_INTERVAL", rewardsFetchIntervalDefault)
	cryptoPricesFetchInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_CRYPTO_PRICES_FETCH_INTERVAL", cryptoPricesFetchInterval)

	// Retrieve update intervals.
	infoUpdateInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_INFO_UPDATE_INTERVAL", infoUpdateIntervalDefault)
	scoreUpdateInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_SCORE_UPDATE_INTERVAL", scoreUpdateIntervalDefault)
	delegatorsUpdateInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_DELEGATORS_UPDATE_INTERVAL", delegatorsUpdateIntervalDefault)
	testStreamUpdateInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_TEST_STREAMS_UPDATE_INTERVAL", testStreamsUpdateIntervalDefault)
	ticketsUpdateInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_TICKETS_UPDATE_INTERVAL", ticketsUpdateIntervalDefault)
	rewardsUpdateInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_REWARDS_UPDATE_INTERVAL", rewardsUpdateIntervalDefault)
	cryptoPricesUpdateInterval := util.GetEnvDuration("LIVEPEER_EXPORTER_CRYPTO_PRICES_UPDATE_INTERVAL", cryptoPricesUpdateIntervalDefault)

	// Setup sub-exporters.
	log.Println("Setting up sub exporters...")
	orchInfoExporter := orch_info_exporter.NewOrchInfoExporter(orchAddr, infoFetchInterval, infoUpdateInterval, orchAddrSecondary)
	orchScoreExporter := orch_score_exporter.NewOrchScoreExporter(orchAddr, scoreFetchInterval, scoreUpdateInterval)
	orchDelegatorsExporter := orch_delegators_exporter.NewOrchDelegatorsExporter(orchAddr, delegatorsFetchInterval, delegatorsUpdateInterval)
	orchTestStreamsExporter := orch_test_streams_exporter.NewOrchTestStreamsExporter(orchAddr, testStreamFetchInterval, testStreamUpdateInterval)
	orchTicketsExporter := orch_tickets_exporter.NewOrchTicketsExporter(orchAddr, ticketsFetchInterval, ticketsUpdateInterval)
	orchRewardsExporter := orch_rewards_exporter.NewOrchRewardsExporter(orchAddr, rewardsFetchInterval, rewardsUpdateInterval)
	cryptoPricesExporter := crypto_prices_exporter.NewCryptoPricesExporter(cryptoPricesFetchInterval, cryptoPricesUpdateInterval)

	// Start sub-exporters.
	log.Println("Starting sub exporters...")
	go orchInfoExporter.Start()
	go orchScoreExporter.Start()
	go orchDelegatorsExporter.Start()
	go orchTestStreamsExporter.Start()
	go orchTicketsExporter.Start()
	go orchRewardsExporter.Start()
	go cryptoPricesExporter.Start()

	// Expose the registered metrics via HTTP.
	log.Println("Exposing metrics via HTTP on port 9153")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9155", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
