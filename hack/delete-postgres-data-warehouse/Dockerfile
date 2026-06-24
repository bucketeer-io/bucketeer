FROM golang:1.25 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o delete-postgres-data-warehouse .

FROM gcr.io/distroless/base

COPY --from=builder /app/delete-postgres-data-warehouse /usr/local/bin/

ENTRYPOINT ["delete-postgres-data-warehouse"]
