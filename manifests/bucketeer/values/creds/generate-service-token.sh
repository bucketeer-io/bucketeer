#!/bin/bash

ISSUER=https://localhost:8000/dex
OUTPUT_DIR=${PWD}

# go get github.com/bucketeer-io/bucketeer

cd ${GOPATH}/src/github.com/bucketeer-io/bucketeer/
go install github.com/golang/protobuf/protoc-gen-go@v1.5.2
make proto-go
go run ./hack/generate-service-token generate \
    --issuer=${ISSUER} \
    --sub=service \
    --audience=bucketeer \
    --email=localenv@bucketeer.io \
    --role=OWNER \
    --key=${OUTPUT_DIR}/${1}/oauth-private.pem \
    --output=${OUTPUT_DIR}/${1}/service.token \
    --no-profile \
    --no-gcp-trace-enabled
