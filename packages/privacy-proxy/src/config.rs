use std::env;
use std::time::Duration;

use alloy::primitives::Address;
use anyhow::{Context, Result};

#[derive(Clone, Debug)]
pub struct Config {
    pub bind_addr: String,
    pub upstream_url: String,
    pub db_url: String,
    pub admin_eoas: Vec<Address>,
    pub challenge_ttl: Duration,
    pub token_ttl: Duration,
    pub domain: String,
}

impl Config {
    pub fn from_env() -> Result<Self> {
        let bind_addr = env::var("BIND_ADDR").unwrap_or_else(|_| "0.0.0.0:8080".to_string());
        let upstream_url =
            env::var("UPSTREAM_URL").context("UPSTREAM_URL env var is required")?;
        let db_url =
            env::var("DATABASE_URL").unwrap_or_else(|_| "sqlite://privacy-proxy.db".to_string());

        let admin_eoas = env::var("ADMIN_EOAS")
            .unwrap_or_default()
            .split(',')
            .map(str::trim)
            .filter(|s| !s.is_empty())
            .map(|s| s.parse::<Address>().with_context(|| format!("ADMIN_EOAS: invalid address `{s}`")))
            .collect::<Result<Vec<_>>>()?;

        let challenge_ttl = Duration::from_secs(
            env::var("CHALLENGE_TTL_SECS")
                .ok()
                .and_then(|s| s.parse().ok())
                .unwrap_or(300),
        );
        let token_ttl = Duration::from_secs(
            env::var("TOKEN_TTL_SECS")
                .ok()
                .and_then(|s| s.parse().ok())
                .unwrap_or(7 * 24 * 60 * 60),
        );
        let domain = env::var("AUTH_DOMAIN").unwrap_or_else(|_| "privacy-proxy".to_string());

        Ok(Self {
            bind_addr,
            upstream_url,
            db_url,
            admin_eoas,
            challenge_ttl,
            token_ttl,
            domain,
        })
    }
}
