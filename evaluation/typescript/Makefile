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
	$(NPM_BIN_DIR)/eslint --fix

.PHONY: copy-genfiles
copy-genfiles:
	@rm -rf $(CURDIR)/lib
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
# Copy proto definition files
####################################

ROOT_DIR := $(CURDIR)/../..
WEB_UI_SRC_DIR := $(ROOT_DIR)/ui/web-v2/src
SRC_DIR := $(CURDIR)/src
NPM_BIN_DIR := $(CURDIR)/node_modules/.bin

.PHONY: copy_proto_generated_code
copy_proto_generated_code: clean_proto
	cp -r $(WEB_UI_SRC_DIR)/proto/feature $(SRC_DIR)/proto/
	cp -r $(WEB_UI_SRC_DIR)/proto/user $(SRC_DIR)/proto/	
	cp -r $(WEB_UI_SRC_DIR)/google $(SRC_DIR)/

.PHONY: clean_proto
clean_proto:
	rm -rf $(SRC_DIR)/proto/**/*.{js,ts}
	rm -rf $(SRC_DIR)/google/**/*.{js,ts}

.PHONY: gen_proto_check
gen_proto_check: gen_proto
	test -z "$$(git diff --name-only | grep ".*pb_service\.d\.ts\|.*pb\.d\.ts\|.*pb_service\.js\|.*\.pb\.js")"
