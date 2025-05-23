#!/bin/bash

# Store process IDs
pids=()

# Function to handle Ctrl+C
cleanup() {
  echo "${pids[0]}"
  echo "Stopping all kvNode instances..."
  for pid in "${pids[@]}"; do
    kill "$pid" 2>/dev/null
  done
  wait
  echo "All processes stopped."
  exit 0
}

# Trap Ctrl+C
trap cleanup SIGINT

go run ./cmd/kvController --config "./config/controller_config.yaml"  &
pids+=($!)

sleep 2
# Start kvNode instances
for i in {1..4}; do
  go run ./cmd/kvNode --config "./config/node_config_${i}.yaml" &
  pids+=($!)
done

# Wait for all processes
wait