#!/usr/bin/env sh
set -e

# Clean up previous config
rm -rf ~/.zrchain

# Run the integration tests
go test -v ./tests/integration/mpc_ecdsa_keygen_test.go -tags=integration
