#############################
# Variables
#############################

# set output directory to CI cache path via environment variables
BZFLAGS =
BUILD_FLAGS =
ifdef IS_CI
	BAZEL_OUTPUT_BASE ?= ../bazel-cache
	BZLFLAGS += --output_base ${BAZEL_OUTPUT_BASE}
	BUILD_FLAGS += --action_env=DOCKER_HOST --remote_cache=${BAZEL_REMOTE_CACHE} --google_credentials=${BAZEL_REMOTE_CACHE_CREDENTIALS}
endif

DELETED_PACKAGES := //proto/external/googleapis/googleapis/83e756a66b80b072bd234abcfe89edf459090974/google/rpc,//proto/external/protocolbuffers/protobuf/v3.18.1/google/protobuf
LOCAL_IMPORT_PATH := github.com/bucketeer-io/bucketeer

#############################
# All
#############################

.PHONY: shutdown
shutdown:
	bazelisk shutdown

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
all: gazelle gofmt-check proto-check update-repos-check lint build test

# protoc-gen-go should be same version as https://github.com/bazelbuild/rules_go/blob/master/go/private/repositories.bzl
.PHONY: local-deps
local-deps:
	mkdir -p ~/go-tools; \
	cd ~/go-tools; \
	if [ ! -e go.mod ]; then go mod init go-tools; fi; \
	go install golang.org/x/tools/cmd/goimports@latest; \
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.2; \
	go install github.com/golang/mock/mockgen@v1.6.0; \
	go install github.com/golang/protobuf/protoc-gen-go@v1.5.2; \
	go install github.com/nilslice/protolock/...@v0.15.0; \
	go get github.com/googleapis/googleapis;

.PHONY: gazelle
gazelle: proto-go
	bazelisk run ${BUILD_FLAGS} //:gazelle

.PHONY: gazelle-check
gazelle-check:
	bazelisk run //:gazelle -- -mode diff

.PHONY: lint
lint:
	golangci-lint run --timeout 3m0s ./cmd/... ./pkg/... ./hack/... ./test/...

.PHONY: build
build:
	bazelisk ${BZLFLAGS} build ${BUILD_FLAGS} --deleted_packages=${DELETED_PACKAGES} --workspace_status_command=$${PWD}/tools/build/status.sh \
		-k -- //cmd/... //pkg/... //proto/... //hack/... //test/...

.PHONY: test
test:
	bazelisk ${BZLFLAGS} test ${BUILD_FLAGS} -- //cmd/... //pkg/... //hack/...

.PHONY: gofmt
gofmt:
	goimports -local ${LOCAL_IMPORT_PATH} -w \
		$$(find . -path "./vendor" -prune -o -path "./proto" -prune -o -path "./bazel-cache" -prune -o -path "./bazel-proto" -prune -o  -type f -name '*.go' -print)

.PHONY: gofmt-check
gofmt-check:
	test -z "$$(goimports -local ${LOCAL_IMPORT_PATH} -d \
		$$(find . -path "./vendor" -prune -o -path "./proto" -prune -o -path "./bazel-cache" -prune -o -path "./bazel-proto" -prune -o  -type f -name '*.go' -print))"

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

.PHONY: proto-go
proto-go:
	make -C proto go

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
	bazelisk run ${BUILD_FLAGS} //:gazelle -- update-repos -from_file=go.mod -to_macro=repositories.bzl%go_repositories -prune=true

.PHONY: update-repos-check
update-repos-check: update-repos diff-check

.PHONY: diff-check
diff-check:
	test -z "$$(git diff --name-only)"

.PHONY: tidy-deps
tidy-deps:
	go mod tidy

#############################
# UI/WEB
#############################

.PHONY: build-ui-web-v2
build-ui-web-v2:
	bazelisk ${BZLFLAGS} build ${BUILD_FLAGS} -k --action_env=RELEASE_CHANNEL=prod -- //ui/web-v2:bundle

.PHONY: build-ui-web-v2-image
build-ui-web-v2-image:
	bazelisk ${BZLFLAGS} run ${BUILD_FLAGS} -k --action_env=RELEASE_CHANNEL=prod -- //ui/web-v2:bundle-image

#############################
# Charts
#############################

.PHONY: build-chart
build-chart: VERSION ?= $(shell git describe --tags --always --dirty --abbrev=7)
build-chart:
	mkdir -p .artifacts
	helm package manifests/bucketeer --version $(VERSION) --app-version $(VERSION) --dependency-update --destination .artifacts

#############################
# Dev tool
#############################

.PHONY: buildifier
buildifier-fix:
	bazelisk run //:buildifier-fix

.PHONY: buildifier-check
buildifier-check:
	bazelisk run //:buildifier-check

#############################
# E2E for backend
#############################

.PHONY: delete-e2e-data-mysql
delete-e2e-data-mysql:
	bazelisk ${BZLFLAGS} run ${BUILD_FLAGS} //hack/delete-e2e-data-mysql:delete-e2e-data-mysql -- delete \
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
	bazelisk ${BZLFLAGS} run ${BUILD_FLAGS} //hack/generate-service-token:generate-service-token -- generate \
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
	bazelisk ${BZLFLAGS} run ${BUILD_FLAGS} //hack/create-api-key:create-api-key -- create \
		--cert=${WEB_GATEWAY_CERT_PATH} \
		--web-gateway=${WEB_GATEWAY_URL}:443 \
		--service-token=${SERVICE_TOKEN_PATH} \
		--name=$$(date +%s) \
		--role=SDK \
		--output=${API_KEY_PATH} \
		--environment-namespace=${ENVIRONMENT_NAMESPACE} \
		--no-profile \
		--no-gcp-trace-enabled

.PHONY: e2e-l4
e2e-l4:
	bazelisk ${BZLFLAGS} test ${BUILD_FLAGS} \
		--cache_test_results=no \
		--test_output=all \
		--test_timeout=500 \
		--verbose_test_summary \
		--test_arg=--web-gateway-addr=${WEB_GATEWAY_URL} \
		--test_arg=--web-gateway-port=443 \
		--test_arg=--web-gateway-cert=${WEB_GATEWAY_CERT_PATH} \
		--test_arg=--api-key=${API_KEY_PATH} \
		--test_arg=--gateway-addr=${GATEWAY_URL} \
		--test_arg=--gateway-port=9000 \
		--test_arg=--gateway-cert=${GATEWAY_CERT_PATH} \
		--test_arg=--service-token=${SERVICE_TOKEN_PATH} \
		--test_arg=--environment-namespace=${ENVIRONMENT_NAMESPACE} \
		--test_arg=--test-id=${TEST_ID} \
		//test/e2e/autoops:go_default_test //test/e2e/environment:go_default_test //test/e2e/feature:go_default_test //test/e2e/experiment:go_default_test //test/e2e/gateway:go_default_test //test/e2e/eventcounter:go_default_test //test/e2e/user:go_default_test //test/e2e/push:go_default_test //test/e2e/notification:go_default_test

.PHONY: e2e
e2e:
	bazelisk ${BZLFLAGS} test ${BUILD_FLAGS} \
		--cache_test_results=no \
		--test_output=all \
		--test_timeout=500 \
		--verbose_test_summary \
		--test_arg=--web-gateway-addr=${WEB_GATEWAY_URL} \
		--test_arg=--web-gateway-port=443 \
		--test_arg=--web-gateway-cert=${WEB_GATEWAY_CERT_PATH} \
		--test_arg=--api-key=${API_KEY_PATH} \
		--test_arg=--gateway-addr=${GATEWAY_URL} \
		--test_arg=--gateway-port=443 \
		--test_arg=--gateway-cert=${GATEWAY_CERT_PATH} \
		--test_arg=--service-token=${SERVICE_TOKEN_PATH} \
		--test_arg=--environment-namespace=${ENVIRONMENT_NAMESPACE} \
		--test_arg=--test-id=${TEST_ID} \
		//test/e2e/autoops:go_default_test //test/e2e/environment:go_default_test //test/e2e/feature:go_default_test //test/e2e/experiment:go_default_test //test/e2e/gateway:go_default_test //test/e2e/eventcounter:go_default_test //test/e2e/user:go_default_test //test/e2e/push:go_default_test //test/e2e/notification:go_default_test

#############################
# Chores
#############################

.PHONY: docker-gen
docker-gen:
	rm -fr bazel-proto
	cp -r $$(bazel info | grep bazel-bin | sed -E 's/bazel-bin: (.+)/\1/')/proto bazel-proto
	docker run -it --rm \
		-v ${PWD}:/go/src/github.com/bucketeer-io/bucketeer \
		-w /go/src/github.com/bucketeer-io/bucketeer \
		--env DIR=/go/src/github.com/bucketeer-io/bucketeer \
		--env DESCRIPTOR_PATH=/go/src/github.com/bucketeer-io/bucketeer/bazel-proto \
		ghcr.io/bucketeer-io/bucketeer-runner:0.1.0 \
		bash tools/gen/gen.sh

.PHONY: remove-bazel-output
remove-bazel-output:
	bazelisk clean --expunge

.PHONY: disable-expired-trial-projects
disable-expired-trial-projects:
	bazelisk ${BZLFLAGS} run ${BUILD_FLAGS} //hack/disable-expired-trial-projects:disable-expired-trial-projects -- disable \
		--cert=${WEB_GATEWAY_CERT_PATH} \
		--service-token=${SERVICE_TOKEN_PATH} \
		--web-gateway=${WEB_GATEWAY_ADDR}:443 \
		--no-profile \
		--no-gcp-trace-enabled

.PHONY: delete-user-data
delete-user-data:
	bazelisk ${BZLFLAGS} run ${BUILD_FLAGS} //hack/delete-user-data:delete-user-data -- delete \
		--mysql-host=${MYSQL_HOST} \
		--mysql-port=${MYSQL_PORT} \
		--mysql-user=${MYSQL_USER} \
		--mysql-pass=${MYSQL_PASS} \
		--mysql-db-name=${MYSQL_DB_NAME} \
		--target-period=${TARGET_PERIOD} \
		--no-profile \
		--no-gcp-trace-enabled
