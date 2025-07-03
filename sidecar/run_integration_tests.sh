#!/bin/bash

# Script to run integration tests for the sidecar
# This script helps developers run integration tests when external dependencies are available

set -e

echo "🚀 Sidecar Integration Tests Runner"
echo "=================================="

# Check if we're in CI
if [ "$CI" = "true" ]; then
    echo "❌ CI environment detected. Integration tests are skipped in CI."
    echo "   Run this script locally with external dependencies available."
    exit 0
fi

# Check if Bitcoin proxy is running
echo "🔍 Checking Bitcoin proxy availability..."
if nc -z 127.0.0.1 1234 2>/dev/null; then
    echo "✅ Bitcoin proxy is running on 127.0.0.1:1234"
else
    echo "❌ Bitcoin proxy not found on 127.0.0.1:1234"
    echo "   Please start the Bitcoin proxy service first."
    echo "   Example: docker run -p 1234:1234 your-bitcoin-proxy"
    exit 1
fi

# Check if Neutrino nodes are available
echo "🔍 Checking Neutrino node availability..."
if [ -d "./neutrino" ]; then
    echo "✅ Neutrino directory found"
else
    echo "⚠️  Neutrino directory not found. Creating..."
    mkdir -p ./neutrino
fi

# Set environment variable to enable integration tests
export RUN_INTEGRATION_TESTS=true

echo "🧪 Running integration tests..."
echo "   This may take some time as Neutrino nodes need to sync..."

# Run the integration tests
go test -v ./zrchain/sidecar -run "Test_ProxyFunctions_Testnet" -timeout 5m

echo ""
echo "✅ Integration tests completed!"
echo ""
echo "📝 Notes:"
echo "   - Tests will be skipped if external dependencies are not available"
echo "   - Set RUN_INTEGRATION_TESTS=true to enable integration tests"
echo "   - Set CI=true to disable integration tests (useful for CI environments)" 