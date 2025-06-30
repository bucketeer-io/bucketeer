TAG := v1.34.1

.PHONY: docker-build
docker-build:
	docker build --platform linux/amd64 . -t bucketeer-envoy:$(TAG)

.PHONY: docker-push-gar
docker-push-gar:
	docker tag bucketeer-envoy:$(TAG) asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-envoy:$(TAG)
	docker push asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-envoy:$(TAG)
