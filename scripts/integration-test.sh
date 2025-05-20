#!/usr/bin/env sh
set -ex  # Add -x for verbose output

# Print environment information
echo "Go version:"
go version

echo "Go environment:"
go env

# Clean up previous config
echo "Cleaning up previous config..."
rm -rf ~/.zrchain

# Run the integration tests with verbose output
echo "Running integration tests..."
go test -v -count=1 ./tests/integration/mpc_ecdsa_keygen_test.go
