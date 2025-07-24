# Example usage:
#
# docker build -f Dockerfile -t zenrockd:0.0.1 .
# docker run -e DOCKER_ENV=true -p 8080:8080 zenrockd:0.0.1

# Use a  golang alpine as the base image
FROM public.ecr.aws/docker/library/golang:1.24.3-alpine3.21 AS go_builder
RUN apk update
RUN apk add make cmake git alpine-sdk linux-headers

# Setup
# Read arguments
ARG ARCH=x86_64
ARG BUILD_DATE
ARG GIT_SHA
ARG VERSION
# Set env variables
ENV arch=$ARCH
ENV build_date=$BUILD_DATE
ENV commit_hash=$GIT_SHA
ENV service_name=zenrockd
ENV version=$VERSION
RUN echo "building service: ${service_name}, version: ${version}, build date: ${build_date}, commit hash: ${commit_hash}, architecture: ${arch}"



# Set the working directory
COPY . /zrchain
WORKDIR /zrchain
ENV BUILD_TAGS=muslc LINK_STATICALLY=true

# Download dependencies
RUN go mod download

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    make build


############################################################################################################

#SSL certs
FROM alpine:3.18.0 AS certs
RUN apk add --no-cache ca-certificates jq
RUN adduser -Ds /bin/bash appuser

# Copy binary to a scratch container. Let's keep our images nice and small!
COPY --from=go_builder /zrchain/build/zenrockd /zenrockd

# Set user
USER appuser

# Expose the port your application will run on
EXPOSE 26656
EXPOSE 26657
EXPOSE 9090

# Run the binary
ENTRYPOINT [ "/zenrockd" ]
