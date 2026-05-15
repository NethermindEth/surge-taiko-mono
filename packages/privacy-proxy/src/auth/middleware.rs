use axum::extract::Request;
use axum::middleware::Next;
use axum::response::Response;
use axum::Extension;
use sha2::{Digest, Sha256};
use sqlx::Row;

use crate::auth::CallerCtx;
use crate::db::now_unix;
use crate::state::AppState;

/// Axum middleware that resolves an `Authorization: Bearer <token>` header
/// into a `CallerCtx` and inserts it into the request extensions. Missing
/// or invalid tokens fall back to `CallerCtx::anonymous()` — actual
/// authorization decisions happen in the handlers, not here.
///
/// Reads `AppState` from a request extension (set by an outer `Extension`
/// layer) to avoid axum's `from_fn_with_state` extractor-tuple inference.
pub async fn caller_ctx_layer(
    Extension(state): Extension<AppState>,
    mut req: Request,
    next: Next,
) -> Response {
    let token = extract_bearer_token(&req);
    let ctx = match token {
        Some(t) => resolve_token(&state, &t).await.unwrap_or_else(|| {
            tracing::debug!("token did not resolve; treating as anonymous");
            CallerCtx::anonymous()
        }),
        None => CallerCtx::anonymous(),
    };
    req.extensions_mut().insert(ctx);
    next.run(req).await
}

fn extract_bearer_token(req: &Request) -> Option<String> {
    let header = req.headers().get(axum::http::header::AUTHORIZATION)?;
    let value = header.to_str().ok()?;
    let token = value
        .strip_prefix("Bearer ")
        .or_else(|| value.strip_prefix("bearer "))?;
    if token.is_empty() {
        None
    } else {
        Some(token.to_string())
    }
}

async fn resolve_token(state: &AppState, token: &str) -> Option<CallerCtx> {
    let token_hash = hex::encode(Sha256::digest(token.as_bytes()));
    let now = now_unix();

    let row = sqlx::query(
        "SELECT u.eoa_address, r.name AS role_name, u.caller_info_json, t.expires_at
         FROM auth_tokens t
         JOIN users u ON u.eoa_address = t.eoa_address
         JOIN roles r ON r.id = u.role_id
         WHERE t.token_hash = ?",
    )
    .bind(&token_hash)
    .fetch_optional(&state.pool)
    .await
    .ok()??;

    let expires_at: i64 = row.try_get("expires_at").ok()?;
    if now > expires_at {
        return None;
    }
    let eoa_str: String = row.try_get("eoa_address").ok()?;
    let role: String = row.try_get("role_name").ok()?;
    let info_str: String = row.try_get("caller_info_json").ok()?;
    let mut caller_info: serde_json::Value =
        serde_json::from_str(&info_str).unwrap_or_else(|_| serde_json::json!({}));
    // The auth layer is the trusted writer of identity fields. Any pre-existing
    // `eoa` value in caller_info_json is overwritten with the token-resolved EOA;
    // a non-object caller_info is replaced with a fresh object so the injection
    // is always available to lambdas.
    if !caller_info.is_object() {
        caller_info = serde_json::json!({});
    }
    if let serde_json::Value::Object(map) = &mut caller_info {
        map.insert(
            "eoa".to_string(),
            serde_json::Value::String(eoa_str.to_ascii_lowercase()),
        );
    }
    let eoa = eoa_str.parse().ok()?;
    Some(CallerCtx {
        eoa: Some(eoa),
        role: Some(role),
        caller_info,
    })
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::Config;
    use crate::state::AppState;
    use std::time::Duration;

    async fn fresh_state() -> AppState {
        let pool = sqlx::sqlite::SqlitePoolOptions::new()
            .max_connections(1)
            .connect("sqlite::memory:")
            .await
            .unwrap();
        sqlx::migrate!("./migrations").run(&pool).await.unwrap();
        let config = Config {
            bind_addr: "127.0.0.1:0".to_string(),
            upstream_url: "http://127.0.0.1:1".to_string(),
            db_url: "sqlite::memory:".to_string(),
            admin_eoas: vec![],
            challenge_ttl: Duration::from_secs(300),
            token_ttl: Duration::from_secs(3600),
            domain: "test".to_string(),
        };
        AppState::new(config, pool)
    }

    async fn seed_user_with_token(
        state: &AppState,
        eoa: &str,
        role: &str,
        caller_info_json: &str,
    ) -> String {
        let now = now_unix();
        sqlx::query(
            "INSERT INTO users (eoa_address, role_id, caller_info_json, created_at)
             VALUES (?, (SELECT id FROM roles WHERE name = ?), ?, ?)",
        )
        .bind(eoa)
        .bind(role)
        .bind(caller_info_json)
        .bind(now)
        .execute(&state.pool)
        .await
        .unwrap();
        use rand::RngCore;
        let mut buf = [0u8; 32];
        rand::thread_rng().fill_bytes(&mut buf);
        let token = hex::encode(buf);
        let token_hash = hex::encode(Sha256::digest(token.as_bytes()));
        sqlx::query(
            "INSERT INTO auth_tokens (token_hash, eoa_address, issued_at, expires_at)
             VALUES (?, ?, ?, ?)",
        )
        .bind(&token_hash)
        .bind(eoa)
        .bind(now)
        .bind(now + 3600)
        .execute(&state.pool)
        .await
        .unwrap();
        token
    }

    #[tokio::test]
    async fn auth_layer_injects_eoa_and_preserves_admin_set_fields() {
        let state = fresh_state().await;
        let eoa = "0x1111111111111111111111111111111111111111";
        let token = seed_user_with_token(&state, eoa, "user", r#"{"kyc":true}"#).await;
        let ctx = resolve_token(&state, &token).await.unwrap();
        assert_eq!(ctx.caller_info["eoa"], serde_json::json!(eoa));
        assert_eq!(ctx.caller_info["kyc"], serde_json::json!(true));
    }

    #[tokio::test]
    async fn auth_layer_overwrites_admin_set_eoa_field() {
        let state = fresh_state().await;
        let eoa = "0x2222222222222222222222222222222222222222";
        let token = seed_user_with_token(
            &state,
            eoa,
            "user",
            r#"{"eoa":"0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"}"#,
        )
        .await;
        let ctx = resolve_token(&state, &token).await.unwrap();
        assert_eq!(ctx.caller_info["eoa"], serde_json::json!(eoa));
    }

    #[tokio::test]
    async fn auth_layer_synthesizes_object_when_caller_info_is_empty_object() {
        let state = fresh_state().await;
        let eoa = "0x3333333333333333333333333333333333333333";
        let token = seed_user_with_token(&state, eoa, "user", "{}").await;
        let ctx = resolve_token(&state, &token).await.unwrap();
        assert_eq!(ctx.caller_info["eoa"], serde_json::json!(eoa));
        assert!(ctx.caller_info.is_object());
    }
}
