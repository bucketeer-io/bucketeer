# Redis Data Copy Tool

This tool is designed to copy data from one Redis instance to another.

## Building

To build the Redis data copy tool, you can use the following make commands:

1. Build the binary:
   ```
   make build
   ```
   This will create the `redis-data-copy` binary in the `bin` directory.

2. Build for Linux:
   ```
   make build-linux
   ```
   This creates a Linux-compatible binary.

3. Build Docker image:
   ```
   make docker-build
   ```
   This builds a Docker image tagged as `redis-data-copy:0.0.1`.

## Running

1. Run the binary directly:
   ```
   make run
   ```
   This builds and runs the binary.

2. Run in Docker:
   After building the Docker image, you can run it with appropriate arguments.

## Deployment

1. Push to GitHub Container Registry:
   ```
   make docker-push-ghcr
   ```
   Note: Ensure you have set the `PAT` and `GITHUB_USER_NAME` environment variables.

2. Deploy to Kubernetes:
   ```
   make kubectl-apply
   ```
   This applies the `pod.yaml` configuration to your current Kubernetes context.

## Cleaning Up

To clean the build artifacts:

```
make clean
```

This removes the binary and other build artifacts.
