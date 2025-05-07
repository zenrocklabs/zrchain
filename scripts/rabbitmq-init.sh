#!/usr/bin/env sh

# Enable stream plugin
rabbitmq-plugins enable rabbitmq_stream rabbitmq_stream_management

# Start rabbitmq
/opt/rabbitmq/sbin/rabbitmq-server &
pid=$!

# wait for it to start
sleep 10

# Create stream
rabbitmqadmin declare queue name="${STREAM_NAME}" queue_type=stream

wait $pid
