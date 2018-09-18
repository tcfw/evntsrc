DOCKER_REGISTRY   ?= asia.gcr.io
IMAGE_PREFIX      ?= evntsrc-io

GO        ?= go
LDFLAGS   := -w -s
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

# Required for globs to work correctly
SHELL=/bin/bash

.PHONY: bootstrap
bootstrap:
ifndef HAS_GLIDE
	${GO} get -u github.com/Masterminds/glide
endif
	glide install --strip-vendor

.DEFAULT_GOAL := all
.PHONY: all
all: storer websocks stsmetrics push

storer:
	docker build -f ./build/storer/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/storer:latest .

websocks:
	docker build -f ./build/websocks/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/websocks:latest .

stsmetrics:
	docker build -f ./build/stsmetrics/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/stsmetrics:latest .

.PHONY: push
push: 
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/storer:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/websocks:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/stsmetrics:latest

protos:
	@scripts/protos.sh