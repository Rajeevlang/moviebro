# =============================================================================
# Stage 1: Build
# =============================================================================
FROM golang:1.27-alpine AS builder

ENV GOTOOLCHAIN=local

WORKDIR /app

# Cache Go dependencies layer
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code
COPY cmd/movie/ ./cmd/movie/
COPY internal/ ./internal/

# Build static Linux binary
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/movie-service ./cmd/movie

# =============================================================================
# Stage 2: Minimal Runtime
# =============================================================================
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/movie-service /app/movie-service

EXPOSE 8081

CMD ["/app/movie-service"]

