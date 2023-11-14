// Package orch_test_streams_exporter implements a Livepeer Orchestrator Test streams exporter that fetches data from the
// https://leaderboard-serverless.vercel.app/api/raw_stats API endpoint and exposes data about the orchestrators test
// streams via Prometheus metrics.
package orch_test_streams_exporter

import (
	"fmt"
	"livepeer-exporter/fetcher"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	orchDelegatorsEndpointTemplate = "https://leaderboard-serverless.vercel.app/api/raw_stats?orchestrator=%s"
)

// TestStreams represents the data structure of the test streams field contained in the API response.
type TestStreams struct {
	Region        string
	Orchestrator  string
	SuccessRate   float64 `json:"success_rate"`
	UploadTime    float64 `json:"upload_time"`
	DownloadTime  float64 `json:"download_time"`
	TranscodeTime float64 `json:"transcode_time"`
	RoundTripTime float64 `json:"round_trip_time"`
}

// OrchTestStreams represents the structure of the data returned by the  API.
type OrchTestStreams struct {
	sync.Mutex

	// Response data.
	FRA []TestStreams
	LAX []TestStreams
	LON []TestStreams
	MDW []TestStreams
	NYC []TestStreams
	PRG []TestStreams
	SAO []TestStreams
	SIN []TestStreams
}

// TestStreamsExporter fetches data from the API endpoint and exposes data about the orchestrator's test streams via Prometheus metrics.
type TestStreamsExporter struct {
	// Metrics.
	SuccessRate   *prometheus.GaugeVec
	UploadTime    *prometheus.GaugeVec
	DownloadTime  *prometheus.GaugeVec
	TranscodeTime *prometheus.GaugeVec
	RoundTripTime *prometheus.GaugeVec

	// Config settings.
	fetchInterval           time.Duration // How often to fetch data.
	updateInterval          time.Duration // How often to update metrics.
	orchTestStreamsEndpoint string        // The endpoint to fetch data from.

	// Data.
	orchTestStreams *OrchTestStreams // The data returned by the API.

	// Fetchers.
	orchTestStreamsFetcher fetcher.Fetcher
}

// initMetrics initializes the orchestrator test streams metrics.
func (m *TestStreamsExporter) initMetrics() {
	m.SuccessRate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "livepeer_orch_test_stream_success_rate",
		Help: "Success rate per region for test streams",
	}, []string{"region", "orchestrator"})
	m.UploadTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "livepeer_orch_test_stream_upload_time",
		Help: "Upload time per region for test streams",
	}, []string{"region", "orchestrator"})
	m.DownloadTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "livepeer_orch_test_stream_download_time",
		Help: "Download time per region for test streams",
	}, []string{"region", "orchestrator"})
	m.TranscodeTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "livepeer_orch_test_stream_transcode_time",
		Help: "Transcode time per region for test streams",
	}, []string{"region", "orchestrator"})
	m.RoundTripTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "livepeer_orch_test_stream_round_trip_time",
		Help: "Round trip time per region for test streams",
	}, []string{"region", "orchestrator"})
}

// registerMetrics registers the orchestrator test streams metrics with Prometheus.
func (m *TestStreamsExporter) registerMetrics() {
	prometheus.MustRegister(
		m.SuccessRate,
		m.UploadTime,
		m.DownloadTime,
		m.TranscodeTime,
		m.RoundTripTime,
	)
}

// updateMetrics updates the metrics with the data fetched from the  'interptr-latest-test-streams' API.
func (m *TestStreamsExporter) updateMetrics() {
	for _, regionData := range []struct {
		Region      string
		TestStreams []TestStreams
	}{
		{"FRA", m.orchTestStreams.FRA},
		{"LAX", m.orchTestStreams.LAX},
		{"LON", m.orchTestStreams.LON},
		{"MDW", m.orchTestStreams.MDW},
		{"NYC", m.orchTestStreams.NYC},
		{"PRG", m.orchTestStreams.PRG},
		{"SAO", m.orchTestStreams.SAO},
		{"SIN", m.orchTestStreams.SIN},
	} {
		for _, orchData := range regionData.TestStreams {
			m.SuccessRate.WithLabelValues(regionData.Region, orchData.Orchestrator).Set(orchData.SuccessRate)
			m.UploadTime.WithLabelValues(regionData.Region, orchData.Orchestrator).Set(orchData.UploadTime)
			m.DownloadTime.WithLabelValues(regionData.Region, orchData.Orchestrator).Set(orchData.DownloadTime)
			m.TranscodeTime.WithLabelValues(regionData.Region, orchData.Orchestrator).Set(orchData.TranscodeTime)
			m.RoundTripTime.WithLabelValues(regionData.Region, orchData.Orchestrator).Set(orchData.RoundTripTime)
		}
	}
}

// NewOrchTestStreamsExporter creates a new TestStreamsExporter.
func NewOrchTestStreamsExporter(orchAddress string, fetchInterval time.Duration, updateInterval time.Duration) *TestStreamsExporter {
	exporter := &TestStreamsExporter{
		fetchInterval:           fetchInterval,
		updateInterval:          updateInterval,
		orchTestStreamsEndpoint: fmt.Sprintf(orchDelegatorsEndpointTemplate, orchAddress),
		orchTestStreams:         &OrchTestStreams{},
	}

	// Initialize fetcher.
	exporter.orchTestStreamsFetcher = fetcher.Fetcher{
		URL:  exporter.orchTestStreamsEndpoint,
		Data: &exporter.orchTestStreams,
	}

	// Initialize metrics.
	exporter.initMetrics()
	exporter.registerMetrics()

	return exporter
}

// Start starts the TestStreamsExporter.
func (m *TestStreamsExporter) Start() {
	// Fetch initial data and update metrics.
	m.orchTestStreamsFetcher.FetchData()
	m.updateMetrics()

	// Start fetcher in a goroutine.
	go func() {
		ticker := time.NewTicker(m.fetchInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchTestStreams.Mutex.Lock()
			m.orchTestStreamsFetcher.FetchData()
			m.orchTestStreams.Mutex.Unlock()
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
