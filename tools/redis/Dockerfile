# Use a specific version of Alpine for the runtime stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY bin/redis-data-copy .

# Command to run the executable
ENTRYPOINT ["./redis-data-copy"]
