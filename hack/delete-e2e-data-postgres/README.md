## Run Command

```
go run ./hack/delete-e2e-data-postgres delete \
  --postgres-user=<POSTGRES_USER> \
  --postgres-pass=<POSTGRES_PASS> \
  --postgres-host=<POSTGRES_HOST> \
  --postgres-port=<POSTGRES_PORT> \
  --postgres-db-name=<POSTGRES_DB_NAME> \
  --test-id=<TEST_ID> \ # optional
  --no-profile \
  --no-gcp-trace-enabled
```

Delete data created by the e2e test whose test_id is `example`.

```
go run ./hack/delete-e2e-data-postgres delete \
  --postgres-user=sample \
  --postgres-pass=${POSTGRES_PASS} \
  --postgres-host=${POSTGRES_HOST} \
  --postgres-port=5432 \
  --postgres-db-name=${DB_NAME} \
  --test-id=example \
  --no-profile \
  --no-gcp-trace-enabled
```

## Create docker image

```
make deps

export PAT=<PERSONAL_ACCESS_TOKEN>
export GITHUB_USER_NAME=<GITHUB_USER_NAME>
export TAG=<TAG>

make docker-build
make docker-push
```

Personal Access Token needs to have `write:packages` permission.
