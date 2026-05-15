pub mod lambdas;
pub mod middleware;
pub mod registry;
pub mod roles;
pub mod users;

use alloy::primitives::Address;
use anyhow::{Context, Result};
use axum::routing::{delete, get, post, put};
use axum::Router;

use crate::auth::format_address;
use crate::db::{now_unix, Pool};
use crate::error::ApiError;
use crate::state::AppState;

pub fn router() -> Router<AppState> {
    Router::new()
        // capability 9 — list in-build lambdas
        .route("/admin/registry/lambdas", get(lambdas::list_lambdas))
        // capability 19 — list synthetic selectors for gated RPC methods
        .route(
            "/admin/registry/synthetic-selectors",
            get(lambdas::list_synthetic_selectors),
        )
        // capabilities 1-8 — access rules + entries
        .route(
            "/admin/registry/rules",
            get(registry::list_rules).post(registry::create_rule),
        )
        .route(
            "/admin/registry/rules/:id",
            get(registry::get_rule)
                .put(registry::replace_rule)
                .delete(registry::delete_rule),
        )
        .route(
            "/admin/registry/rules/:id/entries",
            post(registry::add_entry),
        )
        .route(
            "/admin/registry/rules/:id/entries/:entry_id",
            put(registry::update_entry).delete(registry::delete_entry),
        )
        // capabilities 10-12 — roles
        .route("/admin/roles", get(roles::list_roles).post(roles::create_role))
        .route("/admin/roles/:id", delete(roles::delete_role))
        // capabilities 13-18 — users
        .route("/admin/users", get(users::list_users))
        .route(
            "/admin/users/:eoa",
            get(users::get_user).put(users::upsert_user).delete(users::delete_user),
        )
        .route("/admin/users/:eoa/tokens", delete(users::revoke_tokens))
        .route_layer(axum::middleware::from_fn(middleware::admin_gate))
}

/// On every startup, reconcile the contents of `ADMIN_EOAS` so the seed
/// admins are always promoted in DB. Idempotent.
pub async fn reconcile_seed_admins(pool: &Pool, admin_eoas: &[Address]) -> Result<()> {
    if admin_eoas.is_empty() {
        tracing::warn!("ADMIN_EOAS is empty — no seed admins will exist on this boot");
        return Ok(());
    }
    let now = now_unix();
    for eoa in admin_eoas {
        let addr_hex = format_address(eoa);
        sqlx::query(
            "INSERT INTO users (eoa_address, role_id, caller_info_json, created_at)
             VALUES (?, (SELECT id FROM roles WHERE name = 'admin'), '{}', ?)
             ON CONFLICT(eoa_address) DO UPDATE
             SET role_id = (SELECT id FROM roles WHERE name = 'admin')",
        )
        .bind(&addr_hex)
        .bind(now)
        .execute(pool)
        .await
        .with_context(|| format!("failed to upsert seed admin {addr_hex}"))?;
        tracing::info!("reconciled seed admin: {addr_hex}");
    }
    Ok(())
}

pub(crate) fn normalize_address(s: &str) -> Result<String, ApiError> {
    let addr: Address = s
        .parse()
        .map_err(|_| ApiError::bad_request("invalid address"))?;
    Ok(format_address(&addr))
}

pub(crate) fn normalize_selector(s: &str) -> Result<String, ApiError> {
    // Allow operators to use a gated JSON-RPC method name in place of the
    // synthetic 4-byte selector. The server stores the synthetic value.
    if let Some(method) = crate::rpc::gated_methods::lookup_by_method(s) {
        return Ok(format!("0x{}", hex::encode(method.selector)));
    }
    let trimmed = s.trim_start_matches("0x").to_ascii_lowercase();
    if trimmed.len() != 8 {
        return Err(ApiError::bad_request("selector must be 4 bytes (8 hex chars)"));
    }
    hex::decode(&trimmed).map_err(|_| ApiError::bad_request("selector is not hex"))?;
    Ok(format!("0x{trimmed}"))
}

pub(crate) fn validate_mode(mode: &str) -> Result<&'static str, ApiError> {
    match mode {
        "allow" => Ok("allow"),
        "deny" => Ok("deny"),
        _ => Err(ApiError::bad_request("mode must be 'allow' or 'deny'")),
    }
}
