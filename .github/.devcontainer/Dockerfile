FROM ubuntu:jammy

COPY first-run-notice.txt /tmp/scripts/

ENV LANG="C.UTF-8" \
    DEBIAN_FRONTEND=noninteractive \
    SHELL=/bin/bash \
    DOCKER_BUILDKIT=1

# Restore man command
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && yes | unminimize 2>&1

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

# Install Atlas for MySQL migration
RUN curl -sSf https://atlasgo.sh | sh -s -- -y

# Mount for docker-in-docker
VOLUME [ "/var/lib/docker" ]

# Fire Docker/Moby script if needed
ENTRYPOINT [ "/usr/local/share/docker-init.sh", "/usr/local/share/ssh-init.sh"]
CMD [ "sleep", "infinity" ]

# [Optional] Install debugger for development of Codespaces - Not in resulting image by default
ARG DeveloperBuild
RUN if [ -z $DeveloperBuild ]; then \
    echo "not including debugger" ; \
    else \
    curl -sSL https://aka.ms/getvsdbgsh | bash /dev/stdin -v latest -l /vsdbg ; \
    fi
