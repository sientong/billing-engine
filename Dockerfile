# ─────── STAGE 1: Build ─────────────────────────────────────────────
FROM golang:1.23.4 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Use static binary build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o billing-engine .

# ─────── STAGE 2: Minimal runtime image ─────────────────────────────
FROM alpine:latest

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/billing-engine .
COPY --from=builder /app/db/migrations ./db/migrations 

# Run binary
ENTRYPOINT ["./billing-engine"]