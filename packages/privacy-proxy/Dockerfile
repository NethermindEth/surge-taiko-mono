# syntax=docker/dockerfile:1.6

FROM rust:1.90-slim-bookworm AS builder
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    pkg-config \
    libssl-dev \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Migrations are embedded at compile time via sqlx::migrate!("./migrations"),
# so the migrations directory must be present during the build.
COPY Cargo.toml Cargo.lock ./
COPY src ./src
COPY migrations ./migrations

RUN cargo build --release --bin privacy-proxy

FROM debian:bookworm-slim AS runtime
RUN useradd -m -u 10001 appuser && \
    apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/target/release/privacy-proxy /usr/local/bin/privacy-proxy

# Persistent location for the SQLite database; mount a volume here in prod.
RUN mkdir -p /data && chown appuser:appuser /data
VOLUME ["/data"]

USER appuser
ENV RUST_LOG=info \
    BIND_ADDR=0.0.0.0:8080 \
    DATABASE_URL=sqlite:///data/privacy-proxy.db

EXPOSE 8080
ENTRYPOINT ["privacy-proxy"]
