SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
# .DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

GO_MODULE_BASE := github.com/cgascoig/intersight-webex
DOCKER_IMAGE_ID_BASE := us-central1-docker.pkg.dev/gcp-apjccloudtsadem-nprd-50481/intersight-webex
GO_CMD ?= go
GO_BUILD_CMD := $(GO_CMD) build -v
GIT_COMMIT_SUFFIX := $(shell if [[ -n $$(git status --porcelain) ]]; then echo "-next"; else echo ""; fi)
GIT_COMMIT := $(shell git rev-parse HEAD)$(GIT_COMMIT_SUFFIX)
BUILD_DATETIME := $(shell date '+%F-%T')
GO_BUILD_FLAGS := -ldflags "-X main.commit=$(GIT_COMMIT) -X main.buildDateTime=$(BUILD_DATETIME)"
GO_PATH ?= $(shell go env GOPATH)
SKOPEO_FLAGS := --override-os linux --override-arch amd64

GLOBAL_FILES := go.mod Makefile


all: build/intersight-webex
.PHONY: all

containers: tmp/.intersight-webex-docker-image.sentinel
.PHONY: containers

#####
ISWEBEX_FILES := $(shell find cmd pkg -name \*.go -type f)

build/intersight-webex: $(GLOBAL_FILES) $(WEBEX_BOT_FILES)
> mkdir -p $(@D)
> $(GO_BUILD_CMD) -o "$@" $(GO_BUILD_FLAGS) $(GO_MODULE_BASE)/cmd/webex_bot

build/intersight-webex-linux_amd64: $(GLOBAL_FILES) $(WEBEX_BOT_FILES)
> mkdir -p $(@D)
> GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o "$@" $(GO_BUILD_FLAGS) $(GO_MODULE_BASE)/cmd/webex_bot

tmp/.intersight-webex-docker-image.sentinel: build/intersight-webex-linux_amd64 Makefile build/container/distroless-base
> mkdir -p $(@D)
> mkdir -p build/container
> skopeo $(SKOPEO_FLAGS) copy oci:build/container/distroless-base:latest oci:build/container/intersight-webex:latest
> umoci insert --image build/container/intersight-webex:latest build/intersight-webex-linux_amd64 /intersight-webex
> umoci config --image build/container/intersight-webex:latest --config.cmd /intersight-webex
# > docker build . -f Dockerfile.webex_bot -t "$(DOCKER_IMAGE_ID_BASE)/webex_bot:latest"
> touch $@