// Package orch_score_exporter implements a Livepeer orchestrator score exporter that fetches data
// from the Livepeer orchestrator score API and exposes orchestrator score data via Prometheus metrics.
package orch_score_exporter

import (
	"fmt"
	"livepeer-exporter/constants"
	"livepeer-exporter/fetcher"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	orchScoreEndpointTemplate = "https://explorer.livepeer.org/api/score/%s"
)

// orchScore represents the structure of the data returned by the Livepeer orchestrator score API.
type orchScore struct {
	Mutex sync.Mutex

	// Response data.
	PricePerPixel   float64
	SuccessRates    map[string]float64
	RoundTripScores map[string]float64
	Scores          map[string]float64
}

// OrchScoreExporter fetches data from the Livepeer orchestrator score API and exposes it via Prometheus metrics.
type OrchScoreExporter struct {
	// Metrics.
	PricePerPixel   prometheus.Gauge
	SuccessRates    *prometheus.GaugeVec
	RoundTripScores *prometheus.GaugeVec
	Scores          *prometheus.GaugeVec

	// Config settings.
	fetchInterval    time.Duration // How often to fetch data.
	updateInterval   time.Duration // How often to update metrics.
	orchInfoEndpoint string        // The endpoint to fetch data from.

	// Data.
	orchScore *orchScore // The data returned by the API.

	// Fetchers.
	orchScoreFetcher fetcher.Fetcher
}

// initMetrics initializes the orchestrator score metrics.
func (m *OrchScoreExporter) initMetrics() {
	m.PricePerPixel = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_price_per_pixel",
			Help: "The price per pixel in Wei.",
		},
	)
	m.SuccessRates = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_success_rate",
			Help: "The success rate per region.",
		},
		[]string{"region"},
	)
	m.RoundTripScores = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_round_trip_score",
			Help: "The round trip score per region.",
		},
		[]string{"region"},
	)
	m.Scores = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_total_score",
			Help: "The total score per region.",
		},
		[]string{"region"},
	)
}

// registerMetrics registers the orchestrator score metrics with Prometheus.
func (m *OrchScoreExporter) registerMetrics() {
	prometheus.MustRegister(
		m.PricePerPixel,
		m.SuccessRates,
		m.RoundTripScores,
		m.Scores,
	)
}

// updateMetrics updates the metrics with the data fetched from the Livepeer orchestrator score API.
func (m *OrchScoreExporter) updateMetrics() {
	// Update the PricePerPixel metric
	m.PricePerPixel.Set(m.orchScore.PricePerPixel)

	// Update the SuccessRates metric
	for region, rate := range m.orchScore.SuccessRates {
		m.SuccessRates.WithLabelValues(region).Set(rate)
	}

	// Update the RoundTripScores metric
	for region, score := range m.orchScore.RoundTripScores {
		m.RoundTripScores.WithLabelValues(region).Set(score / 10)
	}

	// Update the Scores metric
	for region, score := range m.orchScore.Scores {
		m.Scores.WithLabelValues(region).Set(score / 10)
	}
}

// NewOrchScoreExporter creates a new OrchScoreExporter.
func NewOrchScoreExporter(orchAddress string, fetchInterval time.Duration, updateInterval time.Duration) *OrchScoreExporter {
	exporter := &OrchScoreExporter{
		fetchInterval:    fetchInterval,
		updateInterval:   updateInterval,
		orchInfoEndpoint: fmt.Sprintf(orchScoreEndpointTemplate, orchAddress),
		orchScore:        &orchScore{},
	}

	// Create request headers.
	headers := map[string][]string{
		"X-Device-ID": {fmt.Sprintf(constants.ClientIDTemplate, orchAddress)},
	}

	// Initialize fetcher.
	exporter.orchScoreFetcher = fetcher.Fetcher{
		URL:     exporter.orchInfoEndpoint,
		Data:    &exporter.orchScore,
		Headers: headers,
	}

	// Initialize metrics.
	exporter.initMetrics()
	exporter.registerMetrics()

	return exporter
}

// Start starts the OrchScoreExporter.
func (m *OrchScoreExporter) Start() {
	// Fetch initial data and update metrics.
	m.orchScoreFetcher.FetchData()
	m.updateMetrics()

	// Start fetcher in a goroutine.
	go func() {
		ticker := time.NewTicker(m.fetchInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchScore.Mutex.Lock()
			m.orchScoreFetcher.FetchData()
			m.orchScore.Mutex.Unlock()
		}
	}()

	// Start metrics updater in a goroutine.
	go func() {
		ticker := time.NewTicker(m.updateInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchScore.Mutex.Lock()
			m.updateMetrics()
			m.orchScore.Mutex.Unlock()
		}
	}()
}
