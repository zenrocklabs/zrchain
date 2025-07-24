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
ENV GOPRIVATE="https://github.com/zenrocklabs/zenrock-avs"
ENV arch=$ARCH
ENV build_date=$BUILD_DATE
ENV commit_hash=$GIT_SHA
ENV service_name=validator_sidecar
ENV version=$VERSION
RUN echo "building service: ${service_name}, version: ${version}, build date: ${build_date}, commit hash: ${commit_hash}, architecture: ${arch}"


# Set the working directory
COPY . /zrchain
WORKDIR /zrchain

# Download dependencies
RUN go mod download
RUN GO111MODULES=on \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build \
    -tags "netgo muslc" \
    -ldflags "-w -linkmode external -extldflags '-Wl,-z,muldefs -static'" \
    -o validator_sidecar ./sidecar

############################################################################################################

#SSL certs
FROM alpine:3.18.0 AS certs
RUN apk add --no-cache ca-certificates
RUN adduser -Ds /bin/bash appuser

# Create writable directory and symlink
RUN mkdir /app && \
    chown appuser:appuser /app

# Copy binary to a scratch container. Let's keep our images nice and small!
COPY --from=go_builder /zrchain/validator_sidecar /app/validator_sidecar
COPY --from=go_builder /zrchain/sidecar/config.yaml.example /app/config.yaml

# Set user and workdir
USER appuser
WORKDIR /app

# Run the binary
ENTRYPOINT [ "/app/validator_sidecar" ]
