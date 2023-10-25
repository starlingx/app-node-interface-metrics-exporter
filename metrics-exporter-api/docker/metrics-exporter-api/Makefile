#
# Copyright (c) 2023 Wind River Systems, Inc.
#
# SPDX-License-Identifier: Apache-2.0
#
# All Rights Reserved.
#

# variables
GO					?= go
GOCOVER             ?= $(GO) tool cover
LABEL				?= metrics-exporter-api
VERSION				?= v0.0.1-stable

# Targets
help:               ## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

apphelp:             ## Show this help.
	$(GO) run . -h

install_dep:        ## install go dependency
	$(GO) mod tidy

run:                ## run app on host machine
	$(GO) run .

test:               ## run go unit test
	$(GO) test ./...

testcov:            ## run go coverage test
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

vet:                ## run go vet
	$(GO) vet

lint:               ## run go lint
	golangci-lint run

build_linux:        ## Build application
	CGO_ENABLED=0 GOOS=linux go build -v -ldflags="-X 'main.Version=$(VERSION)'" -o metrics-api-server

build_image:        ## Build docker image
	docker build -f ../../debian/Dockerfile -t starlingx/metrics-exporter-api ../
