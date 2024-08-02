DEV_DIR := $(CURDIR)/.temp

include .make/tag.mk

APP_NAME := go-text-replacer
APP_EXT := $(if $(filter Windows_NT,$(OS)),.exe)

build:
	@go build -o .bin/$(APP_NAME)$(APP_EXT) example/*

run: build
	@.bin/$(APP_NAME)${APP_EXT}

test:
	@mkdir -p $(DEV_DIR)/.temp
	@CGO_ENABLED=0 go test \
	. \
	-coverprofile=$(DEV_DIR)/.temp/coverage-report.out -covermode=count
	@go tool cover -html=$(DEV_DIR)/.temp/coverage-report.out -o $(DEV_DIR)/.temp/coverage-report.html
	@go tool cover -func=$(DEV_DIR)/.temp/coverage-report.out

.PHONY: build run test
