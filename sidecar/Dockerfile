# Use a  golang alpine as the base image
FROM public.ecr.aws/docker/library/golang:1.23.2-alpine3.20 AS go_builder

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

# Add libwasmvm for musl
ENV WASMVM_VERSION=v2.1.3
ADD https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep faea4e15390e046d2ca8441c21a88dba56f9a0363f92c5d94015df0ac6da1f2d
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 8dab08434a5fe57a6fbbcb8041794bc3c31846d31f8ff5fb353ee74e0fcd3093

RUN cp /lib/libwasmvm_muslc.${arch}.a /lib/libwasmvm_muslc.a


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

# Copy binary to a scratch container. Let's keep our images nice and small!
COPY --from=go_builder /zrchain/validator_sidecar /validator_sidecar
COPY --from=go_builder /zrchain/sidecar/config.yaml /config.yaml

# Set user
USER appuser

# Run the binary
ENTRYPOINT [ "/validator_sidecar" ]
