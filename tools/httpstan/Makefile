.PHONY: docker-build
docker-build:
	docker build --platform linux/amd64 . -t bucketeer-httpstan:latest

.PHONY: docker-push-ghcr
docker-push-ghcr:
	@echo $(PAT) | docker login ghcr.io -u $(GITHUB_USER_NAME) --password-stdin
	docker tag bucketeer-httpstan:latest ghcr.io/bucketeer-io/bucketeer-dex:latest
	docker push ghcr.io/bucketeer-io/bucketeer-httpstan:latest

.PHONY: docker-push-gar
docker-push-gar:
	docker tag bucketeer-httpstan:latest asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-httpstan:latest
	docker push asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-httpstan:latest
