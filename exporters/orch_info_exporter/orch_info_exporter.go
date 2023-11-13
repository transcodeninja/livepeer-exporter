// Package orch_info_exporter implements a Livepeer Orchestrator Info exporter that fetches data from Livepeer's orchestrator info API and exposes info about the orchestrator via Prometheus metrics.
package orch_info_exporter

import (
	"fmt"
	"livepeer-exporter/fetcher"
	"livepeer-exporter/util"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	orchInfoEndpointTemplate       = "https://explorer.livepeer.org/_next/data/xe8lg6V7gubXcRErA1lxB/accounts/%s/orchestrating.json?account=%s"
	delegatingInfoEndpointTemplate = "https://explorer.livepeer.org/_next/data/xe8lg6V7gubXcRErA1lxB/accounts/%s/delegating.json?account=%s"
)

// OrchInfo represents the structure of the data returned by the Livepeer orchestrator info API.
type OrchInfo struct {
	Mutex sync.Mutex

	// Response data.
	PageProps struct {
		Account struct {
			Delegator struct {
				BondedAmount string
				Delegate     struct {
					TotalStake string
				}
				LastClaimRound struct {
					Id string
				}
				StartRound    string
				WithdrawnFees string
			}
			Protocol struct {
				CurrentRound struct {
					Id string
				}
			}
			Transcoder struct {
				ActivationRound string
				Active          bool
				FeeShare        string
				Pools           []struct {
					RewardTokens string
				}
				RewardCut       string
				LastRewardRound struct {
					Id string
				}
				NinetyDayVolumeETH string
				ThirtyDayVolumeETH string
				TotalVolumeETH     string
			}
		}
	}
}

// DelegatingInfo represents the structure of the data returned by the Livepeer delegator info API.
// This is used to fetch extra delegation data for the orchestrator when the `ORCHESTRATOR_ADDRESS_SECONDARY` environment variable is set.
type DelegatingInfo struct {
	Mutex sync.Mutex

	// Response data.
	PageProps struct {
		Account struct {
			Delegator struct {
				BondedAmount string
			}
		}
	}
}

// OrchInfoExporter fetches data from the Livepeer orchestrator info API and exposes info about the orchestrator via Prometheus metrics.
type OrchInfoExporter struct {
	// Metrics.
	BondedAmount       prometheus.Gauge
	TotalStake         prometheus.Gauge
	LastClaimRound     prometheus.Gauge
	StartRound         prometheus.Gauge
	WithdrawnFees      prometheus.Gauge
	CurrentRound       prometheus.Gauge
	ActivationRound    prometheus.Gauge
	Active             prometheus.Gauge
	FeeCut             prometheus.Gauge
	RewardCut          prometheus.Gauge
	LastRewardRound    prometheus.Gauge
	NinetyDayVolumeETH prometheus.Gauge
	ThirtyDayVolumeETH prometheus.Gauge
	TotalVolumeETH     prometheus.Gauge
	TotalReward        prometheus.Gauge
	OrchStake          prometheus.Gauge
	RewardCallRatio    prometheus.Gauge

	// Config settings.
	fetchInterval          time.Duration // How often to fetch data.
	updateInterval         time.Duration // How often to update metrics.
	orchAddressSecondary   string        // The secondary orchestrator address.
	orchInfoEndpoint       string        // The endpoint to fetch data from.
	delegatingInfoEndpoint string        // The endpoint to fetch extra delegation data from.

	// Data.
	orchInfo       *OrchInfo       // The data returned by the orchestrator API.
	delegatingInfo *DelegatingInfo // The data returned by the delegation API.

	// Fetchers.
	orchInfoFetcher       fetcher.Fetcher
	delegatingInfoFetcher fetcher.Fetcher
}

// initMetrics initializes the orchestrator info metrics.
func (m *OrchInfoExporter) initMetrics() {
	m.BondedAmount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_bonded_amount",
			Help: "The amount of LPT bonded to the orchestrator.",
		},
	)
	m.TotalStake = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_total_stake",
			Help: "The total stake of the orchestrator in LPT.",
		},
	)
	m.LastClaimRound = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_last_claim_round",
			Help: "The last round the orchestrator claimed fees.",
		},
	)
	m.StartRound = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_start_round",
			Help: "The round the orchestrator registered.",
		},
	)
	m.WithdrawnFees = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_withdrawn_fees",
			Help: "The amount of fees the orchestrator has withdrawn.",
		},
	)
	m.CurrentRound = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_current_round",
			Help: "The current round.",
		},
	)
	m.ActivationRound = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_activation_round",
			Help: "The round the orchestrator activated.",
		},
	)
	m.Active = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_active",
			Help: "Whether the orchestrator is active.",
		},
	)
	m.FeeCut = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_fee_cut",
			Help: "The proportion of the fees the orchestrator takes.",
		},
	)
	m.RewardCut = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_cut",
			Help: "The proportion of the block reward the orchestrator takes.",
		},
	)
	m.LastRewardRound = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_last_reward_round",
			Help: "The last round the orchestrator received a reward.",
		},
	)
	m.NinetyDayVolumeETH = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_ninety_day_volume_eth",
			Help: "The 90 day volume of ETH.",
		},
	)
	m.ThirtyDayVolumeETH = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_thirty_day_volume_eth",
			Help: "The 30 day volume of ETH.",
		},
	)
	m.TotalVolumeETH = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_total_volume_eth",
			Help: "The total volume of ETH.",
		},
	)
	m.TotalReward = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_total_reward",
			Help: "The total reward of the orchestrator.",
		},
	)
	m.OrchStake = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_stake",
			Help: "The stake provided by the orchestrator.",
		},
	)
	m.RewardCallRatio = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_reward_call_ratio",
			Help: "Ratio of reward calls to total active rounds.",
		},
	)
}

// registerMetrics registers the orchestrator info metrics with Prometheus.
func (m *OrchInfoExporter) registerMetrics() {
	prometheus.MustRegister(
		m.BondedAmount,
		m.TotalStake,
		m.LastClaimRound,
		m.StartRound,
		m.WithdrawnFees,
		m.CurrentRound,
		m.ActivationRound,
		m.Active,
		m.FeeCut,
		m.RewardCut,
		m.LastRewardRound,
		m.NinetyDayVolumeETH,
		m.ThirtyDayVolumeETH,
		m.TotalVolumeETH,
		m.TotalReward,
		m.OrchStake,
		m.RewardCallRatio,
	)
}

// updateMetrics updates the metrics with the data fetched from the Livepeer orchestrator info API.
func (m *OrchInfoExporter) updateMetrics() {
	// Parse and set the current round and activation round.
	currentRound, err := util.StringToFloat64(m.orchInfo.PageProps.Account.Protocol.CurrentRound.Id)
	if err != nil {
		log.Printf("Error parsing CurrentRound: %v", err)
	}
	activationRound, err := util.StringToFloat64(m.orchInfo.PageProps.Account.Transcoder.ActivationRound)
	if err != nil {
		log.Printf("Error parsing Active: %v", err)
	}
	m.CurrentRound.Set(currentRound)
	m.ActivationRound.Set(activationRound)

	// Convert the fee cut and reward cut to fractions.
	feeCut, err := strconv.ParseFloat(m.orchInfo.PageProps.Account.Transcoder.FeeShare, 64)
	if err != nil {
		log.Printf("Error parsing FeeShare: %v", err)
	}
	rewardCut, err := strconv.ParseFloat(m.orchInfo.PageProps.Account.Transcoder.RewardCut, 64)
	if err != nil {
		log.Printf("Error parsing RewardCut: %v", err)
	}
	m.FeeCut.Set(util.Round(1-feeCut*1e-6, 2))
	m.RewardCut.Set(util.Round(rewardCut*1e-6, 2))

	// Calculate and set the total LPT reward received by the orchestrator.
	totalReward := 0.0
	for _, pool := range m.orchInfo.PageProps.Account.Transcoder.Pools {
		rewardTokens, err := strconv.ParseFloat(pool.RewardTokens, 64)
		if err != nil {
			log.Printf("Error parsing RewardTokens: %v", err)
			continue
		}
		totalReward += rewardTokens
	}
	m.TotalReward.Set(totalReward)

	// Calculate the orchestrator's reward call ratio.
	rewardCallRatio := len(m.orchInfo.PageProps.Account.Transcoder.Pools) / int(currentRound-activationRound)
	m.RewardCallRatio.Set(float64(rewardCallRatio))

	// Calculate the total LPT that is staked by the orchestrator.
	// NOTE: This uses the orchestrator bonded amount and the secondary orchestrator bonded amount.
	orchBondedAmount, err := strconv.ParseFloat(m.orchInfo.PageProps.Account.Delegator.BondedAmount, 64)
	if err != nil {
		log.Printf("Error parsing OrchBondedAmount: %v", err)
	}
	var orchSecondaryBondedAmount float64
	if m.orchAddressSecondary != "" {
		orchSecondaryBondedAmount, err = strconv.ParseFloat(m.delegatingInfo.PageProps.Account.Delegator.BondedAmount, 64)
		if err != nil {
			log.Printf("Error parsing OrchSecondaryBondedAmount: %v", err)
		}
	}
	m.OrchStake.Set(orchBondedAmount + orchSecondaryBondedAmount)

	// Update other metrics.
	util.ParseFloatAndSetGauge(m.orchInfo.PageProps.Account.Delegator.BondedAmount, m.BondedAmount)
	util.ParseFloatAndSetGauge(m.orchInfo.PageProps.Account.Delegator.Delegate.TotalStake, m.TotalStake)
	util.ParseFloatAndSetGauge(m.orchInfo.PageProps.Account.Delegator.LastClaimRound.Id, m.LastClaimRound)
	util.ParseFloatAndSetGauge(m.orchInfo.PageProps.Account.Delegator.StartRound, m.StartRound)
	util.ParseFloatAndSetGauge(m.orchInfo.PageProps.Account.Delegator.WithdrawnFees, m.WithdrawnFees)
	m.Active.Set(util.BoolToFloat64(m.orchInfo.PageProps.Account.Transcoder.Active))
	util.ParseFloatAndSetGauge(m.orchInfo.PageProps.Account.Transcoder.LastRewardRound.Id, m.LastRewardRound)
	util.ParseFloatAndSetGauge(m.orchInfo.PageProps.Account.Transcoder.NinetyDayVolumeETH, m.NinetyDayVolumeETH)
	util.ParseFloatAndSetGauge(m.orchInfo.PageProps.Account.Transcoder.ThirtyDayVolumeETH, m.ThirtyDayVolumeETH)
	util.ParseFloatAndSetGauge(m.orchInfo.PageProps.Account.Transcoder.TotalVolumeETH, m.TotalVolumeETH)
}

// NewOrchInfoExporter creates a new OrchInfoExporter.
func NewOrchInfoExporter(orchAddress string, fetchInterval time.Duration, updateInterval time.Duration, orchAddrSecondary string) *OrchInfoExporter {
	exporter := &OrchInfoExporter{
		fetchInterval:          fetchInterval,
		updateInterval:         updateInterval,
		orchAddressSecondary:   orchAddrSecondary,
		orchInfoEndpoint:       fmt.Sprintf(orchInfoEndpointTemplate, orchAddress, orchAddress),
		delegatingInfoEndpoint: fmt.Sprintf(delegatingInfoEndpointTemplate, orchAddrSecondary, orchAddrSecondary),
		orchInfo:               &OrchInfo{},
		delegatingInfo:         &DelegatingInfo{},
	}

	// Initialize fetcher.
	exporter.orchInfoFetcher = fetcher.Fetcher{
		URL:  exporter.orchInfoEndpoint,
		Data: exporter.orchInfo,
	}
	if orchAddrSecondary != "" {
		exporter.delegatingInfoFetcher = fetcher.Fetcher{
			URL:  exporter.delegatingInfoEndpoint,
			Data: exporter.delegatingInfo,
		}
	}

	// Initialize metrics.
	exporter.initMetrics()
	exporter.registerMetrics()

	return exporter
}

// Start starts the OrchInfoExporter.
func (m *OrchInfoExporter) Start() {
	// Fetch initial data and update metrics.
	m.orchInfoFetcher.FetchData()
	if m.orchAddressSecondary != "" {
		m.delegatingInfoFetcher.FetchData()
	}
	m.updateMetrics()

	// Start fetchers in a goroutine.
	go func() {
		ticker := time.NewTicker(m.fetchInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.orchInfo.Mutex.Lock()
			m.orchInfoFetcher.FetchData()
			m.orchInfo.Mutex.Unlock()
			if m.orchAddressSecondary != "" {
				m.delegatingInfo.Mutex.Lock()
				m.delegatingInfoFetcher.FetchData()
				m.delegatingInfo.Mutex.Unlock()
			}
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
