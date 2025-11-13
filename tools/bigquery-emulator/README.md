# BigQuery Emulator Build

This directory contains the build infrastructure for creating multi-architecture (AMD64 + ARM64) Docker images of the [bigquery-emulator](https://github.com/goccy/bigquery-emulator).

## Why This Exists

The upstream `goccy/bigquery-emulator` Docker image only supports AMD64 architecture. This build infrastructure creates multi-arch images that support both AMD64 and ARM64, enabling the emulator to run on Apple Silicon Macs and other ARM64 systems.

## Building Locally

### Single Architecture (Current Platform)

```bash
make docker-build
```

### Specific Architecture

```bash
# Build for AMD64
make docker-build-amd64

# Build for ARM64
make docker-build-arm64
```

### Multi-Architecture Build (Requires Docker Buildx)

```bash
# Build multi-arch images locally (no authentication needed)
make docker-build-multiarch

# Push multi-arch images to GHCR (requires authentication)
PAT=${GITHUB_PERSONAL_ACCESS_TOKEN} \
GITHUB_USER_NAME=${GITHUB_USER_NAME} \
make docker-push-multiarch
```

Requirements:
- Docker Buildx enabled
- For push: `PAT` and `GITHUB_USER_NAME` environment variables

## Version Management

The version can be specified via the `VERSION` variable:

```bash
make docker-build VERSION=0.6.6
```

By default, it uses version `0.6.6` (matching the upstream release).

## Automated Builds

GitHub Actions automatically builds and pushes multi-arch images when:
- Upstream bigquery-emulator releases a new version
- Manual workflow dispatch is triggered

See `.github/workflows/bigquery-emulator-build.yaml` for the automation.

## Image Location

Built images are published to:
- `ghcr.io/bucketeer-io/bigquery-emulator:0.6.6` (versioned)
- `ghcr.io/bucketeer-io/bigquery-emulator:latest` (latest)

## Usage in Kubernetes

The image is used in `manifests/localenv/values.yaml`:

```yaml
bq:
  image:
    repository: ghcr.io/bucketeer-io/bigquery-emulator
    tag: "latest"
```

## Updating to New Upstream Version

1. Update the `VERSION` in the Makefile or pass it as an argument
2. The Dockerfile will automatically checkout the corresponding tag
3. Build and test locally
4. Push via GitHub Actions or manually
