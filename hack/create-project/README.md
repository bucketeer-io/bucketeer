## Run Command

The flag `description` and the `create-environment` are optional.

E.g.

```
go run ./hack/create-project create \
  --cert=full-path-to-certificate \
  --web-gateway=web-gateway-address \
  --service-token=full-path-to-service-token-file \
  --name="Project name" \
  --url-code=url-code \
  --description="Project description" \
  --create-environment=dev \
  --no-profile \
  --no-gcp-trace-enabled
```
