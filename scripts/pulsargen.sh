#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

if ! command -v protoc-gen-go-pulsar &>/dev/null; then
	echo "Error: protoc-gen-go-pulsar is not installed. Please install it before proceeding."
	echo "go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar"
	exit 1
fi

echo "Script executed from: ${PWD}"
project_root_dir=$(git rev-parse --show-toplevel)
echo "Generating proto pulsar code"
proto_root=$project_root_dir/zrchain

pushd .
cd "$proto_root"
buf generate -v --template "$project_root_dir"/zrchain/proto/buf.gen.pulsar.yaml
popd

echo "Pulsar files generated."
