# puppet-github-metrics

This project is responsible for collecting various metrics from GitHubs API and injecting them in to BigQuery.

Data can be presented in Grafana.

## Working locally

### Exporter

For the exporter you'll need to be logged in to GCP with the gcloud cli and have the following environments exported in your current session.

```bash
export GITHUB_TOKEN=xxxxxxxxxx
export BIG_QUERY_DATASET_ID="my_dataset"
```

### Local Grafana stack

First create a new volume where we will store data

```bash
docker volume create --name=grafana-data
```

Then move to the `grafana` directory and start the stack

```bash
cd grafana
docker-compose up
```

Your local stack will be accessible at <http://localhost:9000>
