.PHONY: docker-build
docker-build:
	$(eval TAG := $(shell git describe --tags --always --abbrev=7))
	docker build --platform linux/amd64 . -t bucketeer-dex:$(TAG)

.PHONY: docker-push-ghcr
docker-push-ghcr:
	$(eval TAG := $(shell git describe --tags --always --abbrev=7))
	@echo $(PAT) | docker login ghcr.io -u $(GITHUB_USER_NAME) --password-stdin
	docker tag bucketeer-dex:$(TAG) ghcr.io/bucketeer-io/bucketeer-dex:$(TAG)
	docker push ghcr.io/bucketeer-io/bucketeer-dex:$(TAG)

.PHONY: docker-push-gar
docker-push-gar:
	$(eval TAG := $(shell git describe --tags --always --abbrev=7))
	docker tag bucketeer-dex:$(TAG) asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-dex:$(TAG)
	docker push asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-dex:$(TAG)
