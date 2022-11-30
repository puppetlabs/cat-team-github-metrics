#!/bin/sh

echo $CONNECTION_KEY > /app/credentials.json

export GOOGLE_APPLICATION_CREDENTIALS="/app/credentials.json"
export GITHUB_TOKEN="$GITHUB_TOKEN"
export BIG_QUERY_PROJECT_ID="$BQ_PROJECT_ID"
export REPO_OWNER="$REPO_OWNER"
export REPO_NAME="$REPO_NAME"

COMMAND="$COBRA_COMMAND"

case "$COMMAND" in
    export|stamp) ./collector $COMMAND;;
    *) echo "Invalid command! Should be one of [export, stamp]." && exit 1;;
esac

