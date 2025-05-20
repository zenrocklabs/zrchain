#!/usr/bin/env sh
set -ex  # Add -x for verbose output

# Clean up previous config
echo "Cleaning up previous config..."
rm -rf ~/.zrchain

# Run the integration tests with verbose output
echo "Running integration tests..."
go test -v -count=1 -tags "integration,!wasm" ./tests/integration/keyring_test.go
go test -v -count=1 -tags "integration,!wasm" ./tests/integration/action_test.go
go test -v -count=1 -tags "integration,!wasm" ./tests/integration/mpc_ecdsa_keygen_test.go
