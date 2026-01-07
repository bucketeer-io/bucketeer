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

The version is defined in the `Dockerfile` as the single source of truth:

```dockerfile
ARG VERSION=0.6.6
```

The Makefile and GitHub Actions workflow automatically extract this version.

To override the version temporarily:

```bash
make docker-build VERSION=0.7.0
```

## Automated Builds

GitHub Actions automatically builds and pushes multi-arch images when:
- Changes are pushed to `tools/bigquery-emulator/**` on main branch
- Manual workflow dispatch is triggered (with optional version override)
- Weekly schedule (to catch upstream updates)

See `.github/workflows/bigquery-emulator-build.yaml` for the automation.

## Updating to New Upstream Version

1. Update the `VERSION` ARG in the `Dockerfile`
2. Commit and push to main
3. GitHub Actions will automatically build and push the new image

```bash
# Example: Update to version 0.7.0
sed -i 's/ARG VERSION=.*/ARG VERSION=0.7.0/' Dockerfile
git add Dockerfile
git commit -m "chore: bump bigquery-emulator to 0.7.0"
git push
```

## Image Location

Built images are published to:
- `ghcr.io/bucketeer-io/bigquery-emulator:<version>` (versioned)
- `ghcr.io/bucketeer-io/bigquery-emulator:latest` (latest)

## Usage in Kubernetes

The image is used in `manifests/localenv/dependencies/bq/values.yaml`:

```yaml
image:
  repository: ghcr.io/bucketeer-io/bigquery-emulator
  tag: "0.6.6"
```

And overridden in `manifests/localenv/values.yaml`:

```yaml
bq:
  image:
    repository: ghcr.io/bucketeer-io/bigquery-emulator
    tag: "0.6.6"
```
