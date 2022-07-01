# cat-team-github-metrics

This project is a Relay step that is responsible for collecting various metrics from GitHub and injecting them in to BigQuery.

## Working locally

### Exporter

To use the exporter you'll need to be logged in to GCP with the gcloud cli and have the following environment variables exported in your current session.

```bash
export MODULE_OWNER=""
export MODULE_NAME=""
export GITHUB_TOKEN=xxxxxxxxxx
export BIG_QUERY_DATASET_ID="my_dataset"
```

#### Running the exporter in docker

You can run the exporter in a docker container.
The BigQuery SDK will use the local credential file that is presented to the container.

```bash
make release

docker run \
  -e MODULE_OWNER=$MODULE_OWNER \
  -e MODULE_NAME=$MODULE_NAME \
  -e GITHUB_TOKEN=$GITHUB_TOKEN \
  -e BIG_QUERY_PROJECT_ID=$BIG_QUERY_PROJECT_ID \
  -e GOOGLE_APPLICATION_CREDENTIALS=/tmp/keys/creds.json \
  -v ~/.config/gcloud/application_default_credentials.json:/tmp/keys/creds.json \
  ghcr.io/chelnak/cat-team-github-metrics
```

### Grafana

This repository also comtains a local Grafana stack.
It will be deployed with the BigQuery datasource which will need to be configured.

First create a new volume where we will store data:

```bash
docker volume create --name=grafana-data
```

Then move to the `grafana` directory and start the stack as follows:

```bash
cd grafana
docker-compose up
```
Your local stack will be accessible at <http://localhost:9000>

## Build and release

Builds and releases are handled by goreleaser.
For convenience when working locally use the provided Makefile.

### Release steps
* Ensure that you are on the HEAD of the main branch.
* Create a tag `make tag version=v.0.0.1`. This will also push the tag to the remote.
* The release workflow is triggered by new tags. It will publish a binary and a Docker image.
