# Bucketeer Deployment Guide

This guide covers various deployment options for Bucketeer, from local development to production environments.

## Table of Contents

- [Docker Compose Deployment (Recommended for Small-Medium Companies)](#docker-compose-deployment)
- [Kubernetes Deployment](#kubernetes-deployment)
- [Production Considerations](#production-considerations)
- [Monitoring and Observability](#monitoring-and-observability)
- [Troubleshooting](#troubleshooting)

## Docker Compose Deployment

Docker Compose provides the easiest way to deploy Bucketeer for small to medium-sized companies without requiring Kubernetes expertise.

### Prerequisites

1. **System Requirements**
   - **Minimum**: 4GB RAM, 2 CPU cores, 20GB disk space
   - **Recommended**: 8GB RAM, 4 CPU cores, 50GB disk space
   - Docker Engine 20.10+ and Docker Compose 2.0+

2. **Network Requirements**
   - Ports 80, 443, 3306, 6379 available
   - Internet access for downloading Docker images

3. **Domain Setup** (Optional for production)
   - DNS records pointing to your server IP
   - SSL certificates (or use self-signed for development)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/bucketeer-io/bucketeer.git
   cd bucketeer
   ```

2. **Start all services**
   ```bash
   make docker-compose-up
   ```

   This command will:
   - Generate development certificates if they don't exist
   - Create necessary directories and configuration files
   - Build Bucketeer Docker images with embedded web console
   - Start all services with proper dependencies

3. **Verify deployment**
   ```bash
   # Check service status
   make docker-compose-status

   # View logs
   make docker-compose-logs
   ```

4. **Access the application**
   - **Admin Dashboard**: https://web-gateway.bucketeer.io
   - **API Gateway**: http://localhost (development) or https://api-gateway.bucketeer.io (production)
   - **Health Checks**:
     - Web: https://web-gateway.bucketeer.io/health
     - API: http://localhost/health

### Configuration

#### Environment Variables

Customize deployment by setting environment variables:

```bash
# Use specific versions
export BUCKETEER_VERSION=v1.4.0
export MYSQL_VERSION=8.1
export REDIS_VERSION=7.2-alpine

# Start with custom versions
make docker-compose-up
```

#### Using Environment Files

```bash
# Copy and customize environment file
cp docker-compose/env.default docker-compose/.env

# Edit .env file with your preferred versions
vim docker-compose/.env

# Deploy
make docker-compose-up
```

#### Available Environment Variables

```bash
# Infrastructure versions
MYSQL_VERSION=8.0                    # MySQL version
REDIS_VERSION=7-alpine              # Redis version
NGINX_VERSION=1.25-alpine           # Nginx version

# Bucketeer service versions
BUCKETEER_VERSION=localenv          # Version for all services
BUCKETEER_MIGRATION_VERSION=v0.4.5  # Migration service version
BUCKETEER_WEB_VERSION=localenv      # Web service version
BUCKETEER_API_VERSION=localenv      # API service version
BUCKETEER_BATCH_VERSION=localenv    # Batch service version
BUCKETEER_SUBSCRIBER_VERSION=localenv # Subscriber service version
```

### Production Deployment

For production deployments, follow these additional steps:

#### 1. Security Configuration

**Set up proper SSL certificates:**
```bash
# Replace development certificates with real ones
cp your-ssl-cert.crt tools/dev/cert/tls.crt
cp your-ssl-key.key tools/dev/cert/tls.key
```

**Configure CORS policy:**
Edit `docker-compose/config/nginx/bucketeer.conf`:
```nginx
# Replace localhost with your actual domain
add_header 'Access-Control-Allow-Origin' 'https://your-domain.com' always;
```

**Use Docker secrets for sensitive data:**
```bash
# Create secret files
mkdir -p docker-compose/secrets
echo "your-mysql-password" > docker-compose/secrets/mysql_password.txt
echo "your-mysql-root-password" > docker-compose/secrets/mysql_root_password.txt
```

#### 2. Performance Configuration

**Update resource limits for production:**
The current configuration includes resource limits suitable for medium deployments:
- MySQL: 2GB RAM, 1 CPU
- Web Service: 1.5GB RAM, 1 CPU
- API Service: 1GB RAM, 0.75 CPU

For high-traffic environments, increase these limits in `docker-compose/compose.yml`.

#### 3. Persistence and Backup

**Configure data persistence:**
```bash
# Create backup script
cat > backup-bucketeer.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec bucketeer-mysql mysqldump -u root -p bucketeer > backup_${DATE}.sql
docker run --rm -v mysql_data:/data -v $(pwd):/backup alpine tar czf /backup/mysql_data_${DATE}.tar.gz /data
EOF
chmod +x backup-bucketeer.sh

# Run daily backups via cron
echo "0 2 * * * /path/to/backup-bucketeer.sh" | crontab -
```

#### 4. Domain and Hosts Configuration

**For production with custom domains:**
```bash
# Update /etc/hosts on client machines (development only)
echo "YOUR_SERVER_IP your-bucketeer-domain.com" >> /etc/hosts

# Or configure proper DNS records:
# A record: your-bucketeer-domain.com -> YOUR_SERVER_IP
# A record: web-gateway.your-bucketeer-domain.com -> YOUR_SERVER_IP
# A record: api-gateway.your-bucketeer-domain.com -> YOUR_SERVER_IP
```

### Common Management Tasks

#### Service Management
```bash
# Stop all services
make docker-compose-down

# Restart specific service
docker compose -f docker-compose/compose.yml restart web

# Scale a service (if supported)
docker compose -f docker-compose/compose.yml up -d --scale api=2

# View service logs
docker compose -f docker-compose/compose.yml logs -f web

# Execute command in container
docker exec -it bucketeer-web sh
```

#### Data Management
```bash
# Create MySQL event tables for data warehouse
make docker-compose-create-mysql-event-tables

# Delete E2E test data
make docker-compose-delete-data

# Reset entire environment
make docker-compose-clean
make docker-compose-up
```

#### Monitoring
```bash
# Check resource usage
docker stats

# Monitor service health
make docker-compose-status

# Check nginx access logs
tail -f docker-compose/logs/nginx/access.log
```

### API Key Management

For E2E testing and SDK integration:

```bash
# Create client API key
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
WEB_GATEWAY_PORT=443 \
WEB_GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=$PWD/tools/dev/cert/service-token \
API_KEY_NAME="production-client-$(date +%s)" \
API_KEY_PATH=$PWD/tools/dev/cert/api_key_client \
API_KEY_ROLE=SDK_CLIENT \
ENVIRONMENT_ID=production \
make create-api-key

# Create server API key
WEB_GATEWAY_URL=web-gateway.bucketeer.io \
WEB_GATEWAY_PORT=443 \
WEB_GATEWAY_CERT_PATH=$PWD/tools/dev/cert/tls.crt \
SERVICE_TOKEN_PATH=$PWD/tools/dev/cert/service-token \
API_KEY_NAME="production-server-$(date +%s)" \
API_KEY_PATH=$PWD/tools/dev/cert/api_key_server \
API_KEY_ROLE=SDK_SERVER \
ENVIRONMENT_ID=production \
make create-api-key
```

## Kubernetes Deployment

For larger deployments requiring high availability and scalability, use Kubernetes with Helm:

### Prerequisites
- Kubernetes cluster (1.20+)
- Helm 3.0+
- kubectl configured

### Deployment Steps

```bash
# Add required repositories
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Deploy dependencies
kubectl create namespace bucketeer
helm install mysql bitnami/mysql -n bucketeer
helm install redis bitnami/redis -n bucketeer

# Deploy Bucketeer
helm install bucketeer ./manifests/bucketeer -n bucketeer \
  --values ./manifests/bucketeer/values.prod.yaml
```

For detailed Kubernetes deployment instructions, see the [Minikube development guide](DEVELOPMENT.md#option-1-minikube-setup).

## Production Considerations

### ⚠️ Critical Data Persistence Warning

**IMPORTANT**: Redis in Bucketeer stores critical event data using Redis Streams, not just cache data. **Never configure Redis with memory eviction policies** like `allkeys-lru` or `volatile-lru` as this will cause **permanent data loss**.

The current Docker Compose configuration includes proper persistence settings:
- **AOF (Append Only File)**: `--appendonly yes` ensures all writes are logged
- **RDB Snapshots**: Multiple save points for data durability
- **fsync Policy**: `--appendfsync everysec` balances performance and durability

For production deployments:
```bash
# Add these Redis settings for maximum data safety
--maxmemory-policy noeviction    # CRITICAL: Never evict data
--maxmemory 4gb                 # Set based on your RAM capacity
--appendfsync always            # Maximum durability (slower performance)
```

### High Availability

For production environments, consider:

1. **Database High Availability**
   - Use managed MySQL services (AWS RDS, Google Cloud SQL)
   - Set up MySQL master-slave replication
   - Configure automated backups

2. **Redis High Availability**
   - Use managed Redis services (AWS ElastiCache, Google Memorystore)
   - Set up Redis Sentinel or Cluster mode

3. **Load Balancing**
   - Use external load balancers (AWS ALB, Google Cloud Load Balancer)
   - Configure health checks and auto-scaling

4. **Container Orchestration**
   - Migrate to Kubernetes for production scale
   - Use Docker Swarm for simpler orchestration needs

### Monitoring and Observability

Add monitoring stack to your deployment:

```yaml
# Add to docker-compose.yml
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./config/prometheus:/etc/prometheus
    networks:
      - bucketeer

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana
    networks:
      - bucketeer
```

### Security Hardening

1. **Network Security**
   - Use private networks for internal communication
   - Implement firewall rules
   - Enable TLS everywhere

2. **Access Control**
   - Use OAuth providers for authentication
   - Implement RBAC for user management
   - Regular security audits

3. **Data Protection**
   - Encrypt data at rest
   - Use secrets management systems
   - Regular security updates

## Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check which process is using a port
   lsof -i :3306

   # Solution: Stop conflicting services or change ports
   ```

2. **Certificate Issues**
   ```bash
   # Regenerate development certificates
   make -C tools/dev generate-tls-certificate
   make -C tools/dev generate-oauth
   ```

3. **Service Health Issues**
   ```bash
   # Check service logs
   docker compose -f docker-compose/compose.yml logs service-name

   # Check resource usage
   docker stats

   # Verify network connectivity
   docker exec bucketeer-web ping mysql
   ```

4. **Memory Issues**
   ```bash
   # Check system memory
   free -h

   # Check Docker memory usage
   docker system df

   # Clean up unused resources
   docker system prune -f
   ```

### Performance Tuning

1. **Database Optimization**
   ```sql
   -- Optimize MySQL settings in docker-compose.yml
   command: --default-authentication-plugin=mysql_native_password
           --innodb-buffer-pool-size=1G
           --max-connections=200
   ```

2. **Redis Configuration for Data Persistence**
   ```bash
   # CRITICAL: Redis stores event data via Redis Streams, not just cache
   # Never use memory eviction policies like allkeys-lru for production
   # Current configuration ensures data persistence:
   command: redis-server --appendonly yes --appendfsync everysec
            --save 900 1 --save 300 10 --save 60 10000

   # For production, consider additional settings:
   # --maxmemory-policy noeviction  # Prevent data loss
   # --maxmemory 2gb               # Set appropriate memory limit
   ```

3. **Application Tuning**
   - Adjust worker processes and connection pools
   - Monitor and optimize database queries
   - Implement caching strategies

## Support and Community

- **Documentation**: [Official Docs](https://docs.bucketeer.io)
- **GitHub Issues**: [Report Issues](https://github.com/bucketeer-io/bucketeer/issues)
- **Community**: [Discord](https://discord.gg/bucketeer)

For enterprise support and professional services, contact the Bucketeer team.
