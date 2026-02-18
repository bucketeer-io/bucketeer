TAG=v0.0.9
IMAGE=ghcr.io/bucketeer-io/bucketeer-delete-e2e-data-mysql

.PHONY: deps
deps:
	go mod tidy
	go mod vendor

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./delete-e2e-data-mysql .

.PHONY: build-darwin
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./delete-e2e-data-mysql .

.PHONY: docker-build
docker-build:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE):$(TAG) \
		.

.PHONY: docker-push
docker-push:
	@echo $(PAT) | docker login ghcr.io -u $(GITHUB_USER_NAME) --password-stdin
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-t $(IMAGE):$(TAG) \
		--push .

.PHONY: clean
clean:
	rm -f ./delete-e2e-data-mysql
