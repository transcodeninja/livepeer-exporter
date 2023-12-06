// Package orch_tickets_exporter implements a Livepeer orchestrator tickets exporter that fetches data
// from the Livepeer subgraph GraphQL API endpoint and exposes information about the orchestrator's tickets
// via Prometheus metrics.
package orch_tickets_exporter

import (
	"fmt"
	"livepeer-exporter/fetcher"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	winningTicketRedeemedEventsEndpoint = "https://api.thegraph.com/subgraphs/name/livepeer/arbitrum-one"
)

// graphqlQuery represents the GraphQL query to fetch data from the GraphQL API.
const graphqlQueryTemplate = `
{
	winningTicketRedeemedEvents(
		where: {recipient: "%s"}
	) {
		transaction {
			gasUsed
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
	WinningTicketBlockNumber *prometheus.GaugeVec
	WinningTicketBlockTime   *prometheus.GaugeVec
	WinningTicketRound       *prometheus.GaugeVec

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
			Help: "The amount of fees won by each ticket.",
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
}

// registerMetrics registers the orchestrator tickets metrics with Prometheus.
func (m *OrchTicketsExporter) registerMetrics() {
	prometheus.MustRegister(
		m.WinningTicketAmount,
		m.WinningTicketGasUsed,
		m.WinningTicketBlockNumber,
		m.WinningTicketBlockTime,
		m.WinningTicketRound,
	)
}

// updateMetrics updates the metrics with the data fetched the Livepeer subgraph GraphQL API.
func (m *OrchTicketsExporter) updateMetrics() {
	// Set the metrics for each ticket.
	for _, ticket := range m.orchTickets.Data.WinningTicketRedeemedEvents {
		amount, _ := strconv.ParseFloat(ticket.FaceValue, 64)
		gasUsed, _ := strconv.ParseFloat(ticket.Transaction.GasUsed, 64)
		blockNumber, _ := strconv.ParseFloat(ticket.Transaction.BlockNumber, 64)
		blockTime, _ := strconv.ParseFloat(strconv.Itoa(ticket.Transaction.Timestamp), 64)
		round, _ := strconv.ParseFloat(ticket.Round.ID, 64)

		m.WinningTicketAmount.WithLabelValues(ticket.Transaction.ID).Set(amount)
		m.WinningTicketGasUsed.WithLabelValues(ticket.Transaction.ID).Set(gasUsed)
		m.WinningTicketBlockNumber.WithLabelValues(ticket.Transaction.ID).Set(blockNumber)
		m.WinningTicketBlockTime.WithLabelValues(ticket.Transaction.ID).Set(blockTime * 1000) // Grafana expects milliseconds.
		m.WinningTicketRound.WithLabelValues(ticket.Transaction.ID).Set(round)
	}
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

	// Initialize fetcher.
	exporter.orchTicketsFetcher = fetcher.Fetcher{
		URL:  exporter.orchTicketsEndpoint,
		Data: &exporter.orchTickets,
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
