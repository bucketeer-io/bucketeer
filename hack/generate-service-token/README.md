## Run Command

```
go run ./hack/generate-service-token generate \
  --issuer=issuer \
  --audience=audience \
  --email=email \
  --service=service \
  --ttl=ttl \
  --key=full-path-to-private-key \
  --output=full-path-to-output-file \
  --no-profile \
  --no-gcp-trace-enabled
```

`--service` names the service the token is issued to (`web`, `api`, `batch`,
`subscriber`, or a name chosen for an operator tool). Without it, the token
is issued in the legacy form.

`--ttl` is how long the token is valid, such as `8760h`. The default is one
year with `--service` and 100 years without.

`--key` is the OAuth private key, or the service token private key when the
deployment verifies service tokens with their own key
(`serviceToken.key.secret` in the web and batch charts).
