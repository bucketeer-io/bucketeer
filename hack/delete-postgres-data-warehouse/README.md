# Delete PostgreSQL Data Warehouse

This tool truncates the PostgreSQL data warehouse tables (`evaluation_event` and `goal_event`).

## Usage

```bash
# Build the tool
make build

# Run the tool
./delete-postgres-data-warehouse truncate \
    --postgres-user=bucketeer \
    --postgres-pass=bucketeer \
    --postgres-host=localhost \
    --postgres-port=5432 \
    --postgres-db-name=bucketeer \
    --no-profile \
    --no-gcp-trace-enabled
```

## From the root Makefile

```bash
# For minikube
make delete-postgres-data-warehouse-data

# For docker-compose
make docker-compose-delete-postgres-data-warehouse-data
```

## Release docker image

The image is built and pushed by the
`build-delete-postgres-data-warehouse-image.yaml` GitHub Actions workflow, which
authenticates to ghcr.io with the built-in `GITHUB_TOKEN` (no PAT needed).

1. Bump `TAG` in this directory's `Makefile` to a version greater than the
   newest `vX.Y.Z` tag already published.
2. Merge to `main`. A change under `hack/delete-postgres-data-warehouse/` runs
   the workflow. It publishes a multi-arch image (linux/amd64 + linux/arm64)
   only when that tag is newer than the published release. The push tags the
   image with `<TAG>` and the commit SHA. The same published tag is left
   unchanged. An older tag fails the job.
3. A manual `workflow_dispatch` uses the same check. Leave the tag empty to
   use the Makefile `TAG`.

These packages were first published with a personal access token, so grant
this repository Write under each package's **Manage Actions access** once.
Until that is set, the workflow cannot list the published tags and the job
fails without pushing.

For a local build without pushing (e.g. to verify the Dockerfile), run
`make deps docker-build`. The `deps` target is required: the Docker build
context is this directory only, and `go.mod` replaces `bucketeer/v2` with
`../../`, so the parent module must be vendored in first.
