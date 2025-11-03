#!/usr/bin/env -S uv run
# /// script
# requires-python = ">=3.8"
# dependencies = [
#   "boto3",
# ]
# ///

import boto3

if __name__ == "__main__":
  endpoint_url = "http://localstack:4566"

  kms_client = boto3.client('kms',
                            region_name='eu-west-1',
                            endpoint_url=endpoint_url
                            )

  keys = ["party1", "party2", "party3"]

  for key in keys:
    response = kms_client.create_key(
      Description=f"{key}"
    )
    KeyId = response['KeyMetadata']['KeyId']
    alias = f"alias/{key}"
    response = kms_client.create_alias(
      AliasName=alias,
      TargetKeyId=KeyId
    )
    print(f"{key} key and alias created")
