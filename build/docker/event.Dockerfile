# syntax=docker/dockerfile:1.4
FROM golang:1.24-alpine AS builder
LABEL authors="vizuth"

WORKDIR /build

COPY event/go.* ./event/
COPY common/go.* ./common/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go work init ./common ./event && go mod download

COPY event/ ./event/
COPY common/ ./common/
COPY configs/ ./configs/

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o event-service event/cmd/main.go

FROM gcr.io/distroless/base-debian12 AS runner

WORKDIR /app

COPY --from=builder /build/configs/config.yaml ../configs/config.yaml
COPY --from=builder /build/event-service .

CMD ["./event-service"]