TAG := v1.26.2

.PHONY: docker-build
docker-build:
	docker build --platform linux/amd64 . -t bucketeer-envoy:$(TAG)

.PHONY: docker-push-ghcr
docker-push-ghcr:
	@echo $(PAT) | docker login ghcr.io -u $(GITHUB_USER_NAME) --password-stdin
	docker tag bucketeer-envoy:$(TAG) ghcr.io/bucketeer-io/bucketeer-envoy:$(TAG)
	docker push ghcr.io/bucketeer-io/bucketeer-envoy:$(TAG)

.PHONY: docker-push-gar
docker-push-gar:
	docker tag bucketeer-envoy:$(TAG) asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-envoy:$(TAG)
	docker push asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-envoy:$(TAG)
