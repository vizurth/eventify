export PATH := $(LOCAL_BIN):$(PATH)
LOCAL_BIN:=$(CURDIR)/bin
OPENAPIS_DIR=../third_party/
PROTOC=protoc

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.20.0

gen-openapi:
	$(PROTOC) \
		-I $(OPENAPIS_DIR) \
		-I ../auth/api \
		-I ../event/api \
		-I ../user-interact/api \
		-I . \
		--openapiv2_out ./docs \
		--openapiv2_opt logtostderr=true \
		--openapiv2_opt allow_merge=false \
		../auth/api/auth.proto

	$(PROTOC) \
		-I $(OPENAPIS_DIR) \
		-I ../auth/api \
		-I ../event/api \
		-I ../user-interact/api \
		-I . \
		--openapiv2_out ./docs \
		--openapiv2_opt logtostderr=true \
		--openapiv2_opt allow_merge=false \
		../event/api/event.proto

	$(PROTOC) \
		-I $(OPENAPIS_DIR) \
		-I ../auth/api \
		-I ../event/api \
		-I ../user-interact/api \
		-I . \
		--openapiv2_out ./docs \
		--openapiv2_opt logtostderr=true \
		--openapiv2_opt allow_merge=false \
		../user-interact/api/user-interact.proto
