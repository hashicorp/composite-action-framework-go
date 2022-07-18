SHELL := /usr/bin/env bash -euo pipefail -c

target:
	@echo recipe

MODULE_NAME := $(shell head -n1 go.mod | cut -d' ' -f2)
HEAD_REF    := $(shell git rev-parse HEAD)

module/ref/head:
	@echo "$(MODULE_NAME)@$(HEAD_REF)"

test:
	go test ./...

testv:
	go test -v ./...
