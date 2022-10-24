## Run Command

```
go run ./cmd/delete-e2e-data-mysql delete \
  --mysql-user=<MYSQL_USER> \
  --mysql-pass=<MYSQL_PASS> \
  --mysql-host=<MYSQL_HOST> \
  --mysql-port=<MYSQL_PORT> \
  --mysql-db-name=<MYSQL_DB_NAME> \
  --test-id=<TEST_ID> \ # optional
  --retention-seconds=<RETENTION_SECONDS> \ # optional
  --no-profile \
  --no-gcp-trace-enabled
```

Delete data created by the e2e test whose test_id is `example`.

```
go run ./cmd/delete-e2e-data-mysql delete \
  --mysql-user=sample \
  --mysql-pass=${MYSQL_PASS} \
  --mysql-host=${MYSQL_HOST} \
  --mysql-port=3306 \
  --mysql-db-name=${DB_NAME} \
  --test-id=example \
  --no-profile \
  --no-gcp-trace-enabled
```

Delete data created up to one hour ago.

```
go run ./cmd/delete-e2e-data-mysql delete \
  --mysql-user=sample \
  --mysql-pass=${MYSQL_PASS} \
  --mysql-host=${MYSQL_HOST} \
  --mysql-port=3306 \
  --mysql-db-name=${DB_NAME} \
  --retention-seconds=3600 \
  --no-profile \
  --no-gcp-trace-enabled
```
