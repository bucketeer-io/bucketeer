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
)

show_help() {
    echo "Bucketeer Dev Container Cache Manager"
    echo
    echo "Usage: $0 [COMMAND]"
    echo
    echo "Commands:"
    echo "  status       Show status of all cache volumes"
    echo "  clear        Remove all cache volumes"
    echo "  clear-go     Remove only Go-related cache volumes"
    echo "  clear-js     Remove only JavaScript/Node.js cache volumes"
    echo "  size         Show disk usage of cache volumes"
    echo "  help         Show this help message"
    echo
    echo "Examples:"
    echo "  $0 status                    # Check cache status"
    echo "  $0 size                      # Check cache sizes"
    echo "  $0 clear                     # Clear all caches"
    echo "  $0 clear-js                  # Clear only Node.js caches"
    echo
    echo "Note: Cache volumes will be recreated automatically when the dev container starts."
}

check_volume_exists() {
    local volume_name=$1
    docker volume inspect "$volume_name" &>/dev/null
}

show_status() {
    print_status "Checking cache volume status..."
    echo
    for volume in "${CACHE_VOLUMES[@]}"; do
        if check_volume_exists "$volume"; then
            local size=$(docker system df -v 2>/dev/null | grep "$volume" | awk '{print $3}' || echo "unknown")
            print_success "$volume: exists (size: $size)"
        else
            print_warning "$volume: not found"
        fi
    done
    echo
    print_status "Note: Missing volumes will be created automatically when the dev container starts"
}

show_size() {
    print_status "Cache volume disk usage:"
    echo
    # Show detailed volume information
    docker system df -v 2>/dev/null | head -1
    for volume in "${CACHE_VOLUMES[@]}"; do
        docker system df -v 2>/dev/null | grep "$volume" || echo "$volume: not found"
    done
    echo
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