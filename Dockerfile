from golang:1.18-alpine

WORKDIR /app

COPY ./dist/linux-amd64/collector .
COPY ./step.sh .
CMD ["sh", "-c", "./step.sh"]

LABEL "org.opencontainers.image.title"="cat-team-github-metrics"
LABEL "org.opencontainers.image.description"="A relay step to collect metrics from GitHub and publish them to BigQuery."
LABEL "sh.relay.sdk.version"="v1"
