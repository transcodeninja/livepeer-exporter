// Package orch_rewards_exporter implements a Livepeer orchestrator rewards exporter that fetches
// data from the Livepeer subgraph GraphQL API endpoint and exposes information about the orchestrator's rewards
// via Prometheus metrics.
package orch_rewards_exporter

import (
	"fmt"
	"livepeer-exporter/constants"
	"livepeer-exporter/fetcher"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	rewardEventsEndpoint = constants.LivePeerSubgraphEndpoint
)

// graphqlQuery represents the GraphQL query to fetch data from the GraphQL API.
const graphqlQueryTemplate = `
{
	rewardEvents(where: {delegate: "%s"}) {
		transaction {
			gasUsed
			gasPrice
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
		GasPrice    string
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
	RewardGasPrice    *prometheus.GaugeVec
	RewardGasCost     *prometheus.GaugeVec
	RewardBlockNumber *prometheus.GaugeVec
	RewardBlockTime   *prometheus.GaugeVec
	RewardRound       *prometheus.GaugeVec
	DayRewards        prometheus.Gauge
	WeekRewards       prometheus.Gauge
	ThirtyDaysRewards prometheus.Gauge
	NinetyDaysRewards prometheus.Gauge
	YearRewards       prometheus.Gauge
	TotalRewards      prometheus.Gauge
	DayGasCost        prometheus.Gauge
	WeekGasCost       prometheus.Gauge
	ThirtyDaysGasCost prometheus.Gauge
	NinetyDaysGasCost prometheus.Gauge
	YearGasCost       prometheus.Gauge
	TotalGasCost      prometheus.Gauge

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
	m.RewardGasPrice = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_gas_price",
			Help: "The gas price for each reward transaction in Wei.",
		},
		[]string{"id"},
	)
	m.RewardGasCost = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_gas_cost",
			Help: "The gas cost for each reward transaction in Gwei.",
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
	m.RewardRound = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_round",
			Help: "The round in which each reward was claimed.",
		},
		[]string{"id"},
	)
	m.DayRewards = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_day_rewards",
			Help: "Total rewards claimed by the the orchestrator in the last 24 hours.",
		},
	)
	m.WeekRewards = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_week_rewards",
			Help: "Total rewards claimed by the the orchestrator in the last 7 days.",
		},
	)
	m.ThirtyDaysRewards = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_thirty_day_rewards",
			Help: "Total rewards claimed by the the orchestrator in the last 30 days.",
		},
	)
	m.NinetyDaysRewards = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_ninety_day_rewards",
			Help: "Total rewards claimed by the the orchestrator in the last 90 days.",
		},
	)
	m.YearRewards = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_year_rewards",
			Help: "Total rewards claimed by the the orchestrator in the last 365 days.",
		},
	)
	m.TotalRewards = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_total_rewards",
			Help: "Total rewards claimed by the the orchestrator.",
		},
	)
	m.DayGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_rewards_day_gas_cost",
			Help: "Total gas cost for all reward transactions in the last 24 hours in Gwei.",
		},
	)
	m.WeekGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_rewards_week_gas_cost",
			Help: "Total gas cost for all reward transactions in the last 7 days in Gwei.",
		},
	)
	m.ThirtyDaysGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_rewards_thirty_days_gas_cost",
			Help: "Total gas cost for all reward transactions in the last 30 days in Gwei.",
		},
	)
	m.NinetyDaysGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_rewards_ninety_days_gas_cost",
			Help: "Total gas cost for all reward transactions in the last 90 days in Gwei.",
		},
	)
	m.YearGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_rewards_year_gas_cost",
			Help: "Total gas cost for all reward transactions in the last 365 days in Gwei.",
		},
	)
	m.TotalGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_rewards_total_gas_cost",
			Help: "Total gas cost for all reward transactions in Gwei.",
		},
	)
}

// registerMetrics registers the orchestrator rewards metrics with Prometheus.
func (m *OrchRewardsExporter) registerMetrics() {
	prometheus.MustRegister(
		m.RewardAmount,
		m.RewardGasUsed,
		m.RewardGasPrice,
		m.RewardGasCost,
		m.RewardBlockNumber,
		m.RewardBlockTime,
		m.DayRewards,
		m.WeekRewards,
		m.ThirtyDaysRewards,
		m.NinetyDaysRewards,
		m.YearRewards,
		m.TotalRewards,
		m.DayGasCost,
		m.WeekGasCost,
		m.ThirtyDaysGasCost,
		m.NinetyDaysGasCost,
		m.YearGasCost,
		m.RewardRound,
	)
}

// updateMetrics updates the metrics with the data fetched the Livepeer subgraph GraphQL API.
func (m *OrchRewardsExporter) updateMetrics() {
	// Create required Unix timestamps.
	now := time.Now()
	dayAgo := now.AddDate(0, 0, -1)
	weekAgo := now.AddDate(0, 0, -7)
	ThirtyDaysAgo := now.AddDate(0, -1, 0)
	ninetyDaysAgo := now.AddDate(0, -3, 0)
	yearAgo := now.AddDate(-1, 0, 0)

	// Set the metrics for each reward.
	var totalRewards, totalGasCost float64
	var dayRewards, weekRewards, ThirtyDaysRewards, ninetyDaysRewards, yearRewards float64
	var dayGasCost, weekGasCost, monthGasCost, ninetyDaysGasCost, yearGasCost float64
	for _, reward := range m.orchRewards.Data.RewardEvents {
		amount, _ := strconv.ParseFloat(reward.RewardTokens, 64)
		gasUsed, _ := strconv.ParseFloat(reward.Transaction.GasUsed, 64)
		gasPrice, _ := strconv.ParseFloat(reward.Transaction.GasPrice, 64)
		gasCost := (gasUsed * gasPrice) / 1e9
		blockNumber, _ := strconv.ParseFloat(reward.Transaction.BlockNumber, 64)
		blockTime, _ := strconv.ParseFloat(strconv.Itoa(reward.Transaction.Timestamp), 64)
		round, _ := strconv.ParseFloat(reward.Round.ID, 64)

		m.RewardAmount.WithLabelValues(reward.Transaction.ID).Set(amount)
		m.RewardGasUsed.WithLabelValues(reward.Transaction.ID).Set(gasUsed)
		m.RewardGasPrice.WithLabelValues(reward.Transaction.ID).Set(gasPrice)
		m.RewardGasCost.WithLabelValues(reward.Transaction.ID).Set(gasCost)
		m.RewardBlockNumber.WithLabelValues(reward.Transaction.ID).Set(blockNumber)
		m.RewardBlockTime.WithLabelValues(reward.Transaction.ID).Set(blockTime * 1000) // Grafana expects milliseconds.
		m.RewardRound.WithLabelValues(reward.Transaction.ID).Set(round)

		// Calculate the rewards and gas costs for different periods.
		if blockTime >= float64(dayAgo.Unix()) {
			dayRewards += amount
			dayGasCost += gasCost
		}
		if blockTime >= float64(weekAgo.Unix()) {
			weekRewards += amount
			weekGasCost += gasCost
		}
		if blockTime >= float64(ThirtyDaysAgo.Unix()) {
			ThirtyDaysRewards += amount
			monthGasCost += gasCost
		}
		if blockTime >= float64(ninetyDaysAgo.Unix()) {
			ninetyDaysRewards += amount
			ninetyDaysGasCost += gasCost
		}
		if blockTime >= float64(yearAgo.Unix()) {
			yearRewards += amount
			yearGasCost += gasCost
		}
		totalRewards += amount
		totalGasCost += gasCost
	}

	// Set the period rewards and gas costs.
	m.DayRewards.Set(dayRewards)
	m.WeekRewards.Set(weekRewards)
	m.ThirtyDaysRewards.Set(ThirtyDaysRewards)
	m.NinetyDaysRewards.Set(ninetyDaysRewards)
	m.YearRewards.Set(yearRewards)
	m.TotalRewards.Set(totalRewards)
	m.DayGasCost.Set(dayGasCost)
	m.WeekGasCost.Set(weekGasCost)
	m.ThirtyDaysGasCost.Set(monthGasCost)
	m.NinetyDaysGasCost.Set(ninetyDaysGasCost)
	m.YearGasCost.Set(yearGasCost)
	m.TotalGasCost.Set(totalGasCost)
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

	// Create request headers.
	headers := map[string][]string{
		"X-Device-ID": {fmt.Sprintf(constants.ClientIDTemplate, orchAddress)},
	}

	// Initialize fetcher.
	exporter.orchRewardsFetcher = fetcher.Fetcher{
		URL:     exporter.orchRewardsEndpoint,
		Data:    &exporter.orchRewards,
		Headers: headers,
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
