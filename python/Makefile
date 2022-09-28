PY_FILES := $$(find . -iname '*.py' -not -path "./.venv/*")
GIT_TOP_DIR := $(shell git rev-parse --show-toplevel)
PROTOBUF_INCLUDE_DIR := $(GIT_TOP_DIR)/proto/external/protocolbuffers/protobuf/v3.9.0
PROTO_FOLDERS := $(filter-out $(GIT_TOP_DIR)/proto/external%, $(shell find $(GIT_TOP_DIR)/proto -name '*.proto' -print0 | xargs -0 -n1 dirname | sort --unique))

.PHONY: init
init:
	poetry install

.PHONY: test
test:
	PYTHONPATH=./src poetry run pytest tests

.PHONY: update-deps
update-deps:
	poetry lock
	poetry export -f requirements.txt --output requirements.txt --without-hashes
	poetry export -f requirements.txt --output requirements-dev.txt --without-hashes --dev

.PHONY: fmt
fmt:
	poetry run black ${PY_FILES}

.PHONY: fmt-check
fmt-check:
	poetry run black --check ${PY_FILES}

.PHONY: gen
gen:
	rm -fr ./src/proto
	for f in ${PROTO_FOLDERS}; do \
		poetry run python -m grpc_tools.protoc \
			-I"$(GIT_TOP_DIR)" \
			-I"$(PROTOBUF_INCLUDE_DIR)" \
			-I"${GOPATH}/src/github.com/googleapis/googleapis" \
			--python_out=$(GIT_TOP_DIR)/python/src \
			$$f/*.proto; \
		poetry run python -m grpc_tools.protoc \
			-I"$(GIT_TOP_DIR)" \
			-I"$(PROTOBUF_INCLUDE_DIR)" \
			-I"${GOPATH}/src/github.com/googleapis/googleapis" \
			--grpc_python_out=$(GIT_TOP_DIR)/python/src \
			$$f/*service.proto; \
	done \

.PHONY: docker-build-dev
docker-build-dev:
	docker build -t "bucketeer-python-dev:latest" -f Dockerfile.dev .