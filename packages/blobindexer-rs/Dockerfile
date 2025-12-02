# syntax=docker/dockerfile:1.6

FROM rust:1.90-slim-bookworm AS builder

ARG BUILDKIT_INLINE_CACHE=1

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    pkg-config \
    libssl-dev \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy manifest files first to leverage Docker layer caching
COPY Cargo.toml Cargo.lock ./

# Copy source code
COPY src ./src
COPY bin ./bin
COPY migrations ./migrations

# Build the application
RUN cargo build --release --bin blobindexer && \
    strip /app/target/release/blobindexer

FROM debian:bookworm-slim AS runtime

# Create non-root user
RUN useradd -m -u 10001 appuser && \
    apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary and migrations from builder
COPY --from=builder --chown=appuser:appuser /app/target/release/blobindexer /usr/local/bin/blobindexer
COPY --from=builder --chown=appuser:appuser /app/migrations ./migrations

USER appuser

ENV RUST_LOG=info

ENTRYPOINT ["blobindexer"]
