#!/bin/bash
# Starts the Docker daemon for this dev container, on a firewall it can use.
#
# This used to live at the end of setup.sh (postAttachCommand), which fires
# solely when an editor client attaches -- `gh codespace ssh` sessions ended up
# with no Docker daemon at all, so minikube and every `make` target that talks
# to Docker failed with "Cannot connect to the Docker daemon".
#
# Kept separate from fix-permissions.sh on purpose: that script only ever chowns
# paths, this one only ever deals with the Docker daemon.
set -e

# --- Align the iptables backend with the outer daemon ------------------------
# Codespaces runs this container in the *host* network namespace, so the
# Codespaces VM's own Docker daemon and the daemon we start here write firewall
# rules into the same namespace. That outer daemon is an older release using
# iptables-legacy: it sets `FORWARD policy DROP` and only opens holes for its
# own docker0. If we then use the nft backend, our bridges (minikube's included)
# never get an ACCEPT rule in the table that holds the DROP, so every forwarded
# packet dies there -- image pulls from inside the cluster time out while the
# host itself reaches the registry fine, and nothing in the minikube logs points
# at the host firewall.
#
# Matching the backend makes our daemon add its own ACCEPT rules right next to
# the DROP policy, which resolves it. This is the same fix the upstream
# docker-in-docker devcontainer feature applies, for the same reason; we cannot
# use that feature here because Docker is installed directly in the image.
#
# `ip_tables` being loaded is the signal that something in this namespace
# already uses the legacy backend. The /sys/module check covers kernels that
# build it in rather than as a module. On a host that uses nft, none of this
# matches and the block is a no-op.
if type iptables-legacy > /dev/null 2>&1 \
    && { grep -qE '^ip_tables\b' /proc/modules || [ -d /sys/module/ip_tables ]; } \
    && update-alternatives --list iptables 2>/dev/null | grep -q '/usr/sbin/iptables-legacy'; then

    current=$(update-alternatives --query iptables 2>/dev/null | awk '/^Value:/ {print $2}')
    if [ "$current" != "/usr/sbin/iptables-legacy" ]; then
        echo "🔧 Switching iptables backend to legacy to match the host daemon..."
        sudo update-alternatives --set iptables /usr/sbin/iptables-legacy > /dev/null
        sudo update-alternatives --set ip6tables /usr/sbin/ip6tables-legacy > /dev/null 2>&1 || true
    fi
fi

# --- Start dockerd -----------------------------------------------------------
# Order matters: the switch above only affects a daemon started afterwards. A
# running dockerd has already written its rules to whichever table it picked.
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
