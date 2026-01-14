FROM ubuntu:jammy

COPY first-run-notice.txt /tmp/scripts/

ENV LANG="C.UTF-8" \
    DEBIAN_FRONTEND=noninteractive \
    SHELL=/bin/bash \
    DOCKER_BUILDKIT=1

# Install packages in optimized order for better layer caching
RUN apt-get update \
    && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends \
        # Core essentials (changes rarely)
        ca-certificates \
        curl \
        gettext \
        software-properties-common \
        # Build tools (changes rarely)  
        make \
        build-essential \
        cmake \
        # File utilities
        unzip \
        zip \
        rsync \
        moreutils \
        # Minimal man support
        man-db \
        manpages \
        manpages-dev \
        groff-base \
        less \
        locales \
        # Shell completion support
        bash-completion \
    && add-apt-repository universe \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install development packages in separate layer
RUN apt-get update && apt-get install -y --no-install-recommends \
        # Development tools
        jq \
        vim \
        vim-doc \
        xtail \
        sudo \
        python3-pip \
        python3-dev \
        # Build and debug tools
        swig3.0 \
        cppcheck \
        valgrind \
        clang \
        clang-format \
        lldb \
        llvm \
        gdb \
        # Database drivers
        libpq-dev \
        default-libmysqlclient-dev \
        sqlite3 \
        libsqlite3-dev \
        unixodbc-dev \
        # System libraries
        libssl-dev \
        libc6 \
        libgcc1 \
        libgssapi-krb5-2 \
        libncurses6 \
        liblttng-ust1 \
        libstdc++6 \
        zlib1g \
        libuuid1 \
        libunwind8 \
        tk-dev \
        uuid-dev \
        libgdiplus \
        libsecret-1-dev \
        # GUI libraries (useful for 'puppeteer' project)
        libnss3 \
        libnspr4 \
        libatk-bridge2.0-0 \
        libatk1.0-0 \
        libx11-6 \
        libpangocairo-1.0-0 \
        libx11-xcb1 \
        libcups2 \
        libxcomposite1 \
        libxdamage1 \
        libxfixes3 \
        libpango-1.0-0 \
        libgbm1 \
        libgtk-3-0 \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    # Move first run notice to right spot
    && mkdir -p "/usr/local/etc/vscode-dev-containers/" \
    && mv -f /tmp/scripts/first-run-notice.txt /usr/local/etc/vscode-dev-containers/ \
    && rm -rf /tmp/scripts

# Install Go
RUN ARCH=$(dpkg --print-architecture) \
    && curl -fsSL "https://golang.org/dl/go1.25.1.linux-${ARCH}.tar.gz" | tar -C /usr/local -xzf -

# Install Node.js
RUN ARCH=$(dpkg --print-architecture) \
    && case $ARCH in \
        amd64) NODE_ARCH="x64" ;; \
        arm64) NODE_ARCH="arm64" ;; \
        *) echo "Unsupported architecture: $ARCH" && exit 1 ;; \
    esac \
    && curl -fsSL "https://nodejs.org/dist/v22.18.0/node-v22.18.0-linux-${NODE_ARCH}.tar.xz" | tar -C /usr/local --strip-components=1 -xJ \
    && npm install -g yarn

# Install kubectl, helm, and minikube
# Using fixed versions to avoid network issues fetching "stable.txt" during build
RUN ARCH=$(dpkg --print-architecture) \
    && case $ARCH in \
        amd64) K8S_ARCH="amd64" ;; \
        arm64) K8S_ARCH="arm64" ;; \
        *) echo "Unsupported architecture: $ARCH" && exit 1 ;; \
    esac \
    # Install kubectl (using fixed version for reliable builds)
    && curl -fsSL --retry 3 --retry-delay 5 "https://dl.k8s.io/release/v1.35.0/bin/linux/${K8S_ARCH}/kubectl" -o /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl \
    # Install helm
    && curl -fsSL --retry 3 --retry-delay 5 "https://get.helm.sh/helm-v3.16.4-linux-${K8S_ARCH}.tar.gz" | tar -C /tmp -xzf - \
    && mv "/tmp/linux-${K8S_ARCH}/helm" /usr/local/bin/helm \
    && rm -rf "/tmp/linux-${K8S_ARCH}" \
    # Install minikube (using fixed version for reliable builds)
    && curl -fsSL --retry 3 --retry-delay 5 "https://storage.googleapis.com/minikube/releases/v1.37.0/minikube-linux-${K8S_ARCH}" -o /usr/local/bin/minikube \
    && chmod +x /usr/local/bin/minikube

# Set up Go environment
ENV PATH="/usr/local/go/bin:$PATH"
ENV GOPATH="/go"
ENV CGO_ENABLED=1

# Create go workspace
RUN mkdir -p /go/bin /go/src /go/pkg

# Create codespace user (matching what devcontainer features would create)
RUN groupadd --gid 1000 codespace \
    && useradd --uid 1000 --gid 1000 -m codespace \
    && echo "codespace ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers.d/codespace \
    && chmod 0440 /etc/sudoers.d/codespace

# Install additional development tools and Docker
RUN ARCH=$(dpkg --print-architecture) \
    # Install GitHub CLI
    && curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | gpg --dearmor -o /usr/share/keyrings/githubcli-archive-keyring.gpg \
    && chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg \
    && echo "deb [arch=${ARCH} signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
    && apt-get update \
    && apt-get install -y gh git-lfs \
    # Install protoc
    && case $ARCH in \
        amd64) PROTOC_ARCH="x86_64" ;; \
        arm64) PROTOC_ARCH="aarch_64" ;; \
        *) echo "Unsupported architecture: $ARCH" && exit 1 ;; \
    esac \
    && curl -fsSL "https://github.com/protocolbuffers/protobuf/releases/download/v23.4/protoc-23.4-linux-${PROTOC_ARCH}.zip" -o /tmp/protoc.zip \
    && unzip /tmp/protoc.zip -d /usr/local \
    && rm /tmp/protoc.zip \
    # Install Docker for minikube
    && curl -fsSL https://get.docker.com | sh \
    && usermod -aG docker codespace \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create docker startup script
RUN echo '#!/bin/bash\n\
# Start docker daemon in background\n\
sudo dockerd > /dev/null 2>&1 &\n\
# Wait for docker to be ready\n\
while ! docker info > /dev/null 2>&1; do sleep 1; done\n\
# Execute the command\n\
exec "$@"' > /usr/local/bin/docker-init.sh \
    && chmod +x /usr/local/bin/docker-init.sh

# Install Atlas for MySQL migration
RUN curl -sSf https://atlasgo.sh | sh -s -- -y

# Switch to codespace user
USER codespace

# Mount for docker-in-docker
VOLUME [ "/var/lib/docker" ]

# Fire Docker init script
ENTRYPOINT [ "/usr/local/bin/docker-init.sh" ]
CMD [ "sleep", "infinity" ]

# [Optional] Install debugger for development of Codespaces - Not in resulting image by default
ARG DeveloperBuild
RUN if [ -z $DeveloperBuild ]; then \
    echo "not including debugger" ; \
    else \
    curl -sSL https://aka.ms/getvsdbgsh | bash /dev/stdin -v latest -l /vsdbg ; \
    fi
