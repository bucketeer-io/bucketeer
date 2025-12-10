module github.com/bucketeer-io/bucketeer/v2

go 1.25.1

require (
	cloud.google.com/go/bigquery v1.72.0
	cloud.google.com/go/kms v1.25.0
	cloud.google.com/go/profiler v0.4.3
	cloud.google.com/go/pubsub v1.50.1
	contrib.go.opencensus.io/exporter/stackdriver v0.13.14
	github.com/VividCortex/mysqlerr v1.0.0
	github.com/aws/aws-sdk-go-v2 v1.41.1
	github.com/aws/aws-sdk-go-v2/config v1.32.7
	github.com/aws/aws-sdk-go-v2/credentials v1.19.7
	github.com/aws/aws-sdk-go-v2/service/kms v1.49.5
	github.com/aws/aws-sdk-go-v2/service/sesv2 v1.59.1
	github.com/blang/semver v3.5.1+incompatible
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/coreos/go-oidc/v3 v3.17.0
	github.com/go-gota/gota v0.12.0
	github.com/go-jose/go-jose/v4 v4.1.3
	github.com/go-resty/resty/v2 v2.17.1
	github.com/go-sql-driver/mysql v1.9.3
	github.com/golang/protobuf v1.5.4
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.7.0
	github.com/google/uuid v1.6.0
	github.com/googleapis/gax-go/v2 v2.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.7
	github.com/hashicorp/vault/api v1.22.0
	github.com/improbable-eng/grpc-web v0.15.0
	github.com/jinzhu/copier v0.4.0
	github.com/lib/pq v1.11.1
	github.com/mailersend/mailersend-go v1.6.2
	github.com/mna/redisc v1.4.0
	github.com/nicksnyder/go-i18n/v2 v2.6.1
	github.com/prometheus/client_golang v1.23.2
	github.com/redis/go-redis/v9 v9.17.3
	github.com/sendgrid/sendgrid-go v3.16.1+incompatible
	github.com/slack-go/slack v0.17.3
	github.com/spaolacci/murmur3 v1.1.0
	github.com/stretchr/testify v1.11.1
	github.com/tkuchiki/go-timezone v0.2.3
	go.opencensus.io v0.24.0
	go.uber.org/mock v0.6.0
	go.uber.org/zap v1.27.1
	golang.org/x/crypto v0.47.0
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56
	golang.org/x/oauth2 v0.34.0
	golang.org/x/sync v0.19.0
	golang.org/x/text v0.33.0
	golang.org/x/time v0.14.0
	gonum.org/v1/gonum v0.17.0
	google.golang.org/api v0.264.0
	google.golang.org/genproto v0.0.0-20251202230838-ff82c1b0f217
	google.golang.org/genproto/googleapis/api v0.0.0-20260128011058-8636f8732409
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260128011058-8636f8732409
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.4.0
)

require (
	cloud.google.com/go v0.121.6 // indirect
	cloud.google.com/go/auth v0.18.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/iam v1.5.3 // indirect
	cloud.google.com/go/longrunning v0.7.0 // indirect
	cloud.google.com/go/monitoring v1.24.3 // indirect
	cloud.google.com/go/pubsub/v2 v2.0.0 // indirect
	cloud.google.com/go/trace v1.11.7 // indirect
	filippo.io/edwards25519 v1.1.1 // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/apache/arrow/go/v15 v15.0.2 // indirect
	github.com/aws/aws-sdk-go v1.49.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.4 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.0.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.30.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.35.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.41.6 // indirect
	github.com/aws/smithy-go v1.24.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/desertbit/timer v0.0.0-20180107155436-c41aec40b27f // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/flatbuffers v23.5.26+incompatible // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/pprof v0.0.0-20250602020802-c6617b811d0e // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.11 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.2.0 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.7 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-7 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pierrec/lz4/v4 v4.1.18 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/prometheus/prometheus v0.35.0 // indirect
	github.com/rs/cors v1.8.2 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.61.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 // indirect
	go.opentelemetry.io/otel v1.39.0 // indirect
	go.opentelemetry.io/otel/metric v1.39.0 // indirect
	go.opentelemetry.io/otel/trace v1.39.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/mod v0.31.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/telemetry v0.0.0-20251203150158-8fff8a5912fc // indirect
	golang.org/x/tools v0.40.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	nhooyr.io/websocket v1.8.6 // indirect
)

replace k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191003035328-700b1226c0bd
