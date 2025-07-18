#!/bin/bash
# Cache management helper script for Bucketeer dev container

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

# Docker volumes used for caching
CACHE_VOLUMES=(
    "bucketeer-go-mod-cache"
    "bucketeer-go-tools"
    "bucketeer-dashboard-node-modules"
    "bucketeer-web-v2-node-modules"
    "bucketeer-eval-ts-node-modules"
    "bucketeer-yarn-cache"
    "bucketeer-npm-cache"
    "bucketeer-docker-images"
)

# Helper function to generate consistent cache file names
get_cache_filename() {
    local image="$1"
    echo "$(echo "$image" | tr '/' '_' | tr ':' '_').tar"
}

# Helper function for legacy cache file names (from the old save logic)
get_legacy_cache_filename() {
    local image="$1"
    local cache_file="${image//\//_}_${image##*:}.tar"
    cache_file="${cache_file//:/}_"
    echo "$cache_file"
}

# Helper function to find cache file (tries both new and legacy naming)
find_cache_file() {
    local image="$1"
    local temp_container="$2"
    
    # Try new naming first
    local new_file="$(get_cache_filename "$image")"
    if docker exec "$temp_container" test -f "/cache/$new_file" 2>/dev/null; then
        echo "$new_file"
        return 0
    fi
    
    # Try legacy naming
    local legacy_file="$(get_legacy_cache_filename "$image")"
    if docker exec "$temp_container" test -f "/cache/$legacy_file" 2>/dev/null; then
        echo "$legacy_file"
        return 0
    fi
    
    # File not found
    return 1
}

# Emulator images that are cached
EMULATOR_IMAGES=(
    "ghcr.io/bucketeer-io/bigquery-emulator:latest"
    "gcr.io/google.com/cloudsdktool/google-cloud-cli:449.0.0"
    "docker.io/arigaio/atlas:latest"
    "gcr.io/distroless/base:latest"
    "redis:latest"
    "mysql:8.0"
    "hashicorp/vault:latest"
)

show_help() {
    echo "Bucketeer Dev Container Cache Manager"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Volume Commands:"
    echo "  status       Show status of all cache volumes"
    echo "  clear        Remove all cache volumes"
    echo "  clear-go     Remove only Go-related cache volumes"
    echo "  clear-js     Remove only JavaScript/Node.js cache volumes"
    echo "  size         Show disk usage of cache volumes"
    echo ""
    echo "Docker Image Commands:"
    echo "  images           Show cached Docker images status"
    echo "  pull-images      Pre-pull and cache emulator images"
    echo "  load-images      Load cached images into Docker daemon"
    echo "  load-to-minikube Load cached images into minikube (FULLY AUTOMATED!)"
    echo "  clear-images     Remove cached Docker images"
    echo ""
    echo "General:"
    echo "  help         Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 status                    # Check cache status"
    echo "  $0 images                    # Check Docker image cache"
    echo "  $0 pull-images               # Pre-cache emulator images"
    echo "  $0 load-to-minikube          # Load cached images to minikube (for offline use)"
    echo "  $0 clear                     # Clear all caches"
    echo "  $0 clear-js                  # Clear only Node.js caches"
}

check_volume_exists() {
    local volume_name=$1
    docker volume inspect "$volume_name" &>/dev/null
}

show_status() {
    print_status "Checking cache volume status..."
    echo ""
    
    for volume in "${CACHE_VOLUMES[@]}"; do
        if check_volume_exists "$volume"; then
            local size=$(docker system df -v 2>/dev/null | grep "$volume" | awk '{print $3}' || echo "unknown")
            print_success "$volume: exists (size: $size)"
        else
            print_warning "$volume: not found"
        fi
    done
    echo ""
    print_status "Note: Missing volumes will be created automatically when the dev container starts"
}

show_size() {
    print_status "Cache volume disk usage:"
    echo ""
    
    # Show detailed volume information
    docker system df -v 2>/dev/null | head -1
    for volume in "${CACHE_VOLUMES[@]}"; do
        docker system df -v 2>/dev/null | grep "$volume" || echo "$volume: not found"
    done
    
    echo ""
    local total_cache_size
    total_cache_size=$(docker system df 2>/dev/null | grep "Local Volumes" | awk '{print $3}' || echo "unknown")
    print_status "Total Docker volume space used: $total_cache_size"
}

clear_all_volumes() {
    print_warning "This will remove ALL cache volumes and force a complete reinstall on next dev container start."
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Removing all cache volumes..."
        for volume in "${CACHE_VOLUMES[@]}"; do
            if check_volume_exists "$volume"; then
                docker volume rm "$volume" && print_success "Removed $volume"
            else
                print_warning "$volume already removed"
            fi
        done
        print_success "All cache volumes cleared!"
    else
        print_status "Operation cancelled"
    fi
}

clear_go_volumes() {
    local go_volumes=("bucketeer-go-mod-cache" "bucketeer-go-tools")
    print_warning "This will remove Go-related cache volumes."
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Removing Go cache volumes..."
        for volume in "${go_volumes[@]}"; do
            if check_volume_exists "$volume"; then
                docker volume rm "$volume" && print_success "Removed $volume"
            else
                print_warning "$volume already removed"
            fi
        done
        print_success "Go cache volumes cleared!"
    else
        print_status "Operation cancelled"
    fi
}

clear_js_volumes() {
    local js_volumes=("bucketeer-dashboard-node-modules" "bucketeer-web-v2-node-modules" "bucketeer-eval-ts-node-modules" "bucketeer-yarn-cache" "bucketeer-npm-cache")
    print_warning "This will remove JavaScript/Node.js cache volumes."
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_status "Removing JavaScript cache volumes..."
        for volume in "${js_volumes[@]}"; do
            if check_volume_exists "$volume"; then
                docker volume rm "$volume" && print_success "Removed $volume"
            else
                print_warning "$volume already removed"
            fi
        done
        print_success "JavaScript cache volumes cleared!"
    else
        print_status "Operation cancelled"
    fi
}

# Function to check cached Docker images
show_cached_images() {
    print_status "Checking cached Docker images..."
    echo ""
    
    local cached_count=0
    local total_size=0
    
    if ! check_volume_exists "bucketeer-docker-images"; then
        print_warning "Docker image cache volume not found"
        return 1
    fi
    
    # Create a temporary container to access the volume
    local temp_container="bucketeer-cache-check-$$"
    docker run --rm -d --name "$temp_container" -v bucketeer-docker-images:/cache alpine:latest sleep 60 2>/dev/null || {
        print_error "Could not access Docker image cache"
        return 1
    }
    
    for image in "${EMULATOR_IMAGES[@]}"; do
        if cache_file=$(find_cache_file "$image" "$temp_container"); then
            local size=$(docker exec "$temp_container" stat -c%s "/cache/$cache_file" 2>/dev/null | awk '{printf "%.1fMB", $1/1024/1024}')
            print_success "$image: cached ($size)"
            ((cached_count++))
        else
            print_warning "$image: not cached"
        fi
    done
    
    docker stop "$temp_container" >/dev/null 2>&1
    
    echo ""
    print_status "Summary: $cached_count/${#EMULATOR_IMAGES[@]} emulator images cached"
    
    # Show images currently in Docker daemon
    print_status "Images available in Docker daemon:"
    for image in "${EMULATOR_IMAGES[@]}"; do
        if docker image inspect "$image" >/dev/null 2>&1; then
            local size=$(docker image inspect "$image" --format='{{.Size}}' | awk '{printf "%.1fMB", $1/1024/1024}')
            print_success "$image: available in daemon ($size)"
        else
            print_warning "$image: not in daemon"
        fi
    done
}

# Function to pre-pull and cache emulator images
pull_and_cache_images() {
    print_status "Pre-pulling and caching emulator images..."
    
    if ! curl -s --connect-timeout 5 https://google.com >/dev/null 2>&1; then
        print_error "No internet connection available"
        return 1
    fi
    
    local pulled_count=0
    for image in "${EMULATOR_IMAGES[@]}"; do
        print_status "Pulling $image..."
        # Pull image and filter out gcloud auth warnings from stderr
        local pull_output
        local pull_result=0
        pull_output=$(docker pull "$image" 2>&1) || pull_result=$?
        
        # Show output but filter gcloud warnings  
        echo "$pull_output" | grep -v "gcloud.auth.docker-helper" | grep -v "gcloud auth login" | grep -v "gcloud config set account" | grep -v "Please run:" | grep -v "to obtain new credentials" | grep -v "to select an already authenticated account"
        
        # Check if pull was actually successful (image exists)
        if docker image inspect "$image" >/dev/null 2>&1; then
            ((pulled_count++))
            print_success "Successfully pulled $image"
        else
            print_warning "Failed to pull $image"
        fi
    done
    
    print_success "Successfully pulled $pulled_count/${#EMULATOR_IMAGES[@]} images"
    
    # Now cache them
    save_images_to_cache
}

# Function to save Docker images to cache volume
save_images_to_cache() {
    print_status "Saving images to cache..."
    
    # Create a temporary container to access the cache volume
    local temp_container="bucketeer-cache-writer-$$"
    docker run -d --name "$temp_container" \
        -v bucketeer-docker-images:/cache \
        alpine:latest sleep 3600 >/dev/null
    
    local cached_count=0
    for image in "${EMULATOR_IMAGES[@]}"; do
        if docker image inspect "$image" >/dev/null 2>&1; then
            local cache_file="$(get_cache_filename "$image")"
            
            print_status "Caching $image..."
            # Save image to tar and copy to cache volume
            if docker save "$image" | docker exec -i "$temp_container" sh -c "cat > /cache/$cache_file"; then
                print_success "Cached $image"
                ((cached_count++))
            else
                print_warning "Failed to cache $image"
            fi
        else
            print_warning "$image not available to cache"
        fi
    done
    
    docker stop "$temp_container" >/dev/null 2>&1
    docker rm "$temp_container" >/dev/null 2>&1
    
        print_success "Successfully cached $cached_count/${#EMULATOR_IMAGES[@]} images"
}

# Function to load cached images into Docker daemon
load_cached_images() {
    print_status "Loading cached images into Docker daemon..."
    
    # Try to access the volume directly instead of checking if it exists
    # This works better when run from inside devcontainer
    
    # Create a temporary container to access the cache volume
    local temp_container="bucketeer-cache-reader-$$"
    if ! docker run -d --name "$temp_container" \
        -v bucketeer-docker-images:/cache \
        alpine:latest sleep 3600 >/dev/null 2>&1; then
        print_error "Cannot access Docker image cache volume 'bucketeer-docker-images'"
        print_status "This might happen if:"
        print_status "- Volume doesn't exist (run cache-manager.sh pull-images first)"
        print_status "- Docker permissions issue inside devcontainer"
        print_status "- Docker socket not properly mounted"
        return 1
    fi
    
    local loaded_count=0
    for image in "${EMULATOR_IMAGES[@]}"; do
        # Check if cached image exists (try both naming conventions)
        if cache_file=$(find_cache_file "$image" "$temp_container"); then
            # Check if image is already in daemon
            if ! docker image inspect "$image" >/dev/null 2>&1; then
                print_status "Loading $image from cache..."
                if docker exec "$temp_container" cat "/cache/$cache_file" | docker load >/dev/null; then
                    print_success "Loaded $image"
                    ((loaded_count++))
                else
                    print_warning "Failed to load $image"
                fi
            else
                print_success "$image already available in daemon"
                ((loaded_count++))
            fi
        else
            print_warning "$image not found in cache"
        fi
    done
    
    docker stop "$temp_container" >/dev/null 2>&1
    docker rm "$temp_container" >/dev/null 2>&1
    
    if [ $loaded_count -gt 0 ]; then
        print_success "Successfully loaded/verified $loaded_count/${#EMULATOR_IMAGES[@]} images"
        
        # Also load into minikube if it's running
        if minikube status >/dev/null 2>&1; then
            print_status "Loading images into minikube..."
            local minikube_loaded=0
            for image in "${EMULATOR_IMAGES[@]}"; do
                if docker image inspect "$image" >/dev/null 2>&1; then
                    print_status "Loading $image into minikube..."
                    if minikube image load "$image" >/dev/null 2>&1; then
                        ((minikube_loaded++))
                    fi
                fi
            done
            print_success "Loaded $minikube_loaded images into minikube"
        fi
    else
        print_warning "No images were loaded"
    fi
}

# Function to load cached images into minikube
load_images_to_minikube() {
    print_status "Loading cached images into minikube..."
    
    # Check if minikube is running
    if ! minikube status >/dev/null 2>&1; then
        print_error "Minikube is not running. Start minikube first with: minikube start"
        return 1
    fi
    
    local loaded_count=0
    local total_size=0
    
    # Create temporary container for cache access (reused for all images)
    local temp_restore_container="temp-restore-$$"
    docker run -d --name "$temp_restore_container" -v bucketeer-docker-images:/cache alpine:latest sleep 300 >/dev/null 2>&1
    
    for image in "${EMULATOR_IMAGES[@]}"; do
        # Check if image is available locally
        if ! docker image inspect "$image" >/dev/null 2>&1; then
            print_status "$image not in Docker daemon, trying to restore from cache..."
            
            # Try to load from cache volume
            if cache_file=$(find_cache_file "$image" "$temp_restore_container"); then
                print_status "Restoring $image from cache..."
                
                if docker exec "$temp_restore_container" cat "/cache/$cache_file" | docker load >/dev/null 2>&1; then
                    print_success "Restored $image from cache"
                else
                    print_warning "Failed to restore $image from cache"
                    continue
                fi
            else
                print_warning "$image not available locally or in cache. Run 'pull-images' first."
                continue
            fi
        fi
        
        # Now load into minikube
        if docker image inspect "$image" >/dev/null 2>&1; then
            print_status "Loading $image into minikube..."
            local size=$(docker image inspect "$image" --format='{{.Size}}' | awk '{printf "%.1fMB", $1/1024/1024}')
            
            if minikube image load "$image" >/dev/null 2>&1; then
                print_success "Loaded $image ($size)"
                ((loaded_count++))
                total_size=$(echo "$total_size + $(docker image inspect "$image" --format='{{.Size}}')" | bc 2>/dev/null || echo "$total_size")
            else
                print_warning "Failed to load $image into minikube"
            fi
        fi
    done
    
    # Clean up temporary container
    docker stop "$temp_restore_container" >/dev/null 2>&1
    docker rm "$temp_restore_container" >/dev/null 2>&1
    
    if [ $loaded_count -gt 0 ]; then
        local total_mb=$(echo "scale=1; $total_size/1024/1024" | bc 2>/dev/null || echo "unknown")
        print_success "Successfully loaded $loaded_count/${#EMULATOR_IMAGES[@]} images into minikube (${total_mb}MB total)"
        print_status "Minikube can now use cached images - no internet downloads needed!"
    else
        print_warning "No images were loaded into minikube"
    fi
}

# Function to clear cached Docker images
clear_cached_images() {
    print_warning "This will remove all cached Docker images."
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if check_volume_exists "bucketeer-docker-images"; then
            docker volume rm "bucketeer-docker-images" && print_success "Removed Docker image cache"
        else
            print_warning "Docker image cache volume already removed"
        fi
        print_success "Docker image cache cleared!"
    else
        print_status "Operation cancelled"
    fi
}

# Main command handling
case "${1:-help}" in
    "status")
        show_status
        ;;
    "clear")
        clear_all_volumes
        ;;
    "clear-go")
        clear_go_volumes
        ;;
    "clear-js")
        clear_js_volumes
        ;;
    "size")
        show_size
        ;;
    "images")
        show_cached_images
        ;;
    "pull-images")
        pull_and_cache_images
        ;;
    "load-images")
        load_cached_images
        ;;
    "load-to-minikube")
        load_images_to_minikube
        ;;
    "clear-images")
        clear_cached_images
        ;;
    "help"|"--help"|"-h")
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac 