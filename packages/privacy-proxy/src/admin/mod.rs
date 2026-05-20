pub mod lambdas;
pub mod members;
pub mod middleware;
pub mod registry;
pub mod roles;

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
        .route(
            "/admin/registry/lambdas",
            get(lambdas::list_lambdas).post(lambdas::create_lambda),
        )
        .route(
            "/admin/registry/lambdas/:id",
            get(lambdas::get_lambda).delete(lambdas::delete_lambda),
        )
        .route(
            "/admin/registry/role-attributes",
            get(lambdas::list_role_attributes),
        )
        .route(
            "/admin/registry/synthetic-selectors",
            get(lambdas::list_synthetic_selectors),
        )
        .route("/admin/roles", get(roles::list_roles))
        .route("/admin/members", get(members::list_members))
        .route(
            "/admin/members/:eoa",
            get(members::get_member).put(members::upsert_member).delete(members::delete_member),
        )
        .route("/admin/members/:eoa/tokens", delete(members::revoke_tokens))
        .route_layer(axum::middleware::from_fn(middleware::admin_gate))
}

/// On every startup, ensure every EOA in `ADMIN_EOAS` is promoted to
/// the `admin` role in the DB. Idempotent. If the EOA previously
/// existed as a `user`, its `user_attributes` row is deleted.
///
/// This is **add-only** despite the name: removing an EOA from
/// `ADMIN_EOAS` does **not** demote them on the next boot. To revoke
/// an admin, call `DELETE /admin/members/:eoa` or
/// `PUT /admin/members/:eoa { role: "user" }` — and also drop the
/// EOA from the env var so it isn't re-promoted.
pub async fn reconcile_seed_admins(pool: &Pool, admin_eoas: &[Address]) -> Result<()> {
    if admin_eoas.is_empty() {
        tracing::warn!("ADMIN_EOAS is empty — no seed admins will exist on this boot");
        return Ok(());
    }
    let now = now_unix();
    for eoa in admin_eoas {
        let addr_hex = format_address(eoa);
        let mut tx = pool.begin().await?;
        sqlx::query(
            "INSERT INTO members (eoa_address, role_id, created_at)
             VALUES (?, (SELECT id FROM roles WHERE name = 'admin'), ?)
             ON CONFLICT(eoa_address) DO UPDATE
             SET role_id = (SELECT id FROM roles WHERE name = 'admin')",
        )
        .bind(&addr_hex)
        .bind(now)
        .execute(&mut *tx)
        .await
        .with_context(|| format!("failed to upsert seed admin {addr_hex}"))?;
        sqlx::query("DELETE FROM user_attributes WHERE eoa_address = ?")
            .bind(&addr_hex)
            .execute(&mut *tx)
            .await
            .with_context(|| {
                format!("failed to drop user_attributes for seed admin {addr_hex}")
            })?;
        tx.commit().await?;
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
