APP_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
APP_NAME := $(shell basename $(APP_PATH))
APP_PKG := $(shell go list -m)
SCRIPT_PATH := $(APP_PATH)/scripts
COMPILE_OUT:=$(APP_PATH)/bin/$(APP_NAME)

.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: install test

install: install.go
build: build.ui build.linux build.darwin

install.go:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go get $(GOFLAGS) ./...
	@echo -e "\n"

build.ui:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@cd $(APP_PATH)/webui && yarn install --frozen-lockfile && yarn run build
	@echo -e "\n"

build.linux:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>make $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@export GOOS=linux && $(SCRIPT_PATH)/build/gobuild.sh $(APP_NAME) $(COMPILE_OUT)-$${GOOS} $(APP_PKG)
	@echo -e "\n"

build.darwin:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>make $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@export GOOS=darwin && $(SCRIPT_PATH)/build/gobuild.sh $(APP_NAME) $(COMPILE_OUT)-$${GOOS} $(APP_PKG)
	@echo -e "\n"

test: install
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go test $(GOFLAGS) ./...
	@echo -e "\n"

bench: install
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go test -run=NONE -bench=. $(GOFLAGS) ./...
	@echo -e "\n"

clean:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>making $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@go clean $(GOFLAGS) -i ./...
	@echo -e "\n"

