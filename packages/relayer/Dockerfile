# syntax=docker/dockerfile:1.6

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

ARG PACKAGE=relayer
ARG BUILDKIT_INLINE_CACHE=1
ARG TARGETPLATFORM
ARG BUILDPLATFORM

# Debug: Print platform information
RUN echo "BUILDPLATFORM: $BUILDPLATFORM" && \
    echo "TARGETPLATFORM: $TARGETPLATFORM" && \
    uname -a && \
    go env GOOS GOARCH

# Install build dependencies
RUN apk add --no-cache git make file

WORKDIR /build

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY packages/${PACKAGE} ./packages/${PACKAGE}

# Build the application
WORKDIR /build/packages/${PACKAGE}

# Determine GOARCH from TARGETPLATFORM
RUN case ${TARGETPLATFORM} in \
    "linux/amd64")  GOARCH="amd64"  ;; \
    "linux/arm64")  GOARCH="arm64"  ;; \
    "linux/arm/v7") GOARCH="arm"    ;; \
    *)              GOARCH="amd64"  ;; \
    esac && \
    echo "Building for GOARCH: $GOARCH" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build \
    -ldflags="-w -s" \
    -trimpath \
    -o /build/bin/${PACKAGE} \
    ./cmd/main.go && \
    echo "Binary info:" && \
    file /build/bin/${PACKAGE}

FROM --platform=$TARGETPLATFORM alpine:latest AS runtime

ARG PACKAGE=relayer
ARG TARGETPLATFORM

# Debug: Print platform information
RUN echo "Runtime TARGETPLATFORM: $TARGETPLATFORM" && \
    uname -a

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata file && \
    addgroup -g 10001 -S appuser && \
    adduser -u 10001 -S appuser -G appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /build/bin/${PACKAGE} /usr/local/bin/${PACKAGE}

# Debug: Verify binary architecture matches runtime
RUN echo "Verifying binary architecture:" && \
    file /usr/local/bin/${PACKAGE} && \
    ldd /usr/local/bin/${PACKAGE} 2>&1 || echo "Static binary (no dynamic libs)"

USER appuser

ENTRYPOINT ["/usr/local/bin/${PACKAGE}"]