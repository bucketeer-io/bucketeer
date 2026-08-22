#!/bin/bash
# Makes the mounted cache volumes writable by the container user.
#
# devcontainer.json mounts several named Docker volumes (Go module cache, Go
# tools, node_modules, yarn cache, minikube cache). A named volume that is empty
# on first use is created root-owned, so every tool that writes into one fails
# with EACCES until the ownership is fixed:
#
#   yarn install  -> Error: EACCES: permission denied, mkdir '~/.yarn/cache/v6'
#   go build      -> open /go/pkg/sumdb/sum.golang.org/latest: no such file or directory
#
# This used to run only from setup.sh (postAttachCommand), which fires solely
# when an editor attaches -- `gh codespace ssh` sessions hit the failures above.
# It is wired to postStartCommand instead so every container start is covered.
#
# Kept separate from start-docker.sh on purpose: this script only ever chowns
# paths, that one only ever deals with the Docker daemon.
set -e

if id codespace &> /dev/null; then
    USER_NAME=codespace
else
    USER_NAME=$(whoami)
fi
HOME_DIR=$(getent passwd "$USER_NAME" | cut -d: -f6)

# Directories that must exist and be owned by the user. `go` cannot create
# /go/pkg/sumdb itself because /go/pkg is root-owned, so mkdir -p comes first.
OWNED_DIRS=(
    /go/pkg/mod
    /go/pkg/sumdb
    "$HOME_DIR/go-tools"
    "$HOME_DIR/.yarn"
    "$HOME_DIR/.minikube"
    /workspaces/bucketeer/ui/dashboard/node_modules
    /workspaces/bucketeer/evaluation/typescript/node_modules
)

fixed=0
for dir in "${OWNED_DIRS[@]}"; do
    sudo mkdir -p "$dir"
    # Skip the recursive chown when the top-level owner is already correct --
    # these trees hold tens of thousands of files and this runs on every start.
    if [ "$(stat -c '%U' "$dir")" != "$USER_NAME" ]; then
        echo "🔑 Fixing ownership of $dir..."
        sudo chown -R "$USER_NAME:$USER_NAME" "$dir"
        fixed=$((fixed + 1))
    fi
done

# minikube refuses to start if its cache is not user-writable.
chmod -R u+wrx "$HOME_DIR/.minikube" 2> /dev/null || true

if [ "$fixed" -eq 0 ]; then
    echo "✅ Cache volume permissions already correct"
else
    echo "✅ Fixed permissions on $fixed cache volume(s)"
fi
