module eventify/gateway

go 1.24.0

toolchain go1.24.6

require (
	eventify/auth v0.0.0-00010101000000-000000000000
	eventify/common v0.0.0-00010101000000-000000000000
	eventify/event v0.0.0-00010101000000-000000000000
	eventify/user-interact v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/rs/cors v1.11.1
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.75.1
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250908214217-97024824d090 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250826171959-ef028d996bc1 // indirect
	google.golang.org/protobuf v1.36.9 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace eventify/common => ../common

replace eventify/auth => ../auth

replace eventify/event => ../event

replace eventify/user-interact => ../user-interact
