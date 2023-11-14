// Package orch_rewards_exporter implements a Livepeer Orchestrator Rewards exporter that fetches data from the https://stronk.rocks/api/livepeer/getAllRewardEvents API endpoint and exposes information about the orchestrator's rewards via Prometheus metrics.
package orch_rewards_exporter

import (
	"livepeer-exporter/fetcher"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	getRewardEventsEndpoint = "https://stronk.rocks/api/livepeer/getAllRewardEvents"
)

// RewardTransaction represents the structure of the reward transaction field contained in the API response.
type RewardTransaction struct {
	Address         string
	Amount          float64
	TransactionHash string
	BlockNumber     int
	BlockTime       int
}

// Rewards represents the structure of the data returned by the API.
type Rewards struct {
	sync.Mutex

	// Response data.
	Transactions []RewardTransaction
}

// OrchRewardsExporter fetches data from the API endpoint and exposes data about the orchestrator's rewards via Prometheus metrics.
type OrchRewardsExporter struct {
	// Metrics.
	RewardAmount          *prometheus.GaugeVec
	RewardTransactionHash *prometheus.GaugeVec
	RewardBlockNumber     *prometheus.GaugeVec
	RewardBlockTime       *prometheus.GaugeVec

	// Config settings.
	orchAddress         string        // The orchestrator address to filter rewards by.
	fetchInterval       time.Duration // How often to fetch data.
	updateInterval      time.Duration // How often to update metrics.
	orchRewardsEndpoint string        // The endpoint to fetch data from.

	// Data.
	orchRewards *Rewards // The data returned by the API.

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
	m.RewardTransactionHash = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_transaction_hash",
			Help: "The transaction hash for each rewarded transaction.",
		},
		[]string{"id"},
	)
	m.RewardBlockNumber = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_block_number",
			Help: "The block number for each rewarded transaction.",
		},
		[]string{"id"},
	)
	m.RewardBlockTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_block_time",
			Help: "The block time for each rewarded transaction.",
		},
		[]string{"id"},
	)
}

// registerMetrics registers the orchestrator rewards metrics with Prometheus.
func (m *OrchRewardsExporter) registerMetrics() {
	prometheus.MustRegister(
		m.RewardAmount,
		m.RewardTransactionHash,
		m.RewardBlockNumber,
		m.RewardBlockTime,
	)
}

// updateMetrics updates the metrics with the data fetched from the stronk.rocks rewards API.
func (m *OrchRewardsExporter) updateMetrics() {
	// Filter out rewards that are not for the configured orchestrator.
	var rewards []RewardTransaction
	for _, reward := range m.orchRewards.Transactions {
		if reward.Address == m.orchAddress {
			rewards = append(rewards, reward)
		}
	}

	// Set the metrics for each reward.
	for _, reward := range rewards {
		amount, _ := strconv.ParseFloat(strconv.FormatFloat(reward.Amount, 'f', -1, 64), 64)
		blockNumber, _ := strconv.ParseFloat(strconv.Itoa(reward.BlockNumber), 64)
		blockTime, _ := strconv.ParseFloat(strconv.Itoa(reward.BlockTime), 64)

		m.RewardAmount.WithLabelValues(reward.TransactionHash).Set(amount)
		m.RewardBlockNumber.WithLabelValues(reward.TransactionHash).Set(blockNumber)
		m.RewardBlockTime.WithLabelValues(reward.TransactionHash).Set(blockTime * 1000) // Grafana expects milliseconds.
	}
}

// NewOrchRewardsExporter creates a new OrchRewardsExporter.
func NewOrchRewardsExporter(orchAddress string, fetchInterval time.Duration, updateInterval time.Duration) *OrchRewardsExporter {
	exporter := &OrchRewardsExporter{
		orchAddress:         orchAddress,
		fetchInterval:       fetchInterval,
		updateInterval:      updateInterval,
		orchRewardsEndpoint: getRewardEventsEndpoint,
		orchRewards:         &Rewards{},
	}

	// Initialize fetcher.
	exporter.orchRewardsFetcher = fetcher.Fetcher{
		URL:  exporter.orchRewardsEndpoint,
		Data: &exporter.orchRewards.Transactions,
	}

	// Initialize metrics.
	exporter.initMetrics()
	exporter.registerMetrics()

	return exporter
}

// Start starts the OrchRewardsExporter.
func (m *OrchRewardsExporter) Start() {
	// Fetch initial data and update metrics.
	m.orchRewardsFetcher.FetchDataWithBody(`{"smartUpdate":false}`)
	m.updateMetrics()

	// Start fetcher in a goroutine.
	go func() {
		ticker := time.NewTicker(m.fetchInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchRewards.Mutex.Lock()
			m.orchRewardsFetcher.FetchDataWithBody(`{"smartUpdate":false}`)
			m.orchRewards.Mutex.Unlock()
		}
	}()

	// Start metrics updater in a goroutine.
	go func() {
		ticker := time.NewTicker(m.updateInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.updateMetrics()
		}
	}()
}
