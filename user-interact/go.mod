module eventify/user-interact

go 1.24.0

toolchain go1.24.6

require (
	eventify/common v0.0.0-00010101000000-000000000000
	github.com/Masterminds/squirrel v1.5.4
	github.com/envoyproxy/protoc-gen-validate v1.2.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/jackc/pgx/v5 v5.7.5
	go.uber.org/zap v1.27.0
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff
	google.golang.org/grpc v1.76.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/segmentio/kafka-go v0.4.49 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.42.0 // indirect
	golang.org/x/net v0.45.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251002232023-7c0ddcbb5797 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace eventify/common => ../common
