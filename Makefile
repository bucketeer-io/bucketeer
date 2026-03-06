#############################
# Variables
#############################

LOCAL_IMPORT_PATH := github.com/bucketeer-io/bucketeer

# go applications
GO_APP_DIRS := $(wildcard cmd/*)
GO_APP_BUILD_TARGETS := $(addprefix build-,$(notdir $(GO_APP_DIRS)))
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOTESTSUM_VERSION := v1.13.0

ifeq ($(GOARCH), arm64)
	PLATFORM = linux/arm64
else
	PLATFORM = linux/x86_64
endif

LDFLAGS_PACKAGE := github.com/bucketeer-io/bucketeer/pkg/ldflags
LDFLAGS_VERSION := $(LDFLAGS_PACKAGE).Version
LDFLAGS_HASH := $(LDFLAGS_PACKAGE).Hash
LDFLAGS_BUILDDATE := $(LDFLAGS_PACKAGE).BuildDate

#############################
# Go
#############################

.PHONY: generate-all
generate-all: proto-all mockgen

.PHONY: check-all
check-all: proto-check mockgen update-repos diff-check lint build-go test-go

.PHONY: local-deps
local-deps:
	mkdir -p ~/go-tools; \
	cd ~/go-tools; \
	if [ ! -e go.mod ]; then go mod init go-tools; fi; \
	go install golang.org/x/tools/cmd/goimports@v0.40.0; \
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.7.2; \
	go install go.uber.org/mock/mockgen@v0.4.0; \
	go install github.com/golang/protobuf/protoc-gen-go@v1.5.2; \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0; \
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0; \
	go install github.com/nilslice/protolock/...@v0.15.0; \
	go install github.com/mikefarah/yq/v4@v4.28.2

.PHONY: lint
lint:
	golangci-lint run --timeout 3m0s ./cmd/... ./evaluation/go/... ./pkg/... ./hack/... ./test/...

.PHONY: gofmt
gofmt:
	goimports -local ${LOCAL_IMPORT_PATH} -w \
		$$(find . -path "./vendor" -prune -o -path "./proto" -prune -o -type f -name '*.go' -print)

.PHONY: gofmt-check
gofmt-check: gofmt diff-check

.PHONY: proto-check
proto-check:
	make -C proto check

.PHONY: proto-fmt
proto-fmt:
	make -C proto fmt

.PHONY: proto-fmt-check
proto-fmt-check:
	make -C proto fmt-check

.PHONY: proto-lock-check
proto-lock-check:
	make -C proto lock-check

.PHONY: proto-lock-commit
proto-lock-commit:
	make -C proto lock-commit

.PHONY: proto-lock-commit-force
proto-lock-commit-force:
	make -C proto lock-commit-force

.PHONY: proto-all
proto-all: proto-fmt proto-lock-commit proto-go proto-openapi-gen

.PHONY: proto-go
proto-go:
	make -C proto go

.PHONY: proto-go-check
proto-go-check:
	make -C proto go-check

.PHONY: proto-openapi-gen
proto-openapi-gen:
	make -C proto openapi-api-gen
	make -C proto openapi-web-gen

.PHONY: openapi-ui
proto-openapi-ui:
	make -C proto openapi-api-ui

.PHONY: mockgen
mockgen: proto-go
	find ./pkg -path "*mock*.go" -type f -delete
	go generate -run="mockgen" ./pkg/...
	make gofmt

.PHONY: mockgen-check
mockgen-check: mockgen diff-check

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: update-repos
update-repos: tidy-deps vendor

.PHONY: update-repos-check
update-repos-check: update-repos diff-check

.PHONY: diff-check
diff-check:
	test -z "$$(git diff --name-only)"

.PHONY: migration-validate
migration-validate:
	atlas migrate validate --dir file://migration/mysql

.PHONY: migration-hash-check
migration-hash-check:
	atlas migrate hash --dir file://migration/mysql
	make diff-check

.PHONY: tidy-deps
tidy-deps:
	go mod tidy

.PHONY: $(GO_APP_BUILD_TARGETS)
$(GO_APP_BUILD_TARGETS): build-%:
	$(eval VERSION := $(shell git describe --tags --always --abbrev=7))
	$(eval HASH := $(shell git rev-parse --verify HEAD))
	$(eval BUILDDATE := $(shell date '+%Y/%m/%dT%H:%M:%S%Z'))
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags "-s -w -X $(LDFLAGS_VERSION)=$(VERSION) -X $(LDFLAGS_HASH)=$(HASH) -X $(LDFLAGS_BUILDDATE)=$(BUILDDATE)" \
		-o bin/$* -mod=vendor cmd/$*/$*.go

.PHONY: clean-web-console
clean-web-console:
	rm -rf ui/dashboard/build/*
	touch ui/dashboard/build/DONT-EDIT-FILES-IN-THIS-DIRECTORY

.PHONY: build-web-console
build-web-console:
	rm -rf ui/dashboard/build/*
	make -C ui/dashboard install build

.PHONY: build-go
build-go: $(GO_APP_BUILD_TARGETS)

.PHONY: build-go-embed
build-go-embed: build-web-console $(GO_APP_BUILD_TARGETS) clean-web-console

# Make sure bucketeer-httpstan is already running. If not, run "make start-httpstan".
.PHONY: test-go
test-go:
	TZ=UTC CGO_ENABLED=0 go run gotest.tools/gotestsum@$(GOTESTSUM_VERSION) \
		--format pkgname \
		-- -v ./pkg/... ./evaluation/go/...

.PHONY: start-httpstan
start-httpstan:
	@if docker ps -q -f name=bucketeer-httpstan | grep -q .; then \
		echo "httpstan container is already running"; \
	elif docker ps -aq -f name=bucketeer-httpstan | grep -q .; then \
		echo "Starting existing httpstan container..."; \
		docker start bucketeer-httpstan; \
	else \
		echo "Creating and starting httpstan container..."; \
		docker run --name bucketeer-httpstan -p 8080:8080 -d ghcr.io/bucketeer-io/bucketeer-httpstan:0.0.1; \
	fi

.PHONY: stop-httpstan
stop-httpstan:
	@docker stop bucketeer-httpstan 2>/dev/null || echo "httpstan container is not running"

#############################
# Charts
#############################

.PHONY: build-chart
build-chart: VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7)
build-chart:
	mkdir -p .artifacts
	helm package manifests/bucketeer --version $(VERSION) --app-version $(VERSION) --dependency-update --destination .artifacts

.PHONY: build-migration-chart
build-migration-chart: VERSION ?= $(shell git describe --tags --always --abbrev=7)
build-migration-chart:
	mkdir -p .artifacts
	helm package manifests/bucketeer-migration --version $(VERSION) --app-version $(VERSION) --destination .artifacts

#############################
# E2E for backend
#############################

.PHONY: delete-e2e-data-mysql
delete-e2e-data-mysql:
ifeq ($(GOOS), darwin)
	make -C hack/delete-e2e-data-mysql clean build-darwin
else
	make -C hack/delete-e2e-data-mysql clean build
endif
	./hack/delete-e2e-data-mysql/delete-e2e-data-mysql delete \
		--mysql-user=${MYSQL_USER} \
		--mysql-pass=${MYSQL_PASS} \
		--mysql-host=${MYSQL_HOST} \
		--mysql-port=${MYSQL_PORT} \
		--mysql-db-name=${MYSQL_DB_NAME} \
		--test-id=${TEST_ID} \
		--no-profile \
		--no-gcp-trace-enabled

.PHONY: create-mysql-event-tables
create-mysql-event-tables:
	@echo "Creating MySQL event tables for data warehouse..."
	go run ./hack/create-mysql-event-tables create \
		--mysql-host=${MYSQL_HOST} \
		--mysql-port=${MYSQL_PORT} \
		--mysql-user=${MYSQL_USER} \
		--mysql-pass=${MYSQL_PASS} \
		--mysql-db-name=${MYSQL_DB_NAME} \
		--no-profile \
		--no-gcp-trace-enabled

.PHONY: create-postgres-event-tables
create-postgres-event-tables:
	@echo "Creating Postgres event tables for data warehouse..."
	go run ./hack/create-postgres-event-tables create \
		--postgres-host=${POSTGRES_HOST} \
		--postgres-port=${POSTGRES_PORT} \
		--postgres-user=${POSTGRES_USER} \
		--postgres-pass=${POSTGRES_PASS} \
		--postgres-db-name=${POSTGRES_DB_NAME} \
		--no-profile \
		--no-gcp-trace-enabled \
		--log-level=debug

.PHONY: generate-service-token
generate-service-token:
	go run ./hack/generate-service-token generate \
		--issuer=${ISSUER} \
		--audience=bucketeer \
		--email=${EMAIL} \
		--role=OWNER \
		--key=${OAUTH_KEY_PATH} \
		--output=${SERVICE_TOKEN_PATH} \
		--no-profile \
		--no-gcp-trace-enabled

.PHONY: create-api-key
create-api-key:
	go run ./hack/create-api-key create \
		--cert=${WEB_GATEWAY_CERT_PATH} \
		--web-gateway=${WEB_GATEWAY_URL}:443 \
		--service-token=${SERVICE_TOKEN_PATH} \
		--name=${API_KEY_NAME} \
		--role=${API_KEY_ROLE} \
		--output=${API_KEY_PATH} \
		--environment-id=${ENVIRONMENT_ID} \
		--no-profile \
		--no-gcp-trace-enabled

.PHONY: e2e-l4
e2e-l4:
	TZ=UTC CGO_ENABLED=0 go run gotest.tools/gotestsum@$(GOTESTSUM_VERSION) \
		--format pkgname \
		-- -v ./test/e2e/... -args \
		-web-gateway-addr=${WEB_GATEWAY_URL} \
		-web-gateway-port=443 \
		-web-gateway-cert=${WEB_GATEWAY_CERT_PATH} \
		-api-key=${API_KEY_PATH} \
		-gateway-addr=${GATEWAY_URL} \
		-gateway-port=9000 \
		-gateway-cert=${GATEWAY_CERT_PATH} \
		-service-token=${SERVICE_TOKEN_PATH} \
		-environment-id=${ENVIRONMENT_ID} \
		-test-id=${TEST_ID}

.PHONY: e2e
e2e:
	TZ=UTC CGO_ENABLED=0 go run gotest.tools/gotestsum@$(GOTESTSUM_VERSION) \
		--format pkgname \
		-- -v ./test/e2e/... -args \
		-web-gateway-addr=${WEB_GATEWAY_URL} \
		-web-gateway-port=443 \
		-web-gateway-cert=${WEB_GATEWAY_CERT_PATH} \
		-api-key=${API_KEY_PATH} \
		-api-key-server=${API_KEY_SERVER_PATH} \
		-gateway-addr=${GATEWAY_URL} \
		-gateway-port=443 \
		-gateway-cert=${GATEWAY_CERT_PATH} \
		-service-token=${SERVICE_TOKEN_PATH} \
		-environment-id=${ENVIRONMENT_ID} \
		-organization-id=${ORGANIZATION_ID} \
		-test-id=${TEST_ID}	

.PHONY: delete-dev-container-mysql-data
delete-dev-container-mysql-data:
	MYSQL_USER=bucketeer \
	MYSQL_PASS=bucketeer \
	MYSQL_HOST=$$(minikube ip) \
	MYSQL_PORT=32000 \
	MYSQL_DB_NAME=bucketeer \
	make -C ./ delete-e2e-data-mysql

.PHONY: update-copyright
update-copyright:
	./hack/update-copyright/update-copyright.sh


###################
# Database Migration
###################
.PHONY: create-migration
create-migration:
	# Example: make create-migration NAME=create_table_users USER=root PASS=password HOST=localhost PORT=3306 DB=bucketeer
	atlas migrate diff ${NAME} \
		--dir file://migration/mysql \
		--to mysql://${USER}:${PASS}@${HOST}:${PORT}/${DB} \
		--dev-url docker://mysql/8/${DB}

.PHONY: atlas-set-version
atlas-set-version:
	# Example: make atlas-set-version VERSION=20240311022556 USER=root PASS=password HOST=localhost PORT=3306 DB=bucketeer
	atlas migrate set ${VERSION} \
		--dir file://migration/mysql \
		--url mysql://${USER}:${PASS}@${HOST}:${PORT}/${DB}

.PHONY: apply-migration
check-apply-migration:
	# Example: make check-apply-migration USER=root PASS=password HOST=localhost PORT=3306 DB=bucketeer
	atlas migrate apply \
		--dir file://migration/mysql \
		--url mysql://${USER}:${PASS}@${HOST}:${PORT}/${DB} \
		--dry-run

.PHONY: check-rollback-migration
check-rollback-migration:
	@# Example: make check-rollback-migration USER=root PASS=password HOST=localhost PORT=3306 DB=bucketeer
	@# Example: make check-rollback-migration COUNT=3 USER=root PASS=password HOST=localhost PORT=3306 DB=bucketeer
	@# Example: make check-rollback-migration VERSION=20240815043128 USER=root PASS=password HOST=localhost PORT=3306 DB=bucketeer
	@if [ -n "${VERSION}" ]; then \
		atlas migrate down \
			--dir file://migration/mysql \
			--url mysql://${USER}:${PASS}@${HOST}:${PORT}/${DB} \
			--dev-url docker://mysql/8/${DB} \
			--to-version ${VERSION} \
			--dry-run; \
	else \
		atlas migrate down $${COUNT:-1} \
			--dir file://migration/mysql \
			--url mysql://${USER}:${PASS}@${HOST}:${PORT}/${DB} \
			--dev-url docker://mysql/8/${DB} \
			--dry-run; \
	fi

.PHONY: rollback-migration
rollback-migration:
	@# Example: make rollback-migration USER=root PASS=password HOST=localhost PORT=3306 DB=bucketeer
	@# Example: make rollback-migration COUNT=3 USER=root PASS=password HOST=localhost PORT=3306 DB=bucketeer
	@# Example: make rollback-migration VERSION=20240815043128 USER=root PASS=password HOST=localhost PORT=3306 DB=bucketeer
	@if [ -n "${VERSION}" ]; then \
		echo "Rolling back to version ${VERSION}..."; \
		atlas migrate down \
			--dir file://migration/mysql \
			--url mysql://${USER}:${PASS}@${HOST}:${PORT}/${DB} \
			--dev-url docker://mysql/8/${DB} \
			--to-version ${VERSION}; \
	else \
		echo "Rolling back last $${COUNT:-1} migration(s)..."; \
		atlas migrate down $${COUNT:-1} \
			--dir file://migration/mysql \
			--url mysql://${USER}:${PASS}@${HOST}:${PORT}/${DB} \
			--dev-url docker://mysql/8/${DB}; \
	fi

.PHONY: migration-status
migration-status:
	@# Example: make migration-status USER=root PASS=password HOST=localhost PORT=3306 DB=bucketeer
	atlas migrate status \
		--dir file://migration/mysql \
		--url mysql://${USER}:${PASS}@${HOST}:${PORT}/${DB}

#############################
# dev container
#############################

# build devcontainer locally
.PHONY: build-devcontainer
build-devcontainer:
	devcontainer build --workspace-folder=.github --push=false --image-name="ghcr.io/bucketeer-io/bucketeer-devcontainer:latest"

# start minikube
start-minikube:
	if [ $$(minikube status | grep -c "kubelet: Running") -eq 1 ]; then \
		echo "minikube is already running"; \
		exit 1; \
	elif [ $$(minikube status | grep -c "minikube start") -eq 1 ]; then \
		make -C tools/dev setup-minikube; \
	elif [ $$(minikube status | grep -c "Stopped") -gt 1 ]; then \
		make -C tools/dev start-minikube; \
	fi
	sleep 5
	helm uninstall bucketeer --ignore-not-found
	make -C ./ modify-hosts
	make -C ./ setup-localenv

# modify hosts file to access api-gateway and web-gateway
modify-hosts:
	$(eval MINIKUBE_IP := $(shell minikube ip))
	echo "$(MINIKUBE_IP)   web-gateway.bucketeer.io" | sudo tee -a /etc/hosts
	echo "$(MINIKUBE_IP)   api-gateway.bucketeer.io" | sudo tee -a /etc/hosts

# enable vault transit secret engine
enable-vault-transit:
	kubectl config use-context minikube
	kubectl exec localenv-vault-0 -- vault secrets enable transit

# create bigquery-emulator tables (idempotent - safe to run multiple times)
.PHONY: create-pubsub-topics
create-pubsub-topics:
	@echo "Creating PubSub topics..."
	kubectl config use-context minikube
	go run ./hack/create-pubsub-topics create \
		--pubsub-emulator-host=$$(minikube ip):30089 \
		--project=bucketeer-dev \
		--no-profile \
		--no-gcp-trace-enabled

.PHONY: create-bigquery-emulator-tables
create-bigquery-emulator-tables:
	go run ./hack/create-big-query-table create \
		--project=bucketeer-dev \
		--bigquery-emulator=http://$$(minikube ip):31000 \
		--no-gcp-trace-enabled \
		--no-profile

# create mysql event tables for minikube
create-mysql-event-tables-minikube:
	@echo "Creating MySQL event tables for minikube data warehouse..."
	kubectl config use-context minikube
	@echo "Waiting for MySQL pod to be ready..."
	kubectl wait --for=condition=ready pod localenv-mysql-0 --timeout=300s
	@echo "MySQL pod is ready"
	MYSQL_USER=bucketeer \
	MYSQL_PASS=bucketeer \
	MYSQL_HOST=$$(minikube ip) \
	MYSQL_PORT=32000 \
	MYSQL_DB_NAME=bucketeer \
	make create-mysql-event-tables

.PHONY: delete-mysql-data-warehouse-data
delete-mysql-data-warehouse-data:
ifeq ($(GOOS), darwin)
	make -C hack/delete-mysql-data-warehouse clean build-darwin
else
	make -C hack/delete-mysql-data-warehouse clean build
endif
	./hack/delete-mysql-data-warehouse/delete-mysql-data-warehouse truncate \
		--mysql-user=bucketeer \
		--mysql-pass=bucketeer \
		--mysql-host=$$(minikube ip) \
		--mysql-port=32000 \
		--mysql-db-name=bucketeer \
		--no-profile \
		--no-gcp-trace-enabled

.PHONY: delete-redis-retry-keys
delete-redis-retry-keys:
ifeq ($(GOOS), darwin)
	make -C hack/delete-redis-retry-keys clean build-darwin
else
	make -C hack/delete-redis-retry-keys clean build
endif
	./hack/delete-redis-retry-keys/delete-redis-retry-keys delete \
		--redis-addr=$(if $(REDIS_ADDR),$(REDIS_ADDR),$$(minikube ip):32001) \
		--environment-id=e2e \
		--no-profile \
		--no-gcp-trace-enabled

POSTGRES_ENABLED ?= false
setup-localenv:
	kubectl config use-context minikube
	@echo "Ensuring localenv chart is up to date..."
	helm list | grep -q localenv && helm upgrade localenv manifests/localenv --set postgresql.enabled=$(POSTGRES_ENABLED) || helm install localenv manifests/localenv --set postgresql.enabled=$(POSTGRES_ENABLED)
	@echo "Force restarting infrastructure pods to start fresh..."
	kubectl delete pod -l app.kubernetes.io/name=bq --ignore-not-found=true
	kubectl delete pod -l app.kubernetes.io/name=pubsub --ignore-not-found=true
	kubectl delete pod -l app.kubernetes.io/name=vault --ignore-not-found=true
	kubectl delete pod -l app.kubernetes.io/name=vault-agent-injector --ignore-not-found=true
	kubectl delete pod -l app.kubernetes.io/name=postgresql --ignore-not-found=true
	@echo "Waiting for infrastructure pods to be ready..."
	kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=bq --timeout=300s
	kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=pubsub --timeout=300s
	kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=vault --timeout=300s
	kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=vault-agent-injector --timeout=300s
	if [ "$(POSTGRES_ENABLED)" = "true" ]; then
		kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=postgresql --timeout=300s
	fi
	@echo "Pods are ready"
	@echo "Setting up data warehouse tables..."
	make create-bigquery-emulator-tables
	make create-mysql-event-tables-minikube
	make enable-vault-transit

DEV_IMAGES := gcr.io/distroless/base docker.io/arigaio/atlas:latest
pull-dev-images:
	@echo "Checking and pulling images..."
	@for image in $(DEV_IMAGES); do \
		if docker image inspect $$image >/dev/null 2>&1; then \
			echo "Image $$image already exists. Skipping..."; \
		else \
			echo "Pulling $$image"; \
			docker pull $$image; \
		fi; \
	done

force-pull-dev-images:
	@echo "Force pulling all images..."
	@for image in $(DEV_IMAGES); do \
		echo "Pulling $$image"; \
		docker pull $$image; \
	done

# build go application docker image
# please set the TAG env, eg: TAG=test make build-docker-images
build-docker-images:
	for APP in `ls bin`; do \
		./tools/build/show-dockerfile.sh bin $$APP > Dockerfile-app-$$APP; \
		IMAGE=`./tools/build/show-image-name.sh $$APP`; \
		docker build --platform $(PLATFORM) -f Dockerfile-app-$$APP -t ghcr.io/bucketeer-io/bucketeer-$$IMAGE:${TAG} .; \
		rm Dockerfile-app-$$APP; \
	done
	docker build --platform $(PLATFORM) migration/ -t ghcr.io/bucketeer-io/bucketeer-migration:${TAG}

# copy go application docker image to minikube
# please keep the same TAG env as used in build-docker-images, eg: TAG=test make minikube-load-images
minikube-load-images:
	for APP in $$(ls bin) migration; do \
		IMAGE=`./tools/build/show-image-name.sh $$APP`; \
		docker save ghcr.io/bucketeer-io/bucketeer-$$IMAGE:${TAG} -o $$IMAGE.tar; \
		docker cp $$IMAGE.tar minikube:/home/docker; \
		rm $$IMAGE.tar; \
		minikube ssh "sudo docker load -i /home/docker/$$IMAGE.tar"; \
		minikube ssh "rm /home/docker/$$IMAGE.tar"; \
	done

delete-bucketeer-from-minikube:
	helm uninstall bucketeer --ignore-not-found

# Bucketeer deployment
deploy-bucketeer: delete-bucketeer-from-minikube
	make -C tools/dev service-cert-secret
	make -C tools/dev service-token-secret
	make -C tools/dev oauth-key-secret
	make -C ./ build-go-embed
	make -C ./ pull-dev-images
	TAG=localenv make -C ./ build-docker-images
	TAG=localenv make -C ./ minikube-load-images
	kubectl exec localenv-mysql-0 -- bash -c "mysql -u root -pbucketeer -e 'SET GLOBAL log_bin_trust_function_creators = 1;'"
	@echo "Ensuring BigQuery tables exist (in-memory, may be lost on pod restart)..."
	make -C ./ create-bigquery-emulator-tables
	@echo "Installing Bucketeer services..."
	helm install bucketeer manifests/bucketeer/ --values manifests/bucketeer/values.dev.yaml

#############################
# Docker Compose
#############################

# Detect Docker Compose version and set variables
define detect_docker_compose
	$(eval DOCKER_COMPOSE_CMD := $(shell \
		if docker compose version >/dev/null 2>&1; then \
			echo "docker compose"; \
		elif docker-compose --version >/dev/null 2>&1; then \
			echo "docker-compose"; \
		else \
			echo ""; \
		fi \
	))
	$(eval COMPOSE_FILE := docker-compose/compose.yml)
	@if [ -z "$(DOCKER_COMPOSE_CMD)" ]; then \
		echo "❌ Error: Neither 'docker compose' (v2) nor 'docker-compose' (v1) found"; \
		echo "Please install Docker Compose: https://docs.docker.com/compose/install/"; \
		exit 1; \
	fi
	@echo "🐳 Using Docker Compose command: $(DOCKER_COMPOSE_CMD)"
endef

.PHONY: docker-compose-setup
docker-compose-setup:
	$(call detect_docker_compose)
	@echo "Setting up Docker Compose environment..."
	@if [ ! -d "docker-compose/init-db" ]; then \
		echo "Creating docker-compose/init-db directory..."; \
		mkdir -p docker-compose/init-db; \
	else \
		echo "docker-compose/init-db directory already exists"; \
	fi
	@if [ ! -d "docker-compose/secrets" ]; then \
		echo "Creating docker-compose/secrets directory..."; \
		mkdir -p docker-compose/secrets; \
		echo "Generating MySQL secret files..."; \
		echo "Generating PostgreSQL secret files..."; \
		echo "root" > docker-compose/secrets/mysql_root_password.txt; \
		echo "bucketeer" > docker-compose/secrets/mysql_password.txt; \
		echo "bucketeer" > docker-compose/secrets/postgres_password.txt; \
		chmod 600 docker-compose/secrets/*.txt; \
		echo "MySQL secrets created"; \
		echo "PostgresQL secrets created"; \
	else \
		echo "docker-compose/secrets directory already exists"; \
	fi
	@echo "Docker Compose setup complete"

.PHONY: docker-compose-init-env
docker-compose-init-env:
	@if [ -f docker-compose/.env ]; then \
		echo "docker-compose/.env already exists. Skipping..."; \
		echo "To recreate it, run: rm docker-compose/.env && make docker-compose-init-env"; \
	else \
		echo "Creating docker-compose/.env from template..."; \
		cp docker-compose/env.default docker-compose/.env; \
		echo "Created docker-compose/.env"; \
		echo ""; \
		echo "You can now customize the environment variables in docker-compose/.env"; \
		echo "Then run 'make docker-compose-up' to start the services."; \
	fi

.PHONY: docker-compose-build
docker-compose-build:
	@echo "🔨 Building Bucketeer Docker images..."
	@echo "Building Go applications with embedded web console..."
	GOOS=linux make -C ./ build-go-embed
	@echo "Building Docker images with TAG=localenv..."
	TAG=localenv make -C ./ build-docker-images
	@echo "Docker images built successfully"

# To skip the build step when starting services, run:
# make docker-compose-up SKIP_BUILD=true
.PHONY: docker-compose-up
docker-compose-up: docker-compose-setup docker-compose-init-env
	$(call detect_docker_compose)
	@if [ "$(SKIP_BUILD)" = "true" ]; then \
		echo "⏭️  Skipping build step as requested (SKIP_BUILD=true)."; \
	else \
		make docker-compose-build; \
	fi
	@echo "🚀 Starting Bucketeer services with Docker Compose..."
	@set -a && . docker-compose/.env && set +a && \
	$(DOCKER_COMPOSE_CMD) -f $(COMPOSE_FILE) up -d


.PHONY: docker-compose-down
docker-compose-down:
	$(call detect_docker_compose)
	@echo "Stopping Bucketeer services..."
	$(DOCKER_COMPOSE_CMD) -f $(COMPOSE_FILE) down

.PHONY: docker-compose-logs
docker-compose-logs:
	$(call detect_docker_compose)
	$(DOCKER_COMPOSE_CMD) -f $(COMPOSE_FILE) logs -f

.PHONY: docker-compose-status
docker-compose-status:
	$(call detect_docker_compose)
	$(DOCKER_COMPOSE_CMD) -f $(COMPOSE_FILE) ps

.PHONY: docker-compose-clean
docker-compose-clean:
	$(call detect_docker_compose)
	@echo "Stopping and removing all containers, networks, and volumes..."
	$(DOCKER_COMPOSE_CMD) -f $(COMPOSE_FILE) down -v
	docker system prune -f

.PHONY: docker-compose-regenerate-secrets
docker-compose-regenerate-secrets:
	@echo "Regenerating MySQL and Postgres secrets..."
	@rm -rf docker-compose/secrets
	@mkdir -p docker-compose/secrets
	@echo "root" > docker-compose/secrets/mysql_root_password.txt
	@echo "bucketeer" > docker-compose/secrets/mysql_password.txt
	@echo "bucketeer" > docker-compose/secrets/postgres_password.txt
	@chmod 600 docker-compose/secrets/*.txt
	@echo "MySQL and Postgres secrets regenerated"

.PHONY: docker-compose-delete-data
docker-compose-delete-data:
	@echo "Deleting E2E test data from Docker Compose MySQL..."
	MYSQL_USER=bucketeer \
	MYSQL_PASS=bucketeer \
	MYSQL_HOST=localhost \
	MYSQL_PORT=3306 \
	MYSQL_DB_NAME=bucketeer \
	make -C ./ delete-e2e-data-mysql

.PHONY: docker-compose-create-mysql-event-tables
docker-compose-create-mysql-event-tables:
	@echo "Creating MySQL event tables for Docker Compose data warehouse..."
	MYSQL_USER=bucketeer \
	MYSQL_PASS=bucketeer \
	MYSQL_HOST=localhost \
	MYSQL_PORT=3306 \
	MYSQL_DB_NAME=bucketeer \
	make -C ./ create-mysql-event-tables

.PHONY: docker-compose-create-postgres-event-tables
docker-compose-create-postgres-event-tables:
	@echo "Creating Postgres event tables for Docker Compose data warehouse..."
	POSTGRES_USER=bucketeer \
	POSTGRES_PASS=bucketeer \
	POSTGRES_HOST=localhost \
	POSTGRES_PORT=5432 \
	POSTGRES_DB_NAME=bucketeer \
	make -C ./ create-postgres-event-tables

.PHONY: docker-compose-delete-mysql-data-warehouse-data
docker-compose-delete-mysql-data-warehouse-data:
ifeq ($(GOOS), darwin)
	make -C hack/delete-mysql-data-warehouse clean build-darwin
else
	make -C hack/delete-mysql-data-warehouse clean build
endif
	./hack/delete-mysql-data-warehouse/delete-mysql-data-warehouse truncate \
		--mysql-user=bucketeer \
		--mysql-pass=bucketeer \
		--mysql-host=localhost \
		--mysql-port=3306 \
		--mysql-db-name=bucketeer \
		--no-profile \
		--no-gcp-trace-enabled
