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

GO_MODULE_BASE := github.com/cgascoig/intersight-webhook-notifier
GO_CMD ?= go
GO_BUILD_CMD := $(GO_CMD) build -v
GIT_COMMIT_SUFFIX := $(shell if [[ -n $$(git status --porcelain) ]]; then echo "-next"; else echo ""; fi)
GIT_COMMIT := $(shell git rev-parse HEAD)$(GIT_COMMIT_SUFFIX)
BUILD_DATETIME := $(shell date '+%F-%T')
GO_BUILD_FLAGS := -ldflags "-X main.commit=$(GIT_COMMIT) -X main.buildDateTime=$(BUILD_DATETIME)"
GO_PATH ?= $(shell go env GOPATH)
SKOPEO_FLAGS := --override-os linux --override-arch amd64

GLOBAL_FILES := go.mod Makefile


all: build/intersight-webhook-notifier
.PHONY: all

test:$(GLOBAL_FILES) $(ISWEBEX_FILES)
> $(GO_CMD) test ./cmd/* ./pkg/*
.PHONY: test

containers: tmp/.intersight-webhook-notifier-docker-image.sentinel
.PHONY: containers

build/container/distroless-base:
> mkdir -p $(@D)
> skopeo $(SKOPEO_FLAGS) copy docker://gcr.io/distroless/base:latest oci:$@:latest

containers-push: tmp/.intersight-webhook-notifier-docker-image-push.sentinel
.PHONY: containers-push

#####
ISWEBEX_FILES := $(shell find cmd pkg -name \*.go -type f)

build/intersight-webhook-notifier: $(GLOBAL_FILES) $(ISWEBEX_FILES)
> mkdir -p $(@D)
> $(GO_BUILD_CMD) -o "$@" $(GO_BUILD_FLAGS) $(GO_MODULE_BASE)/cmd/intersight-webhook-notifier

build/intersight-webhook-notifier-linux_amd64: $(GLOBAL_FILES) $(ISWEBEX_FILES)
> mkdir -p $(@D)
> GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o "$@" $(GO_BUILD_FLAGS) $(GO_MODULE_BASE)/cmd/intersight-webhook-notifier

tmp/.intersight-webhook-notifier-docker-image.sentinel: build/intersight-webhook-notifier-linux_amd64 Makefile build/container/distroless-base
> mkdir -p $(@D)
> mkdir -p build/container
> skopeo $(SKOPEO_FLAGS) copy oci:build/container/distroless-base:latest oci:build/container/intersight-webhook-notifier:latest
> umoci insert --image build/container/intersight-webhook-notifier:latest build/intersight-webhook-notifier-linux_amd64 /intersight-webhook-notifier
> umoci config --image build/container/intersight-webhook-notifier:latest --config.cmd /intersight-webhook-notifier
# > docker build . -f Dockerfile.webex_bot -t "$(DOCKER_IMAGE_ID_BASE)/webex_bot:latest"
> touch $@

tmp/.intersight-webhook-notifier-docker-image-push.sentinel: tmp/.intersight-webhook-notifier-docker-image.sentinel
> # always push as :latest
# > docker push "$(DOCKER_IMAGE_ID_BASE)/webex_bot:latest"
> skopeo $(SKOPEO_FLAGS) copy -f v2s2 oci:build/container/intersight-webhook-notifier:latest "docker://$(DOCKER_IMAGE_ID_BASE)/intersight-webhook-notifier:latest"
> gcloud run services update "$(GCLOUD_CLOUDRUN_SERVICE_NAME)" --region us-central1 --image "$(DOCKER_IMAGE_ID_BASE)/intersight-webhook-notifier:latest"
> touch $@