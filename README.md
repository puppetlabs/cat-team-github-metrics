# cat-team-github-metrics

This project is responsible for collecting various metrics from GitHubs API and injecting them in to BigQuery.

Data can be presented in Grafana.

## Working locally

### Exporter

For the exporter you'll need to be logged in to GCP with the gcloud cli and have the following environments exported in your current session.

```bash
export GITHUB_TOKEN=xxxxxxxxxx
export BIG_QUERY_DATASET_ID="my_dataset"

#### Running the exporter locally

You can run locally with go, as follows.

```bash
go run .
```

#### Running the exporter in docker

You can run the exporter in a docker container. The BigQuery SDK will use the credential file that is presented to the container and you
will need to export the required environment variables as described above.

```bash
make build-all

docker run \
  -e GITHUB_TOKEN=$GITHUB_TOKEN \
  -e BIG_QUERY_PROJECT_ID=$BIG_QUERY_PROJECT_ID \
  -e GOOGLE_APPLICATION_CREDENTIALS=/tmp/keys/creds.json \
  -v ~/.config/gcloud/application_default_credentials.json:/tmp/keys/creds.json \
  chelnak/cat-team-github-metrics
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
