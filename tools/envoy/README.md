# bucketeer-envoy

Because Envoy may shut downs before the application we need to run some scripts in the `preStop` hook.
The Envoy image doesn't have wget or curl installed by default, so we need to create an image with it so we can execute the scripts.

## How to image build

```sh
make docker-build
```

## How to image push to Github Container Registry

Only users with appropriate permissions can push.

```sh
PAT=${GITHUB_PERSONAL_ACCESS_TOKEN} \
GITHUB_USER_NAME=${GITHUB_USER_NAME} \
make docker-push-ghcr
```

## How to image push to Google Artifact Registry

Only users with appropriate permissions can push.\
You must set up GAR authentication only once before running it.

```sh
make docker-push-gar
```
