#!/bin/bash

curl -s \
        --header "PRIVATE-TOKEN: $GITLAB_API_TOKEN" \
        $GITLAB_PROJECT_URLS/$CHECKOUT_PROJECT_ID/pipelines




