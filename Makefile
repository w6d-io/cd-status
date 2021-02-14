# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: ci-status

# Run tests
test: fmt vet
	go test -v -coverpkg=./... -coverprofile=cover.out ./...
	@go tool cover -func cover.out | grep total

# Build ci-status binary
ci-status: fmt vet
	go build -o bin/ci-status cmd/ci-status/main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: fmt vet
	go run cmd/ci-status/main.go -config config/tests/config.yaml -log-level 2

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Build the docker image
build: test
	docker build . -t ${IMG}

# Push the docker image
push:
	docker push ${IMG}

