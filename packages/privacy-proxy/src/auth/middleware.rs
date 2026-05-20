use alloy::primitives::Address;
use axum::extract::Request;
use axum::middleware::Next;
use axum::response::Response;
use axum::Extension;
use sha2::{Digest, Sha256};
use sqlx::Row;

use crate::auth::{AdminCallerInfo, CallerAttributes, CallerCtx, UserCallerInfo};
use crate::db::now_unix;
use crate::roles::{ROLE_ADMIN, ROLE_USER};
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
        "SELECT m.eoa_address, r.name AS role_name, t.expires_at
         FROM auth_tokens t
         JOIN members m ON m.eoa_address = t.eoa_address
         JOIN roles r ON r.id = m.role_id
         WHERE t.token_hash = ?",
    )
    .bind(&token_hash)
    .fetch_optional(&state.pool)
    .await
    .ok()??;

    let expires_at: i64 = row.try_get("expires_at").ok()?;
    if now > expires_at {
        // Opportunistic sweep: drop every expired token row we encounter
        // (this one plus any other stale rows in the table). Cheap given the
        // expires_at index path is already warm, and keeps the table from
        // growing without bound.
        let _ = sqlx::query("DELETE FROM auth_tokens WHERE expires_at < ?")
            .bind(now)
            .execute(&state.pool)
            .await;
        return None;
    }
    let eoa_str: String = row.try_get("eoa_address").ok()?;
    let role: String = row.try_get("role_name").ok()?;
    let eoa: Address = eoa_str.parse().ok()?;

    let attributes = match role.as_str() {
        ROLE_ADMIN => CallerAttributes::Admin(AdminCallerInfo { eoa }),
        ROLE_USER => {
            // Missing row would mean data drift (a user without an
            // attribute row); fall back to defaults so the request can
            // still proceed and be gated by lambdas/rules.
            let attrs = sqlx::query(
                "SELECT kyc, blacklisted FROM user_attributes WHERE eoa_address = ?",
            )
            .bind(&eoa_str)
            .fetch_optional(&state.pool)
            .await
            .ok()?;
            let (kyc, blacklisted) = match attrs {
                Some(r) => {
                    let kyc: i64 = r.try_get("kyc").ok()?;
                    let bl: i64 = r.try_get("blacklisted").ok()?;
                    (kyc != 0, bl != 0)
                }
                None => {
                    tracing::warn!(
                        eoa = %eoa_str,
                        "member row exists without user_attributes; using defaults",
                    );
                    (false, false)
                }
            };
            CallerAttributes::User(UserCallerInfo {
                eoa,
                kyc,
                blacklisted,
            })
        }
        other => {
            tracing::warn!(role = %other, "unknown role for resolved token");
            return None;
        }
    };

    Some(CallerCtx {
        eoa: Some(eoa),
        attributes: Some(attributes),
    })
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::config::Config;
    use crate::roles::reconcile_roles;
    use crate::state::AppState;
    use std::time::Duration;

    async fn fresh_state() -> AppState {
        let pool = sqlx::sqlite::SqlitePoolOptions::new()
            .max_connections(1)
            .connect("sqlite::memory:")
            .await
            .unwrap();
        sqlx::migrate!("./migrations").run(&pool).await.unwrap();
        reconcile_roles(&pool).await.unwrap();
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

    async fn seed_member_with_token(
        state: &AppState,
        eoa: &str,
        role: &str,
        kyc: bool,
        blacklisted: bool,
    ) -> String {
        let now = now_unix();
        sqlx::query(
            "INSERT INTO members (eoa_address, role_id, created_at)
             VALUES (?, (SELECT id FROM roles WHERE name = ?), ?)",
        )
        .bind(eoa)
        .bind(role)
        .bind(now)
        .execute(&state.pool)
        .await
        .unwrap();
        if role == ROLE_USER {
            sqlx::query(
                "INSERT INTO user_attributes (eoa_address, kyc, blacklisted) VALUES (?, ?, ?)",
            )
            .bind(eoa)
            .bind(kyc as i64)
            .bind(blacklisted as i64)
            .execute(&state.pool)
            .await
            .unwrap();
        }
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
    async fn resolves_user_with_typed_attributes() {
        let state = fresh_state().await;
        let eoa = "0x1111111111111111111111111111111111111111";
        let token = seed_member_with_token(&state, eoa, "user", true, false).await;
        let ctx = resolve_token(&state, &token).await.unwrap();
        match ctx.attributes.unwrap() {
            CallerAttributes::User(u) => {
                assert_eq!(format!("0x{}", hex::encode(u.eoa.as_slice())), eoa);
                assert!(u.kyc);
                assert!(!u.blacklisted);
            }
            other => panic!("expected User, got {other:?}"),
        }
    }

    #[tokio::test]
    async fn resolves_admin_identity_only() {
        let state = fresh_state().await;
        let eoa = "0x2222222222222222222222222222222222222222";
        let token = seed_member_with_token(&state, eoa, "admin", false, false).await;
        let ctx = resolve_token(&state, &token).await.unwrap();
        match ctx.attributes.unwrap() {
            CallerAttributes::Admin(a) => {
                assert_eq!(format!("0x{}", hex::encode(a.eoa.as_slice())), eoa);
            }
            other => panic!("expected Admin, got {other:?}"),
        }
    }

    #[tokio::test]
    async fn falls_back_to_defaults_when_user_attributes_row_missing() {
        let state = fresh_state().await;
        let eoa = "0x3333333333333333333333333333333333333333";
        // Seed members row without user_attributes (simulating data drift).
        let now = now_unix();
        sqlx::query(
            "INSERT INTO members (eoa_address, role_id, created_at)
             VALUES (?, (SELECT id FROM roles WHERE name = 'user'), ?)",
        )
        .bind(eoa)
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

        let ctx = resolve_token(&state, &token).await.unwrap();
        match ctx.attributes.unwrap() {
            CallerAttributes::User(u) => {
                assert!(!u.kyc);
                assert!(!u.blacklisted);
            }
            other => panic!("expected User, got {other:?}"),
        }
    }
}
