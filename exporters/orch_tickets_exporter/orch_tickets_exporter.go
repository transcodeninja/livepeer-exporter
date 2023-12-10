// Package orch_tickets_exporter implements a Livepeer orchestrator tickets exporter that fetches data
// from the Livepeer subgraph GraphQL API endpoint and exposes information about the orchestrator's tickets
// via Prometheus metrics.
package orch_tickets_exporter

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
	winningTicketRedeemedEventsEndpoint = constants.LivePeerSubgraphEndpoint
)

// graphqlQuery represents the GraphQL query to fetch data from the GraphQL API.
const graphqlQueryTemplate = `
{
	winningTicketRedeemedEvents(where: {recipient: "%s"}) {
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
		faceValue
	}
}
`

// winningTicketRedeemedEvent represents the structure of the winningTicketRedeemedEvent field contained in the GraphQL API response.
type winningTicketRedeemedEvent struct {
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
	FaceValue string
}

// winningTicketRedeemedResponse represents the structure of the GraphQL API response.
type winningTicketRedeemedResponse struct {
	sync.Mutex

	// Response data.
	Data struct {
		WinningTicketRedeemedEvents []winningTicketRedeemedEvent
	}
}

// OrchTicketsExporter fetches data from the API and exposes orchestrator's tickets metrics via Prometheus.
type OrchTicketsExporter struct {
	// Metrics.
	WinningTicketAmount      *prometheus.GaugeVec
	WinningTicketGasUsed     *prometheus.GaugeVec
	WinningTicketGasPrice    *prometheus.GaugeVec
	WinningTicketGasCost     *prometheus.GaugeVec
	WinningTicketBlockNumber *prometheus.GaugeVec
	WinningTicketBlockTime   *prometheus.GaugeVec
	WinningTicketRound       *prometheus.GaugeVec
	DayFees                  prometheus.Gauge
	WeekFees                 prometheus.Gauge
	ThirtyDayFees            prometheus.Gauge
	NinetyDayFees            prometheus.Gauge
	YearFees                 prometheus.Gauge
	TotalFees                prometheus.Gauge
	DayGasCost               prometheus.Gauge
	WeekGasCost              prometheus.Gauge
	ThirtyDayGasCost         prometheus.Gauge
	NinetyDayGasCost         prometheus.Gauge
	YearGasCost              prometheus.Gauge
	TotalGasCost             prometheus.Gauge

	// Config settings.
	orchAddress             string        // The orchestrator address to filter tickets by.
	fetchInterval           time.Duration // How often to fetch data.
	updateInterval          time.Duration // How often to update metrics.
	orchTicketsEndpoint     string        // The endpoint to fetch data from.
	orchTicketsGraphqlQuery string        // The GraphQL query to fetch data from the GraphQL API.

	// Data.
	orchTickets *winningTicketRedeemedResponse // The data returned by the API.

	// Fetchers.
	orchTicketsFetcher fetcher.Fetcher
}

// initMetrics initializes the orchestrator tickets metrics.
func (m *OrchTicketsExporter) initMetrics() {
	m.WinningTicketAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_winning_ticket_amount",
			Help: "The amount of ETH fees won by each ticket.",
		},
		[]string{"id"},
	)
	m.WinningTicketGasUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_winning_ticket_gas_used",
			Help: "The amount of gas used by each ticket.",
		},
		[]string{"id"},
	)
	m.WinningTicketGasPrice = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_winning_ticket_gas_price",
			Help: "The gas price for each ticket in Wei.",
		},
		[]string{"id"},
	)
	m.WinningTicketGasCost = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_winning_ticket_gas_cost",
			Help: "The cost of gas used by each ticket in Gwei.",
		},
		[]string{"id"},
	)
	m.WinningTicketBlockNumber = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_winning_ticket_block_number",
			Help: "The block number for each winning ticket.",
		},
		[]string{"id"},
	)
	m.WinningTicketBlockTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_winning_ticket_block_time",
			Help: "The block time for each winning ticket.",
		},
		[]string{"id"},
	)
	m.WinningTicketRound = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_winning_ticket_round",
			Help: "The round for each winning ticket.",
		},
		[]string{"id"},
	)
	m.DayFees = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_day_fees",
			Help: "The amount of ETH fees won by the orchestrator in the last 24 hours.",
		},
	)
	m.WeekFees = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_week_fees",
			Help: "The amount of ETH fees won by the orchestrator in the last 7 days.",
		},
	)
	m.ThirtyDayFees = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_thirty_day_fees",
			Help: "The amount of ETH fees won by the orchestrator in the last 30 days.",
		},
	)
	m.NinetyDayFees = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_ninety_day_fees",
			Help: "The amount of ETH fees won by the orchestrator in the last 90 days.",
		},
	)
	m.YearFees = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_year_fees",
			Help: "The amount of ETH fees won by the orchestrator in the last 365 days.",
		},
	)
	m.TotalFees = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_total_fees",
			Help: "The total amount of ETH fees won by the orchestrator.",
		},
	)
	m.DayGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_tickets_day_gas_cost",
			Help: "The gas cost for all ticket redeem transactions in the last 24 hours.",
		},
	)
	m.WeekGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_tickets_week_gas_cost",
			Help: "The gas cost for all ticket redeem transactions in the last 7 days.",
		},
	)
	m.ThirtyDayGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_tickets_thirty_day_gas_cost",
			Help: "The gas cost for all ticket redeem transactions in the last 30 days.",
		},
	)
	m.NinetyDayGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_tickets_ninety_day_gas_cost",
			Help: "The gas cost for all ticket redeem transactions in the last 90 days.",
		},
	)
	m.YearGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_tickets_year_gas_cost",
			Help: "The gas cost for all ticket redeem transactions in the last 365 days.",
		},
	)
	m.TotalGasCost = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_tickets_total_gas_cost",
			Help: "The total gas cost for all ticket redeem transactions.",
		},
	)
}

// registerMetrics registers the orchestrator tickets metrics with Prometheus.
func (m *OrchTicketsExporter) registerMetrics() {
	prometheus.MustRegister(
		m.WinningTicketAmount,
		m.WinningTicketGasUsed,
		m.WinningTicketGasPrice,
		m.WinningTicketGasCost,
		m.WinningTicketBlockNumber,
		m.WinningTicketBlockTime,
		m.WinningTicketRound,
		m.DayFees,
		m.WeekFees,
		m.ThirtyDayFees,
		m.NinetyDayFees,
		m.YearFees,
		m.TotalFees,
		m.DayGasCost,
		m.WeekGasCost,
		m.ThirtyDayGasCost,
		m.NinetyDayGasCost,
		m.YearGasCost,
		m.TotalGasCost,
	)
}

// updateMetrics updates the metrics with the data fetched the Livepeer subgraph GraphQL API.
func (m *OrchTicketsExporter) updateMetrics() {
	// Create required Unix timestamps.
	now := time.Now()
	dayAgo := now.AddDate(0, 0, -1)
	weekAgo := now.AddDate(0, 0, -7)
	ThirtyDaysAgo := now.AddDate(0, -1, 0)
	ninetyDaysAgo := now.AddDate(0, -3, 0)
	yearAgo := now.AddDate(-1, 0, 0)

	// Set the metrics for each ticket.
	var totalFees, totalGasCost float64
	var dayFees, weekFees, thirtyDayFees, ninetyDayFees, yearFees float64
	var dayGasCost, weekGasCost, thirtyDayGasCost, ninetyDayGasCost, yearGasCost float64
	for _, ticket := range m.orchTickets.Data.WinningTicketRedeemedEvents {
		amount, _ := strconv.ParseFloat(ticket.FaceValue, 64)
		gasUsed, _ := strconv.ParseFloat(ticket.Transaction.GasUsed, 64)
		gasPrice, _ := strconv.ParseFloat(ticket.Transaction.GasPrice, 64)
		gasCost := (gasUsed * gasPrice) / 1e9
		blockNumber, _ := strconv.ParseFloat(ticket.Transaction.BlockNumber, 64)
		blockTime, _ := strconv.ParseFloat(strconv.Itoa(ticket.Transaction.Timestamp), 64)
		round, _ := strconv.ParseFloat(ticket.Round.ID, 64)

		m.WinningTicketAmount.WithLabelValues(ticket.Transaction.ID).Set(amount)
		m.WinningTicketGasUsed.WithLabelValues(ticket.Transaction.ID).Set(gasUsed)
		m.WinningTicketGasPrice.WithLabelValues(ticket.Transaction.ID).Set(gasPrice)
		m.WinningTicketGasCost.WithLabelValues(ticket.Transaction.ID).Set(gasCost)
		m.WinningTicketBlockNumber.WithLabelValues(ticket.Transaction.ID).Set(blockNumber)
		m.WinningTicketBlockTime.WithLabelValues(ticket.Transaction.ID).Set(blockTime * 1000) // Grafana expects milliseconds.
		m.WinningTicketRound.WithLabelValues(ticket.Transaction.ID).Set(round)

		// Calculate the fees and gas costs for different periods.
		if blockTime >= float64(dayAgo.Unix()) {
			dayFees += amount
			dayGasCost += gasCost
		}
		if blockTime >= float64(weekAgo.Unix()) {
			weekFees += amount
			weekGasCost += gasCost
		}
		if blockTime >= float64(ThirtyDaysAgo.Unix()) {
			thirtyDayFees += amount
			thirtyDayGasCost += gasCost
		}
		if blockTime >= float64(ninetyDaysAgo.Unix()) {
			ninetyDayFees += amount
			ninetyDayGasCost += gasCost
		}
		if blockTime >= float64(yearAgo.Unix()) {
			yearFees += amount
			yearGasCost += gasCost
		}
		totalFees += amount
		totalGasCost += gasCost
	}

	// Set the period fees and gas costs.
	m.DayFees.Set(dayFees)
	m.WeekFees.Set(weekFees)
	m.ThirtyDayFees.Set(thirtyDayFees)
	m.NinetyDayFees.Set(ninetyDayFees)
	m.YearFees.Set(yearFees)
	m.TotalFees.Set(totalFees)
	m.DayGasCost.Set(dayGasCost)
	m.WeekGasCost.Set(weekGasCost)
	m.ThirtyDayGasCost.Set(thirtyDayGasCost)
	m.NinetyDayGasCost.Set(ninetyDayGasCost)
	m.YearGasCost.Set(yearGasCost)
	m.TotalGasCost.Set(totalGasCost)
}

// NewOrchTicketsExporter creates a new OrchTicketsExporter.
func NewOrchTicketsExporter(orchAddress string, fetchInterval time.Duration, updateInterval time.Duration) *OrchTicketsExporter {
	exporter := &OrchTicketsExporter{
		orchAddress:             orchAddress,
		fetchInterval:           fetchInterval,
		updateInterval:          updateInterval,
		orchTicketsEndpoint:     winningTicketRedeemedEventsEndpoint,
		orchTicketsGraphqlQuery: fmt.Sprintf(graphqlQueryTemplate, orchAddress),
		orchTickets:             &winningTicketRedeemedResponse{},
	}

	// Create request headers.
	headers := map[string][]string{
		"X-Device-ID": {fmt.Sprintf(constants.ClientIDTemplate, orchAddress)},
	}

	// Initialize fetcher.
	exporter.orchTicketsFetcher = fetcher.Fetcher{
		URL:     exporter.orchTicketsEndpoint,
		Data:    &exporter.orchTickets,
		Headers: headers,
	}

	// Initialize metrics.
	exporter.initMetrics()
	exporter.registerMetrics()

	return exporter
}

// Start starts the OrchTicketsExporter.
func (m *OrchTicketsExporter) Start() {
	// Fetch initial data and update metrics.
	m.orchTicketsFetcher.FetchGraphQLData(m.orchTicketsGraphqlQuery)
	m.updateMetrics()

	// Start fetcher in a goroutine.
	go func() {
		ticker := time.NewTicker(m.fetchInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchTickets.Mutex.Lock()
			m.orchTicketsFetcher.FetchGraphQLData(m.orchTicketsGraphqlQuery)
			m.orchTickets.Mutex.Unlock()
		}
	}()

	// Start metrics updater in a goroutine.
	go func() {
		ticker := time.NewTicker(m.updateInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchTickets.Mutex.Lock()
			m.updateMetrics()
			m.orchTickets.Mutex.Unlock()
		}
	}()
}
