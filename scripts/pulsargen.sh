#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

# Only proceed if --pulsar flag is provided
if [[ "$*" != *"--pulsar"* ]]; then
	exit 0
fi

if ! command -v protoc-gen-go-pulsar &>/dev/null; then
	echo "Error: protoc-gen-go-pulsar is not installed. Please install it before proceeding."
	echo "go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar"
	exit 1
fi

echo "Script executed from: ${PWD}"
project_root_dir=$(git rev-parse --show-toplevel)
echo "Generating proto pulsar code"

# Run buf generate command
buf generate -v --template "$project_root_dir/proto/buf.gen.pulsar.yaml"

echo "Pulsar files generated."
