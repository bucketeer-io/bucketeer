#!/bin/bash
set -e

echo "🚀 Starting post-attach setup with intelligent caching..."

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

# Function to fix permissions on cache directories
fix_cache_permissions() {
    print_status "Ensuring cache directories have correct permissions..."
    
    # Fix Go modules cache permissions
    sudo mkdir -p /go/pkg/mod /go/pkg/sumdb
    sudo chown -R codespace:codespace /go/pkg/mod /go/pkg/sumdb
    
    # Fix Go tools directory permissions
    if [ -d "/home/codespace/go-tools" ]; then
        sudo chown -R codespace:codespace /home/codespace/go-tools
    fi
    
    # Fix Yarn/NPM cache permissions
    if [ -d "/home/codespace/.yarn" ]; then
        sudo chown -R codespace:codespace /home/codespace/.yarn
    fi
    if [ -d "/home/codespace/.npm" ]; then
        sudo chown -R codespace:codespace /home/codespace/.npm
    fi
    
    # Fix node_modules permissions
    for dir in "ui/dashboard/node_modules" "ui/web-v2/node_modules" "evaluation/typescript/node_modules"; do
        if [ -d "$dir" ]; then
            sudo chown -R codespace:codespace "$dir"
        fi
    done
    
    # Fix Docker image cache permissions
    sudo mkdir -p /home/codespace/.docker-images
    sudo chown -R codespace:codespace /home/codespace/.docker-images
    
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
    
    # Ensure go-tools directory exists and has correct permissions
    sudo mkdir -p /home/codespace/go-tools/bin
    sudo chown -R codespace:codespace /home/codespace/go-tools
    
    cd /home/codespace/go-tools
    if [ ! -e go.mod ]; then go mod init go-tools; fi
    
    # Set additional Go environment variables for better module handling
    export GOSUMDB="sum.golang.org"
    export GOPROXY="https://proxy.golang.org,direct"
    export GONOPROXY=""
    export GONOSUMDB=""
    export GOPRIVATE=""
    
    # Install tools with better error handling
    print_status "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest || {
        print_warning "goimports installation failed, retrying with proxy settings..."
        GOPROXY=direct go install golang.org/x/tools/cmd/goimports@latest
    }
    
    print_status "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
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
    
    # Ensure node_modules directory has correct permissions if it exists
    if [ -d "$dir/node_modules" ]; then
        sudo chown -R codespace:codespace "$dir/node_modules"
    fi
    
    cd "$dir"
    
    # Use optimized yarn flags for faster, deterministic installs
    yarn install --frozen-lockfile --silent || {
        print_error "yarn install failed for $name, continuing..."
        cd /workspaces/bucketeer
        return 1
    }
    
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

# List of critical emulator and development images
EMULATOR_IMAGES=(
    "ghcr.io/bucketeer-io/bigquery-emulator:latest"
    "gcr.io/google.com/cloudsdktool/google-cloud-cli:449.0.0"
    "docker.io/arigaio/atlas:latest"
    "gcr.io/distroless/base:latest"
    "redis:latest"
    "mysql:8.0"
    "hashicorp/vault:latest"
)

# Function to check if internet is available
check_internet() {
    curl -s --connect-timeout 5 https://google.com > /dev/null 2>&1
}

# Function to check which emulator images are already cached
check_cached_images() {
    local cached_images=()
    local missing_images=()
    
    for image in "${EMULATOR_IMAGES[@]}"; do
        local cache_file="/home/codespace/.docker-images/$(echo "$image" | tr '/' '_' | tr ':' '_').tar"
        if [ -f "$cache_file" ]; then
            cached_images+=("$image")
        else
            missing_images+=("$image")
        fi
    done
    
    if [ ${#cached_images[@]} -gt 0 ]; then
        print_success "Cached images: ${cached_images[*]}"
    fi
    
    if [ ${#missing_images[@]} -gt 0 ]; then
        print_warning "Missing cached images: ${missing_images[*]}"
        return 1
    fi
    
    return 0
}

# Function to save Docker image to cache
save_image_to_cache() {
    local image=$1
    local cache_file="/home/codespace/.docker-images/$(echo "$image" | tr '/' '_' | tr ':' '_').tar"
    
    if docker image inspect "$image" >/dev/null 2>&1; then
        print_status "Saving $image to cache..."
        docker save "$image" -o "$cache_file" && {
            print_success "Cached $image"
            return 0
        }
    fi
    return 1
}

# Function to load Docker image from cache
load_image_from_cache() {
    local image=$1
    local cache_file="/home/codespace/.docker-images/$(echo "$image" | tr '/' '_' | tr ':' '_').tar"
    
    if [ -f "$cache_file" ]; then
        print_status "Loading $image from cache..."
        docker load -i "$cache_file" && {
            print_success "Loaded $image from cache"
            return 0
        }
    fi
    return 1
}

# Function to pre-pull and cache emulator images
cache_emulator_images() {
    if ! check_internet; then
        print_warning "No internet connection, skipping image pre-pulling"
        return 0
    fi
    
    print_status "Pre-pulling and caching emulator images..."
    local pulled_count=0
    
    for image in "${EMULATOR_IMAGES[@]}"; do
        local cache_file="/home/codespace/.docker-images/$(echo "$image" | tr '/' '_' | tr ':' '_').tar"
        
        # Skip if already cached
        if [ -f "$cache_file" ]; then
            if [ "${DEBUG_SETUP:-}" = "true" ]; then
                print_success "DEBUG: $image already cached"
            fi
            continue
        fi
        
        # Try to pull and cache the image
        print_status "Pulling $image..."
        if docker pull "$image" 2>/dev/null; then
            if save_image_to_cache "$image"; then
                ((pulled_count++))
            fi
        else
            print_warning "Failed to pull $image (might not exist or network issue)"
        fi
    done
    
    if [ $pulled_count -gt 0 ]; then
        print_success "Successfully cached $pulled_count new emulator images"
    else
        print_success "All emulator images already cached"
    fi
}

# Function to restore cached images to Docker daemon
restore_cached_images() {
    print_status "Restoring cached emulator images to Docker daemon..."
    local restored_count=0
    
    for image in "${EMULATOR_IMAGES[@]}"; do
        # Check if image is already in Docker daemon
        if docker image inspect "$image" >/dev/null 2>&1; then
            if [ "${DEBUG_SETUP:-}" = "true" ]; then
                print_success "DEBUG: $image already in Docker daemon"
            fi
            continue
        fi
        
        # Try to load from cache
        if load_image_from_cache "$image"; then
            ((restored_count++))
        fi
    done
    
    if [ $restored_count -gt 0 ]; then
        print_success "Restored $restored_count images from cache"
    else
        print_success "All cached images already available in Docker daemon"
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
    local need_dashboard_deps=false
    local need_web_v2_deps=false
    local need_eval_ts_deps=false
    
    # Check Go tools
    if ! check_go_tools; then
        need_go_tools=true
    fi
    
    # Check Go dependencies
    if ! check_go_vendor; then
        need_go_deps=true
    fi
    
    # Check Node.js dependencies
    if ! check_node_deps "ui/dashboard" "Dashboard"; then
        need_dashboard_deps=true
    fi
    
    if ! check_node_deps "ui/web-v2" "Web-v2"; then
        need_web_v2_deps=true
    fi
    
    if ! check_node_deps "evaluation/typescript" "Evaluation TypeScript"; then
        need_eval_ts_deps=true
    fi
    
    # Summary of what will be installed
    local tasks_to_run=()
    if [ "$need_go_tools" = true ]; then tasks_to_run+=("Go tools"); fi
    if [ "$need_go_deps" = true ]; then tasks_to_run+=("Go dependencies"); fi
    if [ "$need_dashboard_deps" = true ]; then tasks_to_run+=("Dashboard dependencies"); fi
    if [ "$need_web_v2_deps" = true ]; then tasks_to_run+=("Web-v2 dependencies"); fi
    if [ "$need_eval_ts_deps" = true ]; then tasks_to_run+=("Evaluation TypeScript dependencies"); fi
    
    if [ ${#tasks_to_run[@]} -eq 0 ]; then
        print_success "🎉 All dependencies are cached and up to date! Setup completed in seconds."
        return 0
    fi
    
    print_status "Cache analysis complete. Will install: ${tasks_to_run[*]}"
    
    # Cleanup Docker if needed
    cleanup_docker_if_needed
    
    # Restore cached Docker images first (for offline scenarios)
    restore_cached_images
    
    # Install missing components
    if [ "$need_go_tools" = true ]; then
        install_go_tools
    fi
    
    if [ "$need_go_deps" = true ]; then
        update_go_deps
    fi
    
    if [ "$need_dashboard_deps" = true ]; then
        install_node_deps "ui/dashboard" "Dashboard"
    fi
    
    if [ "$need_web_v2_deps" = true ]; then
        install_node_deps "ui/web-v2" "Web-v2"
    fi
    
    if [ "$need_eval_ts_deps" = true ]; then
        install_node_deps "evaluation/typescript" "Evaluation TypeScript"
    fi
    
    # Pre-pull and cache emulator images if internet is available
    cache_emulator_images
    
    print_success "🎉 Post-attach setup completed!"
    print_status "Cache volumes will persist dependencies for next container restart"
    print_status "Emulator images cached for offline minikube usage"
}

# Run main function
main "$@"