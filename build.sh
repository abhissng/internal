#!/bin/bash

# Exit on error
set -e

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null
then
    echo "golangci-lint could not be found, installing..."
    sudo go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.0
    # curl -sSfL https://github.com/golangci/golangci-lint/releases/latest/download/golangci-lint-$(go env GOOS)-$(go env GOARCH).tar.gz | tar -xzv -C /tmp
    # sudo mv /tmp/golangci-lint-*/golangci-lint /usr/local/bin
fi

# Check if gosec is installed
if ! command -v gosec &> /dev/null
then
    echo "gosec could not be found, installing..."
    go install github.com/securego/gosec/v2/cmd/gosec@latest
fi

# Run golangci-lint
echo "Running golangci-lint..."
golangci-lint run ./...

# Run gosec for security analysis
echo "Running gosec (security analysis)..."
gosec ./...
