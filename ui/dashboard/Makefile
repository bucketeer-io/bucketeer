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
	VITE_DEMO_SIGN_IN_ENABLED=true \
	VITE_DEMO_SIGN_IN_EMAIL=demo@bucketeer.io \
	VITE_DEMO_SIGN_IN_PASSWORD=demo \
	yarn build
