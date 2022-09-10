IMG ?= w6dio/ci-status:latest

export GO111MODULE  := on
export PATH         := ./bin:${PATH}
export NEXT_TAG     ?=

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

ifeq (,$(shell go env GOOS))
GOOS       = $(shell echo $OS)
else
GOOS       = $(shell go env GOOS)
endif

ifeq (,$(shell go env GOARCH))
GOARCH     = $(shell echo uname -p)
else
GOARCH     = $(shell go env GOARCH)
endif

ifeq (gsed not found,$(shell which gsed))
SEDBIN=sed
else
SEDBIN=$(shell which gsed)
endif

ifeq (darwin,$(GOOS))
GOTAGS = "-tags=dynamic"
else
GOTAGS =
endif

export PATH := $(shell pwd)/bin:${PATH}

REF        = $(shell git symbolic-ref --quiet HEAD 2> /dev/null)
VERSION   ?= $(shell basename /$(shell git symbolic-ref --quiet HEAD 2> /dev/null ) )
VCS_REF    = $(shell git rev-parse HEAD)
BUILD_DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

all: ci-status

# go-get-tool will 'go get' any package $2 and install it to $1.
PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
define go-get-tool
@[ -f $(1) ] || { \
set -e ;\
echo "Downloading $(2)" ;\
GOBIN=$(PROJECT_DIR)/bin go install $(2) ;\
rm -rf $$TMP_DIR ;\
}
endef

GOIMPORTS  = $(shell pwd)/bin/goimports
bin/goimports: ## Download goimports locally if necessary
	$(call go-get-tool,$(GOIMPORTS),golang.org/x/tools/cmd/goimports)


# Run tests
test: fmt vet
	go test $(GOTAGS) -v -coverpkg=./... -coverprofile=cover.out ./...
	@go tool cover -func cover.out | grep total

# Build ci-status binary
ci-status: fmt vet
	VERSION=${VERSION/refs\/heads\//}
	go build $(GOTAGS) -ldflags="-X 'main.Version=${VERSION}' -X 'main.Revision=${VCS_REF}' -X 'main.Built=${BUILD_DATE}'" -a -o bin/ci-status cmd/ci-status/main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run $(GOTAGS) cmd/ci-status/main.go -config config/tests/config.yaml -log-level 2

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet $(GOTAGS) ./...

# Formats the code
.PHONY: format
format: bin/goimports
	bin/goimports -w -local gitlab.w6d.io/w6d,github.com/w6d-io internal pkg cmd

# Build the docker image
.PHONY: docker-build
docker-build:
	docker build --build-arg=VERSION=${VERSION} --build-arg=VCS_REF=${VCS_REF} --build-arg=BUILD_DATE=${BUILD_DATE}  -t ${IMG} .

# Push the docker image
.PHONY: docker-push
docker-push:
	docker push ${IMG}

