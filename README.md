[![Docker Build](https://github.com/rickstaa/livepeer-exporter/actions/workflows/docker-build.yml/badge.svg)](https://github.com/rickstaa/livepeer-exporter/actions/workflows/docker-build.yml)
[![Publish Docker image](https://github.com/rickstaa/livepeer-exporter/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/rickstaa/livepeer-exporter/actions/workflows/docker-publish.yml)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/rickstaa/livepeer-exporter?logo=docker)
](https://hub.docker.com/r/rickstaa/livepeer-exporter)
[![Latest Release](https://img.shields.io/github/v/release/rickstaa/livepeer-exporter?label=latest%20release)](https://github.com/rickstaa/livepeer-exporter/releases)

# Livepeer Exporter

![image](https://github.com/rickstaa/livepeer-exporter/assets/17570430/d168d3fc-4e58-4424-9836-d04e425d2991)

Livepeer Exporter is a lightweight tool designed to enhance the monitoring capabilities of [Livepeer](https://livepeer.org/). As a Prometheus exporter, it fetches [various metrics](#metrics) from different Livepeer endpoints and exposes them via an HTTP server, ready for Prometheus to scrape. This tool is the perfect companion to the [Livepeer monitoring service](https://docs.livepeer.org/orchestrators/guides/monitor-metrics), extending the range of Livepeer metrics that can be monitored. By providing deeper insights into Livepeer's performance, Livepeer Exporter helps users optimize their streaming workflows and ensure reliable service delivery.

## Configuration

Before using the Livepeer Exporter, you must configure it using environment variables. These variables allow you to customize the behaviour of the exporter to suit your specific needs. Below, you'll find a list of all the environment variables you can set, a description of what they do, and their default values if they are not specified.

### Required environment variables

- `LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS`: The address of the orchestrator to fetch data from.

### Optional environment variables

- `LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS_SECONDARY`: The address of the secondary orchestrator to fetch data from. Used to calculate the 'livepeer_orch_stake' metric. When set, the LPT stake of this address is added to the LPT stake that the orchestrator bonds.
- `LIVEPEER_EXPORTER_INFO_FETCH_INTERVAL`: How often to fetch general orchestrator information. Defaults to `1m`.
- `LIVEPEER_EXPORTER_SCORE_FETCH_INTERVAL`: How often to fetch score data for the orchestrator. Defaults to `1m`.
- `LIVEPEER_EXPORTER_DELEGATORS_FETCH_INTERVAL`: How often to fetch delegators data for the orchestrator. Defaults to `5m`.
- `LIVEPEER_EXPORTER_TEST_STREAMS_FETCH_INTERVAL`:How often to fetch the test streams data for the orchestrator. Defaults to `1h`.
- `LIVEPEER_EXPORTER_TICKETS_FETCH_INTERVAL`: How often to fetch ticket data for the orchestrator. Defaults to `1h`.
- `LIVEPEER_EXPORTER_REWARDS_FETCH_INTERVAL`: How often to fetch rewards data for the orchestrator. Defaults to `12h`.
- `CRYPTO_PRICES_EXPORTER_FETCH_INTERVAL`: How often to fetch the crypto prices. Defaults to `1m`.
- `LIVEPEER_EXPORTER_INFO_UPDATE_INTERVAL`: How often to update the orchestrator info metrics.  Defaults to `30s`.
- `LIVEPEER_EXPORTER_SCORE_UPDATE_INTERVAL`: How often to update the orchestrator score metrics. Defaults to `30s`.
- `LIVEPEER_EXPORTER_DELEGATORS_UPDATE_INTERVAL`: How often to update the orchestrator delegators metrics.  Defaults to `30s`.
- `LIVEPEER_EXPORTER_TEST_STREAMS_UPDATE_INTERVAL`: How often to update the orchestrator test streams metrics.  Defaults to `30s`.
- `LIVEPEER_EXPORTER_TICKETS_UPDATE_INTERVAL`: How often to update the orchestrator tickets metrics. Defaults to `30s`.
- `LIVEPEER_EXPORTER_REWARDS_UPDATE_INTERVAL`: How often to update the orchestrator rewards metrics.  Defaults to `30s`.
- `CRYPTO_PRICES_EXPORTER_UPDATE_INTERVAL`: How often to update the crypto prices metrics. Defaults to `30s`.

All intervals are specified as a string representation of a duration, e.g., "5m" for 5 minutes, "2h" for 2 hours, etc. See [time#ParseDuration](https://pkg.go.dev/time#ParseDuration) for format details.

> [!NOTE]\
> The `LIVEPEER_EXPORTER_TEST_STREAMS_FETCH_INTERVAL`, `LIVEPEER_EXPORTER_TICKETS_FETCH_INTERVAL`, and `LIVEPEER_EXPORTER_REWARDS_FETCH_INTERVAL`environment variables have higher default values. This adjustment is made considering the nature of the endpoints they fetch—they might be slow or return a substantial amount of data. By setting these values to a higher interval, you can effectively reduce the load on the exporter, ensuring optimal performance.

## Usage

This section explains how to run the Livepeer Exporter. You can run it locally on your own machine, or use Docker for easy setup and teardown.

### Run exporter locally

Running the exporter locally allows you to quickly start and stop it as needed. To do this, you'll need to set the necessary environment variables and then start the exporter. Replace `your-orchestrator-address`, with your own values:

```bash
export LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS=your-orchestrator-address
go run main.go
```

The exporter will be available on port `9153`. Additional [configuration](#configuration) environment variables can be passed to the exporter by adding them to the command above.

### Running the Exporter with Docker

You can run the exporter using the Docker image available on [Docker Hub](https://hub.docker.com/r/rickstaa/livepeer-exporter). To pull and run the exporter from Docker Hub, use the following command:

```bash
docker run --name livepeer-exporter \
    -e "LIVEPEER_EXPORTER_ORCHESTRATOR_ADDRESS=<your-orchestrator-address>" \
    -p 9153:9153 \
    rickstaa/livepeer-exporter:latest
```

Replace `<your-orchestrator-address>` with the address of your orchestrator. This command will start the exporter and expose the metrics on port `9153` for Prometheus to scrape. This command will start the exporter and expose the metrics on port `9153` for Prometheus to scrape. Additional environment variables can be passed to the exporter by adding them to the command above.

> [!NOTE]\
>  This repository also contains a [DockerFile](./Dockerfile) and [docker-compose.yml](./docker-compose.yml) file. These files can be used to build and run the exporter locally. To do this, clone this repository and run `docker compose up` in the repository's root directory.

### Configure Prometheus

For Prometheus to scrape the exporter, add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: livepeer-exporter
    static_configs:
      - targets: ["localhost:9153"]
```

This configuration tells Prometheus to scrape metrics from the Livepeer Exporter running on localhost port `9153`.

## Metrics

This exporter comprises the following sub-exporters, each responsible for fetching specific metrics:

| Sub-Exporter                                                          | Description                                                                                            |
| --------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------ |
| [orch_info_exporter](./exporters/orch_info_exporter/)                 | Collects metrics pertaining to the Livepeer orchestrator.                                              |
| [orch_score_exporter](./exporters/orch_score_exporter/)               | Retrieves metrics concerning the Livepeer orchestrator's score.                                        |
| [orch_delegators_exporter](./exporters/orch_delegators_exporter/)     | Gathers metrics related to the delegators of the designated Livepeer orchestrator.                     |
| [orch_test_streams_exporter](./exporters/orch_test_streams_exporter/) | Procures metrics about the Livepeer orchestrator's test streams.                                       |
| [orch_tickets_exporter](./exporters/orch_tickets_exporter/)           | Fetches metrics about the Livepeer orchestrator's tickets.                                             |
| [orch_reward_exporter](./exporters/orch_reward_exporter/)             | Retrieves metrics about the Livepeer orchestrator's rewards.                                           |
| [crypto_prices_exporter](./exporters/crypto_prices_exporter/)         | Fetches and exposes the prices of different cryptocurrencies used in the Livepeer ecosystem.           |

For enhanced performance, these sub-exporters operate concurrently in separate [goroutines](https://go.dev/tour/concurrency/1). They fetch metrics from various Livepeer endpoints and expose them via the `9153/metrics` endpoint. For detailed information about these sub-exporters and the metrics they provide, refer to the sections below.

### Crypto Prices Exporter

The `crypto_prices_exporter` fetches and exposes the prices of different cryptocurrencies used in the Livepeer ecosystem. They include:

**GaugeVec metrics:**

- `LPT_price`: This metric represents the current price of the LPT token. This GaugeVec includes the label `currency`, representing the currency of the price (e.g., `USD`, `EUR`, etc.)
- `ETH_price`: This metric represents the current price of Ethereum. This GaugeVec includes the label `currency`, representing the currency of the price (e.g., `USD`, `EUR`, etc.)

### orch_delegators_exporter

The `orch_delegators_exporter` fetches metrics about the delegators of the set Livepeer orchestrator from the `https://stronk.rocks/api/livepeer/getOrchestrator/` endpoint. These metrics provide insights into the number and behaviour of the delegators that stake with the orchestrator, including the total number of delegators, the bonded amount of each delegator, and the start round. They include:

**Gauge metrics:**

- `livepeer_orch_delegator_count`: This metric represents the total number of delegators that stake with the Livepeer orchestrator. It can be used to monitor the popularity and trustworthiness of the orchestrator.

**GaugeVec metrics:**

- `livepeer_orch_delegator_bonded_amount`: This metric represents the bonded amount of each delegator address. It can be used to understand the distribution of stakes among delegators. This GaugeVec includes the label `id`, representing the delegator's address.
- `livepeer_orch_delegator_start_round`: This metric represents the start round for each delegator. It can be used to track the longevity and loyalty of delegators. This GaugeVec includes the label `id`, representing the delegator's address.

### orch_info_exporter

The `orch_info_exporter` fetches metrics about the Livepeer orchestrator from the [Livepeer Orchestrator API](https://explorer.livepeer.org/_next/data/xe8lg6V7gubXcRErA1lxB/accounts/%s/orchestrating.json) and [Livepeer Delegating API](https://explorer.livepeer.org/_next/data/xe8lg6V7gubXcRErA1lxB/accounts/%s/delegating.json) endpoints. These metrics provide insights into the orchestrator's performance and behaviour. They include:

**Gauge metrics:**

- `livepeer_orch_bonded_amount`: This metric represents the amount of LPT bonded to the orchestrator.
- `livepeer_orch_total_stake`: This metric represents the total stake of the orchestrator in LPT.
- `livepeer_orch_last_reward_claim_round`: This metric represents the last round in which the orchestrator claimed the reward.
- `livepeer_orch_start_round`: This metric represents the round the orchestrator registered.
- `livepeer_orch_withdrawn_fees`: This metric represents the fees the orchestrator has withdrawn.
- `livepeer_orch_current_round`: This metric represents the current round.
- `livepeer_orch_activation_round`: This metric represents the round the orchestrator activated.
- `livepeer_orch_active`: This metric represents whether the orchestrator is active.
- `livepeer_orch_fee_cut`: This metric represents the proportion of the fees the orchestrator takes.
- `livepeer_orch_reward_cut`: This metric represents the proportion of the block reward the orchestrator takes.
- `livepeer_orch_last_reward_round`: This metric represents the last round in which the orchestrator received rewards while active.
- `livepeer_orch_ninety_day_volume_eth`: This metric represents the 90-day volume of ETH.
- `livepeer_orch_thirty_day_volume_eth`: This metric represents the 30-day volume of ETH.
- `livepeer_orch_total_volume_eth`: This metric represents the total volume of ETH.
- `livepeer_orch_stake`: This metric represents the stake provided by the orchestrator. Includes the orchestrator's bonded stake and the stake of the secondary orchestrator account (if provided).
- `livepeer_orch_thirty_day_reward_claim_ratio`: This metric represents how often an orchestrator claimed rewards in the last thirty rounds, or, if not active for 30 days, the reward claim ratio since activation.

### orch_rewards_exporter

The `orch_rewards_exporter` fetches reward data for the Livepeer orchestrator from the `https://stronk.rocks/api/livepeer/getAllRewardEvents` endpoint and filters it based on the orchestrator ID. These metrics provide insights into the rewards the orchestrator claims, including the total claimed rewards and details about each rewarded transaction. They include:

**Gauge metrics:**

- `livepeer_orch_total_claimed_rewards`: This metric represents the total rewards claimed by the orchestrator. It can be used to track the earnings of the orchestrator from rewards.

**GaugeVec metrics:**

- `livepeer_orch_reward_amount`: This metric represents the rewards earned by each transaction. It can be used to understand the distribution of rewards among transactions. This GaugeVec includes the label `id`, which represents the unique identifier of the transaction in which the ticket was won.
- `livepeer_orch_reward_transaction_hash`: This metric represents the hash of each rewarded transaction. It can track the transactions in which the orchestrator claimed rewards. This GaugeVec includes the label `id`, which represents the unique identifier of the transaction in which the ticket was won.
- `livepeer_orch_reward_block_number`: This metric represents the block number for each rewarded transaction. It can be used to track when the orchestrator claimed rewards. This GaugeVec includes the label `id`, which represents the unique identifier of the transaction in which the ticket was won.
- `livepeer_orch_reward_block_time`: This metric represents the block time for each rewarded transaction. It can be used to understand when the orchestrator claimed rewards. This GaugeVec includes the label `id`, which represents the unique identifier of the transaction in which the ticket was won.

### orch_score_exporter

The `orch_score_exporter` fetches metrics about the Livepeer orchestrator's score from the [Livepeer Score API](https://explorer.livepeer.org/api/score/) endpoint. These metrics provide insights into the performance of the orchestrator, including the price per pixel and success rate, round trip score, and total score per region. They include:

**Gauge metrics:**

- `livepeer_orch_price_per_pixel`: This metric represents the price per pixel. It can be used to understand the cost of using the orchestrator's services.

**GaugeVec metrics:**

- `livepeer_orch_success_rate`: This metric represents the success rate per region. It can monitor the reliability of the orchestrator's services in different regions. This GaugeVec includes the label `region`.
- `livepeer_orch_round_trip_score`: This metric represents the round trip score per region. It can measure the latency of the orchestrator's services in different regions. This GaugeVec includes the label `region`.
- `livepeer_orch_total_score`: This metric represents the total score per region. It can be used to evaluate the overall performance of the orchestrator in different regions. This GaugeVec includes the label `region`.

### orch_test_streams_exporter

The `orch_test_streams_exporter` fetches metrics about the Livepeer orchestrator's test streams from the `https://leaderboard-serverless.vercel.app/api/raw_stats` API endpoint. These metrics provide insights into the performance of the orchestrator's test streams in different regions, including success rate, upload time, download time, transcode time, and round trip time. They include:

**GaugeVec metrics:**

- `livepeer_orch_test_stream_success_rate`: This metric represents the success rate per region for test streams. It can monitor the reliability of the orchestrator's test streams in different regions. This GaugeVec includes the labels `region` and `orchestrator`.
- `livepeer_orch_test_stream_upload_time`: This metric represents the two-segment upload time per region for test streams. It can measure the upload speed of the orchestrator's test streams in different regions. This GaugeVec includes the labels `region` and `orchestrator`.
- `livepeer_orch_test_stream_download_time`: This metric represents the two-segment download time per region for test streams. It can measure the download speed of the orchestrator's test streams in different regions. This GaugeVec includes the labels `region` and `orchestrator`.
- `livepeer_orch_test_stream_transcode_time`: This metric represents the two-segment transcode time per region for test streams. It can measure the transcoding speed of the orchestrator's test streams in different regions. This GaugeVec includes the labels `region` and `orchestrator`.
- `livepeer_orch_test_stream_round_trip_time`: This metric represents the two-segment round trip time per region for test streams. It can measure the overall latency of the orchestrator's test streams in different regions. This GaugeVec includes the labels `region` and `orchestrator`.

### orch_tickets_exporter

The `orch_tickets_exporter` fetches and exposes ticket transaction information for each orchestrator from the `https://stronk.rocks/api/livepeer/getAllRedeemTicketEvents` API endpoint. These metrics provide insights into the winning tickets of the orchestrator, including the amount won, transaction hash, block number, and block time. They include:

**GaugeVec metrics:**

- `livepeer_orch_winning_ticket_amount`: This metric represents the fees won by each winning orchestrator ticket. It can be used to track the earnings of the orchestrator from winning tickets. This GaugeVec includes the label `id`, which represents the unique identifier of each ticket.
- `livepeer_orch_winning_ticket_transaction_hash`: This metric represents the transaction hash for each winning ticket. It can track the transactions in which the orchestrator won tickets. This GaugeVec includes the label `id`, which represents the unique identifier of the transaction in which the ticket was won.
- `livepeer_orch_winning_ticket_block_number`: This metric represents the block number for each winning ticket. It can be used to track when the orchestrator won tickets. This GaugeVec includes the label `id`, which represents the unique identifier of the transaction in which the ticket was won.
- `livepeer_orch_winning_ticket_block_time`: This metric represents the block time for each winning ticket. It can be used to understand when the orchestrator won tickets. This GaugeVec includes the label `id`, which represents the unique identifier of the transaction in which the ticket was won.

## Contributing

Feel free to open an issue if you have ideas on how to make this repository better or if you want to report a bug! All contributions are welcome. :rocket: Please consult the [contribution guidelines](CONTRIBUTING.md) for more information.
