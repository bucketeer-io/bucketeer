PROTO_FOLDERS := $(filter-out ./external%, $(shell find . -name '*.proto' -print0 | xargs -0 -n1 dirname | sort --unique))
PROTO_WITH_DESCRIPTOR_FOLDERS := ./account ./auditlog ./auth ./autoops ./environment ./eventcounter ./experiment ./feature ./gateway ./migration ./notification ./push ./user
PROTO_WITHOUT_DESCRIPTOR_FOLDERS := $(filter-out $(PROTO_WITH_DESCRIPTOR_FOLDERS), $(PROTO_FOLDERS))

GIT_TOP_DIR := $(shell cd .. ; pwd)
PROTOBUF_INCLUDE_DIR := ./external/protocolbuffers/protobuf/v3.18.1
GOOGLEAPIS := ./external/googleapis/googleapis/83e756a66b80b072bd234abcfe89edf459090974

.PHONY: go
go: remove-go
	for f in ${PROTO_WITHOUT_DESCRIPTOR_FOLDERS}; do \
		protoc \
			--proto_path=${GIT_TOP_DIR} \
			--proto_path=${PROTOBUF_INCLUDE_DIR} \
			--proto_path=${GOOGLEAPIS} \
			--go_out=plugins=grpc:${GIT_TOP_DIR} \
			--go_opt=paths=source_relative \
			${GIT_TOP_DIR}/proto/$$f/*.proto; \
	done; \
	for f in ${PROTO_WITH_DESCRIPTOR_FOLDERS}; do \
		protoc \
			--proto_path=${GIT_TOP_DIR} \
			--proto_path=${PROTOBUF_INCLUDE_DIR} \
			--proto_path=${GOOGLEAPIS} \
			--go_out=plugins=grpc:${GIT_TOP_DIR} \
			--go_opt=paths=source_relative \
			--include_imports \
			--descriptor_set_out=${GIT_TOP_DIR}/proto/$$f/proto_descriptor.pb \
			${GIT_TOP_DIR}/proto/$$f/*.proto; \
	done

.PHONY: remove-go
remove-go:
	find . -name "*.pb.go" -type f -delete; \
	find . -name "proto_descriptor.pb" -type f -delete

.PHONY: go-descriptor
go-descriptor:
	DIR=$(GIT_TOP_DIR) DESCRIPTOR_PATH=$(GIT_TOP_DIR)/proto bash $(GIT_TOP_DIR)/tools/gen/descriptor.sh

.PHONY: go-descriptor-check
go-descriptor-check: go-descriptor
	$(eval api_gateway_values_path := manifests/bucketeer/charts/api-gateway/values.yaml)
	$(eval web_gateway_values_path := manifests/bucketeer/charts/web-gateway/values.yaml)
	test -z "$$(git diff --name-only | grep -e "$(api_gateway_values_path)\|$(web_gateway_values_path)")"

.PHONY: check
check: fmt-check lock-check

.PHONY: fmt
fmt:
	find . -name "*.proto" | grep -v external | xargs clang-format -i

.PHONY: fmt-check
fmt-check:
	test -z "$$(find . -name "*.proto" | grep -v external | xargs clang-format -i -output-replacements-xml | grep "<replacement ")"

.PHONY: lock-init
lock-init:
	protolock init --ignore ./external

.PHONY: lock-commit
lock-commit:
	protolock commit --ignore ./external

# This can be used when conflict errors occur
# It will rewite the proto.lock with current tree
.PHONY: lock-commit-force
lock-commit-force:
	protolock commit --force --ignore ./external

.PHONY: lock-check
lock-check:
	protolock status --uptodate --ignore ./external
