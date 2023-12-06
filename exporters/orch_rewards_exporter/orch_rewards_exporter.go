// Package orch_rewards_exporter implements a Livepeer orchestrator rewards exporter that fetches
// data from the Livepeer subgraph GraphQL API endpoint and exposes information about the orchestrator's rewards
// via Prometheus metrics.
package orch_rewards_exporter

import (
	"fmt"
	"livepeer-exporter/fetcher"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	rewardEventsEndpoint = "https://api.thegraph.com/subgraphs/name/livepeer/arbitrum-one"
)

// graphqlQuery represents the GraphQL query to fetch data from the GraphQL API.
const graphqlQueryTemplate = `
{
	rewardEvents(where: {delegate: "%s"}) {
		transaction {
			gasUsed
			blockNumber
			timestamp
			id
		}
		round {
			id
		}
		rewardTokens
	}
}
`

// rewardEvent represents the structure of the rewardEvent field contained in the GraphQL API response.
type rewardEvent struct {
	Transaction struct {
		GasUsed     string
		BlockNumber string
		Timestamp   int
		ID          string
	}
	Round struct {
		ID string
	}
	RewardTokens string
}

// rewardEventResponse represents the structure of the GraphQL API response.
type rewardEventResponse struct {
	sync.Mutex

	// Response data.
	Data struct {
		RewardEvents []rewardEvent
	}
}

// OrchRewardsExporter fetches data from the API and exposes orchestrator's rewards metrics via Prometheus.
type OrchRewardsExporter struct {
	// Metrics.
	RewardAmount      *prometheus.GaugeVec
	RewardGasUsed     *prometheus.GaugeVec
	RewardBlockNumber *prometheus.GaugeVec
	RewardBlockTime   *prometheus.GaugeVec
	TotalReward       prometheus.Gauge
	RewardRound       *prometheus.GaugeVec

	// Config settings.
	orchAddress             string        // The orchestrator address to filter rewards by.
	fetchInterval           time.Duration // How often to fetch data.
	updateInterval          time.Duration // How often to update metrics.
	orchRewardsEndpoint     string        // The endpoint to fetch data from.
	orchRewardsGraphqlQuery string        // The GraphQL query to fetch data from the GraphQL API.

	// Data.
	orchRewards *rewardEventResponse // The data returned by the API.

	// Fetchers.
	orchRewardsFetcher fetcher.Fetcher
}

// initMetrics initializes the orchestrator rewards metrics.
func (m *OrchRewardsExporter) initMetrics() {
	m.RewardAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_amount",
			Help: "The amount of rewards earned by each transaction.",
		},
		[]string{"id"},
	)
	m.RewardGasUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_gas_used",
			Help: "The amount of gas used by each reward transaction.",
		},
		[]string{"id"},
	)
	m.RewardBlockNumber = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_block_number",
			Help: "The block number for each reward transaction.",
		},
		[]string{"id"},
	)
	m.RewardBlockTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_block_time",
			Help: "The block time for each reward transaction.",
		},
		[]string{"id"},
	)
	m.TotalReward = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_total_claimed_rewards",
			Help: "Total rewards claimed by the the orchestrator.",
		},
	)
	m.RewardRound = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_round",
			Help: "The round in which each reward was claimed.",
		},
		[]string{"id"},
	)
}

// registerMetrics registers the orchestrator rewards metrics with Prometheus.
func (m *OrchRewardsExporter) registerMetrics() {
	prometheus.MustRegister(
		m.RewardAmount,
		m.RewardGasUsed,
		m.RewardBlockNumber,
		m.RewardBlockTime,
		m.TotalReward,
		m.RewardRound,
	)
}

// updateMetrics updates the metrics with the data fetched the Livepeer subgraph GraphQL API.
func (m *OrchRewardsExporter) updateMetrics() {
	// Set the metrics for each reward.
	var total float64
	for _, reward := range m.orchRewards.Data.RewardEvents {
		amount, _ := strconv.ParseFloat(reward.RewardTokens, 64)
		gasUsed, _ := strconv.ParseFloat(reward.Transaction.GasUsed, 64)
		blockNumber, _ := strconv.ParseFloat(reward.Transaction.BlockNumber, 64)
		blockTime, _ := strconv.ParseFloat(strconv.Itoa(reward.Transaction.Timestamp), 64)
		round, _ := strconv.ParseFloat(reward.Round.ID, 64)

		m.RewardAmount.WithLabelValues(reward.Transaction.ID).Set(amount)
		m.RewardGasUsed.WithLabelValues(reward.Transaction.ID).Set(gasUsed)
		m.RewardBlockNumber.WithLabelValues(reward.Transaction.ID).Set(blockNumber)
		m.RewardBlockTime.WithLabelValues(reward.Transaction.ID).Set(blockTime * 1000) // Grafana expects milliseconds.
		m.RewardRound.WithLabelValues(reward.Transaction.ID).Set(round)

		// Calculate the total rewards.
		total += amount
	}

	// Set the total rewards metric.
	m.TotalReward.Set(total)
}

// NewOrchRewardsExporter creates a new OrchRewardsExporter.
func NewOrchRewardsExporter(orchAddress string, fetchInterval time.Duration, updateInterval time.Duration) *OrchRewardsExporter {
	exporter := &OrchRewardsExporter{
		orchAddress:             orchAddress,
		fetchInterval:           fetchInterval,
		updateInterval:          updateInterval,
		orchRewardsEndpoint:     rewardEventsEndpoint,
		orchRewardsGraphqlQuery: fmt.Sprintf(graphqlQueryTemplate, orchAddress),
		orchRewards:             &rewardEventResponse{},
	}

	// Initialize fetcher.
	exporter.orchRewardsFetcher = fetcher.Fetcher{
		URL:  exporter.orchRewardsEndpoint,
		Data: &exporter.orchRewards,
	}

	// Initialize metrics.
	exporter.initMetrics()
	exporter.registerMetrics()

	return exporter
}

// Start starts the OrchRewardsExporter.
func (m *OrchRewardsExporter) Start() {
	// Fetch initial data and update metrics.
	m.orchRewardsFetcher.FetchGraphQLData(m.orchRewardsGraphqlQuery)
	m.updateMetrics()

	// Start fetcher in a goroutine.
	go func() {
		ticker := time.NewTicker(m.fetchInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchRewards.Mutex.Lock()
			m.orchRewardsFetcher.FetchGraphQLData(m.orchRewardsGraphqlQuery)
			m.orchRewards.Mutex.Unlock()
		}
	}()

	// Start metrics updater in a goroutine.
	go func() {
		ticker := time.NewTicker(m.updateInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchRewards.Mutex.Lock()
			m.updateMetrics()
			m.orchRewards.Mutex.Unlock()
		}
	}()
}
