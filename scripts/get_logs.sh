#!/bin/bash

job_id=$1
curl -s \
        --header "PRIVATE-TOKEN: $GITLAB_API_TOKEN" \
        $GITLAB_PROJECT_URLS/$CHECKOUT_PROJECT_ID/jobs/$job_id/trace
