module github.com/Zenrock-Foundation/zrchain/v4

go 1.23.2

replace (
	// replace broken sdk lib to ensure problematic version v0.12.0 is not used
	cosmossdk.io/core => cosmossdk.io/core v0.11.0
	github.com/Layr-Labs/eigensdk-go => github.com/Layr-Labs/eigensdk-go v0.1.7-0.20240425202952-954cd7661775
	github.com/cosmos/cosmos-sdk => github.com/zenrocklabs/cosmos-sdk v0.50.10-zenrock
	// fix upstream GHSA-h395-qcrw-5vmq vulnerability.
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.0
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.3.2
	// replace broken goleveldb
	github.com/syndtr/goleveldb => github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
)

require (
	cosmossdk.io/api v0.7.5
	cosmossdk.io/client/v2 v2.0.0-beta.4
	cosmossdk.io/collections v0.4.0
	cosmossdk.io/core v0.12.1-0.20240725072823-6a2d039e1212
	cosmossdk.io/depinject v1.0.0
	cosmossdk.io/errors v1.0.1
	cosmossdk.io/log v1.4.1
	cosmossdk.io/math v1.3.0
	cosmossdk.io/store v1.1.1
	cosmossdk.io/tools/confix v0.1.2
	cosmossdk.io/x/circuit v0.1.1
	cosmossdk.io/x/evidence v0.1.1
	cosmossdk.io/x/feegrant v0.1.1
	cosmossdk.io/x/nft v0.1.1
	cosmossdk.io/x/tx v0.13.5
	cosmossdk.io/x/upgrade v0.1.4
	github.com/CosmWasm/wasmd v0.53.0
	github.com/CosmWasm/wasmvm/v2 v2.1.2
	github.com/Layr-Labs/eigensdk-go v0.1.12
	github.com/btcsuite/btcd v0.24.2
	github.com/btcsuite/btcd/btcec/v2 v2.3.4
	github.com/btcsuite/btcd/btcutil v1.1.6
	github.com/btcsuite/btcd/chaincfg/chainhash v1.1.0
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f
	github.com/btcsuite/btcwallet/walletdb v1.4.4
	github.com/bufbuild/buf v1.28.1
	github.com/cometbft/cometbft v0.38.12
	github.com/cosmos/cosmos-db v1.0.2
	github.com/cosmos/cosmos-proto v1.0.0-beta.5
	github.com/cosmos/cosmos-sdk v0.50.9
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/gogoproto v1.7.0
	github.com/cosmos/ibc-go/modules/capability v1.0.1
	github.com/cosmos/ibc-go/v8 v8.5.0
	github.com/ethereum/go-ethereum v1.14.11
	github.com/fxamacker/cbor/v2 v2.7.0
	github.com/gagliardetto/binary v0.8.0
	github.com/gagliardetto/solana-go v1.11.0
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.4
	github.com/gorilla/mux v1.8.1
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0
	github.com/hashicorp/go-metrics v0.5.3
	github.com/lightninglabs/neutrino v0.16.0
	github.com/near/borsh-go v0.3.1
	github.com/prometheus/client_golang v1.20.3
	github.com/smartcontractkit/chainlink/v2 v2.17.0
	github.com/spf13/cast v1.7.0
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
	github.com/zenrocklabs/zenbtc v0.0.3
	github.com/zenrocklabs/zenrock-avs v1.0.1
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0
	golang.org/x/tools v0.25.0
	google.golang.org/genproto/googleapis/api v0.0.0-20241007155032-5fefd90f89a9
	google.golang.org/grpc v1.67.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.35.1
	gopkg.in/yaml.v3 v3.0.1
	gotest.tools/v3 v3.5.1
)

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.31.0-20231106192134-1baebb0a1518.2 // indirect
	buf.build/gen/go/bufbuild/registry/protocolbuffers/go v1.31.0-20231111212044-1119bf4b707e.2 // indirect
	cloud.google.com/go v0.115.1 // indirect
	cloud.google.com/go/auth v0.9.3 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.4 // indirect
	cloud.google.com/go/compute/metadata v0.5.0 // indirect
	cloud.google.com/go/iam v1.2.0 // indirect
	cloud.google.com/go/storage v1.43.0 // indirect
	connectrpc.com/connect v1.12.0 // indirect
	connectrpc.com/otelconnect v0.6.0 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4 // indirect
	github.com/99designs/keyring v1.2.2 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/DataDog/datadog-go v4.8.3+incompatible // indirect
	github.com/DataDog/zstd v1.5.6 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/aead/siphash v1.0.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.0 // indirect
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/bgentry/speakeasy v0.2.0 // indirect
	github.com/bits-and-blooms/bitset v1.14.2 // indirect
	github.com/blendle/zapdriver v1.3.1 // indirect
	github.com/btcsuite/btcwallet/wtxmgr v1.5.0 // indirect
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd // indirect
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792 // indirect
	github.com/bufbuild/protocompile v0.9.0 // indirect
	github.com/bufbuild/protovalidate-go v0.4.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/fifo v0.0.0-20240816210425-c5d0cb0b6fc0 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.2 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/cometbft/cometbft-db v0.14.0 // indirect
	github.com/consensys/bavard v0.1.13 // indirect
	github.com/consensys/gnark-crypto v0.14.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.15.1 // indirect
	github.com/cosmos/btcutil v1.0.5 // indirect
	github.com/cosmos/gogogateway v1.2.0 // indirect
	github.com/cosmos/iavl v1.2.0 // indirect
	github.com/cosmos/ics23/go v0.11.0 // indirect
	github.com/cosmos/ledger-cosmos-go v0.13.3 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20240223125850-b1e8a79f509c // indirect
	github.com/crate-crypto/go-kzg-4844 v1.1.0 // indirect
	github.com/creachadair/atomicfile v0.3.5 // indirect
	github.com/creachadair/tomledit v0.0.26 // indirect
	github.com/danieljoos/wincred v1.2.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/decred/dcrd/crypto/blake256 v1.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/decred/dcrd/lru v1.0.0 // indirect
	github.com/desertbit/timer v1.0.1 // indirect
	github.com/dgraph-io/badger/v4 v4.3.0 // indirect
	github.com/dgraph-io/ristretto v0.1.2-0.20240116140435-c67e07994f91 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/docker/cli v24.0.7+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker v25.0.5+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.8.0 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/dvsekhvalnov/jose2go v1.7.0 // indirect
	github.com/emicklei/dot v1.6.2 // indirect
	github.com/ethereum/c-kzg-4844 v1.0.3 // indirect
	github.com/ethereum/go-verkle v0.1.1-0.20240829091221-dffa7562dbe9 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gagliardetto/treeout v0.1.4 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-chi/chi/v5 v5.0.10 // indirect
	github.com/go-kit/kit v0.13.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/godbus/dbus v0.0.0-20190726142602-4481cbc300e2 // indirect
	github.com/gofrs/uuid/v5 v5.0.0 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.3 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/glog v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/btree v1.1.3 // indirect
	github.com/google/cel-go v0.18.2 // indirect
	github.com/google/flatbuffers v24.3.25+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-containerregistry v0.16.1 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/orderedcode v0.0.1 // indirect
	github.com/google/pprof v0.0.0-20240424215950-a892ee059fd6 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.3 // indirect
	github.com/googleapis/gax-go/v2 v2.13.0 // indirect
	github.com/gorilla/handlers v1.5.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/gsterjov/go-libsecret v0.0.0-20161001094733-a6f4afe4910c // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.7.6 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-plugin v1.6.2-0.20240829161738-06afb6d7ae99 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/hdevalence/ed25519consensus v0.2.0 // indirect
	github.com/holiman/uint256 v1.3.1 // indirect
	github.com/huandu/skiplist v1.2.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/improbable-eng/grpc-web v0.15.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jdx/go-netrc v1.0.0 // indirect
	github.com/jhump/protoreflect v1.15.6 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jmhodges/levigo v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kkdai/bstream v0.0.0-20161212061736-f391b8402d23 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/lightninglabs/neutrino/cache v1.1.0 // indirect
	github.com/lightningnetwork/lnd/clock v1.0.1 // indirect
	github.com/lightningnetwork/lnd/queue v1.0.1 // indirect
	github.com/lightningnetwork/lnd/ticker v1.0.0 // indirect
	github.com/linxGnu/grocksdb v1.9.3 // indirect
	github.com/lmittmann/tint v1.0.4 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/manifoldco/promptui v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mmcloughlin/addchain v0.4.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mostynb/zstdpool-freelist v0.0.0-20201229113212-927304c0c3b1 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/mtibben/percent v0.2.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oasisprotocol/curve25519-voi v0.0.0-20230904125328-1f23a7beb09a // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/petermattis/goid v0.0.0-20240813172612-4fcff4a6cae7 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pkg/profile v1.7.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.59.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/rs/zerolog v1.33.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sasha-s/go-deadlock v0.3.5 // indirect
	github.com/shamaton/msgpack/v2 v2.2.1 // indirect
	github.com/shirou/gopsutil v3.21.11+incompatible // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/streamingfast/logging v0.0.0-20230608130331-f22c91403091 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/supranational/blst v0.3.13 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/tetratelabs/wazero v1.5.0 // indirect
	github.com/tidwall/btree v1.7.0 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/urfave/cli v1.22.14 // indirect
	github.com/vbatts/tar-split v0.11.5 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zondax/hid v0.9.2 // indirect
	github.com/zondax/ledger-go v0.15.0 // indirect
	go.etcd.io/bbolt v1.4.0-alpha.1 // indirect
	go.mongodb.org/mongo-driver v1.16.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.54.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel v1.31.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.31.0 // indirect
	go.opentelemetry.io/otel/metric v1.31.0 // indirect
	go.opentelemetry.io/otel/sdk v1.31.0 // indirect
	go.opentelemetry.io/otel/trace v1.31.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/ratelimit v0.3.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/oauth2 v0.23.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/term v0.25.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	golang.org/x/time v0.6.0 // indirect
	google.golang.org/api v0.196.0 // indirect
	google.golang.org/genproto v0.0.0-20240903143218-8af14fe29dc1 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	nhooyr.io/websocket v1.8.17 // indirect
	pgregory.net/rapid v1.1.0 // indirect
	rsc.io/tmplfunc v0.0.3 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)