#!/bin/bash
# Starts the SSH daemon so `gh codespace ssh` (and other SSH clients) can connect.
#
# The Dockerfile's own ENTRYPOINT (docker-init.sh) also starts sshd, but Codespaces/
# devcontainer tooling overrides the container's entrypoint at container-start time,
# so that logic never runs there. This script is wired to postStartCommand so sshd
# comes up at every container start — with or without an editor attach.
set -e

if pgrep -x sshd > /dev/null 2>&1; then
    echo "✅ SSH daemon already running"
    exit 0
fi

echo "🔑 Starting SSH daemon..."
# Regenerate host keys if missing (never baked into the shared image)
sudo ssh-keygen -A > /dev/null 2>&1
# Privilege-separation dir; /run can be tmpfs, wiping the build-time mkdir
sudo mkdir -p /run/sshd
sudo /usr/sbin/sshd
echo "✅ SSH daemon started successfully"
