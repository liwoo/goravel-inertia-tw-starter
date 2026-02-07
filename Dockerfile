# =============================================================================
# Multi-stage Dockerfile for Goravel (Go + React/Inertia) Application
# =============================================================================
# Build arguments for flexibility
ARG GO_VERSION=1.24
ARG NODE_VERSION=20
ARG ALPINE_VERSION=3.20

# =============================================================================
# Stage 1: Go Backend Builder
# =============================================================================
FROM golang:${GO_VERSION}-alpine AS go-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    upx

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY main.go ./
COPY bootstrap ./bootstrap
COPY app ./app
COPY config ./config
COPY routes ./routes
COPY database ./database
COPY docs ./docs

# Build with optimizations
# -ldflags="-w -s" strips debug info and symbol table
# -trimpath removes file system paths from binary
ARG BUILD_DATE
ARG VCS_REF
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} \
    go build -a -trimpath \
    -ldflags="-w -s -extldflags '-static' -X main.BuildDate=${BUILD_DATE} -X main.GitCommit=${VCS_REF}" \
    -o goravel-app .

# Compress binary with UPX (optional, reduces size by ~50%)
RUN upx --best --lzma goravel-app || true

# =============================================================================
# Stage 2: Node.js Frontend Builder
# =============================================================================
FROM node:${NODE_VERSION}-alpine AS node-builder

WORKDIR /app

# Copy package files for better caching
COPY package.json ./

# Install dependencies using fresh npm install to get correct platform binaries
# (npm ci uses lockfile which may have wrong platform-specific deps)
RUN npm install --legacy-peer-deps

# Copy frontend source files
COPY tsconfig.json ./
COPY vite.config.ts ./
COPY vitest.config.ts* ./
COPY postcss.config.js ./
COPY tailwind.config.js ./
COPY components.json ./
COPY resources ./resources

# Copy existing public directory (static images, etc.) BEFORE vite build
# Vite will add its output to this directory
COPY public ./public

# Build frontend assets (skip tsc type-checking, just build with vite)
# Type checking is done in CI, here we just need the production bundle
RUN npx vite build

# =============================================================================
# Stage 3: Runtime Image
# =============================================================================
FROM alpine:${ALPINE_VERSION} AS runtime

# Labels for container metadata
LABEL org.opencontainers.image.title="Goravel Blog" \
      org.opencontainers.image.description="Goravel application with React/Inertia frontend" \
      org.opencontainers.image.vendor="Books" \
      org.opencontainers.image.source="https://github.com/Tiyeni/books-database" \
      org.opencontainers.image.licenses="MIT"

# Install only runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    curl \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1001 goravel && \
    adduser -D -u 1001 -G goravel goravel

WORKDIR /app

# Copy the binary from go-builder
COPY --from=go-builder --chown=goravel:goravel /app/goravel-app .

# Copy built frontend assets
COPY --from=node-builder --chown=goravel:goravel /app/public ./public

# Copy view templates (Go templates, not built by Node)
COPY --chown=goravel:goravel resources/views ./resources/views

# Create necessary directories with proper permissions
RUN mkdir -p storage/logs storage/app/public storage/framework/cache \
             storage/framework/sessions storage/framework/views \
             database/migrations database/seeders && \
    chown -R goravel:goravel storage database && \
    chmod -R 755 storage

# Copy database files (migrations run at deploy time, not embedded in binary)
COPY --chown=goravel:goravel database/migrations ./database/migrations
COPY --chown=goravel:goravel database/seeders ./database/seeders

# Switch to non-root user
USER goravel

# Expose application port
EXPOSE 3000

# Health check - verify the app responds on /health endpoint
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

# Environment variables
ENV APP_ENV=production \
    APP_HOST=0.0.0.0 \
    APP_PORT=3000 \
    TZ=UTC

# Use exec form for better signal handling
ENTRYPOINT ["/app/goravel-app"]
