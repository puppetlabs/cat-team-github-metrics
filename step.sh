#!/bin/sh

echo $(ni get -p {.connection.serviceAccountKey}) > /app/credentials.json

export GOOGLE_APPLICATION_CREDENTIALS="credentials.json"
export GITHUB_TOKEN=$(ni get -p {.github_token})
export BIG_QUERY_PROJECT_ID=$(ni get -p {.big_query_project_id})
export REPO_OWNER=$(ni get -p {.repo_owner})
export REPO_NAME=$(ni get -p {.repo_name})

./collector
