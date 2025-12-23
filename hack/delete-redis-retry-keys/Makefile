.PHONY: deps
deps:
	go mod tidy
	go mod vendor

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./delete-redis-retry-keys .

.PHONY: build-darwin
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./delete-redis-retry-keys .

.PHONY: docker-build
docker-build: clean build
	docker build . -t bucketeer-delete-redis-retry-keys:$(TAG)

.PHONY: docker-push
docker-push:
	@echo $(PAT) | docker login ghcr.io -u $(GITHUB_USER_NAME) --password-stdin
	docker tag bucketeer-delete-redis-retry-keys:$(TAG) ghcr.io/bucketeer-io/bucketeer-delete-redis-retry-keys:$(TAG)
	docker push ghcr.io/bucketeer-io/bucketeer-delete-redis-retry-keys:$(TAG)

.PHONY: clean
clean:
	rm -f ./delete-redis-retry-keys

