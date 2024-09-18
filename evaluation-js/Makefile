.PHONY: all
all: install lint copy_proto_generated_code build

####################################
# Yarn
####################################

.PHONY: install
install:
	yarn install

.PHONY: build
build:
	yarn build

.PHONY: lint
lint:
	yarn lint

####################################
# Generate proto definition files
####################################

ROOT_DIR := $(CURDIR)/..
WEB_UI_SRC_DIR := $(ROOT_DIR)/ui/web-v2/src
SRC_DIR := $(CURDIR)/src
NPM_BIN_DIR := $(CURDIR)/node_modules/.bin

.PHONY: copy_proto_generated_code
copy_proto_generated_code: clean_proto
	cp -r $(WEB_UI_SRC_DIR)/proto/feature $(SRC_DIR)/proto/feature
	cp -r $(WEB_UI_SRC_DIR)/proto/user $(SRC_DIR)/proto/user	
	cp -r $(WEB_UI_SRC_DIR)/google $(SRC_DIR)/google
	yarn lintfix

.PHONY: clean_proto
clean_proto:
	rm -rf $(SRC_DIR)/proto/**/*.{js,ts}
	rm -rf $(SRC_DIR)/google/**/*.{js,ts}

.PHONY: gen_proto_check
gen_proto_check: gen_proto
	test -z "$$(git diff --name-only | grep ".*pb_service\.d\.ts\|.*pb\.d\.ts\|.*pb_service\.js\|.*\.pb\.js")"
