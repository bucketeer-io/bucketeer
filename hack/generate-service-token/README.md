## Run Command

```
bazelisk run //hack/generate-service-token:generate-service-token -- generate \
  --issuer=dex-issuer \
  --sub=sub \
  --audience=client-id-set-in-dex-config \
  --email=email \
  --role=role \
  --key=full-path-to-private-key \
  --output=full-path-to-output-file \
  --no-profile \
  --no-gcp-trace-enabled
```
