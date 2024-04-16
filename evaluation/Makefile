LOCAL_IMPORT_PATH := github.com/bucketeer-io/bucketeer/evaluation

.PHONY: update-deps
update-deps:
	go mod tidy && go mod vendor

.PHONY: update-deps-check
update-deps-check: update-deps diff-check

.PHONY: diff-check
diff-check:
	test -z "$$(git diff --name-only)"

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: lint
lint:
	golangci-lint run --timeout 3m0s ./

.PHONY: gofmt
gofmt:
	goimports -local ${LOCAL_IMPORT_PATH} -w \
		$$(find . -path "./vendor" -prune -o -type f -name '*.go' -print)

.PHONY: gofmt-check
gofmt-check:
	test -z "$$(goimports -local ${LOCAL_IMPORT_PATH} -d \
		$$(find . -path "./vendor" -prune -o -type f -name '*.go' -print))"

.PHONY: build
build:
	go build

.PHONY: test
test:
	TZ=UTC CGO_ENABLED=0 go test -v ./...
