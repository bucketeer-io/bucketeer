# Evaluation module for Node.JS

## Development

### Setup

```sh
make init
make gen_proto
```

### Build

```sh
make build
```

### Unit tests

```sh
make test
```

### Lint

```sh
make lint
```

## Release to NPM

Publishing is run manually with the
[`Publish TypeScript evaluation to npm`](https://github.com/bucketeer-io/bucketeer/actions/workflows/publish-evaluation-ts.yaml)
GitHub Actions workflow.

Before the first release, configure an npm trusted publisher for
`@bucketeer/evaluation` with:

- Organization or user: `bucketeer-io`
- Repository: `bucketeer`
- Workflow filename: `publish-evaluation-ts.yaml`

To release a new version:

1. Update the version in `package.json`.
2. Merge the change into `main`.
3. Run the workflow from the `main` branch. Enable `dry_run` to verify the
   package without publishing it.

To verify another branch, select that branch when starting the workflow and
enable `dry_run`. Alternatively, run it from `main` and set `checkout_ref` to
the branch, tag, or commit SHA to test.

The workflow builds and publishes the package with npm trusted publishing
(OIDC), so an npm token is not required.
