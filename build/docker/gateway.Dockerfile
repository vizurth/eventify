# syntax=docker/dockerfile:1.4
FROM golang:1.24-alpine AS builder
LABEL authors="vizuth"

WORKDIR /build

COPY gateway/go.* ./gateway/
COPY common/go.* ./common/
COPY auth/go.* ./auth/
COPY event/go.* ./event/
COPY user-interact/go.* ./user-interact/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go work init ./common ./gateway && go mod download

COPY gateway/ ./gateway/
COPY auth/ ./auth/
COPY event/ ./event/
COPY user-interact/ ./user-interact/
COPY common/ ./common/
COPY configs/ ./configs/

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o gateway-service gateway/cmd/main.go

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app

COPY --from=builder /build/configs/config.yaml ../configs/config.yaml
COPY --from=builder /build/gateway-service .

CMD ["./gateway-service"]