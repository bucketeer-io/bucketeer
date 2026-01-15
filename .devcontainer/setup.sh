#!/bin/bash
set -e

echo "🚀 Starting post-attach setup with intelligent caching..."

# ===== PROJECT CONFIGURATION =====
# Configure which projects to cache dependencies for
# To add/remove projects, just update these arrays
#
# Example: To add a new Node.js project at "ui/new-app":
# 1. Add "ui/new-app" to NODE_PROJECTS array
# 2. Add "New App" to NODE_PROJECTS_NAMES array
# 3. Add volume mount in devcontainer.json
# 4. Update cache-manager.sh volume arrays
# The setup script will handle the rest automatically!

# Node.js projects (directories containing package.json)
declare -a NODE_PROJECTS=(
    "ui/dashboard"
    "evaluation/typescript"
)

# Display names for Node.js projects (must match array order)
declare -a NODE_PROJECTS_NAMES=(
    "Dashboard"
    "Evaluation TypeScript"
)

# Go projects (directories containing go.mod)
# Note: "." represents the root directory
declare -a GO_PROJECTS=(
    "."
)

# Display names for Go projects (must match array order)
declare -a GO_PROJECTS_NAMES=(
    "Main Go"
)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper function for colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Add go-tools to PATH if not already there
export PATH="/home/codespace/go-tools/bin:$PATH"
export GOBIN="/home/codespace/go-tools/bin"

# Make PATH change persistent by adding to shell profile if not already there
if ! grep -q "/home/codespace/go-tools/bin" ~/.bashrc; then
    echo 'export PATH="/home/codespace/go-tools/bin:$PATH"' >> ~/.bashrc
    echo 'export GOBIN="/home/codespace/go-tools/bin"' >> ~/.bashrc
    print_status "Added Go tools to PATH in ~/.bashrc for persistence"
fi

# Add kubectl autocomplete if kubectl is available
# Note: bash-completion package is pre-installed in the devcontainer image
if command -v kubectl &> /dev/null; then
    if ! grep -q "kubectl completion bash" ~/.bashrc; then
        echo '' >> ~/.bashrc
        echo '# kubectl autocomplete' >> ~/.bashrc
        echo '# Ensure bash-completion is loaded first (required for kubectl completion)' >> ~/.bashrc
        echo 'if [ -f /usr/share/bash-completion/bash_completion ]; then' >> ~/.bashrc
        echo '    source /usr/share/bash-completion/bash_completion' >> ~/.bashrc
        echo 'elif [ -f /etc/bash_completion ]; then' >> ~/.bashrc
        echo '    source /etc/bash_completion' >> ~/.bashrc
        echo 'fi' >> ~/.bashrc
        echo 'source <(kubectl completion bash)' >> ~/.bashrc
        echo 'alias k=kubectl' >> ~/.bashrc
        echo 'complete -o default -F __start_kubectl k' >> ~/.bashrc
        print_status "Added kubectl autocomplete to ~/.bashrc"
    fi
fi

# Function to fix permissions on cache directories
fix_cache_permissions() {
    print_status "Ensuring cache directories have correct permissions..."

    # Check if codespace user exists
    if id "codespace" &>/dev/null; then
        USER_NAME="codespace"
    else
        USER_NAME=$(whoami)
        print_status "codespace user not found, using current user: $USER_NAME"
    fi

    # Fix Go modules cache permissions
    sudo mkdir -p /go/pkg/mod /go/pkg/sumdb
    sudo chown -R $USER_NAME:$USER_NAME /go/pkg/mod /go/pkg/sumdb

    # Fix Go tools directory permissions
    if [ -d "/home/$USER_NAME/go-tools" ]; then
        sudo chown -R $USER_NAME:$USER_NAME /home/$USER_NAME/go-tools
    fi

    # Fix Yarn cache permissions
    if [ -d "/home/$USER_NAME/.yarn" ]; then
        sudo chown -R $USER_NAME:$USER_NAME /home/$USER_NAME/.yarn
    fi

    # Fix node_modules permissions
    for project in "${NODE_PROJECTS[@]}"; do
        if [ -d "$project/node_modules" ]; then
            sudo chown -R $USER_NAME:$USER_NAME "$project/node_modules"
        fi
    done

    # Fix minikube cache permissions
    if [ -d "/home/$USER_NAME/.minikube" ]; then
        sudo chown -R $USER_NAME:$USER_NAME /home/$USER_NAME/.minikube
        chmod -R u+wrx /home/$USER_NAME/.minikube
        print_status "Fixed minikube cache permissions"
    fi

    print_success "Cache permissions fixed"
}

# Function to check if Go tools are installed
check_go_tools() {
    local tools=("goimports" "golangci-lint" "mockgen" "protoc-gen-go" "protoc-gen-openapiv2" "protoc-gen-grpc-gateway" "protolock" "yq")
    local missing_tools=()

    # Ensure PATH includes go-tools directory for checking
    export PATH="/home/codespace/go-tools/bin:$PATH"

    # Debug: Show where we're looking for tools
    if [ "${DEBUG_SETUP:-}" = "true" ]; then
        print_status "DEBUG: Checking for Go tools in PATH and /home/codespace/go-tools/bin/"
        print_status "DEBUG: Current PATH: $PATH"
        if [ -d "/home/codespace/go-tools/bin" ]; then
            print_status "DEBUG: Contents of /home/codespace/go-tools/bin: $(ls -la /home/codespace/go-tools/bin/ 2>/dev/null || echo 'empty')"
        fi
    fi

    for tool in "${tools[@]}"; do
        # Check both PATH and direct file existence
        if command -v "$tool" &> /dev/null || [ -f "/home/codespace/go-tools/bin/$tool" ]; then
            # Tool is available
            if [ "${DEBUG_SETUP:-}" = "true" ]; then
                print_success "DEBUG: Found $tool"
            fi
            continue
        else
            if [ "${DEBUG_SETUP:-}" = "true" ]; then
                print_warning "DEBUG: Missing $tool"
            fi
            missing_tools+=("$tool")
        fi
    done

    if [ ${#missing_tools[@]} -eq 0 ]; then
        print_success "All Go tools are already installed"
        return 0
    else
        print_warning "Missing Go tools: ${missing_tools[*]}"
        return 1
    fi
}

# Function to install Go tools
install_go_tools() {
    print_status "Installing Go development tools..."

    # Check if codespace user exists
    if id "codespace" &>/dev/null; then
        USER_NAME="codespace"
    else
        USER_NAME=$(whoami)
    fi

    # Ensure go-tools directory exists and has correct permissions
    sudo mkdir -p /home/$USER_NAME/go-tools/bin
    sudo chown -R $USER_NAME:$USER_NAME /home/$USER_NAME/go-tools

    cd /home/$USER_NAME/go-tools
    if [ ! -e go.mod ]; then go mod init go-tools; fi

    # Set additional Go environment variables for better module handling
    export GOSUMDB="sum.golang.org"
    export GOPROXY="https://proxy.golang.org,direct"
    export GONOPROXY=""
    export GONOSUMDB=""
    export GOPRIVATE=""

    # Install tools with better error handling
    print_status "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@v0.40.0 || {
        print_warning "goimports installation failed, retrying with proxy settings..."
        GOPROXY=direct go install golang.org/x/tools/cmd/goimports@latest
    }

    print_status "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.2

    print_status "Installing mockgen..."
    go install go.uber.org/mock/mockgen@v0.4.0

    print_status "Installing protoc-gen-go..."
    go install github.com/golang/protobuf/protoc-gen-go@v1.5.2

    print_status "Installing protoc-gen-openapiv2..."
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0

    print_status "Installing protoc-gen-grpc-gateway..."
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0

    print_status "Installing protolock..."
    go install github.com/nilslice/protolock/...@v0.15.0

    print_status "Installing yq..."
    go install github.com/mikefarah/yq/v4@v4.28.2

    cd /workspaces/bucketeer
    print_success "Go tools installed successfully"
}

# Function to check if Go vendor directory is up to date
check_go_vendor() {
    if [ ! -d "vendor" ]; then
        print_warning "vendor/ directory missing"
        return 1
    fi

    if [ "go.mod" -nt "vendor" ] || [ "go.sum" -nt "vendor" ]; then
        print_warning "go.mod/go.sum newer than vendor/ directory"
        return 1
    fi

    print_success "Go vendor directory is up to date"
    return 0
}

# Function to update Go dependencies
update_go_deps() {
    print_status "Updating Go dependencies..."
    make update-repos || {
        print_error "make update-repos failed, continuing..."
        return 1
    }
    print_success "Go dependencies updated"
}

# Function to check if Node.js dependencies are up to date
check_node_deps() {
    local dir=$1
    local name=$2

    if [ ! -d "$dir/node_modules" ]; then
        print_warning "$name: node_modules/ directory missing"
        return 1
    fi

    if [ "$dir/yarn.lock" -nt "$dir/node_modules" ]; then
        print_warning "$name: yarn.lock newer than node_modules/"
        return 1
    fi

    print_success "$name: Node.js dependencies are up to date"
    return 0
}

# Function to install Node.js dependencies
install_node_deps() {
    local dir=$1
    local name=$2

    print_status "Installing $name dependencies..."

    # Check if codespace user exists
    if id "codespace" &>/dev/null; then
        USER_NAME="codespace"
    else
        USER_NAME=$(whoami)
    fi

    # Ensure node_modules directory has correct permissions if it exists
    if [ -d "$dir/node_modules" ]; then
        sudo chown -R $USER_NAME:$USER_NAME "$dir/node_modules"
    fi

    cd "$dir"

    # Handle .npmrc file if it requires NPM_TOKEN but token is not set
    local npmrc_backup=""
    if [ -f ".npmrc" ] && grep -q "\${NPM_TOKEN}" .npmrc 2>/dev/null; then
        if [ -z "${NPM_TOKEN:-}" ]; then
            print_status "$name: NPM_TOKEN not set, temporarily disabling .npmrc for installation"
            npmrc_backup=".npmrc.backup"
            mv .npmrc "$npmrc_backup" || true
        fi
    fi

    # Use optimized yarn flags for faster, deterministic installs
    yarn install --frozen-lockfile --silent || {
        print_error "yarn install failed for $name, continuing..."
        # Restore .npmrc if it was backed up
        if [ -n "$npmrc_backup" ] && [ -f "$npmrc_backup" ]; then
            mv "$npmrc_backup" .npmrc || true
        fi
        cd /workspaces/bucketeer
        return 1
    }

    # Update node_modules timestamp to prevent false warnings with Docker volumes
    touch node_modules

    # Restore .npmrc if it was backed up
    if [ -n "$npmrc_backup" ] && [ -f "$npmrc_backup" ]; then
        mv "$npmrc_backup" .npmrc || true
    fi

    cd /workspaces/bucketeer
    print_success "$name dependencies installed"
}

# Function to check disk space and cleanup if needed
cleanup_docker_if_needed() {
    local available_space_gb
    available_space_gb=$(df /var/lib/docker 2>/dev/null | awk 'NR==2 {print int($4/1024/1024)}' || echo "10")

    if [ "$available_space_gb" -lt 5 ]; then
        print_warning "Low disk space (${available_space_gb}GB), cleaning up Docker..."
        docker image prune -f || print_error "Docker image prune failed"
        docker builder prune -f || print_error "Docker builder prune failed"
        print_success "Docker cleanup completed"
    else
        print_success "Sufficient disk space available (${available_space_gb}GB)"
    fi
}

# Main setup logic
main() {
    # Fix any permission issues with mounted cache volumes first
    fix_cache_permissions

    print_status "Checking cache status..."

    # Track what needs to be done
    local need_go_tools=false
    local need_go_deps=false

    # Check Go tools
    if ! check_go_tools; then
        need_go_tools=true
    fi

    # Check Go dependencies
    if ! check_go_vendor; then
        need_go_deps=true
    fi

    # Check Node.js dependencies
    declare -a need_node_deps=()
    for i in "${!NODE_PROJECTS[@]}"; do
        local project="${NODE_PROJECTS[$i]}"
        local name="${NODE_PROJECTS_NAMES[$i]}"
        if ! check_node_deps "$project" "$name"; then
            need_node_deps+=("$i")
        fi
    done

    # Summary of what will be installed
    local tasks_to_run=()
    if [ "$need_go_tools" = true ]; then tasks_to_run+=("Go tools"); fi
    if [ "$need_go_deps" = true ]; then tasks_to_run+=("Go dependencies"); fi
    for i in "${need_node_deps[@]}"; do
        local name="${NODE_PROJECTS_NAMES[$i]}"
        tasks_to_run+=("$name dependencies")
    done

    if [ ${#tasks_to_run[@]} -eq 0 ]; then
        print_success "🎉 All dependencies are cached and up to date! Setup completed in seconds."
        return 0
    fi

    print_status "Cache analysis complete. Will install: ${tasks_to_run[*]}"

    # Cleanup Docker if needed
    cleanup_docker_if_needed

    # Install missing components
    if [ "$need_go_tools" = true ]; then
        install_go_tools
    fi

    if [ "$need_go_deps" = true ]; then
        update_go_deps
    fi

    for i in "${need_node_deps[@]}"; do
        local project="${NODE_PROJECTS[$i]}"
        local name="${NODE_PROJECTS_NAMES[$i]}"
        install_node_deps "$project" "$name"
    done

    print_success "🎉 Post-attach setup completed!"
    print_status "Cache volumes will persist dependencies for next container restart"
}

# Run main function
main "$@"

# Start Docker daemon automatically after main setup is complete
if ! docker info > /dev/null 2>&1; then
    echo "🐳 Starting Docker daemon..."
    nohup sudo dockerd > /tmp/dockerd.log 2>&1 < /dev/null &
    # Wait for Docker to be ready
    while ! docker info > /dev/null 2>&1; do
        sleep 1
    done
    echo "✅ Docker daemon started successfully"
    # Reset cursor position for clean terminal state
    printf "\r"
else
    echo "✅ Docker daemon already running"
fi
