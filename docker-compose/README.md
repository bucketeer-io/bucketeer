# Bucketeer Docker Compose Setup

This directory contains Docker Compose configuration for running Bucketeer locally in a containerized environment.

## Overview

The Docker Compose setup provides an alternative to the Minikube deployment for local development and testing. It includes all necessary services to run Bucketeer with minimal external dependencies.

## Architecture

```
┌─────────────────┐    ┌─────────────────┐
│     MySQL       │    │     Redis       │
│   (Database)    │    │ (Cache/PubSub)  │
│   Port: 3306    │    │   Port: 6379    │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────────┬───────────┘
                     │
┌─────────────────────────────────────────┐
│              Migration                  │
│         (Schema Setup)                  │
└─────────────────────────────────────────┘
                     │
    ┌───────────────┴────────────────┐
    │                                │
┌─────────────────┐    ┌─────────────────┐
│   Web Service   │    │   API Service   │
│ (gRPC internal) │    │ (gRPC internal) │
└─────────────────┘    └─────────────────┘
         │                       │
         └───────────┬───────────┘
                     │
           ┌─────────────────┐
           │                 │
    ┌─────────────────┐    ┌─────────────────┐
    │ Batch Service   │    │ Subscriber Svc  │
    │ (gRPC internal) │    │ (gRPC internal) │
    └─────────────────┘    └─────────────────┘
                     │
                     |
┌─────────────────────────────────────────┐
│                   Nginx                 │
│              (Reverse Proxy)            │
│  Ports: 80 (HTTP), 443 (HTTPS)          │
│                                         │
│ - web-gateway.bucketeer.io              │
│ - api-gateway.bucketeer.io              │
└─────────────────────────────────────────┘
                    │
                    |
    ┌───────────────┴────────────────┐
    │                                │
┌─────────────────┐    ┌─────────────────┐
│   Client/SDK    │    │   Admin UI      │
│                 │    │                 │
└─────────────────┘    └─────────────────┘
```

## Services

### Infrastructure Services
- **MySQL 8.0**: Primary database (port 3306)
- **Redis 7**: Cache and pub/sub messaging (port 6379)
- **Nginx**: Reverse proxy for routing (ports 80, 443)

### Application Services
- **Migration**: Database schema migration (runs once)
- **Web**: Core backend with gRPC gateway (frontend + admin APIs, internal ports)
- **API**: Gateway service with gRPC gateway (SDK client communication, internal ports)
- **Batch**: Background job processing (internal ports)
- **Subscriber**: Event processing service (internal ports)

## Prerequisites

1. **Docker & Docker Compose**: Ensure you have Docker and Docker Compose installed
2. **Development Certificates**: The setup requires TLS certificates in `../tools/dev/cert/`
3. **Hosts file**: You need to add the following entries to your `/etc/hosts` file:
   ```
   127.0.0.1 web-gateway.bucketeer.io
   127.0.0.1 api-gateway.bucketeer.io
   ```

## Quick Start

1. **Start all services**:
   ```shell
   make docker-compose-up
   ```

2. **Check service status**:
   ```shell
   make docker-compose-status
   ```

3. **View logs**:
   ```shell
   make docker-compose-logs
   ```

4. **Stop services**:
   ```shell
   make docker-compose-down
   ```

## Configuration

### Docker Image Versions

You can customize the Docker image versions used by setting environment variables:

#### Using Default Versions
```shell
# Use the default versions specified in compose.yml
make docker-compose-up
```

#### Using Environment File
```shell
# Option 1: Source the default environment file
source env.default
make docker-compose-up

# Option 2: Copy and modify the default environment file
cp env.default .env
# Edit .env with your preferred versions
# Then start services
make docker-compose-up
```

#### Using Environment Variables
```shell
# Set versions for all Bucketeer services
export BUCKETEER_VERSION=v1.4.0
make docker-compose-up

# Or set versions for specific services
export BUCKETEER_WEB_VERSION=v1.4.0
export BUCKETEER_API_VERSION=v1.3.0
make docker-compose-up

# Set infrastructure versions
export MYSQL_VERSION=8.1
export REDIS_VERSION=7.2-alpine
make docker-compose-up
```

#### Available Environment Variables
- `MYSQL_VERSION` (default: 8.0)
- `REDIS_VERSION` (default: 7-alpine)
- `NGINX_VERSION` (default: 1.25-alpine)
- `BUCKETEER_VERSION` (default: localenv) - Sets version for all Bucketeer services
- `BUCKETEER_MIGRATION_VERSION` - Override migration service version (default: v0.4.5)
- `BUCKETEER_WEB_VERSION` - Override web service version
- `BUCKETEER_API_VERSION` - Override API service version
- `BUCKETEER_BATCH_VERSION` - Override batch service version
- `BUCKETEER_SUBSCRIBER_VERSION` - Override subscriber service version

### Subscriber Service Configuration

The subscriber service uses JSON configuration files located in `config/subscriber-config/`:

- **`subscribers.json`**: Main subscriber configurations for different event types
- **`onDemandSubscribers.json`**: On-demand subscriber configurations
- **`processors.json`**: Processor configurations with flush settings
- **`onDemandProcessors.json`**: On-demand processor configurations with BigQuery settings

### Environment Variables

All services are configured with appropriate environment variables for:
- Database connections (MySQL)
- Cache connections (Redis)
- Inter-service communication
- Pub/sub messaging (Redis Streams)
- TLS certificates and tokens

## Service Dependencies

The services start in the following order:
1. MySQL and Redis (parallel)
2. Migration (depends on MySQL health)
3. Web service (depends on migration completion)
4. API, Batch, and Subscriber services (depend on Web service)
5. Nginx (depends on Web and API services)

## Service Access

After starting the services, you can access:

- **Web Dashboard**: https://web-gateway.bucketeer.io
  - Admin interface for managing feature flags, experiments, etc.
  - Routes internally to Web service gRPC Gateway.

- **API Gateway**: https://api-gateway.bucketeer.io
  - SDK and client API endpoints.
  - Routes internally to API service gRPC Gateway.

- **Health Check**:
  - Web: https://web-gateway.bucketeer.io/health
  - API: https://api-gateway.bucketeer.io/health

**Note**: All Bucketeer services communicate using HTTPS/TLS internally. Nginx handles SSL termination and routing to the appropriate service ports.

### Internal Service Ports

The Docker Compose setup uses nginx as a reverse proxy. Internal service architecture:

**API Service:**
- Main gRPC: 9090
- gRPC Gateway: 9089 (used by nginx)
- Health Check: 9090 (same as main gRPC)

**Web Service:**
- gRPC Gateway: 9089 (used by nginx)
- Health Check: 8000
- Individual services: 9091-9107 (internal routing)

**Batch & Subscriber Services:**
- Main gRPC: 9000
- gRPC Gateway: 9089

**Communication Flow:**
1. Client/Admin UI → nginx (80/443) via `web-gateway.bucketeer.io` or `api-gateway.bucketeer.io`
2. nginx → API/Web services (gRPC Gateway: 9089, Web health: 8000, API health: 9090)
3. Services communicate internally via gRPC.

### Direct Database Access
- **MySQL**: localhost:3306 (bucketeer/bucketeer)
- **Redis**: localhost:6379

## Development Commands

### Docker Compose Management
```shell
# Start all services
make start-docker-compose

# Stop all services
make docker-compose-down

# View service status
make docker-compose-status

# View logs from all services
make docker-compose-logs

# Clean up everything (containers, networks, volumes)
make docker-compose-clean

# Create MySQL event tables for data warehouse functionality
make docker-compose-create-mysql-event-tables
```

### Individual Service Management
```shell
# View logs for specific service (Docker Compose v2)
docker compose -f docker-compose/compose.yml logs -f web
# Or for Docker Compose v1
docker-compose -f docker-compose/compose.yml logs -f web

# Restart a specific service (Docker Compose v2)
docker compose -f docker-compose/compose.yml restart api
# Or for Docker Compose v1
docker-compose -f docker-compose/compose.yml restart api

# Scale a service (if supported) - Docker Compose v2
docker compose -f docker-compose/compose.yml up -d --scale batch=2
# Or for Docker Compose v1
docker-compose -f docker-compose/compose.yml up -d --scale batch=2
```

## Running E2E Tests

When using Docker Compose, you can run E2E tests against the local services:

### Create API Keys
```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
WEB_GATEWAY_PORT=443 \
WEB_GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=$PWD/tools/dev/cert/service-token \
API_KEY_NAME="e2e-test-$(date +%s)-client" \
API_KEY_PATH=$PWD/tools/dev/cert/api_key_client \
API_KEY_ROLE=SDK_CLIENT \
ENVIRONMENT_ID=e2e \
make create-api-key

WEB_GATEWAY_URL=web-gateway.bucketeer.io \
WEB_GATEWAY_PORT=443 \
WEB_GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=$PWD/tools/dev/cert/service-token \
API_KEY_NAME="e2e-test-$(date +%s)-server" \
API_KEY_PATH=$PWD/tools/dev/cert/api_key_server \
API_KEY_ROLE=SDK_SERVER \
ENVIRONMENT_ID=e2e \
make create-api-key
```

### Run E2E Tests
```shell
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
GATEWAY_URL=api-gateway.bucketeer.io \
WEB_GATEWAY_PORT=443 \
GATEWAY_PORT=443 \
WEB_GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=$PWD/tools/dev/cert/service-token \
API_KEY_PATH=$PWD/tools/dev/cert/api_key_client \
API_KEY_SERVER_PATH=$PWD/tools/dev/cert/api_key_server \
ENVIRONMENT_ID=e2e \
ORGANIZATION_ID=default \
make e2e
```

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 80, 443, 3306, 6379 are available on your host machine.
2. **Certificate issues**: Make sure development certificates exist in `../tools/dev/cert/`.
3. **`/etc/hosts` not configured**: Ensure `web-gateway.bucketeer.io` and `api-gateway.bucketeer.io` are mapped to `127.0.0.1`.
4. **Service startup order**: Services have dependencies; wait for each tier to be healthy.

### Debugging Commands

```shell
# Check service health (Docker Compose v2)
docker compose -f docker-compose/compose.yml ps
# Or for Docker Compose v1
docker-compose -f docker-compose/compose.yml ps

# View detailed logs for a service (Docker Compose v2)
docker compose -f docker-compose/compose.yml logs -f --tail=100 web
# Or for Docker Compose v1
docker-compose -f docker-compose/compose.yml logs -f --tail=100 web

```

### Reset Environment

```shell
# Complete cleanup and restart
make docker-compose-clean
make start-docker-compose
```

## Differences from Minikube Setup

- **No Kubernetes**: Uses Docker Compose instead of Minikube/Kubernetes
- **Simpler networking**: Direct container-to-container communication
- **Local storage**: Uses Docker volumes instead of persistent volumes
- **Redis Streams**: Uses Redis Streams for pub/sub instead of Google Pub/Sub emulator
- **No BigQuery emulator**: Events are processed but not stored in BigQuery
- **TLS between services**: All internal gRPC communication between services is secured with TLS. Nginx proxy also connects to backends using TLS.

## File Structure

```
docker-compose/
├── README.md                 # This file
├── compose.yml               # Main Docker Compose configuration
├── env.default               # Default environment variables
├── .gitignore                # Git ignore rules
└── config/
    ├── nginx/                # Nginx configuration
    │   ├── nginx.conf
    │   └── bucketeer.conf
    ├── datawarehouse.yaml    # Data warehouse configuration
    ├── oauth-config.json     # OAuth configuration
    └── subscriber-config/    # Subscriber service configuration
        ├── subscribers.json
        ├── onDemandSubscribers.json
        ├── processors.json
        └── onDemandProcessors.json
```

## Contributing

When making changes to the Docker Compose setup:

1. Test with `make docker-compose-up`
2. Verify all services start correctly
3. Run E2E tests to ensure functionality
4. Update this README if adding new services or configuration options 