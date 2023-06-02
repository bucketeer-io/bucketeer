TAG := 4.10.1

.PHONY: docker-build
docker-build:
	docker build --build-arg TAG=$(TAG) --platform linux/amd64 . -t bucketeer-httpstan:$(TAG)

.PHONY: docker-push-ghcr
docker-push-ghcr:
	@echo $(PAT) | docker login ghcr.io -u $(GITHUB_USER_NAME) --password-stdin
	docker tag bucketeer-httpstan:$(TAG) ghcr.io/bucketeer-io/bucketeer-httpstan:$(TAG)
	docker push ghcr.io/bucketeer-io/bucketeer-httpstan:$(TAG)

.PHONY: docker-push-gar
docker-push-gar:
	docker tag bucketeer-httpstan:$(TAG) asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-httpstan:$(TAG)
	docker push asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-httpstan:$(TAG)

.PHONY: run-httpstan-container
run-httpstan-container:
	docker run --name httpstan -p 8080:8080 -d bucketeer-httpstan:$(TAG)