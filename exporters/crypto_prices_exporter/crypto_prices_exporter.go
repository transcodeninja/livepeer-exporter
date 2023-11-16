// Package crypto_prices_exporter implements a crypto prices exporter that fetches data from the https://api.coinbase.com/v2/exchange-rates?currency=USD API endpoint and exposes information about several crypto currencies that
// are relevant to Livepeer.
package crypto_prices_exporter

import (
	"livepeer-exporter/fetcher"
	"livepeer-exporter/util"
	"log"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	getCryptoPricesEndpoint = "https://api.coinbase.com/v2/exchange-rates?currency=USD"
)

// CryptoPricesResponse represents the structure of the data returned by the API.
type CryptoPricesResponse struct {
	sync.Mutex

	Data struct {
		Currency string
		Rates    map[string]string `json:"rates"`
	}
}

// cryptoPrices represents the structure of the data returned by the API, parsed into a struct.
type cryptoPrices struct {
	LPTUSDPrice float64
	ETHUSDPrice float64
	LPTEURPrice float64
	ETHEURPrice float64
}

// CryptoPricesExporter fetches data from the  API endpoint and exposes data about the crypto prices via Prometheus metrics.
type CryptoPricesExporter struct {
	// Metrics.
	LPTPrice *prometheus.GaugeVec
	ETHPrice *prometheus.GaugeVec

	// Config settings.
	fetchInterval        time.Duration // How often to fetch data.
	updateInterval       time.Duration // How often to update metrics.
	cryptoPricesEndpoint string        // The endpoint to fetch data from.

	// Data.
	cryptoPricesResponse *CryptoPricesResponse // The data returned by the API.
	cryptoPrices         *cryptoPrices         // The data returned by the  API, parsed into a struct.

	// Fetchers.
	cryptoPricesFetcher fetcher.Fetcher
}

// initMetrics initializes the crypto prices metrics.
func (m *CryptoPricesExporter) initMetrics() {
	m.LPTPrice = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "LPT_price",
		Help: "LPT token price.",
	}, []string{"currency"})
	m.ETHPrice = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ETH_price",
		Help: "Ethereum  price.",
	}, []string{"currency"})
}

// registerMetrics registers the crypto prices metrics with Prometheus.
func (m *CryptoPricesExporter) registerMetrics() {
	prometheus.MustRegister(
		m.LPTPrice,
		m.ETHPrice,
	)
}

// parseMetrics parses the values from the cryptoResponse and populates the cryptoPricesResponse struct.
func (m *CryptoPricesExporter) parseMetrics() {
	// Retrieve dollar prices.
	LPTUSDPrice, err := util.StringToFloat64(m.cryptoPricesResponse.Data.Rates["LPT"])
	if err != nil {
		log.Printf("Error trying to parse LPT price: %v", err)
		return
	}
	ETHUSDPrice, err := util.StringToFloat64(m.cryptoPricesResponse.Data.Rates["ETH"])
	if err != nil {
		log.Printf("Error trying to parse ETH price: %v", err)
		return
	}
	m.cryptoPrices.LPTUSDPrice = 1 / LPTUSDPrice
	m.cryptoPrices.ETHUSDPrice = 1 / ETHUSDPrice

	// Calculate prices in euros.
	USDToEUR, err := util.StringToFloat64(m.cryptoPricesResponse.Data.Rates["EUR"])
	if err != nil {
		log.Printf("Error trying to parse USD to EUR conversion rate: %v", err)
		return
	}
	m.cryptoPrices.LPTEURPrice = m.cryptoPrices.LPTUSDPrice * USDToEUR
	m.cryptoPrices.ETHEURPrice = m.cryptoPrices.ETHUSDPrice * USDToEUR
}

// updateMetrics updates the metrics with the data fetched from the Coinbase exchange-rates API.
func (m *CryptoPricesExporter) updateMetrics() {
	// Parse the metrics from the response data.
	m.parseMetrics()

	// Set the metrics.
	m.LPTPrice.WithLabelValues("USD").Set(m.cryptoPrices.LPTUSDPrice)
	m.LPTPrice.WithLabelValues("EUR").Set(m.cryptoPrices.LPTEURPrice)
	m.ETHPrice.WithLabelValues("USD").Set(m.cryptoPrices.ETHUSDPrice)
	m.ETHPrice.WithLabelValues("EUR").Set(m.cryptoPrices.ETHEURPrice)
}

// NewCryptoPricesExporter creates a new CryptoPricesExporter.
func NewCryptoPricesExporter(fetchInterval time.Duration, updateInterval time.Duration) *CryptoPricesExporter {
	exporter := &CryptoPricesExporter{
		fetchInterval:        fetchInterval,
		updateInterval:       updateInterval,
		cryptoPricesEndpoint: getCryptoPricesEndpoint,
		cryptoPricesResponse: &CryptoPricesResponse{},
		cryptoPrices:         &cryptoPrices{},
	}

	// Initialize fetcher.
	exporter.cryptoPricesFetcher = fetcher.Fetcher{
		URL:  exporter.cryptoPricesEndpoint,
		Data: &exporter.cryptoPricesResponse,
	}

	// Initialize metrics.
	exporter.initMetrics()
	exporter.registerMetrics()

	return exporter
}

func (m *CryptoPricesExporter) Start() {
	// Fetch initial data and update metrics.
	m.cryptoPricesFetcher.FetchData()
	m.updateMetrics()

	// Start fetchers in a goroutine.
	go func() {
		ticker := time.NewTicker(m.fetchInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.cryptoPricesResponse.Mutex.Lock()
			m.cryptoPricesFetcher.FetchData()
			m.cryptoPricesResponse.Mutex.Unlock()
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
