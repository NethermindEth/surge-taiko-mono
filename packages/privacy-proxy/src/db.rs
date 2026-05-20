use std::time::{SystemTime, UNIX_EPOCH};

use anyhow::{Context, Result};
use sqlx::sqlite::{SqliteConnectOptions, SqlitePoolOptions};
use sqlx::{ConnectOptions, SqlitePool};

pub type Pool = SqlitePool;

pub async fn init_pool(url: &str) -> Result<Pool> {
    let mut options: SqliteConnectOptions = url
        .parse::<SqliteConnectOptions>()
        .with_context(|| format!("invalid DATABASE_URL `{url}`"))?
        .create_if_missing(true)
        .foreign_keys(true);

    // Statement logging is silenced by default to avoid leaking row contents
    // into prod logs. Set `PROXY_LOG_SQL=1` (or any non-empty value) when
    // debugging locally.
    if std::env::var("PROXY_LOG_SQL").ok().filter(|v| !v.is_empty()).is_none() {
        options = options.disable_statement_logging();
    }

    let pool = SqlitePoolOptions::new()
        .max_connections(8)
        .connect_with(options)
        .await
        .context("failed to open sqlite database")?;

    sqlx::migrate!("./migrations")
        .run(&pool)
        .await
        .context("failed to apply migrations")?;

    Ok(pool)
}

pub fn now_unix() -> i64 {
    SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map(|d| d.as_secs() as i64)
        .unwrap_or(0)
}
