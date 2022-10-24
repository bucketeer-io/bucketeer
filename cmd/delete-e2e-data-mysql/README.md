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
