## Delete Redis Retry Keys

A tool to delete Redis goal event retry keys for a specific environment. This is useful for cleaning up e2e test data from Redis.

The tool deletes keys matching the pattern `{environment-id}:goal_event_retry:*`.

## Run Command

```bash
go run ./hack/delete-redis-retry-keys delete \
  --redis-addr=<REDIS_HOST:PORT> \
  --environment-id=<ENVIRONMENT_ID> \
  --redis-password=<REDIS_PASSWORD> \ # optional
  --scan-count=<SCAN_COUNT> \ # optional, defaults to 100
  --no-profile \
  --no-gcp-trace-enabled
```

### Example: Delete retry keys for e2e environment (minikube)

```bash
./hack/delete-redis-retry-keys/delete-redis-retry-keys delete \
  --redis-addr=$(minikube ip):32001 \
  --environment-id=e2e \
  --no-profile \
  --no-gcp-trace-enabled
```

### Example: Delete retry keys for a specific environment

```bash
./hack/delete-redis-retry-keys/delete-redis-retry-keys delete \
  --redis-addr=localhost:6379 \
  --environment-id=production \
  --no-profile \
  --no-gcp-trace-enabled
```

## Build

```bash
# Build for Linux
make -C hack/delete-redis-retry-keys build

# Build for macOS
make -C hack/delete-redis-retry-keys build-darwin
```

## Create Docker Image

```bash
make -C hack/delete-redis-retry-keys deps

export PAT=<PERSONAL_ACCESS_TOKEN>
export GITHUB_USER_NAME=<GITHUB_USER_NAME>
export TAG=<TAG>

make -C hack/delete-redis-retry-keys docker-build
make -C hack/delete-redis-retry-keys docker-push
```

Personal Access Token needs to have `write:packages` permission.

## Using in GitHub Workflows

The tool can be used in GitHub Actions workflows to clean up Redis retry keys after e2e tests.

### Using the Makefile with REDIS_ADDR variable

```yaml
- name: Delete Redis retry keys
  run: make delete-redis-retry-keys REDIS_ADDR=${{ secrets.REDIS_ADDR }}
```

### Running directly

```yaml
- name: Delete Redis retry keys
  run: |
    make -C hack/delete-redis-retry-keys build
    ./hack/delete-redis-retry-keys/delete-redis-retry-keys delete \
      --redis-addr=${{ secrets.REDIS_ADDR }} \
      --redis-password=${{ secrets.REDIS_PASSWORD }} \
      --environment-id=e2e \
      --no-profile \
      --no-gcp-trace-enabled
```

### Using the Docker image

```yaml
- name: Delete Redis retry keys
  run: |
    docker run --rm \
      ghcr.io/bucketeer-io/bucketeer-delete-redis-retry-keys:latest \
      delete \
      --redis-addr=${{ secrets.REDIS_ADDR }} \
      --redis-password=${{ secrets.REDIS_PASSWORD }} \
      --environment-id=e2e \
      --no-profile \
      --no-gcp-trace-enabled
```
