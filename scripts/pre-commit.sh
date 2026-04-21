#!/usr/bin/env bash
set -e

echo "====================================="
echo "Running pre-commit hooks..."
echo "====================================="

echo "1. Formatting code..."
go fmt ./...

echo "2. Running linters..."
if command -v golangci-lint >/dev/null 2>&1; then
    make lint
else
    echo "Warning: golangci-lint is not installed. Skipping linter."
    echo "To install: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin"
fi

echo "3. Running tests..."
make test

echo "====================================="
echo "Pre-commit checks passed! ✨"
echo "====================================="
