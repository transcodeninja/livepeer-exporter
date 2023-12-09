// Package orch_info_exporter implements a Livepeer orchestrator info exporter that fetches data
// from the Livepeer subgraph GraphQL API endpoint and exposes info about the orchestrator via Prometheus metrics.
package orch_info_exporter

import (
	"fmt"
	"livepeer-exporter/constants"
	"livepeer-exporter/fetcher"
	"livepeer-exporter/util"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	orchInfoEndpoint = constants.LivePeerSubgraphEndpoint

	// Global variables to track whether a warning has already been logged for a invalid delegator address.
	hasLoggedNoDelegator bool
)

// graphqlQuery represents the GraphQL query to fetch data from the GraphQL API.
const graphqlQueryTemplate = `
{
	transcoder(id: "%s") {
		delegator {
			bondedAmount
			withdrawnFees
			lastClaimRound {
			id
			}
			startRound
		}
		totalStake
		lastRewardRound {
			id
		}
		activationRound
		active
		feeShare
		pools {
			rewardTokens
			round {
				id
			}
		}
		rewardCut
		lastRewardRound {
			id
		}
		ninetyDayVolumeETH
		thirtyDayVolumeETH
		totalVolumeETH
		delegators (where:{id: "%s"}){
			bondedAmount
		}
	}
	protocol(id: "0") {
		currentRound {
			id
		}
	}
}
`

// delegatingInfoResponse represents the structure of the pools field contained in the GraphQL API response.
type pool struct {
	RewardTokens string
	Round        struct {
		ID string
	}
}

// transcoderResponse represents the structure of the GraphQL API response.
type transcoderResponse struct {
	sync.Mutex

	// Response data.
	Data struct {
		Transcoder struct {
			Delegator struct {
				BondedAmount   string
				WithdrawnFees  string
				LastClaimRound struct {
					ID string
				}
				StartRound string
			}
			TotalStake      string
			LastRewardRound struct {
				ID string
			}
			ActivationRound    string
			Active             bool
			FeeShare           string
			Pools              []pool
			RewardCut          string
			NinetyDayVolumeETH string
			ThirtyDayVolumeETH string
			TotalVolumeETH     string
			Delegators         []struct {
				BondedAmount string
			}
		}
		Protocol struct {
			CurrentRound struct {
				ID string
			}
		}
	}
}

// orchInfo represents the parsed data from the the Livepeer subgraph GraphQL API.
type orchInfo struct {
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
	OrchStake          float64
	RewardCallRatio    float64
}

// getRewardCallRatio calculates the ratio of rounds in the last 30 days that the orchestrator claimed rewards.
func getRewardCallRatio(pools []pool, currentRound, activationRound int) float64 {
	// Calculate the round 30 days back
	thirtyDaysBack := currentRound - 30

	// If the activation round is less than 30 days back, use it instead
	if activationRound > thirtyDaysBack {
		thirtyDaysBack = activationRound
	}

	// Create a map of all round IDs in the pools
	poolRounds := make(map[int]bool)
	for _, pool := range pools {
		roundID, _ := strconv.Atoi(pool.Round.ID)
		poolRounds[roundID] = true
	}

	// Count the rounds from the current round to 30 days back that exist in the pools
	rewardedRounds := 0
	totalRounds := currentRound - thirtyDaysBack + 1
	for round := currentRound; round >= thirtyDaysBack; round-- {
		if poolRounds[round] {
			rewardedRounds++
		}
	}

	// Calculate and return the ratio
	return float64(rewardedRounds) / float64(totalRounds)
}

// OrchInfoExporter fetches data from the API and exposes orchestrator info via Prometheus.
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
	OrchStake          prometheus.Gauge
	RewardCallRatio    prometheus.Gauge

	// Config settings.
	fetchInterval        time.Duration // How often to fetch data.
	updateInterval       time.Duration // How often to update metrics.
	orchAddressSecondary string        // The secondary orchestrator address.
	orchInfoEndpoint     string        // The endpoint to fetch data from.
	orchInfoGraphqlQuery string        // The GraphQL query to fetch data from the GraphQL API.

	// Data.
	transcoderResponse *transcoderResponse // The data returned by the API.
	orchInfo           *orchInfo           // The data returned by the orchestrator API, parsed into a struct.

	// Fetchers.
	orchInfoFetcher fetcher.Fetcher
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
			Help: "The total amount of LPT that is staked to the orchestrator.",
		},
	)
	m.LastClaimRound = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_last_reward_claim_round",
			Help: "The last round in which the orchestrator claimed the reward.",
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
			Help: "The amount of ETH fees the orchestrator has withdrawn.",
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
			Help: "The last round the orchestrator received rewards while active.",
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
			Help: "The stake personally contributed by the orchestrator.",
		},
	)
	m.RewardCallRatio = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "livepeer_orch_thirty_day_reward_claim_ratio",
			Help: "How often an orchestrator claimed rewards in the last thirty rounds.",
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
	)
}

// parseMetrics parses the values from the transcoderResponse and delegatingInfoResponse and populates the orchInfo struct.
func (m *OrchInfoExporter) parseMetrics() {
	// Parse and set the orchestrator info.
	util.SetFloatFromStr(&m.orchInfo.BondedAmount, m.transcoderResponse.Data.Transcoder.Delegator.BondedAmount)
	util.SetFloatFromStr(&m.orchInfo.TotalStake, m.transcoderResponse.Data.Transcoder.TotalStake)
	util.SetFloatFromStr(&m.orchInfo.LastClaimRound, m.transcoderResponse.Data.Transcoder.Delegator.LastClaimRound.ID)
	util.SetFloatFromStr(&m.orchInfo.StartRound, m.transcoderResponse.Data.Transcoder.Delegator.StartRound)
	util.SetFloatFromStr(&m.orchInfo.WithdrawnFees, m.transcoderResponse.Data.Transcoder.Delegator.WithdrawnFees)
	util.SetFloatFromStr(&m.orchInfo.CurrentRound, m.transcoderResponse.Data.Protocol.CurrentRound.ID)
	util.SetFloatFromStr(&m.orchInfo.ActivationRound, m.transcoderResponse.Data.Transcoder.ActivationRound)
	m.orchInfo.Active = util.BoolToFloat64(m.transcoderResponse.Data.Transcoder.Active)
	util.SetFloatFromStr(&m.orchInfo.LastRewardRound, m.transcoderResponse.Data.Transcoder.LastRewardRound.ID)
	util.SetFloatFromStr(&m.orchInfo.NinetyDayVolumeETH, m.transcoderResponse.Data.Transcoder.NinetyDayVolumeETH)
	util.SetFloatFromStr(&m.orchInfo.ThirtyDayVolumeETH, m.transcoderResponse.Data.Transcoder.ThirtyDayVolumeETH)
	util.SetFloatFromStr(&m.orchInfo.TotalVolumeETH, m.transcoderResponse.Data.Transcoder.TotalVolumeETH)
	m.orchInfo.RewardCallRatio = getRewardCallRatio(m.transcoderResponse.Data.Transcoder.Pools, int(m.orchInfo.CurrentRound), int(m.orchInfo.ActivationRound))

	// Calculate and set reward and fee cut proportions.
	feeShare, err := util.StringToFloat64(m.transcoderResponse.Data.Transcoder.FeeShare)
	if err != nil {
		log.Printf("Error parsing fee share: %v", err)
	} else {
		m.orchInfo.FeeCut = util.Round(1-feeShare*1e-6, 2)
	}
	rewardCut, err := util.StringToFloat64(m.transcoderResponse.Data.Transcoder.RewardCut)
	if err != nil {
		log.Printf("Error parsing reward cut: %v", err)
	} else {
		m.orchInfo.RewardCut = util.Round(rewardCut*1e-6, 2)
	}

	// Calculate and set the orchestrator stake.
	// NOTE: If the orchestrator has a secondary address, we need to add the stake from the secondary address to the stake from the primary address.
	util.SetFloatFromStr(&m.orchInfo.OrchStake, m.transcoderResponse.Data.Transcoder.Delegator.BondedAmount)
	if m.orchAddressSecondary != "" {
		var secondaryStake float64
		if len(m.transcoderResponse.Data.Transcoder.Delegators) > 0 {
			util.SetFloatFromStr(&secondaryStake, m.transcoderResponse.Data.Transcoder.Delegators[0].BondedAmount)
		} else {
			secondaryStake = 0
			if !hasLoggedNoDelegator {
				log.Printf("No delegator account found for secondary address '%s'", m.orchAddressSecondary)
				hasLoggedNoDelegator = true
			}
		}
		m.orchInfo.OrchStake += secondaryStake
	}
}

// updateMetrics updates the metrics with the data fetched from the Livepeer subgraph GraphQL API.
func (m *OrchInfoExporter) updateMetrics() {
	// Parse the metrics from the response data.
	m.transcoderResponse.Mutex.Lock()
	m.parseMetrics()
	m.transcoderResponse.Mutex.Unlock()

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
}

// NewOrchInfoExporter creates a new OrchInfoExporter.
func NewOrchInfoExporter(orchAddress string, fetchInterval time.Duration, updateInterval time.Duration, orchAddrSecondary string) *OrchInfoExporter {
	exporter := &OrchInfoExporter{
		fetchInterval:        fetchInterval,
		updateInterval:       updateInterval,
		orchAddressSecondary: orchAddrSecondary,
		orchInfoEndpoint:     orchInfoEndpoint,
		orchInfoGraphqlQuery: fmt.Sprintf(graphqlQueryTemplate, orchAddress, orchAddrSecondary),
		transcoderResponse:   &transcoderResponse{},
		orchInfo:             &orchInfo{},
	}

	// Create request headers.
	headers := map[string][]string{
		"X-Device-ID": {fmt.Sprintf(constants.ClientIDTemplate, orchAddress)},
	}

	// Initialize fetcher.
	exporter.orchInfoFetcher = fetcher.Fetcher{
		URL:     exporter.orchInfoEndpoint,
		Data:    &exporter.transcoderResponse,
		Headers: headers,
	}

	// Initialize metrics.
	exporter.initMetrics()
	exporter.registerMetrics()

	return exporter
}

// Start starts the OrchInfoExporter.
func (m *OrchInfoExporter) Start() {
	// Fetch initial data and update metrics.
	m.orchInfoFetcher.FetchGraphQLData(m.orchInfoGraphqlQuery)
	m.updateMetrics()

	// Start fetchers in a goroutine.
	go func() {
		ticker := time.NewTicker(m.fetchInterval)
		defer ticker.Stop()

		for range ticker.C {
			m.transcoderResponse.Mutex.Lock()
			m.orchInfoFetcher.FetchGraphQLData(m.orchInfoGraphqlQuery)
			m.transcoderResponse.Mutex.Unlock()
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
