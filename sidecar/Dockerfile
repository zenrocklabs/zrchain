# Use a  golang alpine as the base image
FROM public.ecr.aws/docker/library/golang:1.23.0-alpine3.20 AS go_builder

RUN apk update
RUN apk add make cmake git alpine-sdk linux-headers

# Setup
# Read arguments
ARG ARCH=x86_64
ARG BUILD_DATE
ARG GIT_SHA
ARG VERSION
ARG GITHUB_TOKEN

# Set env variables
ENV GOPRIVATE="https://github.com/zenrocklabs/zenrock-avs"
ENV arch=$ARCH
ENV build_date=$BUILD_DATE
ENV commit_hash=$GIT_SHA
ENV service_name=zenrockd
ENV version=$VERSION
ENV github_token=$GITHUB_TOKEN
RUN echo "building service: ${service_name}, version: ${version}, build date: ${build_date}, commit hash: ${commit_hash}, architecture: ${arch}"

# Set the working directory
COPY zrchain /zrchain
COPY avs /avs
COPY tools /tools
WORKDIR /avs/validator_sidecar

# Git config
RUN git config --global \
    url."https://git:${github_token}@github.com/zenrocklabs/zenrock-avs".insteadOf \
    "https://github.com/zenrocklabs/zenrock-avs"

# Download dependencies
RUN go mod download
RUN GO111MODULES=on \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags "-w -linkmode external -extldflags '-static'" \
    -o validator_sidecar

############################################################################################################

#SSL certs
FROM alpine:3.18.0 AS certs
RUN apk add --no-cache ca-certificates
RUN adduser -Ds /bin/bash appuser

# Copy binary to a scratch container. Let's keep our images nice and small!
COPY --from=go_builder /avs/validator_sidecar/validator_sidecar /validator_sidecar
COPY --from=go_builder /avs/validator_sidecar/config.yaml /config.yaml

# Set user
USER appuser

# Run the binary
ENTRYPOINT [ "/validator_sidecar" ]
