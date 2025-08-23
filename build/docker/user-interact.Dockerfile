# syntax=docker/dockerfile:1.4
FROM golang:1.24-alpine AS builder
LABEL authors="vizuth"

WORKDIR /build

COPY user-interact/go.* ./user-interact/
COPY common/go.* ./common/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go work init ./common ./user-interact && go mod download

COPY user-interact/ ./user-interact/
COPY common/ ./common/
COPY configs/ ./configs/

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o user-interact-service user-interact/cmd/main.go

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app

COPY --from=builder /build/configs/config.yaml ../configs/config.yaml
COPY --from=builder /build/user-interact-service .

CMD ["./user-interact-service"]