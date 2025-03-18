PROTO_FOLDERS := $(filter-out ./external%, $(shell find . -name '*.proto' -not -path "./openapi/web_default_settings.proto" -print0 | xargs -0 -n1 dirname | sort --unique))
PROTO_WITH_DESCRIPTOR_FOLDERS := ./account ./auditlog ./auth ./autoops ./environment ./eventcounter ./experiment ./feature ./gateway ./notification ./push ./batch ./tag ./coderef
PROTO_WITHOUT_DESCRIPTOR_FOLDERS := $(filter-out $(PROTO_WITH_DESCRIPTOR_FOLDERS), $(PROTO_FOLDERS))

GIT_TOP_DIR := $(shell cd .. ; pwd)
PROTOBUF_INCLUDE_DIR := ./external/protocolbuffers/protobuf/v23.4
GOOGLEAPIS := ./external/googleapis/googleapis/83e756a66b80b072bd234abcfe89edf459090974
OPENAPI := ./external/openapiv2

.PHONY: go
go: remove-go
	for f in ${PROTO_WITHOUT_DESCRIPTOR_FOLDERS}; do \
		protoc \
			--proto_path=${GIT_TOP_DIR} \
			--proto_path=${PROTOBUF_INCLUDE_DIR} \
			--proto_path=${GOOGLEAPIS} \
			--proto_path=${OPENAPI} \
			--go_out=plugins=grpc:${GIT_TOP_DIR} \
			--go_opt=paths=source_relative \
			${GIT_TOP_DIR}/proto/$$f/*.proto; \
	done; \
	for f in ${PROTO_WITH_DESCRIPTOR_FOLDERS}; do \
		protoc \
			--proto_path=${GIT_TOP_DIR} \
			--proto_path=${PROTOBUF_INCLUDE_DIR} \
			--proto_path=${GOOGLEAPIS} \
			--proto_path=${OPENAPI} \
			--go_out=plugins=grpc:${GIT_TOP_DIR} \
			--go_opt=paths=source_relative \
			--include_imports \
			--descriptor_set_out=${GIT_TOP_DIR}/proto/$$f/proto_descriptor.pb \
			${GIT_TOP_DIR}/proto/$$f/*.proto; \
	done
	# Because there are old projects using empty in the `environment_id`
	# we must remove the `omitempty` option from the generated files.
	# Otherwise, it will omit the field when converting it to json when saving the data in MySQL.
	@for f in ${PROTO_FOLDERS}; do \
		if [ ${shell uname} = "Darwin" ]; then \
			sed -i '' 's/,omitempty//g' ${GIT_TOP_DIR}/proto/$$f/*.pb.go; \
		else \
			sed -i 's/,omitempty//g' ${GIT_TOP_DIR}/proto/$$f/*.pb.go; \
		fi \
	done

.PHONY: openapi-api-gen
openapi-api-gen:
	protoc \
		--proto_path=${GIT_TOP_DIR} \
		--proto_path=${PROTOBUF_INCLUDE_DIR} \
		--proto_path=${GOOGLEAPIS} \
		--proto_path=${OPENAPI} \
		--proto_path=${PROTOC_GEN_OPENAPIV2} \
		--openapiv2_out=${GIT_TOP_DIR}/api-description \
		--openapiv2_opt output_format=yaml \
		--openapiv2_opt allow_merge=true \
		--openapiv2_opt=disable_service_tags=true \
		${GIT_TOP_DIR}/proto/gateway/*.proto; \

PROTO_OPENAPI_WEB_TARGETS := ${GIT_TOP_DIR}/proto/openapi/web_default_settings.proto ${GIT_TOP_DIR}/proto/auth/service.proto ${GIT_TOP_DIR}/proto/environment/service.proto ${GIT_TOP_DIR}/proto/account/service.proto ${GIT_TOP_DIR}/proto/notification/service.proto ${GIT_TOP_DIR}/proto/push/service.proto ${GIT_TOP_DIR}/proto/feature/service.proto ${GIT_TOP_DIR}/proto/experiment/service.proto ${GIT_TOP_DIR}/proto/tag/service.proto ${GIT_TOP_DIR}/proto/coderef/service.proto ${GIT_TOP_DIR}/proto/autoops/service.proto ${GIT_TOP_DIR}/proto/eventcounter/service.proto ${GIT_TOP_DIR}/proto/auditlog/service.proto
.PHONY: openapi-web-gen
openapi-web-gen:
	protoc \
		--proto_path=${GIT_TOP_DIR} \
		--proto_path=${PROTOBUF_INCLUDE_DIR} \
		--proto_path=${GOOGLEAPIS} \
		--proto_path=${OPENAPI} \
		--proto_path=${PROTOC_GEN_OPENAPIV2} \
		--openapiv2_out=${GIT_TOP_DIR}/api-description \
		--openapiv2_opt output_format=yaml \
		--openapiv2_opt allow_merge=true \
		--openapiv2_opt merge_file_name=web-api \
		--openapiv2_opt=disable_service_tags=true \
		${PROTO_OPENAPI_WEB_TARGETS};

.PHONY: openapi-ui
openapi-api-ui:
	docker run -p 80:8080 -e SWAGGER_JSON=/openapi/apidocs.swagger.yaml -v ${GIT_TOP_DIR}/api-description/:/openapi swaggerapi/swagger-ui

go-check: go
	test -z "$$(git diff --name-only | grep ".*\.pb.go\|.*\.pb")"

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
