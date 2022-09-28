## Run Command

```
bazelisk run //hack/create-environment:create-environment -- create \
  --cert=full-path-to-certificate \
  --web-gateway=web-gateway-address \
  --service-token=full-path-to-service-token-file \
  --id=environment-id \
  --description=optional-environment-description \
  --project-id=project-id
```
