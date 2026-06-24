FROM golang:1.25 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o delete-e2e-data-postgres .

FROM gcr.io/distroless/base

COPY --from=builder /app/delete-e2e-data-postgres /usr/local/bin/

ENTRYPOINT ["delete-e2e-data-postgres"]
