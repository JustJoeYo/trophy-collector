# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download dependencies first (cached layer unless go.mod changes)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a statically linked binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
# Distroless = no shell, no package manager, minimal attack surface
FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["/app/server"]
