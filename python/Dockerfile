FROM index.docker.io/library/python:3.7.6-buster as builder

RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.3.1/grpc_health_probe-linux-amd64

FROM index.docker.io/library/python:3.7.6-slim-buster as runner

RUN apt update \
  && apt install -y build-essential python3-dev

COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe
RUN chmod +x /bin/grpc_health_probe

WORKDIR /opt/app

COPY requirements.txt /opt/app

RUN pip install --upgrade pip

RUN pip install -r requirements.txt

COPY src /opt/app

ENV PYTHONPATH "${PYTHONPATH}:/opt/app"

CMD [ "python", "cmd/calculator/main.py"]
