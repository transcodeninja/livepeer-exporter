// Package orch_tickets_exporter implements a Livepeer orchestrator tickets exporter that fetches data
// from the https://stronk.rocks/api/livepeer/getAllRedeemTicketEvents API endpoint and exposes
// information about the orchestrator's tickets via Prometheus metrics.
package orch_tickets_exporter

import (
	"livepeer-exporter/fetcher"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	getRedeemedTicketsEndpoint = "https://stronk.rocks/api/livepeer/getAllRedeemTicketEvents"
)

// ticketTransaction represents the structure of the ticket transaction field contained in the API response.
type ticketTransaction struct {
	Address         string
	Amount          float64
	TransactionHash string
	BlockNumber     int
	BlockTime       int
}

// tickets represents the structure of the data returned by the API.
type tickets struct {
	sync.Mutex

	// Response data.
	Transactions []ticketTransaction
}

// OrchTicketsExporter fetches data from the API and exposes orchestrator's tickets metrics via Prometheus.
type OrchTicketsExporter struct {
	// Metrics.
	WinningTicketAmount          *prometheus.GaugeVec
	WinningTicketTransactionHash *prometheus.GaugeVec
	WinningTicketBlockNumber     *prometheus.GaugeVec
	WinningTicketBlockTime       *prometheus.GaugeVec

	// Config settings.
	orchAddress         string        // The orchestrator address to filter tickets by.
	fetchInterval       time.Duration // How often to fetch data.
	updateInterval      time.Duration // How often to update metrics.
	orchTicketsEndpoint string        // The endpoint to fetch data from.

	// Data.
	orchTickets *tickets // The data returned by the API.

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
	m.WinningTicketTransactionHash = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_winning_ticket_transaction_hash",
			Help: "The transaction hash for each winning ticket.",
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
}

// registerMetrics registers the orchestrator tickets metrics with Prometheus.
func (m *OrchTicketsExporter) registerMetrics() {
	prometheus.MustRegister(
		m.WinningTicketAmount,
		m.WinningTicketTransactionHash,
		m.WinningTicketBlockNumber,
		m.WinningTicketBlockTime,
	)
}

// updateMetrics updates the metrics with the data fetched from the stonk.rocks tickets API.
func (m *OrchTicketsExporter) updateMetrics() {
	// Filter out tickets that are not for the configured orchestrator.
	var tickets []ticketTransaction
	for _, ticket := range m.orchTickets.Transactions {
		if ticket.Address == m.orchAddress {
			tickets = append(tickets, ticket)
		}
	}

	// Set the metrics for each ticket.
	for _, ticket := range tickets {
		amount, _ := strconv.ParseFloat(strconv.FormatFloat(ticket.Amount, 'f', -1, 64), 64)
		blockNumber, _ := strconv.ParseFloat(strconv.Itoa(ticket.BlockNumber), 64)
		blockTime, _ := strconv.ParseFloat(strconv.Itoa(ticket.BlockTime), 64)

		m.WinningTicketAmount.WithLabelValues(ticket.TransactionHash).Set(amount)
		m.WinningTicketBlockNumber.WithLabelValues(ticket.TransactionHash).Set(blockNumber)
		m.WinningTicketBlockTime.WithLabelValues(ticket.TransactionHash).Set(blockTime * 1000) // Grafana expects milliseconds.
	}
}

// NewOrchTicketsExporter creates a new OrchTicketsExporter.
func NewOrchTicketsExporter(orchAddress string, fetchInterval time.Duration, updateInterval time.Duration) *OrchTicketsExporter {
	exporter := &OrchTicketsExporter{
		orchAddress:         orchAddress,
		fetchInterval:       fetchInterval,
		updateInterval:      updateInterval,
		orchTicketsEndpoint: getRedeemedTicketsEndpoint,
		orchTickets:         &tickets{},
	}

	// Initialize fetcher.
	exporter.orchTicketsFetcher = fetcher.Fetcher{
		URL:  exporter.orchTicketsEndpoint,
		Data: &exporter.orchTickets.Transactions,
	}

	// Initialize metrics.
	exporter.initMetrics()
	exporter.registerMetrics()

	return exporter
}

// Start starts the OrchTicketsExporter.
func (m *OrchTicketsExporter) Start() {
	// Fetch initial data and update metrics.
	m.orchTicketsFetcher.FetchDataWithBody(`{"smartUpdate":false}`)
	m.updateMetrics()

	// Start fetcher in a goroutine.
	go func() {
		ticker := time.NewTicker(m.fetchInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchTickets.Mutex.Lock()
			m.orchTicketsFetcher.FetchDataWithBody(`{"smartUpdate":false}`)
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
