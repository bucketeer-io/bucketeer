PROTO_FOLDERS := $(filter-out ./external%, $(shell find . -name '*.proto' -print0 | xargs -0 -n1 dirname | sort --unique))
GIT_TOP_DIR := $(shell cd .. ; pwd)
PROTOBUF_INCLUDE_DIR := ./external/protocolbuffers/protobuf/v3.18.1
GOOGLEAPIS := ./external/googleapis/googleapis/83e756a66b80b072bd234abcfe89edf459090974

.PHONY: go
go: remove-go
	for f in ${PROTO_FOLDERS}; do \
		protoc \
			--proto_path=${GIT_TOP_DIR} \
			--proto_path=${PROTOBUF_INCLUDE_DIR} \
			--proto_path=${GOOGLEAPIS} \
			--go_out=plugins=grpc:${GIT_TOP_DIR} \
			--go_opt=paths=source_relative \
			${GIT_TOP_DIR}/proto/$$f/*.proto; \
	done

.PHONY: remove-go
remove-go:
	find . -name "*.pb.go" -type f -delete

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
