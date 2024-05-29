module github.com/bucketeer-io/bucketeer/evaluation

go 1.20

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/bucketeer-io/bucketeer/proto v0.0.0-20240529105832-0d897d36a2f3
	github.com/golang/protobuf v1.5.4
	github.com/stretchr/testify v1.8.4
	go.uber.org/mock v0.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240412170617-26222e5d3d56 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/bucketeer-io/bucketeer/proto v0.0.0-20240529105832-0d897d36a2f3 => ../proto
