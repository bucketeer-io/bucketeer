#!/bin/bash

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
    
    echo "DEBUG: Looking for image: $image"
    
    # Try new naming first
    local new_file="$(get_cache_filename "$image")"
    echo "DEBUG: Trying new file: $new_file"
    if docker exec "$temp_container" test -f "/cache/$new_file" 2>/dev/null; then
        echo "DEBUG: Found new file!"
        echo "$new_file"
        return 0
    fi
    
    # Try legacy naming
    local legacy_file="$(get_legacy_cache_filename "$image")"
    echo "DEBUG: Trying legacy file: $legacy_file"
    if docker exec "$temp_container" test -f "/cache/$legacy_file" 2>/dev/null; then
        echo "DEBUG: Found legacy file!"
        echo "$legacy_file"
        return 0
    fi
    
    echo "DEBUG: No file found!"
    return 1
}

# Test just one image
test_single_image() {
    print_status "Testing cache file detection for one image..."
    
    local temp_container="debug-cache-$$"
    if ! docker run -d --name "$temp_container" \
        -v bucketeer-docker-images:/cache \
        alpine:latest sleep 120 >/dev/null 2>&1; then
        print_error "Cannot access Docker image cache volume"
        return 1
    fi
    
    local image="ghcr.io/bucketeer-io/bigquery-emulator:latest"
    print_status "Testing image: $image"
    
    echo "=== Files in cache ==="
    docker exec "$temp_container" ls -la /cache/ | grep bigquery
    echo "======================"
    
    if cache_file=$(find_cache_file "$image" "$temp_container"); then
        print_success "Found cache file: $cache_file"
    else
        print_warning "Cache file not found"
    fi
    
    docker stop "$temp_container" >/dev/null 2>&1
    docker rm "$temp_container" >/dev/null 2>&1
}

test_single_image 