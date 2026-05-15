use alloy::primitives::{Address, Signature};
use axum::extract::State;
use axum::Json;
use rand::RngCore;
use serde::{Deserialize, Serialize};
use sha2::{Digest, Sha256};
use sqlx::Row;

use crate::auth::{build_signin_message, format_address};
use crate::db::now_unix;
use crate::error::{ApiError, ApiResult};
use crate::state::AppState;

#[derive(Deserialize)]
pub struct VerifyRequest {
    pub address: String,
    pub signature: String,
}

#[derive(Serialize)]
pub struct VerifyResponse {
    pub token: String,
    pub expires_at: i64,
}

/// `POST /auth/verify` — completes the sign-in by verifying the wallet
/// signature against the pending challenge and issuing an auth token.
pub async fn handler(
    State(state): State<AppState>,
    Json(req): Json<VerifyRequest>,
) -> ApiResult<Json<VerifyResponse>> {
    let address: Address = req
        .address
        .parse()
        .map_err(|_| ApiError::bad_request("invalid address"))?;
    let addr_hex = format_address(&address);

    let sig_hex = req.signature.trim_start_matches("0x");
    let sig_bytes = hex::decode(sig_hex)
        .map_err(|_| ApiError::bad_request("signature is not valid hex"))?;
    if sig_bytes.len() != 65 {
        return Err(ApiError::bad_request("signature must be 65 bytes"));
    }
    let signature = Signature::try_from(sig_bytes.as_slice())
        .map_err(|_| ApiError::bad_request("invalid signature encoding"))?;

    let now = now_unix();

    let challenge = sqlx::query("SELECT nonce, expires_at FROM challenges WHERE eoa_address = ?")
        .bind(&addr_hex)
        .fetch_optional(&state.pool)
        .await?;
    let Some(row) = challenge else {
        return Err(ApiError::not_found("no pending challenge for this address"));
    };
    let nonce: String = row.try_get("nonce")?;
    let expires_at: i64 = row.try_get("expires_at")?;
    if now > expires_at {
        // best-effort cleanup, then reject
        let _ = sqlx::query("DELETE FROM challenges WHERE eoa_address = ?")
            .bind(&addr_hex)
            .execute(&state.pool)
            .await;
        return Err(ApiError::bad_request("challenge expired"));
    }

    let message = build_signin_message(&state.config.domain, &address, &nonce);
    let recovered = signature
        .recover_address_from_msg(message.as_bytes())
        .map_err(|_| ApiError::bad_request("signature recovery failed"))?;
    if recovered != address {
        return Err(ApiError::unauthorized("signature does not match address"));
    }

    // ensure a users row exists. New EOAs default to role = 'user'.
    sqlx::query(
        "INSERT OR IGNORE INTO users (eoa_address, role_id, caller_info_json, created_at)
         VALUES (?, (SELECT id FROM roles WHERE name = 'user'), '{}', ?)",
    )
    .bind(&addr_hex)
    .bind(now)
    .execute(&state.pool)
    .await?;

    let mut token_bytes = [0u8; 32];
    rand::thread_rng().fill_bytes(&mut token_bytes);
    let token = hex::encode(token_bytes);
    let token_hash = hex::encode(Sha256::digest(token.as_bytes()));
    let token_expires = now + state.config.token_ttl.as_secs() as i64;

    sqlx::query(
        "INSERT INTO auth_tokens (token_hash, eoa_address, issued_at, expires_at)
         VALUES (?, ?, ?, ?)",
    )
    .bind(&token_hash)
    .bind(&addr_hex)
    .bind(now)
    .bind(token_expires)
    .execute(&state.pool)
    .await?;

    sqlx::query("DELETE FROM challenges WHERE eoa_address = ?")
        .bind(&addr_hex)
        .execute(&state.pool)
        .await?;

    Ok(Json(VerifyResponse {
        token,
        expires_at: token_expires,
    }))
}
