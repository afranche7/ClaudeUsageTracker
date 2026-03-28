#!/usr/bin/env bash
# Build the Claude Usage Tracker (browser mode on Linux/macOS)
set -e

echo "Building Claude Usage Tracker..."

CGO_ENABLED=0 go build -ldflags="-s -w" -o claude-usage-tracker ./cmd/widget

echo "Build successful: ./claude-usage-tracker"
echo ""
echo "Run with:"
echo "  ./claude-usage-tracker              (auto browser mode on Linux/macOS)"
echo "  ./claude-usage-tracker --port 8080  (custom port)"
echo "  ./claude-usage-tracker --interval 15 (refresh every 15s)"
