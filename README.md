[![Docker Build](https://github.com/rickstaa/livepeer-exporter/actions/workflows/docker-build.yml/badge.svg)](https://github.com/rickstaa/livepeer-exporter/actions/workflows/docker-build.yml)
[![Publish Docker image](https://github.com/rickstaa/livepeer-exporter/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/rickstaa/livepeer-exporter/actions/workflows/docker-publish.yml)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/rickstaa/livepeer-exporter?logo=docker)
](https://hub.docker.com/r/rickstaa/livepeer-exporter)
![Latest Release](https://img.shields.io/github/v/release/rickstaa/livepeer-exporter?label=latest%20release)


# Livepeer Exporter

Livepeer Exporter is a lightweight tool designed to enhance the monitoring capabilities of [Livepeer](https://livepeer.org/). As a Prometheus exporter, it fetches various metrics from different Livepeer endpoints and exposes them via an HTTP server, ready for Prometheus to scrape. This tool is the perfect companion to the [Livepeer monitoring service](https://docs.livepeer.org/orchestrators/guides/monitor-metrics), extending the range of Livepeer metrics that can be monitored. By providing deeper insights into Livepeer's performance, Livepeer Exporter helps users optimize their streaming workflows and ensure reliable service delivery.

## Metrics

This exporter comprises the following sub-exporters, each responsible for fetching specific metrics:

- [orch_delegators_exporter](./exporters/orch_delegators_exporter/): Gathers metrics related to the delegators of the designated Livepeer orchestrator.
- [orch_info_exporter](./exporters/orch_info_exporter/): Collects metrics pertaining to the Livepeer orchestrator.
- [orch_score_exporter](./exporters/orch_score_exporter/): Retrieves metrics concerning the Livepeer orchestrator's score.
- [orch_test_streams_exporter](./exporters/orch_test_streams_exporter/): Procures metrics about the Livepeer orchestrator's test streams.
- [orch_tickets_exporter](./exporters/orch_tickets_exporter/): Fetches metrics about the Livepeer orchestrator's tickets.

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
- `livepeer_orch_stake`: Stake provided by the orchestrator.
- `livepeer_orch_reward_call_ratio`: Ratio of reward calls to total active rounds.
- `livepeer_orch_total_reward`: Total reward of the orchestrator.

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

### orch_tickets_exporter

Fetches and exposes ticket transaction information for each orchestrator from the https://stronk.rocks/api/livepeer/getAllRedeemTicketEvents API endpoint. The exposed metrics include:

**GaugeVec metrics:**

- `livepeer_orch_winning_ticket_amount`: Fees won by each winning orchestrator ticket. It contains the label `id`, which represents the unique identifier of each ticket.
- `livepeer_orch_winning_ticket_transaction_hash`: Transaction hash for each winning ticket. The `id` label is a unique identifier of the transaction in which the ticket was won.
- `livepeer_orch_winning_ticket_block_number`: Block number for each winning ticket. The `id` label is a unique identifier of the transaction in which the ticket was won.
- `livepeer_orch_winning_ticket_block_time`: Block time for each winning ticket. The `id` label is a unique identifier of the transaction in which the ticket was won.
  
## Configuration

The exporter is configured using the following environment variables:

- **LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS (Required):** Address of the primary orchestrator to fetch data for.
- **LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS_SECONDARY (Optional):** Address of the secondary orchestrator to fetch data for. This is used to calculate the `livepeer_orch_stake` metric.
- **LIVEPEER_EXPORTER_FETCH_INTERVAL (Optional, default: 5m):** How often to fetch general orchestrator data. For example, if set to `5m`, the exporter fetches data every 5 minutes. See [time#ParseDuration](https://pkg.go.dev/time#ParseDuration) for format details.
- **LIVEPEER_EXPORTER_FETCH_TEST_STREAMS_INTERVAL (Optional, default: 15m):** How often to fetch test streams data for the orchestrator. For example, if set to `5m`, the exporter fetches test data every 5 minutes. See [time#ParseDuration](https://pkg.go.dev/time#ParseDuration) for format details.
- **LIVEPEER_EXPORTER_TICKETS_FETCH_INTERVAL (Optional, default: 5m):** How often to fetch ticket data for the orchestrator. For example, if set to `5m`, the exporter fetches ticket data every 5 minutes. See [time#ParseDuration](https://pkg.go.dev/time#ParseDuration) for format details.
- **LIVEPEER_EXPORTER_UPDATE_INTERVAL (Optional, default: 30s):** How often to update Prometheus metrics. For example, if set to `5m`, the exporter updates metrics every 5 minutes. See [time#ParseDuration](https://pkg.go.dev/time#ParseDuration) for format details.

## Usage

### Run exporter

#### Run exporter locally

To run the exporter, set the necessary environment variables and start the exporter:

```bash
export LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS=your-orchestrator-address
export LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS_SECONDARY=your-secondary-orchestrator-address
export LIVEPEER_EXPORTER_FETCH_INTERVAL=your-fetch-interval
export LIVEPEER_EXPORTER_FETCH_TEST_STREAMS_INTERVAL=your-test-streams-fetch-interval
export LIVEPEER_EXPORTER_TICKETS_FETCH_INTERVAL=your-tickets-fetch-interval
export LIVEPEER_EXPORTER_UPDATE_INTERVAL=your-update-interval
go run main.go
```

The exporter will be available on port `9153`.

#### Running the Exporter with Docker

You can run the exporter using the Docker image available on [Docker Hub](https://hub.docker.com/r/rickstaa/livepeer-exporter).To pull and run the exporter from Docker Hub, use the following command:

```bash
docker run --name livepeer-exporter \
    -e "LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS=<your-orchestrator-address>" \
    -p 9153:9153 \
    rickstaa/livepeer-exporter:latest
```

This command will start the exporter and expose the metrics on port `9153` for Prometheus to scrape. Additional environment variables can be passed to the exporter by adding them to the command above.

> [!NOTE]
> This repository provides a [Dockerfile](./Dockerfile) and a [docker-compose](./docker-compose.yml) file to facilitate running the exporter with Docker. To utilize these, first configure the necessary environment variables within the docker-compose file. Subsequently, initiate the exporter using the command `docker-compose up`.

### Configure Prometheus

For Prometheus to scrape the exporter, add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: livepeer-exporter
    static_configs:
      - targets: ["localhost:9153"]
```

## Contributing

Feel free to open an issue if you have ideas on how to make this repository better or if you want to report a bug! All contributions are welcome. :rocket: Please consult the [contribution guidelines](CONTRIBUTING.md) for more information.
