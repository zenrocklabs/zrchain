#!/usr/bin/env bash

# Dependencies:
# - yq
# - docker

set -eo pipefail

if [ -f "/etc/alpine-release" ]; then
  apk add yq
fi

echo "Script executed from: ${PWD}"

echo "Generating proto code"
proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  proto_files=$(find "${dir}" -maxdepth 1 -name '*.proto')
  for file in $proto_files; do
    # Check if the go_package in the file is pointing to zenrock
    if grep -q "option go_package.*zenrock" "$file"; then
      echo "Processing file ${file}"
      buf generate --template proto/buf.gen.gogo.yaml "$file"
      # Only generate Python files if --python flag is provided
      if [[ "$*" == *"--python"* ]]; then
        buf generate --template proto/buf.gen.python.yaml "$file"
      fi
    fi
  done
done

# Always show proto files generated message
echo "Proto files generated."

# Show Python files message if they were generated
if [[ "$*" == *"--python"* ]]; then
  echo "Python files generated."
fi

# move proto files to the right places
if [ -d "github.com" ]; then
  # Find the correct directory structure
  GEN_DIR=$(find github.com -type d -name "zrchain" | head -n 1)
  if [ -n "$GEN_DIR" ]; then
    echo "Copying generated files from $GEN_DIR"
    cp -a "$GEN_DIR"/* ./
    rm -rf ./github.com
  else
    echo "Warning: Could not find generated files in expected location"
  fi
else
  echo "Warning: No generated files found in github.com directory"
fi

