## Run Command

```
bazelisk run //hack/create-api-key:create-api-key -- create \
  --cert=full-path-to-certificate \
  --web-gateway=web-gateway-address \
  --service-token=full-path-to-service-token-file \
  --name=key-name \
  --role=key-role \
  --output=full-path-to-output-file \
  --environment-namespace=environment-namespace \
  --no-profile \
  --no-gcp-trace-enabled
```
