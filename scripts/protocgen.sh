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
      buf generate --template proto/buf.gen.python.yaml "$file"
    fi
  done
done
echo "All files done"

# Generate dependencies for Python
# yq '.deps[]' proto/buf.yaml | while read -r dep; do \
#   echo "Generate python dependencies for ${dep}"
#   buf generate --template proto/buf.gen.python.yaml "$dep"; \
# done
# echo "All dependencies done"

# restore user permissions
chown -R 1000:1000 proto/python
chown -R 1000:1000 github.com/zenrocklabs

# move proto files to the right places
cp -ar github.com/Zenrock-Foundation/zrchain/v4/* ./
rm -rf ./github.com

