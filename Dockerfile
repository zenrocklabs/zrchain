# Example usage:
#
# docker build -f Dockerfile -t zenrockd:0.0.1 .
# docker run -e DOCKER_ENV=true -p 8080:8080 zenrockd:0.0.1

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
ENV arch=$ARCH
ENV build_date=$BUILD_DATE
ENV commit_hash=$GIT_SHA
ENV service_name=zenrockd
ENV version=$VERSION
RUN echo "building service: ${service_name}, version: ${version}, build date: ${build_date}, commit hash: ${commit_hash}, architecture: ${arch}"


# Add libwasmvm for musl
ENV WASMVM_VERSION=v2.1.2
ADD https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/$WASMVM_VERSION/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 0881c5b463e89e229b06370e9e2961aec0a5c636772d5142c68d351564464a66
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 58e1f6bfa89ee390cb9abc69a5bc126029a497fe09dd399f38a82d0d86fe95ef

RUN cp /lib/libwasmvm_muslc.${arch}.a /lib/libwasmvm_muslc.a

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
RUN apk add --no-cache ca-certificates
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
