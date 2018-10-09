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
all: storer websocks stsmetrics streams passport users apigw bridge streamauth push

storer:
	docker build -f ./build/storer/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/storer:latest .

websocks:
	docker build -f ./build/websocks/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/websocks:latest .

stsmetrics:
	docker build -f ./build/stsmetrics/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/stsmetrics:latest .

apigw:
	docker build -f ./build/apigw/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/apigw:latest .

streams:
	docker build -f ./build/streams/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/streams:latest .

passport:
	docker build -f ./build/passport/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/passport:latest .

users:
	docker build -f ./build/users/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/users:latest .

bridge:
	docker build -f ./build/bridge/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/bridge:latest .

streamauth:
	docker build -f ./build/streamauth/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/streamauth:latest .

.PHONY: push
push: 
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/storer:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/websocks:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/stsmetrics:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/streams:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/passport:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/users:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/apigw:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/bridge:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/streamauth:latest


protos:
	@scripts/protos.sh