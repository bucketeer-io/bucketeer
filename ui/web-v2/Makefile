.PHONY: all
all: install lint gen_proto build

####################################
# Yarn
####################################

.PHONY: install
install:
	yarn install

.PHONY: build
build:
	yarn build

.PHONY: start
start:
	yarn start

.PHONY: lint
lint:
	yarn lint

####################################
# Generate proto definition files
####################################

ROOT_DIR := $(CURDIR)/../..
SRC_DIR := $(CURDIR)/apps/admin/src
PROTOBUF_INCLUDE_DIR := $(ROOT_DIR)/proto/external/protocolbuffers/protobuf/v23.4
GOOGLEAPIS := $(ROOT_DIR)/proto/external/googleapis/googleapis/83e756a66b80b072bd234abcfe89edf459090974
NPM_BIN_DIR := $(CURDIR)/node_modules/.bin

.PHONY: gen_proto
gen_proto: clean_proto gen_google_api
	protoc \
	--plugin=protoc-gen-ts=$(NPM_BIN_DIR)/protoc-gen-ts \
	--plugin=protoc-gen-js=$(NPM_BIN_DIR)/protoc-gen-js \
	--js_out=import_style=commonjs,binary:$(SRC_DIR) \
	--ts_out=service=grpc-web:$(SRC_DIR) \
	-I $(ROOT_DIR) \
	-I $(PROTOBUF_INCLUDE_DIR) \
	-I $(GOOGLEAPIS) \
	$(shell find $(ROOT_DIR)/proto -type f -name "*.proto" -not -path "**/gateway/*.proto" -not -path "**/google/protobuf/*.proto" -not -path "**/google/api/*.proto")

.PHONY: gen_google_api
gen_google_api:
	protoc \
	--plugin=protoc-gen-ts=$(NPM_BIN_DIR)/protoc-gen-ts \
	--plugin=protoc-gen-js=$(NPM_BIN_DIR)/protoc-gen-js \
	--js_out=import_style=commonjs,binary:$(SRC_DIR) \
	--ts_out=service=grpc-web:$(SRC_DIR) \
	-I $(PROTOBUF_INCLUDE_DIR) \
	-I $(GOOGLEAPIS) \
	$(shell find $(GOOGLEAPIS) -type f -name "*.proto")

.PHONY: clean_proto
clean_proto:
	rm -rf $(SRC_DIR)/proto/**/*.{js,ts}
	rm -rf $(SRC_DIR)/google/**/*.{js,ts}

.PHONY: gen_proto_check
gen_proto_check: gen_proto
	test -z "$$(git diff --name-only | grep ".*pb_service\.d\.ts\|.*pb\.d\.ts\|.*pb_service\.js\|.*\.pb\.js")"
