
FROM relaysh/core:latest

WORKDIR /app

COPY collector .
COPY step.sh .
CMD ["sh", "-c", "./step.sh"]

LABEL "org.opencontainers.image.title"="cat-team-github-metrics"
LABEL "org.opencontainers.image.description"="A relay step to collect metrics from GitHub and publish them to BigQuery."
LABEL "sh.relay.sdk.version"="v1"
