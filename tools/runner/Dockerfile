FROM ubuntu:18.04

ENV GO_VERSION go1.17.2
ENV PYTHON2_VERSION 2.7.17
ENV PYTHON3_VERSION 3.7.6

ENV GOLANGCILINT_VERSION v1.39.0
ENV MOCKGEN_VERSION v1.3.1

ENV PROTO_VERSION 3.18.1
ENV PROTOLOCK_VERSION v0.15.0
ENV GOPROTOGEN_VERSION v1.5.2
ENV GOOGLEAPIS_VERSION 83e756a66b80b072bd234abcfe89edf459090974

ENV BAZELISK_VERSION v1.3.0
ENV BUILDIFIER_VERSION 3.3.0

ENV GCLOUD_VERSION 366.0.0
ENV DOCKER_VERSION 5:18.09.3~3-0~ubuntu-bionic

RUN apt update && apt -qq -y install \
  pkg-config \
  zip \
  g++ \
  zlib1g-dev \
  wget \
  git \
  curl \
  jq \
  gcc \
  make \
  whois \
  clang-format-9 \
  build-essential \
  libffi-dev \
  libssl-dev \
  zlib1g-dev \
  liblzma-dev \
  libbz2-dev  \
  libreadline-dev \
  libsqlite3-dev
RUN mv /usr/bin/clang-format-9 /usr/bin/clang-format

RUN apt install -y -qq lsb-core software-properties-common
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
RUN add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
RUN apt update
RUN apt install -y docker-ce=${DOCKER_VERSION}
RUN curl -sSL https://get.docker.com/ | sh

RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTO_VERSION}/protoc-${PROTO_VERSION}-linux-x86_64.zip
RUN unzip protoc-${PROTO_VERSION}-linux-x86_64.zip -d protoc3
RUN mv protoc3/bin/* /usr/local/bin/
RUN mv protoc3/include/* /usr/local/include/

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN wget -P /tmp https://dl.google.com/go/$GO_VERSION.linux-amd64.tar.gz \
  && tar -C /usr/local -xzf /tmp/$GO_VERSION.linux-amd64.tar.gz \
  && rm /tmp/$GO_VERSION.linux-amd64.tar.gz \
  && mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN git clone https://github.com/googleapis/googleapis.git $GOPATH/src/github.com/googleapis/googleapis
RUN cd $GOPATH/src/github.com/googleapis/googleapis && \
  git checkout ${GOOGLEAPIS_VERSION}

RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCILINT_VERSION}
RUN go install go.uber.org/mock/mockgen@${MOCKGEN_VERSION}
RUN go install github.com/golang/protobuf/protoc-gen-go@${GOPROTOGEN_VERSION}
RUN go install github.com/nilslice/protolock/...@${PROTOLOCK_VERSION}
RUN go install github.com/bazelbuild/buildtools/buildifier@${BUILDIFIER_VERSION}
RUN go install github.com/bazelbuild/bazelisk@${BAZELISK_VERSION}

ENV PYENV_ROOT /root/.pyenv
ENV PATH $PYENV_ROOT/bin:$PATH
RUN git clone https://github.com/pyenv/pyenv.git ~/.pyenv
RUN pyenv init -
RUN pyenv install -f ${PYTHON2_VERSION}
RUN pyenv install -f ${PYTHON3_VERSION}
RUN ln -sf /root/.pyenv/versions/${PYTHON2_VERSION}/bin/python /usr/bin/python
RUN ln -sf /root/.pyenv/versions/${PYTHON3_VERSION}/bin/python /usr/bin/python3
RUN ln -sf /root/.pyenv/versions/${PYTHON2_VERSION}/bin/pip /usr/bin/pip
RUN ln -sf /root/.pyenv/versions/${PYTHON3_VERSION}/bin/pip /usr/bin/pip3

RUN rm /usr/bin/lsb_release
RUN pip3 install yq pyyaml requests
RUN curl -sSL https://raw.githubusercontent.com/python-poetry/poetry/master/install-poetry.py | python3 -
ENV PATH /root/.local/bin:$PATH

ENV PATH /google-cloud-sdk/bin:$PATH
RUN curl -Lo /tmp/google-cloud-sdk.tar.gz https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.tar.gz && \
  tar -C / -xzf /tmp/google-cloud-sdk.tar.gz && \
  rm /tmp/google-cloud-sdk.tar.gz && \
  /google-cloud-sdk/install.sh --usage-reporting=true && \
  gcloud --quiet components update --version=$GCLOUD_VERSION

RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
RUN echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null
RUN apt update -y && apt install -y gh
