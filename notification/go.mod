module eventify/notification

go 1.24.0

toolchain go1.24.6

replace eventify/common => ../common

require (
	eventify/common v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.3
	github.com/ilyakaznacheev/cleanenv v1.5.0
	go.uber.org/zap v1.27.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/segmentio/kafka-go v0.4.49 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
