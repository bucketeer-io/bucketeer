#############################
# Variables
#############################

LOCAL_IMPORT_PATH := github.com/bucketeer-io/bucketeer

# go applications
GO_APP_DIRS := $(wildcard cmd/*)
GO_APP_BUILD_TARGETS := $(addprefix build-,$(notdir $(GO_APP_DIRS)))

ifndef GOOS
	GOOS := $(shell go env GOOS)
endif

ifndef GOARCH
	GOARCH := $(shell go env GOARCH)
endif

LDFLAGS_PACKAGE := github.com/bucketeer-io/bucketeer/pkg/ldflags
LDFLAGS_VERSION := $(LDFLAGS_PACKAGE).Version
LDFLAGS_HASH := $(LDFLAGS_PACKAGE).Hash
LDFLAGS_BUILDDATE := $(LDFLAGS_PACKAGE).BuildDate

#############################
# Run make commands on docker container
#############################

# E.g. make docker-run CMD=proto-go
RUNNER_IMAGE = ghcr.io/bucketeer-io/bucketeer-runner:0.1.0
DOCKER_REPO_PATH = /go/src/github.com/bucketeer-io/bucketeer
DOCKER_RUN_CMD = docker run -it --rm -v ${PWD}:${DOCKER_REPO_PATH} -w ${DOCKER_REPO_PATH} ${RUNNER_IMAGE}
.PHONY: docker-run
docker-run:
	eval ${DOCKER_RUN_CMD} make $$CMD

#############################
# Go
#############################

.PHONY: all
all: gofmt-check proto-check update-repos-check lint build-go test-go

.PHONY: local-deps
local-deps:
	mkdir -p ~/go-tools; \
	cd ~/go-tools; \
	if [ ! -e go.mod ]; then go mod init go-tools; fi; \
	go install golang.org/x/tools/cmd/goimports@latest; \
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.2; \
	go install go.uber.org/mock/mockgen@v0.1.0; \
	go install github.com/golang/protobuf/protoc-gen-go@v1.5.2; \
	go install github.com/nilslice/protolock/...@v0.15.0;
	go install github.com/mikefarah/yq/v4@v4.28.2

.PHONY: lint
lint:
	golangci-lint run --timeout 3m0s ./cmd/... ./pkg/... ./hack/... ./test/...

.PHONY: gofmt
gofmt:
	goimports -local ${LOCAL_IMPORT_PATH} -w \
		$$(find . -path "./vendor" -prune -o -path "./proto" -prune -o -type f -name '*.go' -print)

.PHONY: gofmt-check
gofmt-check:
	test -z "$$(goimports -local ${LOCAL_IMPORT_PATH} -d \
		$$(find . -path "./vendor" -prune -o -path "./proto" -prune -o -type f -name '*.go' -print))"

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
proto-all:
	make -C proto go
	make -C ui/web-v2 gen_proto

.PHONY: proto-go
proto-go:
	make -C proto go

.PHONY: proto-go-check
proto-go-check:
	make -C proto go-check

.PHONY: proto-go-descriptor
proto-go-descriptor:
	make -C proto go-descriptor

.PHONY: proto-go-descriptor-check
proto-go-descriptor-check:
	make -C proto go-descriptor-check

.PHONY: mockgen
mockgen: proto-go
	find ./pkg -path "*mock*.go" -type f -delete
	go generate -run="mockgen" ./pkg/...
	make gofmt

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

.PHONY: build-go
build-go: $(GO_APP_BUILD_TARGETS)

# Make sure bucketeer-httpstan is already running. If not, run "make run-httpstan".
.PHONY: test-go
test-go:
	TZ=UTC CGO_ENABLED=0 go test -v ./pkg/...

.PHONY: start-httpstan
start-httpstan:
	docker run --name bucketeer-httpstan -p 8080:8080 -d ghcr.io/bucketeer-io/bucketeer-httpstan:0.0.1

.PHONY: stop-httpstan
stop-httpstan:
	docker stop bucketeer-httpstan

#############################
# UI/WEB
#############################

.PHONY: build-ui-web-v2
build-ui-web-v2:
	make -C ui/web-v2 build

.PHONY: build-ui-web-v2-prod
build-ui-web-v2-prod:
	RELEASE_CHANNEL=prod make -C ui/web-v2 build

#############################
# Charts
#############################

.PHONY: build-chart
build-chart: VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7)
build-chart:
	mkdir -p .artifacts
	helm package manifests/bucketeer --version $(VERSION) --app-version $(VERSION) --dependency-update --destination .artifacts

#############################
# E2E for backend
#############################

.PHONY: delete-e2e-data-mysql
delete-e2e-data-mysql:
	go run ./hack/delete-e2e-data-mysql delete \
		--mysql-user=${MYSQL_USER} \
		--mysql-pass=${MYSQL_PASS} \
		--mysql-host=mysql-${ENV}.bucketeer.private \
		--mysql-port=3306 \
		--mysql-db-name=master \
		--test-id=${TEST_ID} \
		--no-profile \
		--no-gcp-trace-enabled

.PHONY: generate-service-token
generate-service-token:
	go run ./hack/generate-service-token generate \
		--issuer=${ISSUER} \
		--sub=service \
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
		--name=$$(date +%s) \
		--role=SDK_CLIENT \
		--output=${API_KEY_PATH} \
		--environment-namespace=${ENVIRONMENT_NAMESPACE} \
		--no-profile \
		--no-gcp-trace-enabled

.PHONY: e2e-l4
e2e-l4:
	go test -v ./test/e2e/... -args \
		-web-gateway-addr=${WEB_GATEWAY_URL} \
		-web-gateway-port=443 \
		-web-gateway-cert=${WEB_GATEWAY_CERT_PATH} \
		-api-key=${API_KEY_PATH} \
		-gateway-addr=${GATEWAY_URL} \
		-gateway-port=9000 \
		-gateway-cert=${GATEWAY_CERT_PATH} \
		-service-token=${SERVICE_TOKEN_PATH} \
		-environment-namespace=${ENVIRONMENT_NAMESPACE} \
		-test-id=${TEST_ID}

.PHONY: e2e
e2e:
	go test -v ./test/e2e/... -args \
		-web-gateway-addr=${WEB_GATEWAY_URL} \
		-web-gateway-port=443 \
		-web-gateway-cert=${WEB_GATEWAY_CERT_PATH} \
		-api-key=${API_KEY_PATH} \
		-gateway-addr=${GATEWAY_URL} \
		-gateway-port=443 \
		-gateway-cert=${GATEWAY_CERT_PATH} \
		-service-token=${SERVICE_TOKEN_PATH} \
		-environment-namespace=${ENVIRONMENT_NAMESPACE} \
		-test-id=${TEST_ID}

.PHONY: update-copyright
update-copyright:
	./hack/update-copyright/update-copyright.sh


#############################
# dev container
#############################

# start minikube
start-minikube: 
	if [ $$(minikube status | grep -c "minikube start") -eq 1 ]; then \
		make -C tools/dev setup-minikube; \
	elif [ $$(minikube status | grep -c "Stopped") -gt 1 ]; then \
		make -C tools/dev start-minikube; \
	fi
	sleep 5
	make -C ./ modify-hosts
	make -C ./ setup-bigquery-vault

# modify hosts file to access api-gateway and web-gateway
modify-hosts:
	$(eval MINIKUBE_IP := $(shell minikube ip))
	echo "$(MINIKUBE_IP)   web-gateway.bucketeer.org" | sudo tee -a /etc/hosts
	echo "$(MINIKUBE_IP)   api-gateway.bucketeer.org" | sudo tee -a /etc/hosts

# enable vault transit secret engine
enable-vault-transit:
	kubectl exec localenv-vault-0 -- vault secrets enable transit

# create bigquery-emulator tables
create-bigquery-emulator-tables:
	go run ./hack/create-big-query-table create \
		--bigquery-emulator=http://$$(minikube ip):31000 \
		--no-gcp-trace-enabled \
		--no-profile

setup-bigquery-vault:
	while [ "$$(kubectl get pods | grep localenv-bq | awk '{print $$3}')" != "Running" ] || [ "$$(kubectl get pods | grep localenv-vault-0 | awk '{print $$3}')" != "Running" ]; \
	do \
		sleep 5; \
	done; \
	make create-bigquery-emulator-tables
	make enable-vault-transit

# generate tls certificate
generate-tls-certificate:
	make -C tools/dev generate-tls-certificate

# generate oauth key
generate-oauth:
	make -C tools/dev generate-oauth

# create service cert secret in minikube
service-cert-secret:
	make -C tools/dev service-cert-secret

# create oauth key secret in minikube
oauth-key-secret:
	make -C tools/dev oauth-key-secret

# create github token secret in minikube
generate-github-token:
	make -C tools/dev generate-github-token


# build go application docker image
# please set the TAG env, eg: TAG=test make build-docker-images
build-docker-images:
	for APP in `ls bin`; do \
		./tools/build/show-dockerfile.sh bin $$APP > Dockerfile-app-$$APP; \
		IMAGE=`./tools/build/show-image-name.sh $$APP`; \
		docker build -f Dockerfile-app-$$APP -t ghcr.io/bucketeer-io/bucketeer-$$IMAGE:${TAG} .; \
		rm Dockerfile-app-$$APP; \
	done


# copy go application docker image to minikube
# please keep the same TAG env as used in build-docker-images, eg: TAG=test make minikube-load-images
minikube-load-images:
	for APP in `ls bin`; do \
		IMAGE=`./tools/build/show-image-name.sh $$APP`; \
		docker save  ghcr.io/bucketeer-io/bucketeer-$$IMAGE:${TAG} -o $$IMAGE.tar; \
		docker cp $$IMAGE.tar minikube:/home/docker; \
		rm $$IMAGE.tar; \
		minikube ssh "docker load -i /home/docker/$$IMAGE.tar"; \
		minikube ssh "rm /home/docker/$$IMAGE.tar"; \
	done

SERVICES := api-gateway auditlog-persister backend batch event-persister-evaluation-events-dwh event-persister-evaluation-events-evaluation-count event-persister-evaluation-events-ops event-persister-goal-events-dwh event-persister-goal-events-ops experiment-calculator metrics-event-persister push-sender user-persister web-gateway web dex

# Deploy Bucketeer to minikube
deploy-service-to-minikube:
	helm install ${SERVICE} manifests/bucketeer/charts/${SERVICE}/ --values manifests/bucketeer/charts/${SERVICE}/values.dev.yaml \
	--set serviceToken.token=$$(cat tools/dev/cert/service-token)

# Delete all the services from Minikube
delete-all-services-from-minikube:
	$(foreach var,$(SERVICES),helm uninstall $(var);)

# Deploy All the services to minikube
deploy-all-services-to-minikube:
	$(foreach var,$(SERVICES),SERVICE=$(var) make deploy-service-to-minikube;)

# bucketeer deploy
deploy-bucketeer:
	make -C ./ generate-tls-certificate
	make -C ./ generate-oauth
	make -C ./ service-cert-secret
	make -C ./ oauth-key-secret
	GITHUB_TOKEN=$(GITHUB_TOKEN) make -C ./ generate-github-token
	ISSUER=$(ISSUER) \
	EMAIL=$(EMAIL) \
	OAUTH_KEY_PATH=/workspaces/bucketeer/tools/dev/cert/oauth-private.pem \
	SERVICE_TOKEN_PATH=/workspaces/bucketeer/tools/dev/cert/service-token \
	make generate-service-token
	make -C ./ build-go
	TAG=$(TAG) make -C ./ build-docker-images
	TAG=$(TAG) make -C ./ minikube-load-images
	make -C ./ deploy-all-services-to-minikube