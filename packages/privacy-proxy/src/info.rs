use axum::extract::State;
use axum::Json;
use serde::Serialize;

use crate::error::ApiResult;
use crate::state::AppState;

const PROXY_NAME: &str = "privacy-proxy";
const PROXY_VERSION: &str = env!("CARGO_PKG_VERSION");

#[derive(Serialize)]
pub struct InfoResponse {
    pub name: &'static str,
    pub version: &'static str,
    pub chain_id: u64,
    pub domain: String,
    pub auth: AuthInfo,
}

#[derive(Serialize)]
pub struct AuthInfo {
    pub scheme: &'static str,
    pub challenge_path: &'static str,
    pub verify_path: &'static str,
}

/// `GET /info` — public identification endpoint. Wallets probe this to
/// detect that an RPC URL is a privacy-proxy and learn the auth scheme.
/// No authentication required.
pub async fn handler(State(state): State<AppState>) -> ApiResult<Json<InfoResponse>> {
    let chain_id = state.upstream.chain_id().await?;

    Ok(Json(InfoResponse {
        name: PROXY_NAME,
        version: PROXY_VERSION,
        chain_id,
        domain: state.config.domain.clone(),
        auth: AuthInfo {
            scheme: "bearer",
            challenge_path: "/auth/challenge",
            verify_path: "/auth/verify",
        },
    }))
}
