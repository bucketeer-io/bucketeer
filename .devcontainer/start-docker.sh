#!/bin/bash
# Starts the Docker daemon for this dev container.
#
# This used to live at the end of setup.sh (postAttachCommand), which fires
# solely when an editor client attaches -- `gh codespace ssh` sessions ended up
# with no Docker daemon at all, so minikube and every `make` target that talks
# to Docker failed with "Cannot connect to the Docker daemon".
#
# Kept separate from fix-permissions.sh on purpose: that script only ever chowns
# paths, this one only ever deals with the Docker daemon.
set -e

if docker info > /dev/null 2>&1; then
    echo "✅ Docker daemon already running"
    exit 0
fi

echo "🐳 Starting Docker daemon..."
# `service` keeps the daemon log at /var/log/docker.log. The previous
# `nohup sudo dockerd > /tmp/dockerd.log` left no trace once /tmp was cleared.
sudo service docker start > /dev/null

for _ in $(seq 1 60); do
    if docker info > /dev/null 2>&1; then
        echo "✅ Docker daemon started successfully"
        exit 0
    fi
    sleep 1
done

echo "❌ Docker daemon did not become ready; see /var/log/docker.log" >&2
exit 1
