#!/bin/bash
set -e

echo "Starting post-attach setup..."

# Clean up Docker
echo "Cleaning up Docker images and builders..."
docker image prune -f || echo "Docker image prune failed, continuing..."
docker builder prune -f || echo "Docker builder prune failed, continuing..."

# Install dependencies
echo "Installing local dependencies..."
make local-deps || echo "make local-deps failed, continuing..."

# Update repositories
echo "Updating repositories..."
make update-repos || echo "make update-repos failed, continuing..."

# Install UI dependencies
echo "Installing UI dependencies..."
cd ui/dashboard && yarn || echo "yarn install failed, continuing..."

echo "Post-attach setup completed!"