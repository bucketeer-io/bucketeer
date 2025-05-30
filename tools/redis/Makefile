# Makefile for building redis-data-copy

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary name
BINARY_NAME=redis-data-copy

# Build directory
BUILD_DIR=bin

# Main package path
MAIN_PACKAGE=../../hack/redis-data-copy

TAG := 0.0.1

all: build

build:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PACKAGE)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

run:
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PACKAGE)
	./$(BUILD_DIR)/$(BINARY_NAME)

deps:
	$(GOGET) -v -d ./...

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v $(MAIN_PACKAGE)

docker-build:
	docker build -t redis-data-copy:$(TAG) .

docker-push-ghcr:
	@echo $(PAT) | docker login ghcr.io -u $(GITHUB_USER_NAME) --password-stdin
	docker tag redis-data-copy:$(TAG) ghcr.io/bucketeer-io/redis-data-copy:$(TAG)
	docker push ghcr.io/bucketeer-io/redis-data-copy:$(TAG)
kubectl-apply:
	kubectl apply -f pod.yaml

.PHONY: all build test clean run deps build-linux docker-build
