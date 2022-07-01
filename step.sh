#!/bin/sh

export GITHUB_TOKEN=$(ni get -p {.github_token})
export BIG_QUERY_PROJECT_ID=$(ni get -p {.big_query_project_id})
export MODULE_OWNER=$(ni get -p {.module_owner})
export MODULE_NAME=$(ni get -p {.module_namena})

./collector
