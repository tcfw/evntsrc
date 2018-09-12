DOCKER_REGISTRY   ?= gcr.io
IMAGE_PREFIX      ?= evntsrc

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

