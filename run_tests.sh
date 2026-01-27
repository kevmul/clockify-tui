#!/bin/bash

echo "Running Clockify App Tests..."
echo "=============================="

# Run all tests with coverage
echo "Running tests with coverage..."
go test ./... -v -cover

# Generate coverage report
echo -e "\nGenerating coverage report..."
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

echo -e "\nTest coverage report generated: coverage.html"
echo "Open coverage.html in your browser to view detailed coverage"
