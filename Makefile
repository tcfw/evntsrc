DOCKER_REGISTRY   ?= asia.gcr.io
IMAGE_PREFIX      ?= evntsrc-io

GO        ?= go
LDFLAGS   := -w -s
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin

# Required for globs to work correctly
SHELL=/bin/bash

DOCKER=docker build --compress

.PHONY: bootstrap
	GO111MODULE=on go get

.DEFAULT_GOAL := changed
.PHONY: all
all: storer websocks stsmetrics streams passport users apigw bridge streamauth ingress billing wui emails metrics interconnect ttlscheduler ttlworker push

storer:
	${DOCKER} -f ./build/storer/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/storer:latest .

websocks:
	${DOCKER} -f ./build/websocks/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/websocks:latest .

stsmetrics:
	${DOCKER} -f ./build/stsmetrics/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/stsmetrics:latest .

apigw:
	${DOCKER} -f ./build/apigw/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/apigw:latest .

streams:
	${DOCKER} -f ./build/streams/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/streams:latest .

passport:
	${DOCKER} -f ./build/passport/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/passport:latest .

users:
	${DOCKER} -f ./build/users/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/users:latest .

bridge:
	${DOCKER} -f ./build/bridge/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/bridge:latest .

streamauth:
	${DOCKER} -f ./build/streamauth/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/streamauth:latest .

ingress:
	${DOCKER} -f ./build/ingress/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/ingress:latest .

billing:
	${DOCKER} -f ./build/billing/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/billing:latest .

wui:
	${DOCKER} -f ./build/wui/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/wui:latest .

emails:
	${DOCKER} -f ./build/emails/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/emails:latest .

metrics:
	${DOCKER} -f ./build/metrics/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/metrics:latest .

interconnect:
	${DOCKER} -f ./build/interconnect/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/interconnect:latest .

ttlscheduler:
	${DOCKER} -f ./build/ttlscheduler/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/ttlscheduler:latest .

ttlworker:
	${DOCKER} -f ./build/ttlworker/Dockerfile -t ${DOCKER_REGISTRY}/${IMAGE_PREFIX}/ttlworker:latest .

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