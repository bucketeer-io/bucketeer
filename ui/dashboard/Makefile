# Makefile for the dashboard UI

# Install dependencies
.PHONY: install
install:
	yarn install

# Run the development server
.PHONY: dev
dev:
	yarn start

# Build for production
.PHONY: build
build:
	VITE_RELEASE_CHANNEL=prod \
	yarn build
