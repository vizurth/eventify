include configs/.env

# Пути к docker-compose файлов
COMPOSE_PROJECT=build/docker/docker-compose.yaml
# Запуск всех сервисов
up:
	docker-compose -f $(COMPOSE_PROJECT) up -d --build

# Остановка всех сервисов
down:
	docker-compose -f $(COMPOSE_PROJECT) down

# Перезапуск
restart: down up

# Просмотр логов всех сервисов
logs:
	docker-compose -f $(COMPOSE_ORDER) logs -f &
	docker-compose -f $(COMPOSE_MIGRATION) logs -f &

LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

local-migration-status:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

vendor-proto:
		@if [ ! -d third_party/google ]; then \
			git clone https://github.com/googleapis/googleapis third_party/googleapis &&\
			mkdir -p  third_party/google/ &&\
			mv third_party/googleapis/google/api third_party/google &&\
			rm -rf third_party/googleapis ;\
		fi
		@if [ ! -d third_party/protoc-gen-openapiv2 ]; then \
			mkdir -p third_party/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway third_party/openapiv2 &&\
			mv third_party/openapiv2/protoc-gen-openapiv2/options/*.proto third_party/protoc-gen-openapiv2/options &&\
			rm -rf third_party/openapiv2 ;\
		fi


#go mod edit -replace eventify/common=../common
#go get eventify/common

#mkdir -p third_party && cd third_party
#git clone https://github.com/googleapis/googleapis.git
#cd ..
