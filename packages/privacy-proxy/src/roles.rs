use anyhow::{Context, Result};

use crate::db::Pool;

pub const ROLE_ADMIN: &str = "admin";
pub const ROLE_USER: &str = "user";

/// Every role recognized by this build. Adding a role is a code change:
/// add the name here, declare its attribute type (and table, if any),
/// extend `CallerAttributes`, and add its lambda registry under
/// `src/acl/lambdas/<role>/`.
pub const ROLES: &[&str] = &[ROLE_ADMIN, ROLE_USER];

/// Reconcile `ROLES` into the `roles` table at boot. Idempotent.
pub async fn reconcile_roles(pool: &Pool) -> Result<()> {
    for name in ROLES {
        sqlx::query("INSERT OR IGNORE INTO roles (name) VALUES (?)")
            .bind(name)
            .execute(pool)
            .await
            .with_context(|| format!("failed to reconcile role `{name}`"))?;
    }
    Ok(())
}
