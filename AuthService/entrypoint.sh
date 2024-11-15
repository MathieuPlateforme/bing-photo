#!/bin/sh
set -e

# Vérifier si la variable RUN_TESTS est définie à "true"
if [ "$RUN_TESTS" = "true" ]; then
  if command -v go >/dev/null 2>&1; then
    echo "Running tests..."
    go test -v ./...
  else
    echo "Go environment not available, skipping tests."
  fi
else
  echo "RUN_TESTS is not set to 'true', skipping tests."
fi

echo "Starting the application..."
exec "$@"
