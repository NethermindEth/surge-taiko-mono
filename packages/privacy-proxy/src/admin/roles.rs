use axum::extract::{Path, State};
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use serde::{Deserialize, Serialize};
use sqlx::Row;

use crate::error::{ApiError, ApiResult};
use crate::state::AppState;

#[derive(Serialize)]
pub struct Role {
    pub id: i64,
    pub name: String,
}

#[derive(Deserialize)]
pub struct CreateRoleReq {
    pub name: String,
}

/// Capability 10: `GET /admin/roles`
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

/// Capability 11: `POST /admin/roles { name }`
pub async fn create_role(
    State(state): State<AppState>,
    Json(req): Json<CreateRoleReq>,
) -> ApiResult<impl IntoResponse> {
    let name = req.name.trim();
    if name.is_empty() {
        return Err(ApiError::bad_request("role name must not be empty"));
    }
    sqlx::query("INSERT INTO roles (name) VALUES (?)")
        .bind(name)
        .execute(&state.pool)
        .await?;
    let row = sqlx::query("SELECT id, name FROM roles WHERE name = ?")
        .bind(name)
        .fetch_one(&state.pool)
        .await?;
    Ok((
        StatusCode::CREATED,
        Json(Role {
            id: row.get("id"),
            name: row.get("name"),
        }),
    ))
}

/// Capability 12: `DELETE /admin/roles/:id`
pub async fn delete_role(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> ApiResult<StatusCode> {
    // Reject if any user or rule entry references the role.
    let user_refs: i64 = sqlx::query("SELECT COUNT(*) AS c FROM users WHERE role_id = ?")
        .bind(id)
        .fetch_one(&state.pool)
        .await?
        .get("c");
    let entry_refs: i64 =
        sqlx::query("SELECT COUNT(*) AS c FROM access_rule_entries WHERE role_id = ?")
            .bind(id)
            .fetch_one(&state.pool)
            .await?
            .get("c");
    if user_refs > 0 || entry_refs > 0 {
        return Err(ApiError::conflict(
            "role is referenced by users or rule entries",
        ));
    }
    let res = sqlx::query("DELETE FROM roles WHERE id = ?")
        .bind(id)
        .execute(&state.pool)
        .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("role"));
    }
    Ok(StatusCode::NO_CONTENT)
}
