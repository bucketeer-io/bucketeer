NPM_BIN_DIR := $(CURDIR)/node_modules/.bin
GENFILES_DIR := $(CURDIR)/__lib

GIT_REVISION := $(shell git rev-parse --verify HEAD)

export PACKAGE_NAME := $(shell node -p "require('./package.json').name")
export CURRENT_VERSION := $(shell npm view $(PACKAGE_NAME) version 2>/dev/null || echo 0.0.0)
export LOCAL_VERSION := $(shell node -p "require('./package.json').version")

.PHONY: init
init:
	yarn

.PHONY: clean
clean:
	rm -rf $(CURDIR)/__test $(CURDIR)/__lib

.PHONY: build
build: clean-build tsc rename-js

.PHONY: clean-build
clean-build:
	rm -rf $(GENFILES_DIR)

.PHONY: tsc
tsc: clean
	$(NPM_BIN_DIR)/tsc --project tsconfig.json
	cp -r $(SRC_DIR)/proto __lib/

.PHONY: build_test
build_test: clean
	$(NPM_BIN_DIR)/tsc --project tsconfig.test.json
	cp -r $(SRC_DIR)/proto __test/

.PHONY: rename-js
rename-js:
	$(NPM_BIN_DIR)/rename '$(GENFILES_DIR)/**/*.js' '{{f}}.mjs'

.PHONY: test
test: build_test
	$(NPM_BIN_DIR)/ava --config ava.config.mjs

.PHONY: fmt
fmt: tsfmt

.PHONY: tsfmt
tsfmt:
	$(NPM_BIN_DIR)/prettier --write '$(CURDIR)/src/**/*.ts'

.PHONY: lint
lint:
	$(NPM_BIN_DIR)/eslint

.PHONY: fixLintErrors
fixLintErrors:
	$(NPM_BIN_DIR)/eslint . --fix

.PHONY: copy-genfiles
copy-genfiles:
	@rm -rf $(CURDIR)/lib # Clear cache
	@mkdir -p $(CURDIR)/lib  # Ensure lib directory exists
	$(NPM_BIN_DIR)/cpx '$(GENFILES_DIR)/**/*.{mjs,d.ts}' $(CURDIR)/lib
	$(NPM_BIN_DIR)/cpx '$(GENFILES_DIR)/bucketeer.*' $(CURDIR)/lib
	@find $(CURDIR)/lib -type f -exec chmod 644 {} +
	$(NPM_BIN_DIR)/rename '$(CURDIR)/lib/**/*.js' '{{f}}.mjs'
	$(NPM_BIN_DIR)/babel lib --extensions '.mjs' --config-file "$(CURDIR)/babel.config.js" --out-dir "lib"

.PHONY: publish-dry
publish-dry: copy-genfiles
	npm publish --dry-run

.PHONY: publish
publish: copy-genfiles
ifeq ($(shell $(NPM_BIN_DIR)/semver -r ">$(CURRENT_VERSION)" $(LOCAL_VERSION) ),$(LOCAL_VERSION))
	npm publish --access public
else
	@echo "$(LOCAL_VERSION) exists. skip publish."
endif

####################################
# Generate proto definition files
####################################

ROOT_DIR := $(CURDIR)/../..
SRC_DIR := $(CURDIR)/src
PROTOBUF_INCLUDE_DIR := $(ROOT_DIR)/proto/external/protocolbuffers/protobuf/v23.4
NPM_BIN_DIR := $(CURDIR)/node_modules/.bin

PROTO_TARGETS := \
	$(ROOT_DIR)/proto/feature \
	$(ROOT_DIR)/proto/user \
	$(ROOT_DIR)/proto/event/client
	
.PHONY: gen_proto
gen_proto: clean_proto
	protoc \
	--plugin=protoc-gen-ts=$(NPM_BIN_DIR)/protoc-gen-ts \
	--plugin=protoc-gen-js=$(NPM_BIN_DIR)/protoc-gen-js \
	--js_out=import_style=commonjs,binary:$(SRC_DIR) \
	--ts_out=service=grpc-web:$(SRC_DIR) \
	-I $(ROOT_DIR) \
	-I $(PROTOBUF_INCLUDE_DIR) \
	$(shell find $(PROTO_TARGETS) -type f -name "*.proto" -not -path "*/feature/service.proto")

.PHONY: clean_proto
clean_proto:
	rm -rf $(SRC_DIR)/proto/