#!/bin/bash

DIR=$1
BINARY=$2

cat <<EOF
FROM gcr.io/distroless/base

COPY $DIR/$BINARY /usr/local/bin/

ENTRYPOINT ["$BINARY"]
EOF
