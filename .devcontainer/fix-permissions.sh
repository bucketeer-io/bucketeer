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

# USER_NAME, HOME_DIR and USER_WRITABLE_DIRS come from here.
source "$(dirname "${BASH_SOURCE[0]}")/paths.sh"

# Same output format as setup.sh, which also runs this script on attach.
BLUE='\033[0;34m'
GREEN='\033[0;32m'
NC='\033[0m'
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}
print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

fixed=0
for dir in "${USER_WRITABLE_DIRS[@]}"; do
    # mkdir -p first: `go` cannot create /go/pkg/sumdb itself because /go/pkg
    # is root-owned.
    sudo mkdir -p "$dir"
    # Skip the recursive chown when nothing in the tree is foreign-owned --
    # these trees hold tens of thousands of files and this runs on every start.
    # Testing "$dir" itself is not enough: the drift that actually happens is a
    # root-owned file appearing *inside* a tree whose root still looks correct,
    # which is exactly what the setup.sh re-attach call is meant to repair.
    # `-print -quit` stops at the first offender, and unlike chown -R this walk
    # only reads metadata, so a healthy cache is never rewritten.
    if [ -n "$(sudo find "$dir" ! -user "$USER_NAME" -print -quit)" ]; then
        print_status "Fixing ownership of $dir..."
        sudo chown -R "$USER_NAME:$USER_NAME" "$dir"
        fixed=$((fixed + 1))
    fi
done

# minikube refuses to start if its cache is not user-writable.
chmod -R u+wrx "$MINIKUBE_DIR" 2> /dev/null || true

if [ "$fixed" -eq 0 ]; then
    print_success "Cache permissions already correct"
else
    print_success "Cache permissions fixed"
fi
