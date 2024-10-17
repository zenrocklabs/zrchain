#!/usr/bin/env bash

ignite generate hooks

dir=$(cd "$(dirname "$0")"; pwd -P)
zrchain_dir="$(dirname "$dir")"
root_dir="$(dirname "$zrchain_dir")"

# replace ts client with the newly generated one
rm -rf $root_dir/web/ts-client
cp -r $root_dir/zrchain/ts-client $root_dir/web
rm -rf $root_dir/zrchain/ts-client

# replace react hooks with the newly generated ones
rm -rf $root_dir/web/react/src/hooks
cp -r $root_dir/zrchain/react/src/hooks $root_dir/web/react/src
rm -rf $root_dir/zrchain/react
