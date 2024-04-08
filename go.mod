module github.com/tofutf/tofutf

go 1.21

replace github.com/r3labs/sse/v2 => github.com/leg100/sse/v2 v2.0.0-20220910081853-79ffbd7c2fad

replace github.com/jaschaephraim/lrserver => github.com/sowiner/lrserver v0.0.0-20230123160823-795409868576

replace github.com/natefinch/atomic => github.com/sdassow/atomic v0.0.0-20220219102542-174b5d2a3ea6

require (
	cloud.google.com/go/profiler v0.4.0
	cloud.google.com/go/pubsub v1.38.0
	github.com/DataDog/jsonapi v0.8.6
	github.com/Masterminds/sprig/v3 v3.2.3
	github.com/ProtonMail/go-crypto v1.0.0
	github.com/allegro/bigcache v1.2.1
	github.com/antchfx/htmlquery v1.3.1
	github.com/bradleyfalzon/ghinstallation/v2 v2.10.0
	github.com/buildkite/terminal-to-html v3.2.0+incompatible
	github.com/cenkalti/backoff/v4 v4.3.0
	github.com/chromedp/cdproto v0.0.0-20240524221637-55927c2a4565
	github.com/chromedp/chromedp v0.9.5
	github.com/coreos/go-oidc/v3 v3.10.0
	github.com/fatih/color v1.17.0
	github.com/felixge/httpsnoop v1.0.4
	github.com/fsnotify/fsnotify v1.7.0
	github.com/gfleury/go-bitbucket-v1 v0.0.0-20240131155556-0b41d7863037
	github.com/gobwas/glob v0.2.3
	github.com/gomarkdown/markdown v0.0.0-20240419095408-642f0ee99ae2
	github.com/google/go-github/v41 v41.0.0
	github.com/google/go-github/v55 v55.0.0
	github.com/google/goexpect v0.0.0-20210430020637-ab937bf7fd6f
	github.com/google/uuid v1.6.0
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/schema v1.3.0
	github.com/hashicorp/go-retryablehttp v0.7.6
	github.com/hashicorp/go-tfe v1.53.0
	github.com/hashicorp/hcl/v2 v2.20.1
	github.com/hashicorp/terraform-config-inspect v0.0.0-20240509232506-4708120f8f30
	github.com/iancoleman/strcase v0.3.0
	github.com/jackc/pgconn v1.14.3
	github.com/jackc/pgtype v1.14.3
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackc/pgx/v4 v4.18.2
	github.com/jaschaephraim/lrserver v0.0.0-00010101000000-000000000000
	github.com/leg100/surl v0.0.6
	github.com/lestrrat-go/jwx/v2 v2.0.21
	github.com/mitchellh/iochan v1.0.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/natefinch/atomic v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/pressly/goose/v3 v3.20.0
	github.com/prometheus/client_golang v1.19.1
	github.com/redis/go-redis/v9 v9.5.1
	github.com/spf13/cobra v1.8.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.9.0
	github.com/xanzy/go-gitlab v0.105.0
	golang.org/x/exp v0.0.0-20240525044651-4c93da0ed11d
	golang.org/x/mod v0.17.0
	golang.org/x/net v0.25.0
	golang.org/x/oauth2 v0.20.0
	golang.org/x/sync v0.7.0
	google.golang.org/api v0.181.0
	gorm.io/driver/postgres v1.5.7
	gorm.io/gen v0.3.26
	gorm.io/gorm v1.25.10
)

require (
	cloud.google.com/go v0.113.0 // indirect
	cloud.google.com/go/auth v0.4.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.2 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/iam v1.1.7 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/antchfx/xpath v1.3.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/cloudflare/circl v1.3.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-jose/go-jose/v4 v4.0.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.3.2 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-github/v60 v60.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/goterm v0.0.0-20190703233501-fc88cf888a3f // indirect
	github.com/google/pprof v0.0.0-20230602150820-91b7bce49751 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.4 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-slug v0.15.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl v0.0.0-20170504190234-a4b07c25de5f // indirect
	github.com/hashicorp/jsonapi v1.3.1 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/imdario/mergo v0.3.11 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/itchyny/base58-go v0.2.1 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9 // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle v1.3.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.5 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mfridman/interpolate v0.0.2 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/zclconf/go-cty v1.14.4 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
	google.golang.org/genproto v0.0.0-20240401170217-c3f982113cda // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240506185236-b8a5c65736ae // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240513163218-0867130af1f8 // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/datatypes v1.1.1-0.20230130040222-c43177d3cf8c // indirect
	gorm.io/driver/mysql v1.4.4 // indirect
	gorm.io/hints v1.1.0 // indirect
	gorm.io/plugin/dbresolver v1.5.0 // indirect
)
