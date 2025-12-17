# Makefile for building bigquery-emulator multi-arch images
# Extract VERSION from Dockerfile as the single source of truth
VERSION ?= $(shell grep -E '^ARG VERSION=' Dockerfile | sed 's/ARG VERSION=//')
LOCAL_IMAGE := bucketeer-bigquery-emulator:$(VERSION)
GHCR_IMAGE := ghcr.io/bucketeer-io/bigquery-emulator:$(VERSION)
GAR_IMAGE := asia-docker.pkg.dev/bucketeer-io/bucketeer/bucketeer-bigquery-emulator:$(VERSION)

.PHONY: docker-build
docker-build:
	@echo "Building bigquery-emulator for current platform..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		-t $(LOCAL_IMAGE) \
		.

.PHONY: docker-build-amd64
docker-build-amd64:
	@echo "Building bigquery-emulator for linux/amd64..."
	docker build \
		--platform linux/amd64 \
		--build-arg VERSION=$(VERSION) \
		-t $(LOCAL_IMAGE)-amd64 \
		-t $(GHCR_IMAGE)-amd64 \
		.

.PHONY: docker-build-arm64
docker-build-arm64:
	@echo "Building bigquery-emulator for linux/arm64..."
	docker build \
		--platform linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		-t $(LOCAL_IMAGE)-arm64 \
		-t $(GHCR_IMAGE)-arm64 \
		.

.PHONY: docker-build-multiarch
docker-build-multiarch:
	@echo "Building bigquery-emulator for both architectures in parallel..."
	@echo "Note: Requires sufficient memory (16GB+ recommended)"
	@echo "Note: Builds images locally - use 'make docker-push-multiarch' to push after authentication"
	@echo "Note: Use 'make docker-build-multiarch-sequential' if you encounter OOM issues"
	docker buildx create --use --name multiarch-builder 2>/dev/null || docker buildx use multiarch-builder
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		-t $(GHCR_IMAGE)-amd64 \
		-t $(GHCR_IMAGE)-arm64 \
		--load \
		.
	@echo "Multi-arch build complete! Images are loaded locally."
	@echo "Run 'make docker-push-multiarch' to push after authenticating to GHCR."

.PHONY: docker-build-multiarch-sequential
docker-build-multiarch-sequential:
	@echo "Building bigquery-emulator for both architectures sequentially..."
	@echo "Use this if parallel build fails due to OOM"
	docker buildx create --use --name multiarch-builder 2>/dev/null || docker buildx use multiarch-builder
	@echo "Building AMD64 architecture..."
	docker buildx build \
		--platform linux/amd64 \
		--build-arg VERSION=$(VERSION) \
		-t $(GHCR_IMAGE)-amd64 \
		--push \
		.
	@echo "Building ARM64 architecture..."
	docker buildx build \
		--platform linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		-t $(GHCR_IMAGE)-arm64 \
		--push \
		.
	@echo "Creating and pushing multi-arch manifest..."
	@make docker-manifest-create
	@make docker-manifest-push
	@echo "Tagging latest..."
	docker manifest create ghcr.io/bucketeer-io/bigquery-emulator:latest \
		--amend $(GHCR_IMAGE)-amd64 \
		--amend $(GHCR_IMAGE)-arm64
	docker manifest push ghcr.io/bucketeer-io/bigquery-emulator:latest
	@echo "Multi-arch build complete!"

.PHONY: docker-push-multiarch
docker-push-multiarch:
	@echo "Pushing multi-arch images to GitHub Container Registry..."
	@if [ -z "$(PAT)" ]; then \
		echo "Error: PAT environment variable is required"; \
		echo "Usage: PAT=\$$GITHUB_PERSONAL_ACCESS_TOKEN GITHUB_USER_NAME=\$$GITHUB_USER_NAME make docker-push-multiarch"; \
		exit 1; \
	fi
	@if [ -z "$(GITHUB_USER_NAME)" ]; then \
		echo "Error: GITHUB_USER_NAME environment variable is required"; \
		echo "Usage: PAT=\$$GITHUB_PERSONAL_ACCESS_TOKEN GITHUB_USER_NAME=\$$GITHUB_USER_NAME make docker-push-multiarch"; \
		exit 1; \
	fi
	@echo $(PAT) | docker login ghcr.io -u $(GITHUB_USER_NAME) --password-stdin
	@echo "Pushing AMD64 image..."
	docker push $(GHCR_IMAGE)-amd64
	@echo "Pushing ARM64 image..."
	docker push $(GHCR_IMAGE)-arm64
	@echo "Creating and pushing multi-arch manifest..."
	@make docker-manifest-create
	@make docker-manifest-push
	@echo "Multi-arch images pushed successfully!"

.PHONY: docker-push-ghcr
docker-push-ghcr:
	@echo "Pushing to GitHub Container Registry..."
	@if [ -z "$(PAT)" ]; then \
		echo "Error: PAT environment variable is required"; \
		exit 1; \
	fi
	@echo $(PAT) | docker login ghcr.io -u $(GITHUB_USER_NAME) --password-stdin
	docker tag $(LOCAL_IMAGE) $(GHCR_IMAGE)
	docker push $(GHCR_IMAGE)

.PHONY: docker-push-gar
docker-push-gar:
	@echo "Pushing to Google Artifact Registry..."
	docker tag $(LOCAL_IMAGE) $(GAR_IMAGE)
	docker push $(GAR_IMAGE)

.PHONY: docker-manifest-create
docker-manifest-create:
	@echo "Creating multi-arch manifest..."
	@echo "Removing old local manifests if they exist..."
	-docker manifest rm $(GHCR_IMAGE) 2>/dev/null || true
	-docker manifest rm ghcr.io/bucketeer-io/bigquery-emulator:latest 2>/dev/null || true
	@echo "Extracting image digests from architecture-specific manifest lists..."
	@AMD64_DIGEST=$$(docker manifest inspect $(GHCR_IMAGE)-amd64 2>/dev/null | jq -r '.manifests[] | select(.platform.architecture == "amd64") | .digest' | head -1); \
	ARM64_DIGEST=$$(docker manifest inspect $(GHCR_IMAGE)-arm64 2>/dev/null | jq -r '.manifests[] | select(.platform.architecture == "arm64") | .digest' | head -1); \
	if [ -z "$$AMD64_DIGEST" ] || [ -z "$$ARM64_DIGEST" ]; then \
		echo "Error: Could not extract digests. Trying alternative method..."; \
		AMD64_DIGEST=$$(docker manifest inspect $(GHCR_IMAGE)-amd64 2>/dev/null | grep -A 5 '"amd64"' | grep '"digest"' | head -1 | sed 's/.*"digest": "\([^"]*\)".*/\1/'); \
		ARM64_DIGEST=$$(docker manifest inspect $(GHCR_IMAGE)-arm64 2>/dev/null | grep -A 5 '"arm64"' | grep '"digest"' | head -1 | sed 's/.*"digest": "\([^"]*\)".*/\1/'); \
	fi; \
	if [ -z "$$AMD64_DIGEST" ] || [ -z "$$ARM64_DIGEST" ]; then \
		echo "Error: Failed to get image digests"; \
		exit 1; \
	fi; \
	echo "Using digests: AMD64=$$AMD64_DIGEST, ARM64=$$ARM64_DIGEST"; \
	docker manifest create $(GHCR_IMAGE) \
		--amend ghcr.io/bucketeer-io/bigquery-emulator@$$AMD64_DIGEST \
		--amend ghcr.io/bucketeer-io/bigquery-emulator@$$ARM64_DIGEST
	@AMD64_DIGEST=$$(docker manifest inspect $(GHCR_IMAGE)-amd64 2>/dev/null | jq -r '.manifests[] | select(.platform.architecture == "amd64") | .digest' | head -1); \
	ARM64_DIGEST=$$(docker manifest inspect $(GHCR_IMAGE)-arm64 2>/dev/null | jq -r '.manifests[] | select(.platform.architecture == "arm64") | .digest' | head -1); \
	if [ -z "$$AMD64_DIGEST" ] || [ -z "$$ARM64_DIGEST" ]; then \
		AMD64_DIGEST=$$(docker manifest inspect $(GHCR_IMAGE)-amd64 2>/dev/null | grep -A 5 '"amd64"' | grep '"digest"' | head -1 | sed 's/.*"digest": "\([^"]*\)".*/\1/'); \
		ARM64_DIGEST=$$(docker manifest inspect $(GHCR_IMAGE)-arm64 2>/dev/null | grep -A 5 '"arm64"' | grep '"digest"' | head -1 | sed 's/.*"digest": "\([^"]*\)".*/\1/'); \
	fi; \
	docker manifest create ghcr.io/bucketeer-io/bigquery-emulator:latest \
		--amend ghcr.io/bucketeer-io/bigquery-emulator@$$AMD64_DIGEST \
		--amend ghcr.io/bucketeer-io/bigquery-emulator@$$ARM64_DIGEST

.PHONY: docker-manifest-push
docker-manifest-push:
	@echo "Pushing multi-arch manifest..."
	docker manifest push $(GHCR_IMAGE)
	docker manifest push ghcr.io/bucketeer-io/bigquery-emulator:latest

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  docker-build          - Build for current platform"
	@echo "  docker-build-amd64    - Build for linux/amd64"
	@echo "  docker-build-arm64   - Build for linux/arm64"
	@echo "  docker-build-multiarch - Build multi-arch images locally (faster, needs 16GB+ RAM)"
	@echo "  docker-push-multiarch - Push multi-arch images to GHCR (requires PAT and GITHUB_USER_NAME)"
	@echo "  docker-build-multiarch-sequential - Build architectures one at a time (safer, slower)"
	@echo "  docker-push-ghcr     - Push to GHCR (requires PAT and GITHUB_USER_NAME)"
	@echo "  docker-push-gar      - Push to GAR (Google Artifact Registry)"
	@echo "  docker-manifest-create - Create multi-arch manifest"
	@echo "  docker-manifest-push  - Push multi-arch manifest"
