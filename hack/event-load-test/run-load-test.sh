#!/bin/bash

# Exit on error
set -e

# Default values
GATEWAY_ADDR="localhost:8443"
API_KEY_PATH=""
NUM_EVAL_EVENTS="10000"
NUM_GOAL_EVENTS="5000"
CONCURRENT_WORKERS="10"
ENVIRONMENT_ID=""

# Process command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --gateway-addr)
      GATEWAY_ADDR="$2"
      shift 2
      ;;
    --api-key-path)
      API_KEY_PATH="$2"
      shift 2
      ;;
    --num-eval-events)
      NUM_EVAL_EVENTS="$2"
      shift 2
      ;;
    --num-goal-events)
      NUM_GOAL_EVENTS="$2"
      shift 2
      ;;
    --concurrent-workers)
      CONCURRENT_WORKERS="$2"
      shift 2
      ;;
    --environment-id)
      ENVIRONMENT_ID="$2"
      shift 2
      ;;
    --help)
      echo "Usage: $0 [options]"
      echo ""
      echo "Options:"
      echo "  --gateway-addr ADDR        Gateway address (default: localhost:8443)"
      echo "  --api-key-path PATH        Path to API key (REQUIRED)"
      echo "  --num-eval-events NUM      Number of evaluation events (default: 10000)"
      echo "  --num-goal-events NUM      Number of goal events (default: 5000)"
      echo "  --concurrent-workers NUM   Number of concurrent workers (default: 10)"
      echo "  --environment-id ID        Environment ID"
      echo "  --help                     Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      echo "Use --help for usage information"
      exit 1
      ;;
  esac
done

# Check if go is installed
if ! command -v go &> /dev/null; then
  echo "Error: go is not installed or not in the PATH"
  exit 1
fi

# Check for required API key path
if [ -z "$API_KEY_PATH" ]; then
  echo "Error: --api-key-path is required"
  exit 1
fi

# Validate API key file exists
if [ ! -f "$API_KEY_PATH" ]; then
  echo "Error: API key file not found at $API_KEY_PATH"
  exit 1
fi

# Print run configuration
echo "Running event load test with the following configuration:"
echo "Gateway Address:     $GATEWAY_ADDR"
echo "API Key Path:        $API_KEY_PATH"
echo "Evaluation Events:   $NUM_EVAL_EVENTS"
echo "Goal Events:         $NUM_GOAL_EVENTS"
echo "Concurrent Workers:  $CONCURRENT_WORKERS"
if [ -n "$ENVIRONMENT_ID" ]; then
  echo "Environment ID:      $ENVIRONMENT_ID"
else
  echo "Environment ID:      (not specified)"
fi
echo ""

# Build environment ID parameter if provided
ENV_PARAM=""
if [ -n "$ENVIRONMENT_ID" ]; then
  ENV_PARAM="--environment-id=$ENVIRONMENT_ID"
fi

# Run the load test
echo "Starting load test..."
go run main.go \
  --gateway-addr="$GATEWAY_ADDR" \
  --api-key-path="$API_KEY_PATH" \
  --num-eval-events="$NUM_EVAL_EVENTS" \
  --num-goal-events="$NUM_GOAL_EVENTS" \
  --concurrent-workers="$CONCURRENT_WORKERS" \
  $ENV_PARAM

echo "Load test completed." 