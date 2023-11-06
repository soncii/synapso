#!/bin/bash

# Keep this script running to ensure it restarts the Go application
while true; do
  ./app/main
  exit_status=$?
  if [ $exit_status -ne 0 ]; then
    echo "Error detected. Restarting..."
  else
    echo "Application exited gracefully. Exiting..."
    exit 0
  fi
  sleep 2
done

