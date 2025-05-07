#!/bin/bash

echo "Creating initial tables"
sql_commands=(
  "CREATE TABLE IF NOT EXISTS $POSTGRES_KEYS_TABLE1 (
     id VARCHAR(64) PRIMARY KEY,
     data VARCHAR(10240) NOT NULL
   )"
  "CREATE TABLE IF NOT EXISTS $POSTGRES_KEYS_TABLE2 (
     id VARCHAR(64) PRIMARY KEY,
     data VARCHAR(10240) NOT NULL
   )"
  "CREATE TABLE IF NOT EXISTS $POSTGRES_KEYS_TABLE3 (
     id VARCHAR(64) PRIMARY KEY,
     data VARCHAR(10240) NOT NULL
   )"
  "CREATE TABLE $POSTGRES_REQUESTS_TABLE (
    id INTEGER NOT NULL,
    type VARCHAR(10) NOT NULL CHECK (type IN ('key', 'sign')),
    key_type VARCHAR(10) NOT NULL CHECK (key_type IN ('ecdsa', 'eddsa')),
    status VARCHAR(20) NOT NULL CHECK (status IN ('pending', 'in_progress')),
    created_at TIMESTAMP NOT NULL,
    next_attempt_time TIMESTAMP NOT NULL,
    attempt_count INTEGER NOT NULL DEFAULT 0,
    zrchain_request BYTEA,

    PRIMARY KEY (id, type)
  )"
)

# Loop through and execute each SQL command
for sql_command in "${sql_commands[@]}"; do
  cmd=(
    psql
    -U "$POSTGRES_USER"
    -d "$POSTGRES_DB"
    -c "$sql_command"
  )
  echo "${cmd[@]}"
  "${cmd[@]}"
done
