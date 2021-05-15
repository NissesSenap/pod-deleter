SHELL := /bin/bash

TAG ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
IMG ?= poddeleter:$(TAG)
DATE_FMT = +%Y-%m-%d
BUILD_DATE = $(shell date "$(DATE_FMT)")

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

build:
	go build -o bin/poddeleter ./main.go

.SILENT:
container/build:
	docker build --build-arg BUILD_DATE=$(BUILD_DATE) --build-arg VERSION=$(TAG) . -t $(IMG)

.SILENT:
container/push:
	docker push $(IMG)
