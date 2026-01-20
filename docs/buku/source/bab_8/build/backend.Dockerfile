# syntax=docker/dockerfile:1

# Stage 1: Build
FROM golang:1.25-alpine AS build

WORKDIR /src

# Download dependencies first (caching layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
# -buildvcs=false is often needed in CI/Docker if .git is missing or mismatched
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /out/api ./cmd/api

# Stage 2: Runtime
FROM alpine:3.20

# Install CA certificates for HTTPS calls (e.g. to Packing Service)
RUN apk add --no-cache ca-certificates \
    && adduser -D -h /app app

# Copy binary from build stage
COPY --from=build /out/api /usr/local/bin/api

# Run as non-root user
USER app
WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/api"]
