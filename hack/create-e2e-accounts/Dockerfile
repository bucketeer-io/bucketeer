FROM golang:1.25 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o create-e2e-accounts .

FROM gcr.io/distroless/base

COPY --from=builder /app/create-e2e-accounts /usr/local/bin/

ENTRYPOINT ["create-e2e-accounts"]
