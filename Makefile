DOCKER_REGISTRY   ?= asia.gcr.io
IMAGE_PREFIX      ?= evntsrc-io

GO        ?= go
LDFLAGS   := -w -s
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

# Required for globs to work correctly
SHELL=/bin/bash

.PHONY: bootstrap
	GO111MODULE=on go get

.DEFAULT_GOAL := changed
.PHONY: all
all: storer websocks stsmetrics streams passport users apigw bridge streamauth ingress billing wui emails metrics interconnect ttlscheduler ttlworker adapter push

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

ingress:
	docker build -f ./build/ingress/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/ingress:latest .

billing:
	docker build -f ./build/billing/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/billing:latest .

wui:
	docker build -f ./build/wui/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/wui:latest .

adapter:
	docker build -f ./build/adapter/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/adapter:latest .

emails:
	docker build -f ./build/emails/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/emails:latest .

metrics:
	docker build -f ./build/metrics/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/metrics:latest .

interconnect:
	docker build -f ./build/interconnect/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/interconnect:latest .

ttlscheduler:
	docker build -f ./build/ttlscheduler/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/ttlscheduler:latest .

ttlworker:
	docker build -f ./build/ttlworker/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/ttlworker:latest .

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
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/ingress:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/billing:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/wui:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/adapter:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/emails:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/metrics:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/interconnect:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/ttlscheduler:latest
	docker push ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/ttlworker:latest

.PHONY: changed
changed:
	@[ "${COMMIT}" ] || ( echo ">> COMMIT env variable is not set"; exit 1 )
	@COMMIT_RANGE=${COMMIT}..HEAD scripts/changed_services.sh 

protos:
	@scripts/protos.sh