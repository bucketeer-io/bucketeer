PROTO_FOLDERS := $(filter-out ./external%, $(shell find . -name '*.proto' -print0 | xargs -0 -n1 dirname | sort --unique))
GIT_TOP_DIR := $(shell realpath ..)
PROTOBUF_INCLUDE_DIR := $(GIT_TOP_DIR)/proto/external/protocolbuffers/protobuf/v3.18.1

.PHONY: go
go: remove-go
	for f in ${PROTO_FOLDERS}; do \
		protoc -I"$(GIT_TOP_DIR)" \
			-I"$(PROTOBUF_INCLUDE_DIR)" \
			-I"${GOPATH}/src/github.com/googleapis/googleapis" \
			--go_out=plugins=grpc:${GOPATH}/src \
			$(GIT_TOP_DIR)/proto/$$f/*.proto; \
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
