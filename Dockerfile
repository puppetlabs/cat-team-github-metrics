
FROM ubuntu:jammy

RUN apt-get update && apt-get install -y pcre2-utils ca-certificates && update-ca-certificates

WORKDIR /app

COPY linux-amd64/collector .
COPY step.sh .
CMD ["sh", "-c", "./step.sh"]

LABEL "org.opencontainers.image.title"="cat-team-github-metrics"
LABEL "org.opencontainers.image.description"="A step to collect metrics from GitHub and publish them to BigQuery."
