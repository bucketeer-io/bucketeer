## Run Command

```
go run ./hack/delete-e2e-data-mysql delete \
  --mysql-user=<MYSQL_USER> \
  --mysql-pass=<MYSQL_PASS> \
  --mysql-host=<MYSQL_HOST> \
  --mysql-port=<MYSQL_PORT> \
  --mysql-db-name=<MYSQL_DB_NAME> \
  --test-id=<TEST_ID> \ # optional
  --no-profile \
  --no-gcp-trace-enabled
```

Delete data created by the e2e test whose test_id is `example`.

```
go run ./hack/delete-e2e-data-mysql delete \
  --mysql-user=sample \
  --mysql-pass=${MYSQL_PASS} \
  --mysql-host=${MYSQL_HOST} \
  --mysql-port=3306 \
  --mysql-db-name=${DB_NAME} \
  --test-id=example \
  --no-profile \
  --no-gcp-trace-enabled
```

## Create docker image

```
make deps
make docker-build

export PAT=<PERSONAL_ACCESS_TOKEN>
export GITHUB_USER_NAME=<GITHUB_USER_NAME>
export TAG=<TAG>

make docker-push
```

Personal Access Token needs to have `write:packages` permission.
