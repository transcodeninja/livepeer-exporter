// Package orch_delegators_exporter implements a Livepeer orchestrator elegators exporter that
// fetches data from the https://stronk.rocks/api/livepeer/getOrchestrator/ API endpoint and exposes
// information about the orchestrator's delegators via Prometheus metrics.
package orch_delegators_exporter

import (
	"fmt"
	"livepeer-exporter/fetcher"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	orchDelegatorsEndpointTemplate = "https://stronk.rocks/api/livepeer/getOrchestrator/%s"
)

// delegator represents the structure of the delegators field contained in the API response.
type delegator struct {
	ID           string
	BondedAmount string
	StartRound   string
}

// orchDelegators represents the structure of the data returned by the API.
type orchDelegators struct {
	sync.Mutex

	// Response data.
	Delegators []delegator
}

// OrchDelegatorsExporter fetches data from the API and exposes orchestrator's delegators metrics via Prometheus.
type OrchDelegatorsExporter struct {
	// Metrics.
	BondedAmount   *prometheus.GaugeVec
	StartRound     *prometheus.GaugeVec
	DelegatorCount prometheus.Gauge

	// Config settings.
	fetchInterval          time.Duration // How often to fetch data.
	updateInterval         time.Duration // How often to update metrics.
	orchDelegatorsEndpoint string        // The endpoint to fetch data from.

	// Data.
	orchDelegators *orchDelegators // The data returned by the API.

	// Fetchers.
	orchDelegatorsFetcher fetcher.Fetcher
}

// initMetrics initializes the orchestrator delegators metrics.
func (m *OrchDelegatorsExporter) initMetrics() {
	m.BondedAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_delegator_bonded_amount",
			Help: "The bonded amount for each delegator.",
		},
		[]string{"id"},
	)
	m.StartRound = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_delegator_start_round",
			Help: "The start round for each delegator.",
		},
		[]string{"id"},
	)
	m.DelegatorCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_delegator_count",
			Help: "The total number of delegators that are staked with the orchestrator.",
		},
	)
}

// registerMetrics registers the orchestrator delegators metrics with Prometheus.
func (m *OrchDelegatorsExporter) registerMetrics() {
	prometheus.MustRegister(
		m.BondedAmount,
		m.StartRound,
		m.DelegatorCount,
	)
}

// updateMetrics updates the metrics with the data fetched from the stonk.rocks orchestrator API.
func (m *OrchDelegatorsExporter) updateMetrics() {
	// Set the DelegatorCount metric by counting the length of the Delegators slice.
	m.DelegatorCount.Set(float64(len(m.orchDelegators.Delegators)))

	// Set the BondedAmount and StartRound metrics for each delegator.
	for _, delegator := range m.orchDelegators.Delegators {
		bondedAmount, _ := strconv.ParseFloat(delegator.BondedAmount, 64)
		startRound, _ := strconv.ParseFloat(delegator.StartRound, 64)

		m.BondedAmount.WithLabelValues(delegator.ID).Set(bondedAmount)
		m.StartRound.WithLabelValues(delegator.ID).Set(startRound)
	}
}

// NewOrchDelegatorsExporter creates a new OrchDelegatorsExporter.
func NewOrchDelegatorsExporter(orchAddress string, fetchInterval time.Duration, updateInterval time.Duration) *OrchDelegatorsExporter {
	exporter := &OrchDelegatorsExporter{
		fetchInterval:          fetchInterval,
		updateInterval:         updateInterval,
		orchDelegatorsEndpoint: fmt.Sprintf(orchDelegatorsEndpointTemplate, orchAddress),
		orchDelegators:         &orchDelegators{},
	}

	// Initialize fetcher.
	exporter.orchDelegatorsFetcher = fetcher.Fetcher{
		URL:  exporter.orchDelegatorsEndpoint,
		Data: &exporter.orchDelegators,
	}

	// Initialize metrics.
	exporter.initMetrics()
	exporter.registerMetrics()

	return exporter
}

// Start starts the OrchDelegatorsExporter.
func (m *OrchDelegatorsExporter) Start() {
	// Fetch initial data and update metrics.
	m.orchDelegatorsFetcher.FetchData()
	m.updateMetrics()

	// Start fetcher in a goroutine.
	go func() {
		ticker := time.NewTicker(m.fetchInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchDelegators.Mutex.Lock()
			m.orchDelegatorsFetcher.FetchData()
			m.orchDelegators.Mutex.Unlock()
		}
	}()

	// Start metrics updater in a goroutine.
	go func() {
		ticker := time.NewTicker(m.updateInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchDelegators.Mutex.Lock()
			m.updateMetrics()
			m.orchDelegators.Mutex.Unlock()
		}
	}()
}
