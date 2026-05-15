use alloy::primitives::Address;
use axum::extract::{Query, State};
use axum::Json;
use rand::RngCore;
use serde::{Deserialize, Serialize};

use crate::auth::{build_signin_message, format_address};
use crate::db::now_unix;
use crate::error::{ApiError, ApiResult};
use crate::state::AppState;

#[derive(Deserialize)]
pub struct ChallengeQuery {
    address: String,
}

#[derive(Serialize)]
pub struct ChallengeResponse {
    pub message: String,
    pub expires_at: i64,
}

/// `GET /auth/challenge?address=0x...` — generates a fresh nonce for the
/// supplied EOA, persists it, and returns the EIP-191 message to sign.
pub async fn handler(
    State(state): State<AppState>,
    Query(q): Query<ChallengeQuery>,
) -> ApiResult<Json<ChallengeResponse>> {
    let address: Address = q
        .address
        .parse()
        .map_err(|_| ApiError::bad_request("invalid address"))?;
    let addr_hex = format_address(&address);

    let mut nonce_bytes = [0u8; 16];
    rand::thread_rng().fill_bytes(&mut nonce_bytes);
    let nonce = hex::encode(nonce_bytes);

    let now = now_unix();
    let expires_at = now + state.config.challenge_ttl.as_secs() as i64;

    sqlx::query(
        "INSERT INTO challenges (eoa_address, nonce, expires_at) VALUES (?, ?, ?)
         ON CONFLICT(eoa_address) DO UPDATE SET nonce = excluded.nonce, expires_at = excluded.expires_at",
    )
    .bind(&addr_hex)
    .bind(&nonce)
    .bind(expires_at)
    .execute(&state.pool)
    .await?;

    let message = build_signin_message(&state.config.domain, &address, &nonce);
    Ok(Json(ChallengeResponse {
        message,
        expires_at,
    }))
}
