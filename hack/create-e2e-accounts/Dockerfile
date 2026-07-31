FROM golang:1.26.5 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o create-e2e-accounts .

FROM gcr.io/distroless/base

# Links the ghcr.io package to this repository, so packages created by the
# image-build workflow automatically grant the repo's GITHUB_TOKEN write access.
LABEL org.opencontainers.image.source=https://github.com/bucketeer-io/bucketeer

COPY --from=builder /app/create-e2e-accounts /usr/local/bin/

ENTRYPOINT ["create-e2e-accounts"]
