# Shared location definitions for the dev container scripts.
#
# Sourced (not executed) by setup.sh and fix-permissions.sh so both agree on
# where projects and caches live. Only defines variables -- no side effects.
#
# To add a new Node.js project at "ui/new-app":
# 1. Add "ui/new-app" to NODE_PROJECTS
# 2. Add "New App" to NODE_PROJECTS_NAMES
# 3. Add its node_modules volume mount in devcontainer.json
# 4. Update cache-manager.sh volume arrays
# setup.sh installs its dependencies and fix-permissions.sh keeps its
# node_modules volume writable from then on.

# The image (.github/.devcontainer/Dockerfile) always creates and runs as this
# user, and devcontainer.json mounts the caches under /home/codespace.
USER_NAME=codespace
HOME_DIR=/home/$USER_NAME

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Node.js projects (directories containing package.json), relative to REPO_ROOT
declare -a NODE_PROJECTS=(
    "ui/dashboard"
    "evaluation/typescript"
)

# Display names for Node.js projects (must match array order)
declare -a NODE_PROJECTS_NAMES=(
    "Dashboard"
    "Evaluation TypeScript"
)

# Go projects (directories containing go.mod), relative to REPO_ROOT
# Note: "." represents the root directory
declare -a GO_PROJECTS=(
    "."
)

# Display names for Go projects (must match array order)
declare -a GO_PROJECTS_NAMES=(
    "Main Go"
)

GO_MOD_CACHE_DIR=/go/pkg/mod
# Not a mounted volume, but /go/pkg is root-owned so `go` cannot create it.
GO_SUMDB_DIR=/go/pkg/sumdb
GO_TOOLS_DIR=$HOME_DIR/go-tools
YARN_DIR=$HOME_DIR/.yarn
MINIKUBE_DIR=$HOME_DIR/.minikube

# Directories the container user writes into but that can start out
# root-owned (fresh named volumes, or under the root-owned /go/pkg): each must
# exist and be owned by USER_NAME. /var/lib/docker is deliberately absent --
# dockerd owns it as root.
declare -a USER_WRITABLE_DIRS=(
    "$GO_MOD_CACHE_DIR"
    "$GO_SUMDB_DIR"
    "$GO_TOOLS_DIR"
    "$YARN_DIR"
    "$MINIKUBE_DIR"
)
for project in "${NODE_PROJECTS[@]}"; do
    USER_WRITABLE_DIRS+=("$REPO_ROOT/$project/node_modules")
done
unset project
