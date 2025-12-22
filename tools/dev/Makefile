NAMESPACE := default

# For minikube image caching documentation, see: MINIKUBE_CACHING.md

generate-tls-certificate:
	mkdir -p $(CURDIR)/cert
	openssl req -x509 -nodes -days 36500 -newkey rsa:2048 \
		-keyout $(CURDIR)/cert/tls.key \
		-out $(CURDIR)/cert/tls.crt \
		-subj "/CN=*.bucketeer.io" \
		-config service.config

generate-oauth:
	mkdir -p ${CURDIR}/$*
	openssl genrsa -out ${CURDIR}/cert/oauth-private.pem 4096
	openssl rsa -in ${CURDIR}/cert/oauth-private.pem -pubout -out ${CURDIR}/cert/oauth-public.pem

service-cert-secret:
	kubectl config use-context minikube
	kubectl delete secret bucketeer-service-cert --namespace $(NAMESPACE) --ignore-not-found
	kubectl create secret tls bucketeer-service-cert --key ${CURDIR}/cert/tls.key --cert ${CURDIR}/cert/tls.crt --namespace $(NAMESPACE)

service-token-secret:
	kubectl config use-context minikube
	kubectl delete secret bucketeer-service-token --namespace $(NAMESPACE) --ignore-not-found
	kubectl create secret generic bucketeer-service-token --from-file token=${CURDIR}/cert/service-token --namespace $(NAMESPACE)

oauth-key-secret:
	kubectl config use-context minikube
	kubectl delete secret bucketeer-oauth-key --namespace $(NAMESPACE) --ignore-not-found
	kubectl create secret generic bucketeer-oauth-key --from-file public.pem=${CURDIR}/cert/oauth-public.pem --from-file private.pem=${CURDIR}/cert/oauth-private.pem --namespace $(NAMESPACE)

# Pre-cache Docker images to persist them across minikube deletions
pre-cache-docker-images:
	@echo "Pre-caching Docker images to survive minikube deletions..."
	@echo "Pulling base minikube image..."
	docker pull gcr.io/k8s-minikube/kicbase:v0.0.47 || echo "Failed to pull kicbase, will be downloaded later"
	@echo "Pulling Kubernetes preload image..."
	docker pull registry.k8s.io/pause:3.10 || echo "Failed to pull pause image"
	@echo "Pulling ingress controller images..."
	docker pull registry.k8s.io/ingress-nginx/controller:v1.12.2 || echo "Failed to pull ingress controller"
	docker pull registry.k8s.io/ingress-nginx/kube-webhook-certgen:v1.5.3 || echo "Failed to pull webhook certgen"
	@echo "Pulling storage provisioner..."
	docker pull gcr.io/k8s-minikube/storage-provisioner:v5 || echo "Failed to pull storage provisioner"
	@echo "Pulling application images..."
	docker pull mysql:8.0 || echo "Failed to pull mysql"
	docker pull postgres:18 || echo "Failed to pull postgres"
	docker pull ghcr.io/bucketeer-io/bigquery-emulator:0.6.6 || echo "Failed to pull bigquery emulator"
	docker pull gcr.io/google.com/cloudsdktool/google-cloud-cli:545.0.0 || echo "Failed to pull cloud-sdk"
	@echo "Docker image pre-caching completed. These images will persist across minikube deletions."

# Smart minikube setup that uses cached images when available
setup-minikube-cached:
	@echo "Setting up minikube with cached images..."
	@$(MAKE) pre-cache-docker-images
	minikube start --memory max --cpus max --cache-images
	minikube addons enable ingress
	kubectl config use-context minikube
	kubectl patch deployment -n ingress-nginx ingress-nginx-controller --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value":"--enable-ssl-passthrough"}]'
	helm install localenv ../../manifests/localenv

setup-minikube:
	minikube start --memory max --cpus max
	minikube addons enable ingress
	kubectl config use-context minikube
	kubectl patch deployment -n ingress-nginx ingress-nginx-controller --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value":"--enable-ssl-passthrough"}]'
	helm install localenv ../../manifests/localenv

start-minikube:
	minikube start --cache-images
	# ensure ingress addon is enabled
	minikube addons enable ingress
	# used to enable ssl-passthrough for nginx ingress container
	kubectl config use-context minikube
	kubectl patch deployment -n ingress-nginx ingress-nginx-controller --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value":"--enable-ssl-passthrough"}]'
	# update helm dependencies
	cd ../../manifests/localenv && helm dependency update
	# Delete old Bitnami StatefulSets if they exist (migration from Bitnami charts)
	# This is needed because selector labels changed from app.kubernetes.io/name: localenv to app.kubernetes.io/name: redis/mysql/postgres
	kubectl delete statefulset localenv-mysql localenv-redis localenv-postgresql -n default --ignore-not-found=true
	# Install or upgrade localenv
	helm list | grep -q localenv && helm upgrade localenv ../../manifests/localenv || helm install localenv ../../manifests/localenv

# Clean Docker cache and volumes (use when troubleshooting or freeing space)
clean-cache:
	@echo "Cleaning Docker cache and persistent volumes..."
	@echo "This will remove all cached images and force fresh downloads next time."
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	@echo "Stopping minikube..."
	-minikube stop 2>/dev/null || echo "Minikube wasn't running"
	@echo "Removing Docker images..."
	docker system prune -a -f
	@echo "Removing persistent volumes..."
	-docker volume rm bucketeer-docker-cache bucketeer-minikube-cache 2>/dev/null || echo "Volumes didn't exist"
	@echo "Cache cleaned! Run 'make pre-cache-docker-images' to rebuild cache."

# Light cleanup - only remove unused Docker images (keeps volumes)
clean-cache-light:
	@echo "Light cleanup: removing unused Docker images only..."
	docker image prune -f
	@echo "Unused images removed. Volumes and cached images preserved."

# Show cache status and disk usage
cache-status:
	@echo "Cache Status"
	@echo "=================="
	@echo "\nDocker Images:"
	docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | head -20
	@echo "\nPersistent Volumes:"
	docker volume ls | grep bucketeer || echo "No bucketeer volumes found"
	@echo "\nDisk Usage:"
	docker system df
	@if command -v minikube >/dev/null 2>&1 && minikube status >/dev/null 2>&1; then \
		echo "\nMinikube Cache:"; \
		minikube cache list || echo "No minikube cache"; \
	fi
