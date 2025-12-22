.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./delete-mysql-data-warehouse .

.PHONY: build-darwin
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./delete-mysql-data-warehouse .

.PHONY: clean
clean:
	rm -f ./delete-mysql-data-warehouse

