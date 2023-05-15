# bucketeer-dex

Bucketeer uses dex for authentication.\
We manage our own dex image to add some customizations.

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
You must set up GAR authentication only once before run

```sh
make docker-push-gar
```
