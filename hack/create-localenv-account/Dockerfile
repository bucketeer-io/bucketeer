FROM golang:1.25 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o create-localenv-account .

FROM gcr.io/distroless/base

COPY --from=builder /app/create-localenv-account /usr/local/bin/

ENTRYPOINT ["create-localenv-account"]
