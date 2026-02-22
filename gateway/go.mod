module eventify/gateway

go 1.24.0

toolchain go1.24.6

require (
	eventify/auth v0.0.0-00010101000000-000000000000
	eventify/common v0.0.0-00010101000000-000000000000
	eventify/event v0.0.0-00010101000000-000000000000
	eventify/user-interact v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.3
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/rs/cors v1.11.1
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.76.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/alecthomas/kingpin v2.2.6+incompatible // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/cheggaaa/pb v1.0.29 // indirect
	github.com/codesenberg/bombardier v1.2.6 // indirect
	github.com/codesenberg/concurrent v0.0.0-20180531114123-64560cfcf964 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/juju/ratelimit v1.0.2 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.46.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.45.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251007200510-49b9836ed3ff // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251002232023-7c0ddcbb5797 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace eventify/common => ../common

replace eventify/auth => ../auth

replace eventify/event => ../event

replace eventify/user-interact => ../user-interact
