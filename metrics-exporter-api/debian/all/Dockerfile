#
# Copyright (c) 2023-2026 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#
# All Rights Reserved.
#

# ARG BASE
# FROM ${BASE}

# Build Stage
FROM golang:1.24-alpine AS buildstage

# Set destination for COPY
WORKDIR /app

COPY metrics-exporter-api /app/

RUN go mod download

ARG VERSION="v1.0.0-stable"

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-X main.Version=$VERSION" -o /metrics-api-server

# Deploy binary which will make image size smaller
FROM alpine:3.22

# Set workdir context of current path wrt image
WORKDIR /

COPY --from=buildstage /metrics-api-server /metrics-api-server

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the
# docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 9110

ENTRYPOINT ["/metrics-api-server"]
