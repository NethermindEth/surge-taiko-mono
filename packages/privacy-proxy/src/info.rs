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
    pub chain_id: Option<u64>,
    pub upstream_ok: bool,
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
/// No authentication required. Returns the static fields plus
/// `chain_id: null, upstream_ok: false` when the upstream is unreachable,
/// so probers can tell "not a privacy-proxy" from "upstream offline".
pub async fn handler(State(state): State<AppState>) -> ApiResult<Json<InfoResponse>> {
    let (chain_id, upstream_ok) = match state.upstream.chain_id().await {
        Ok(id) => (Some(id), true),
        Err(e) => {
            tracing::warn!("upstream chain_id unreachable for /info: {e}");
            (None, false)
        }
    };

    Ok(Json(InfoResponse {
        name: PROXY_NAME,
        version: PROXY_VERSION,
        chain_id,
        upstream_ok,
        domain: state.config.domain.clone(),
        auth: AuthInfo {
            scheme: "bearer",
            challenge_path: "/auth/challenge",
            verify_path: "/auth/verify",
        },
    }))
}
