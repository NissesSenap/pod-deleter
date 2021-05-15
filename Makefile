SHELL := /bin/bash

TAG = dev
IMG ?= azad-kube-proxy:$(TAG)
TEST_ENV_FILE = tmp/test_env
VERSION ?= "v0.0.0-dev"
REVISION ?= ""
CREATED ?= ""
K8DASH_DIR ?= ${PWD}/pkg/dashboard/static/k8dash


ifneq (,$(wildcard $(TEST_ENV_FILE)))
    include $(TEST_ENV_FILE)
    export
endif

.SILENT:
all: tidy lint fmt vet gosec go-test cover

.SILENT:
lint:
	golangci-lint run

.SILENT:
fmt:
	go fmt ./...

.SILENT:
tidy:
	go mod tidy

.SILENT:
vet:
	go vet ./...

.SILENT:
go-test:
	mkdir -p tmp
	go test -timeout 1m ./... -cover

.SILENT:
gosec:
	gosec ./...

.SILENT:
cover:
	go test -timeout 1m ./... -coverprofile=tmp/coverage.out                                                                                                                                                                                         16:10:38
	go tool cover -html=tmp/coverage.out
