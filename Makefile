APP_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
APP_NAME := $(shell basename $(APP_PATH))
APP_PKG := $(shell go list -m)
SCRIPT_PATH := $(APP_PATH)/scripts
COMPILE_OUT:=$(APP_PATH)/bin/$(APP_NAME)

.PHONY: all test clean build install

GOFLAGS ?= $(GOFLAGS:)

all: install test

build:build.linux build.darwin

build.linux:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>make $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@export GOOS=linux && $(SCRIPT_PATH)/build/gobuild.sh $(APP_NAME) $(COMPILE_OUT)-$${GOOS} $(APP_PKG)
	@echo -e "\n"

build.darwin:
	@echo ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>make $@<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<"
	@export GOOS=darwin && $(SCRIPT_PATH)/build/gobuild.sh $(APP_NAME) $(COMPILE_OUT)-$${GOOS} $(APP_PKG)
	@echo -e "\n"

install:
	go get $(GOFLAGS) ./...

test: install
	go test $(GOFLAGS) ./...

bench: install
	go test -run=NONE -bench=. $(GOFLAGS) ./...

clean:
	go clean $(GOFLAGS) -i ./...

