# Livepeer Exporter

Livepeer Exporter is a Prometheus exporter for [Livepeer](https://livepeer.org/) metrics. It fetches various metrics from different Livepeer endpoints and exposes them via an HTTP server for Prometheus to scrape.

## Metrics

This exporter comprises the following sub-exporters, each responsible for fetching specific metrics:

- [orch_delegators_exporter](./exporters/orch_delegators_exporter/): Gathers metrics related to the delegators of the designated Livepeer orchestrator.
- [orch_info_exporter](./exporters/orch_info_exporter/): Collects metrics pertaining to the Livepeer orchestrator.
- [orch_score_exporter](./exporters/orch_score_exporter/): Retrieves metrics concerning the Livepeer orchestrator's score.
- [orch_test_streams_exporter](./exporters/orch_test_streams_exporter/): Procures metrics about the Livepeer orchestrator's test streams.

These sub-exporters operate concurrently in separate [goroutines](https://go.dev/tour/concurrency/1) for enhanced performance. They fetch metrics from various Livepeer endpoints and expose them via the `9153/metrics` endpoint. For detailed information about these sub-exporters and the metrics they provide, refer to the sections below.

### orch_delegators_exporter

Fetches metrics about the delegators of the set Livepeer orchestrator from the https://stronk.rocks/api/livepeer/getOrchestrator/ endpoint. It exposes the following metrics:

**Gauge metrics:**

- `livepeer_orch_delegator_count`: Total number of delegators that stake with the Livepeer orchestrator.

**GaugeVec metrics:**

- `livepeer_orch_delegator_bonded_amount`: Bonded amount of a delegator address. This GaugeVec contains the label `id`.
- `livepeer_orch_delegator_start_round`: Start round of a delegator address. This GaugeVec contains the label `id`.

### orch_info_exporter

Fetches metrics about the Livepeer orchestrator from the [Livepeer Orchestrator API](https://explorer.livepeer.org/_next/data/xe8lg6V7gubXcRErA1lxB/accounts/%s/orchestrating.json) and [Livepeer Delegating API](https://explorer.livepeer.org/_next/data/xe8lg6V7gubXcRErA1lxB/accounts/%s/delegating.json) endpoints. It exposes the following metrics:

**Gauge metrics:**

- `livepeer_orch_bonded_amount`: Amount of LPT bonded to the orchestrator.
- `livepeer_orch_total_stake`: Total stake of the orchestrator in LPT.
- `livepeer_orch_last_claim_round`: Last round the orchestrator claimed fees.
- `livepeer_orch_start_round`: Round the orchestrator registered.
- `livepeer_orch_withdrawn_fees`: Amount of fees the orchestrator has withdrawn.
- `livepeer_orch_current_round`: Current round.
- `livepeer_orch_activation_round`: Round the orchestrator activated.
- `livepeer_orch_active`: Whether the orchestrator is active.
- `livepeer_orch_fee_cut`: Proportion of the fees the orchestrator takes.
- `livepeer_orch_reward_cut`: Proportion of the block reward the orchestrator takes.
- `livepeer_orch_last_reward_round`: Last round the orchestrator received a reward.
- `livepeer_orch_ninety_day_volume_eth`: 90-day volume of ETH.
- `livepeer_orch_thirty_day_volume_eth`: 30-day volume of ETH.
- `livepeer_orch_total_volume_eth`: Total volume of ETH.
- `livepeer_orch_total_reward`: Total reward of the orchestrator.
- `livepeer_orch_stake`: Stake provided by the orchestrator.

### orch_score_exporter

Fetches metrics about the Livepeer orchestrator's score from the [Livepeer Score API](https://explorer.livepeer.org/api/score/) endpoint. It exposes the following metrics:

**Gauge metrics:**

- `livepeer_orch_price_per_pixel`: Price per pixel.

**GaugeVec metrics:**

- `livepeer_orch_success_rates`: Success rates per region. This GaugeVec contains the label `region`.
- `livepeer_orch_round_trip_scores`: Round trip scores per region. This GaugeVec contains the label `region`.
- `livepeer_orch_scores`: Scores per region. This GaugeVec contains the label `region`.

### orch_test_streams_exporter

Fetches metrics about the LivePeer orchestrator's test streams from the https://leaderboard-serverless.vercel.app/api/raw_stats API endpoint. It exposes the following metrics:

**GaugeVec metrics:**

- `livepeer_orch_test_stream_success_rate`: Success rate per region for test streams. This GaugeVec contains the labels `region` and `orchestrator`.
- `livepeer_orch_test_stream_upload_time`: Upload time per region for test streams. This GaugeVec contains the labels `region` and `orchestrator`.
- `livepeer_orch_test_stream_download_time`: Download time per region for test streams. This GaugeVec contains the labels `region` and `orchestrator`.
- `livepeer_orch_test_stream_transcode_time`: Transcode time per region for test streams. This GaugeVec contains the labels `region` and `orchestrator`.
- `livepeer_orch_test_stream_round_trip_time`: Round trip time per region for test streams. This GaugeVec contains the labels `region` and `orchestrator`.

## Configuration

The exporter is configured via environment variables:

- `ORCHESTRATOR_ADDRESS`: Address of the orchestrator to fetch data from. **Required**
- `ORCHESTRATOR_ADDRESS_SECONDARY`: Address of the secondary orchestrator to fetch data from. This is used to calculate the `livepeer_orch_stake` metric. **Optional**
- `FETCH_INTERVAL`: How often to fetch data from the orchestrator. For example, if this is set to `5m`, the exporter will fetch data from the orchestrator every 5 minutes. See https://pkg.go.dev/time#ParseDuration for more information about the accepted format. **Optional** (default: `5m`)
- `FETCH_TEST_STREAMS_INTERVAL`: How often to fetch test streams data from the orchestrator. For example, if this is set to `5m`, the exporter will fetch test data from the orchestrator every 5 minutes. See https://pkg.go.dev/time#ParseDuration for more information about the accepted format. **Optional** (default: `5m`)
- `UPDATE_INTERVAL`: How often to update metrics. For example, if this is set to `5m`, the exporter will update metrics every 5 minutes. See https://pkg.go.dev/time#ParseDuration for more information about the accepted format. **Optional** (default: `30s`)

## Usage

### Run exporter

#### Run exporter locally

To run the exporter, set the necessary environment variables and start the exporter:

```bash
export ORCHESTRATOR_ADDRESS=your-orchestrator-address
export ORCHESTRATOR_ADDRESS_SECONDARY=your-secondary-orchestrator-address
export FETCH_INTERVAL=your-fetch-interval
export FETCH_TEST_STREAMS_INTERVAL=your-test-streams-fetch-interval
export UPDATE_INTERVAL=your-update-interval
go run main.go
```

The exporter will be available on port `9153`.

#### Run exporter with Docker

This repository also contains a [Dockerfile](./Dockerfile) and [docker-compose](./docker-compose.yml) file to run the exporter with Docker. To run the exporter with Docker, set the necessary environment variables in the docker-compose file and start the exporter:

```bash
docker-compose up
```

The exporter will be available on port `9153`.

### Configure Prometheus

For Prometheus to scrape the exporter, add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: livepeer
    static_configs:
      - targets: ['localhost:9153']
```
