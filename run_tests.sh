#!/bin/bash

echo "Running Clockify App Tests..."
echo "=============================="

# Run tests with coverage
go test -v -cover ./internal/api ./internal/utils ./internal/config ./internal/models

echo ""
echo "Test Summary:"
echo "============="
go test ./internal/api ./internal/utils ./internal/config ./internal/models | grep -E "(PASS|FAIL)"
