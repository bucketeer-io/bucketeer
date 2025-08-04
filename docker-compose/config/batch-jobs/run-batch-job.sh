#!/bin/sh
#
# Bucketeer Batch Job Execution Script
# This script makes HTTP calls to the Bucketeer batch service to execute jobs
#

# Check if job ID is provided
if [ -z "$1" ]; then
    echo "Error: Job ID is required"
    echo "Usage: $0 <JOB_ID>"
    exit 1
fi

JOB_ID="$1"
ENDPOINT="${BATCH_SERVICE_ADDRESS}/bucketeer.batch.BatchService/ExecuteBatchJob"
TOKEN_FILE="/usr/local/service-token/token"
CERT_FILE="/usr/local/certs/service/tls.crt"

# Check if required files exist
if [ ! -f "$TOKEN_FILE" ]; then
    echo "Error: Service token file not found at $TOKEN_FILE"
    exit 1
fi

if [ ! -f "$CERT_FILE" ]; then
    echo "Error: Certificate file not found at $CERT_FILE"
    exit 1
fi

# Read the service token
TOKEN=$(cat "$TOKEN_FILE")

if [ -z "$TOKEN" ]; then
    echo "Error: Service token is empty"
    exit 1
fi

echo "$(date): Starting batch job: $JOB_ID"

# Make the HTTP request to execute the batch job
RES=$(curl -X POST \
    -m 3600 \
    --cacert "$CERT_FILE" \
    -d '{"job":"'$JOB_ID'"}' \
    -H "authorization: bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -s -o /dev/null -w '%{http_code}\n' \
    "$ENDPOINT")

echo "$(date): Batch job $JOB_ID result: HTTP $RES"

# Check response status
case "$RES" in
    200)
        echo "$(date): Batch job $JOB_ID completed successfully"
        exit 0
        ;;
    503)
        echo "$(date): Batch job $JOB_ID service unavailable (expected for some jobs)"
        exit 0
        ;;
    000)
        echo "$(date): Batch job $JOB_ID connection issue (network error)"
        exit 0
        ;;
    *)
        echo "$(date): Batch job $JOB_ID failed with HTTP status: $RES"
        exit 1
        ;;
esac 