.PHONY: all test test-local install-deps lint fmt vet

REPO_NAME = go-plugins-helpers
REPO_OWNER = docker
PKG_NAME = github.com/${REPO_OWNER}/${REPO_NAME}
IMAGE = golang:1.21

all: test

test-local: install-deps fmt lint vet
	@echo "+ $@"
	@go test -v ./...

test:
	@docker run --rm -v ${shell pwd}:/go/src/${PKG_NAME} -w /go/src/${PKG_NAME} ${IMAGE} make test-local

install-deps:
	@echo "+ $@"
	@go install golang.org/x/lint/golint@latest

lint:
	@echo "+ $@"
	@test -z "$$(golint ./... | tee /dev/stderr)"

fmt:
	@echo "+ $@"
	@test -z "$$(gofmt -s -l . | tee /dev/stderr)"

vet:
	@echo "+ $@"
	@go vet ./...

