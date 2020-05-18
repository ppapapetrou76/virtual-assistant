SHELL := /bin/bash

# VARIABLES USEd
export GO111MODULE ?= on
export GOBIN = $(shell pwd)/bin

# COLORS
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20


## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)


## Build project and generate binary
build:
	@ go build -o action ./cmd

## Clean project
clean:
	@ rm -rf bin && rm action

## Install project dependencies
deps:
	@ echo "-> Installing project dependencies..."
	@ GO111MODULE=off go get -u github.com/myitcv/gobin
	@ $(GOBIN)/gobin golang.org/x/lint/golint
	@ curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(GOBIN) v1.22.2
	@ echo "-> Done."

## Formats code and fixes as many as possible linter errors
format: deps
	@ echo "-> Formatting/auto-fixing Go files..."
	@ $(GOBIN)/golangci-lint run --fix
	@ echo "-> Done."

## Runs various checks
lint: deps
	@ echo "-> Running linters..."
	@ $(GOBIN)/golint -set_exit_status ./...
	@ $(GOBIN)/golangci-lint run
	@ echo "-> Done."

## Run unit tests
test:
	@ echo "-> Running unit tests..."
	@ go test -timeout 10s -p 4 -race -count=1 ./...
	@ echo "-> Done."

