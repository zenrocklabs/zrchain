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

# Add libwasmvm for musl
ENV WASMVM_VERSION=v2.1.5
ADD https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 1bad0e3f9b72603082b8e48307c7a319df64ca9e26976ffc7a3c317a08fe4b1a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep c6612d17d82b0997696f1076f6d894e339241482570b9142f29b0d8f21b280bf

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

# Copy config file (use config.yaml if it exists, otherwise use example)
COPY --from=go_builder /zrchain/sidecar/config.yaml* /tmp/
RUN if [ -f /tmp/config.yaml ]; then \
        cp /tmp/config.yaml /config.yaml; \
    else \
        cp /tmp/config.yaml.example /config.yaml; \
    fi && \
    rm -rf /tmp/config.yaml*

# Set user
USER appuser

# Run the binary
ENTRYPOINT [ "/validator_sidecar" ]
