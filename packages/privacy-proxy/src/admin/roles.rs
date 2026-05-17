use axum::extract::State;
use axum::Json;
use serde::Serialize;
use sqlx::Row;

use crate::error::ApiResult;
use crate::state::AppState;

#[derive(Serialize)]
pub struct Role {
    pub id: i64,
    pub name: String,
}

/// Capability 11: `GET /admin/roles`. Enumerates the set of role
/// names declared in `src/roles.rs::ROLES`. Used by clients/UIs to
/// populate role selectors on other endpoints.
pub async fn list_roles(State(state): State<AppState>) -> ApiResult<Json<Vec<Role>>> {
    let rows = sqlx::query("SELECT id, name FROM roles ORDER BY id")
        .fetch_all(&state.pool)
        .await?;
    let out = rows
        .into_iter()
        .map(|r| Role {
            id: r.get("id"),
            name: r.get("name"),
        })
        .collect();
    Ok(Json(out))
}
