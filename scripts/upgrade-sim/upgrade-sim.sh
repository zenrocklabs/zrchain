#!/bin/bash

echo "starting docker init"

echo "checking if docker is down"

docker-compose down

echo "remove volumes"

docker volume rm upgrade-sim_zenrock-mainnet

echo "remove images"

docker rmi upgrade-sim-init-mainnet

docker rmi zenrockd

echo "build new zenrockd image"

cd ../../

echo "pwd: $(pwd)"

docker login

rm zenrockd

echo "build new zenrockd image"

docker build -t zenrockd -f Dockerfile . --progress=plain --no-cache

echo "copy zenrockd to current directory"

docker run --entrypoint=/bin/cp --rm -v $(pwd):/workspace zenrockd "/zenrockd" "/workspace/zenrockd"

echo "copy done"

cd scripts/upgrade-sim

echo "starting init-mainnet, wait 60 seconds"

docker-compose up -d init-mainnet

sleep 60

echo "starting zenrock-mainnet, wait 300 seconds"

docker-compose up -d zenrock-mainnet

sleep 200

echo "stopping zenrock-mainnet"
    
docker stop zenrock-mainnet

echo "comment out the start command and uncomment the tail command"

# Comment out the start command and uncomment the tail command
# Simple text replacement approach
# Create a backup
cp docker-compose.yaml docker-compose.yaml.backup

# Use a simple sed command to replace the command section
# First, let's use a different approach - create a temporary file with the new content
cat > docker-compose.yaml.tmp << 'EOF'
services:
  init-mainnet:
    build:
      context: ./
      dockerfile: Dockerinit
    environment:
      COSMOVISOR_VERSION: 1.6.0
      # update this version to the one currently running on mainnet
      ZENROCKD_GENESIS_VERSION: 6.25.0
      MONIKER: zenrock-mainnet
    volumes:
      - zenrock-mainnet:/root-data

  zenrock-mainnet:
    image: alpine:3.18.0
    #platform: linux/x86_64
    container_name: zenrock-mainnet
    environment:
      VALIDATOR_NAME: zenrock-mainnet
      COSMOVISOR_COLOR_LOGS: true
      COSMOVISOR_DISABLE_LOGS: false
      COSMOVISOR_DISABLE_RECASE: false
      COSMOVISOR_TIMEFORMAT_LOGS: kitchen
      DAEMON_ALLOW_DOWNLOAD_BINARIES: true
      DAEMON_DOWNLOAD_MUST_HAVE_CHECKSUM: false
      DAEMON_HOME: /root-data
      DAEMON_NAME: zenrockd
      DAEMON_POLL_INTERVAL: 5s
      DAEMON_PREUPGRADE_MAX_RETRIES: 0
      DAEMON_RESTART_AFTER_UPGRADE: true
      DAEMON_RESTART_DELAY: 5s
      DAEMON_SHUTDOWN_GRACE: 5s
      UNSAFE_SKIP_BACKUP: true
    volumes:
      - zenrock-mainnet:/root-data
    command:
      [
        "tail",
        "-f",
        "/dev/null"
      ]
    ports:
      - 26656:26656 #peer
      - 26657:26657 #rpc
      - 9790:9790 #grpc
      - 1717:1717 #api
      - 26660:26660 #metrics
    restart: unless-stopped

volumes:
  zenrock-mainnet:
EOF

# Replace the original file with the modified version
mv docker-compose.yaml.tmp docker-compose.yaml

echo "restarting zenrock-mainnet without cosmovisor"

docker-compose up -d zenrock-mainnet

echo "copy new zenrockd to zenrock-mainnet root"

docker cp ../../zenrockd zenrock-mainnet:/root-data/ 

docker exec zenrock-mainnet /root-data/zenrockd export --output-document export.json --log_format json --home /root-data 

docker cp zenrock-mainnet:/export.json .

cp docker-compose.yaml docker-compose.yaml.backup

# Use a simple sed command to replace the command section
# First, let's use a different approach - create a temporary file with the new content
cat > docker-compose.yaml.tmp << 'EOF'
services:
  init-mainnet:
    build:
      context: ./
      dockerfile: Dockerinit
    environment:
      COSMOVISOR_VERSION: 1.6.0
      # update this version to the one currently running on mainnet
      ZENROCKD_GENESIS_VERSION: 6.25.0
      MONIKER: zenrock-mainnet
    volumes:
      - zenrock-mainnet:/root-data

  zenrock-mainnet:
    image: alpine:3.18.0
    #platform: linux/x86_64
    container_name: zenrock-mainnet
    environment:
      VALIDATOR_NAME: zenrock-mainnet
      COSMOVISOR_COLOR_LOGS: true
      COSMOVISOR_DISABLE_LOGS: false
      COSMOVISOR_DISABLE_RECASE: false
      COSMOVISOR_TIMEFORMAT_LOGS: kitchen
      DAEMON_ALLOW_DOWNLOAD_BINARIES: true
      DAEMON_DOWNLOAD_MUST_HAVE_CHECKSUM: false
      DAEMON_HOME: /root-data
      DAEMON_NAME: zenrockd
      DAEMON_POLL_INTERVAL: 5s
      DAEMON_PREUPGRADE_MAX_RETRIES: 0
      DAEMON_RESTART_AFTER_UPGRADE: true
      DAEMON_RESTART_DELAY: 5s
      DAEMON_SHUTDOWN_GRACE: 5s
      UNSAFE_SKIP_BACKUP: true
    volumes:
      - zenrock-mainnet:/root-data
    command:
      [
        "/root-data/cosmovisor/bin/cosmovisor",
        "run",
        "start",
        "--home",
        "/root-data",
      ]
    ports:
      - 26656:26656 #peer
      - 26657:26657 #rpc
      - 9790:9790 #grpc
      - 1717:1717 #api
      - 26660:26660 #metrics
    restart: unless-stopped

volumes:
  zenrock-mainnet:
EOF

# Replace the original file with the modified version
mv docker-compose.yaml.tmp docker-compose.yaml