#!/bin/bash

echo "Generating tests for all .go files in the project..."

# Find all .go files except those ending in _test.go and generate tests
find . -name '*.go' -not -name '*_test.go' -exec gotests -w -all {} +

echo "Test generation complete."
