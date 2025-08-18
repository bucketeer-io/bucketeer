FROM ubuntu:jammy

COPY first-run-notice.txt /tmp/scripts/

ENV LANG="C.UTF-8" \
    DEBIAN_FRONTEND=noninteractive \
    SHELL=/bin/bash \
    DOCKER_BUILDKIT=1

# Install all packages in a single RUN command to reduce layers
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends \
        # Essential build tools
        make \
        build-essential \
        cmake \
        curl \
        git \
        unzip \
        zip \
        rsync \
        jq \
        xz-utils \
        # Database drivers (only what's actually needed)
        libpq-dev \
        default-libmysqlclient-dev \
        sqlite3 \
        libsqlite3-dev \
        # Development tools
        python3-dev \
        python3-pip \
        vim \
        protobuf-compiler \
        # Essential libraries
        libssl-dev \
        libgcc1 \
        libstdc++6 \
        zlib1g \
        libuuid1 \
        uuid-dev \
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/* \
    # Move first run notice to right spot
    && mkdir -p "/usr/local/etc/vscode-dev-containers/" \
    && mv -f /tmp/scripts/first-run-notice.txt /usr/local/etc/vscode-dev-containers/ \
    && rm -rf /tmp/scripts

# Install Go
RUN curl -fsSL https://golang.org/dl/go1.24.6.linux-$(dpkg --print-architecture).tar.gz | tar -C /usr/local -xzf -

# Install Node.js and npm manually (more reliable for multi-platform builds)
RUN ARCH=$(dpkg --print-architecture) \
    && case $ARCH in \
        amd64) NODE_ARCH="x64" ;; \
        arm64) NODE_ARCH="arm64" ;; \
        *) echo "Unsupported architecture: $ARCH" && exit 1 ;; \
    esac \
    && curl -fsSL "https://nodejs.org/dist/v22.11.0/node-v22.11.0-linux-${NODE_ARCH}.tar.xz" | tar -C /usr/local --strip-components=1 -xJ \
    && npm install -g yarn

# Set up Go environment
ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=1

# Create go workspace
RUN mkdir -p /go/bin /go/src /go/pkg

# Install Atlas for MySQL migration
RUN curl -sSf https://atlasgo.sh | sh -s -- -y

# Mount for docker-in-docker
VOLUME [ "/var/lib/docker" ]

# Simple entrypoint for dev container
CMD [ "sleep", "infinity" ]

# [Optional] Install debugger for development of Codespaces - Not in resulting image by default
ARG DeveloperBuild
RUN if [ -z $DeveloperBuild ]; then \
    echo "not including debugger" ; \
    else \
    curl -sSL https://aka.ms/getvsdbgsh | bash /dev/stdin -v latest -l /vsdbg ; \
    fi
