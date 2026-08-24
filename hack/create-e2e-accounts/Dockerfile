# Build on the host's native architecture and cross-compile for the target:
# running the Go compiler under QEMU emulation takes ~10x longer.
FROM --platform=$BUILDPLATFORM golang:1.26.7 AS builder
ARG TARGETOS TARGETARCH

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o create-e2e-accounts .

FROM gcr.io/distroless/base

# Links the ghcr.io package to this repository, so packages created by the
# image-build workflow automatically grant the repo's GITHUB_TOKEN write access.
LABEL org.opencontainers.image.source=https://github.com/bucketeer-io/bucketeer

COPY --from=builder /app/create-e2e-accounts /usr/local/bin/

ENTRYPOINT ["create-e2e-accounts"]
