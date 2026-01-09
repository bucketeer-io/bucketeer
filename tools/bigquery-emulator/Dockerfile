# Multi-arch Dockerfile for bigquery-emulator
# Builds from upstream goccy/bigquery-emulator source with ARM64 support
# Note: Base image (go-zetasql) only supports AMD64, so ARM64 builds use emulation
# This is expected and harmless - the final binary will be native ARM64
# Pin to specific version for reproducible builds
FROM ghcr.io/goccy/go-zetasql:0.5.5 AS builder

ARG VERSION=0.6.6
ARG TARGETARCH
ARG TARGETOS

# Configure Go to use multiple proxies and increase timeout
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org
ENV GOPRIVATE=

# Set Go build target architecture (required for cross-compilation)
# TARGETARCH is set by Docker buildx (amd64, arm64, etc.)
# Convert Docker arch names to Go arch names
ENV GOARCH=${TARGETARCH}
ENV GOOS=${TARGETOS:-linux}

WORKDIR /work

# Install git if not present
RUN which git >/dev/null 2>&1 || (apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*)

# Install cross-compilation toolchain for ARM64 builds
# When building ARM64 on AMD64 base, we need ARM64 cross-compilers for CGO
# Try clang first (better compatibility with go-zetasql), fallback to gcc
RUN if [ "$TARGETARCH" = "arm64" ] && [ "$(uname -m)" != "aarch64" ]; then \
        apt-get update && \
        apt-get install -y \
            clang \
            gcc-aarch64-linux-gnu \
            g++-aarch64-linux-gnu \
            binutils-aarch64-linux-gnu && \
        rm -rf /var/lib/apt/lists/*; \
    fi

# Clone upstream bigquery-emulator repository with shallow clone
# Use shallow clone to speed up and reduce size
RUN git clone --depth 1 --branch v${VERSION} https://github.com/goccy/bigquery-emulator.git . || \
    (git clone --depth 1 https://github.com/goccy/bigquery-emulator.git . && \
     git checkout v${VERSION} 2>/dev/null || git checkout ${VERSION} 2>/dev/null || git checkout $(git describe --tags --abbrev=0 2>/dev/null || echo "v0.6.6"))

# Download dependencies first (with retry logic) - this helps with caching and timeout issues
RUN go mod download || \
    (sleep 5 && go mod download) || \
    (sleep 10 && go mod download) || \
    (sleep 15 && go mod download)

# Verify dependencies are downloaded
RUN go mod verify

# Build the emulator
# Explicitly set GOARCH/GOOS to ensure correct target architecture
# This is critical when building ARM64 on AMD64 base image with emulation
# For ARM64 cross-compilation, use clang with cross-compilation target
# CGO_ENABLED=1 is scoped to the build command for explicit behavior
# -trimpath removes local paths for reproducible builds
RUN if [ "$TARGETARCH" = "arm64" ] && [ "$(uname -m)" != "aarch64" ]; then \
        CC="clang -target aarch64-linux-gnu" \
        CXX="clang++ -target aarch64-linux-gnu" \
        CGO_CFLAGS="-target aarch64-linux-gnu" \
        CGO_CXXFLAGS="-target aarch64-linux-gnu" \
        GOARCH=${GOARCH} GOOS=${GOOS} \
        CGO_ENABLED=1 \
        go build -trimpath -o bigquery-emulator \
        -ldflags='-s -w -X main.version=${VERSION}' \
        ./cmd/bigquery-emulator; \
    else \
        CXX=clang++ GOARCH=${GOARCH} GOOS=${GOOS} \
        CGO_ENABLED=1 \
        go build -trimpath -o bigquery-emulator \
        -ldflags='-s -w -X main.version=${VERSION}' \
        ./cmd/bigquery-emulator; \
    fi

# Final stage
# Use bookworm (Debian 12) to match GLIBC version from builder
# The builder uses a newer base with GLIBC 2.34+, so we need matching runtime
FROM debian:bookworm-slim

COPY --from=builder /work/bigquery-emulator /bin/bigquery-emulator

WORKDIR /work

ENTRYPOINT ["/bin/bigquery-emulator"]
