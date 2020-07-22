#!/bin/bash
pipeline_id=$(./get_pipelines.sh | grep failed | head -n1 | awk '{print $1}')

job_id=$(./get_jobs.sh "$pipeline_id" | grep failed | head -n1 | awk '{print $1}')

./get_logs.sh "$job_id"
