// Package orch_info_exporter implements a Livepeer Orchestrator Info exporter that fetches data from Livepeer's orchestrator info API and exposes
// info about the orchestrator via Prometheus metrics.
package orch_info_exporter

import (
	"fmt"
	"livepeer-exporter/fetcher"
	"livepeer-exporter/util"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	orchInfoEndpointTemplate       = "https://explorer.livepeer.org/_next/data/xe8lg6V7gubXcRErA1lxB/accounts/%s/orchestrating.json?account=%s"
	delegatingInfoEndpointTemplate = "https://explorer.livepeer.org/_next/data/xe8lg6V7gubXcRErA1lxB/accounts/%s/delegating.json?account=%s"
)

// OrchInfoResponse represents the structure of the data returned by the Livepeer orchestrator info API.
type OrchInfoResponse struct {
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

// OrchInfo represents the parsed data from the Livepeer orchestrator info API.
type OrchInfo struct {
	BondedAmount       float64
	TotalStake         float64
	LastClaimRound     float64
	StartRound         float64
	WithdrawnFees      float64
	CurrentRound       float64
	ActivationRound    float64
	Active             float64
	FeeCut             float64
	RewardCut          float64
	LastRewardRound    float64
	NinetyDayVolumeETH float64
	ThirtyDayVolumeETH float64
	TotalVolumeETH     float64
	TotalReward        float64
	OrchStake          float64
	RewardCallRatio    float64
}

// DelegationInfoResponse represents the structure of the data returned by the Livepeer delegator info API. This is used to fetch extra delegation
// data for the orchestrator when the `ORCHESTRATOR_ADDRESS_SECONDARY` environment variable is set.
type DelegationInfoResponse struct {
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
	orchInfoResponse       *OrchInfoResponse       // The data returned by the orchestrator API.
	delegatingInfoResponse *DelegationInfoResponse // The data returned by the delegation API.
	orchInfo               *OrchInfo               // The data returned by the orchestrator API, parsed into a struct.

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
	m.TotalReward = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_total_reward",
			Help: "The total reward of the orchestrator.",
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
		m.OrchStake,
		m.RewardCallRatio,
		m.TotalReward,
	)
}

// parseMetrics parses the metrics from the orchInfoResponse and delegatingInfoResponse and populates the orchInfo struct.
func (m *OrchInfoExporter) parseMetrics() {
	// Parse and set the orchestrator info.
	util.SetFloatFromStr(&m.orchInfo.BondedAmount, m.orchInfoResponse.PageProps.Account.Delegator.BondedAmount, -1)
	util.SetFloatFromStr(&m.orchInfo.TotalStake, m.orchInfoResponse.PageProps.Account.Delegator.Delegate.TotalStake, -1)
	util.SetFloatFromStr(&m.orchInfo.LastClaimRound, m.orchInfoResponse.PageProps.Account.Delegator.LastClaimRound.Id, -1)
	util.SetFloatFromStr(&m.orchInfo.StartRound, m.orchInfoResponse.PageProps.Account.Delegator.StartRound, -1)
	util.SetFloatFromStr(&m.orchInfo.WithdrawnFees, m.orchInfoResponse.PageProps.Account.Delegator.WithdrawnFees, -1)
	util.SetFloatFromStr(&m.orchInfo.CurrentRound, m.orchInfoResponse.PageProps.Account.Protocol.CurrentRound.Id, -1)
	util.SetFloatFromStr(&m.orchInfo.ActivationRound, m.orchInfoResponse.PageProps.Account.Transcoder.ActivationRound, -1)
	m.orchInfo.Active = util.BoolToFloat64(m.orchInfoResponse.PageProps.Account.Transcoder.Active)
	util.SetFloatFromStr(&m.orchInfo.FeeCut, m.orchInfoResponse.PageProps.Account.Transcoder.FeeShare, 2)
	util.SetFloatFromStr(&m.orchInfo.RewardCut, m.orchInfoResponse.PageProps.Account.Transcoder.RewardCut, 2)
	util.SetFloatFromStr(&m.orchInfo.LastRewardRound, m.orchInfoResponse.PageProps.Account.Transcoder.LastRewardRound.Id, -1)
	util.SetFloatFromStr(&m.orchInfo.NinetyDayVolumeETH, m.orchInfoResponse.PageProps.Account.Transcoder.NinetyDayVolumeETH, -1)
	util.SetFloatFromStr(&m.orchInfo.ThirtyDayVolumeETH, m.orchInfoResponse.PageProps.Account.Transcoder.ThirtyDayVolumeETH, -1)
	util.SetFloatFromStr(&m.orchInfo.TotalVolumeETH, m.orchInfoResponse.PageProps.Account.Transcoder.TotalVolumeETH, -1)

	// Calculate and set the orchestrator stake.
	// NOTE: If the orchestrator has a secondary address, we need to add the stake from the secondary address to the stake from the primary address.
	util.SetFloatFromStr(&m.orchInfo.OrchStake, m.orchInfoResponse.PageProps.Account.Delegator.BondedAmount, -1)
	if m.orchAddressSecondary != "" {
		var secondaryStake float64
		util.SetFloatFromStr(&secondaryStake, m.delegatingInfoResponse.PageProps.Account.Delegator.BondedAmount, -1)
		m.orchInfo.OrchStake += secondaryStake
	}

	// Calculate and set the total reward and reward call ratio.
	totalReward := 0.0
	for _, pool := range m.orchInfoResponse.PageProps.Account.Transcoder.Pools {
		rewardTokens, err := util.StringToFloat64(pool.RewardTokens)
		if err == nil {
			totalReward += rewardTokens
		}
	}
	m.orchInfo.TotalReward = totalReward
	if m.orchInfo.CurrentRound > m.orchInfo.ActivationRound {
		m.orchInfo.RewardCallRatio = float64(len(m.orchInfoResponse.PageProps.Account.Transcoder.Pools)) / float64(int(m.orchInfo.CurrentRound-m.orchInfo.ActivationRound))
	}
}

// updateMetrics updates the metrics with the data fetched from the Livepeer orchestrator info API.
func (m *OrchInfoExporter) updateMetrics() {
	// Parse the metrics from the response data.
	m.parseMetrics()

	// Set the metrics.
	m.BondedAmount.Set(m.orchInfo.BondedAmount)
	m.TotalStake.Set(m.orchInfo.TotalStake)
	m.LastClaimRound.Set(m.orchInfo.LastClaimRound)
	m.StartRound.Set(m.orchInfo.StartRound)
	m.WithdrawnFees.Set(m.orchInfo.WithdrawnFees)
	m.CurrentRound.Set(m.orchInfo.CurrentRound)
	m.ActivationRound.Set(m.orchInfo.ActivationRound)
	m.Active.Set(m.orchInfo.Active)
	m.FeeCut.Set(m.orchInfo.FeeCut)
	m.RewardCut.Set(m.orchInfo.RewardCut)
	m.LastRewardRound.Set(m.orchInfo.LastRewardRound)
	m.NinetyDayVolumeETH.Set(m.orchInfo.NinetyDayVolumeETH)
	m.ThirtyDayVolumeETH.Set(m.orchInfo.ThirtyDayVolumeETH)
	m.TotalVolumeETH.Set(m.orchInfo.TotalVolumeETH)
	m.OrchStake.Set(m.orchInfo.OrchStake)
	m.RewardCallRatio.Set(m.orchInfo.RewardCallRatio)
	m.TotalReward.Set(m.orchInfo.TotalReward)
}

// NewOrchInfoExporter creates a new OrchInfoExporter.
func NewOrchInfoExporter(orchAddress string, fetchInterval time.Duration, updateInterval time.Duration, orchAddrSecondary string) *OrchInfoExporter {
	exporter := &OrchInfoExporter{
		fetchInterval:          fetchInterval,
		updateInterval:         updateInterval,
		orchAddressSecondary:   orchAddrSecondary,
		orchInfoEndpoint:       fmt.Sprintf(orchInfoEndpointTemplate, orchAddress, orchAddress),
		delegatingInfoEndpoint: fmt.Sprintf(delegatingInfoEndpointTemplate, orchAddrSecondary, orchAddrSecondary),
		orchInfoResponse:       &OrchInfoResponse{},
		delegatingInfoResponse: &DelegationInfoResponse{},
		orchInfo:               &OrchInfo{},
	}

	// Initialize fetcher.
	exporter.orchInfoFetcher = fetcher.Fetcher{
		URL:  exporter.orchInfoEndpoint,
		Data: exporter.orchInfoResponse,
	}
	if orchAddrSecondary != "" {
		exporter.delegatingInfoFetcher = fetcher.Fetcher{
			URL:  exporter.delegatingInfoEndpoint,
			Data: exporter.delegatingInfoResponse,
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
			m.orchInfoResponse.Mutex.Lock()
			m.orchInfoFetcher.FetchData()
			m.orchInfoResponse.Mutex.Unlock()
			if m.orchAddressSecondary != "" {
				m.delegatingInfoResponse.Mutex.Lock()
				m.delegatingInfoFetcher.FetchData()
				m.delegatingInfoResponse.Mutex.Unlock()
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
