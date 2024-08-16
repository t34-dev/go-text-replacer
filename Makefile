DEV_DIR := $(CURDIR)

include .make/tag.mk
include .make/test.mk

APP_NAME := go-text-replacer
APP_EXT := $(if $(filter Windows_NT,$(OS)),.exe)

build:
	@go build -o .bin/$(APP_NAME)$(APP_EXT) example/*

run: build
	@.bin/$(APP_NAME)${APP_EXT}



.PHONY: build run
