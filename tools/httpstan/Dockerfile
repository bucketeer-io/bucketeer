FROM python:3.9-slim

WORKDIR /httpstan

ARG TAG

RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    software-properties-common \
    git \
    && rm -rf /var/lib/apt/lists/*

RUN git clone --branch $TAG --depth=1 https://github.com/stan-dev/httpstan.git .

# Build shared libraries (Stan C++)
RUN make

# Build the httpstan wheel on your system
RUN python3 -m pip install poetry
RUN python3 -m poetry build

# Install the wheel
RUN python3 -m pip install dist/*.whl

EXPOSE 8080

CMD ["python3", "-m", "httpstan", "--host", "0.0.0.0"]