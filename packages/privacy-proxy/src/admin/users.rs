use axum::extract::{Path, Query, State};
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use serde::{Deserialize, Serialize};
use sqlx::Row;

use crate::admin::normalize_address;
use crate::db::now_unix;
use crate::error::{ApiError, ApiResult};
use crate::state::AppState;

#[derive(Serialize)]
pub struct UserView {
    pub eoa_address: String,
    pub role: String,
    pub caller_info: serde_json::Value,
    pub created_at: i64,
}

#[derive(Deserialize)]
pub struct ListUsersQuery {
    pub role: Option<String>,
    pub limit: Option<i64>,
    pub offset: Option<i64>,
}

#[derive(Deserialize)]
pub struct UpsertUserReq {
    pub role: String,
    #[serde(default)]
    pub caller_info: Option<serde_json::Value>,
}

/// Capability 13: `GET /admin/users`
pub async fn list_users(
    State(state): State<AppState>,
    Query(q): Query<ListUsersQuery>,
) -> ApiResult<Json<Vec<UserView>>> {
    let limit = q.limit.unwrap_or(100).clamp(1, 1000);
    let offset = q.offset.unwrap_or(0).max(0);

    let rows = if let Some(role) = &q.role {
        sqlx::query(
            "SELECT u.eoa_address, r.name AS role, u.caller_info_json, u.created_at
             FROM users u JOIN roles r ON r.id = u.role_id
             WHERE r.name = ?
             ORDER BY u.created_at DESC
             LIMIT ? OFFSET ?",
        )
        .bind(role)
        .bind(limit)
        .bind(offset)
        .fetch_all(&state.pool)
        .await?
    } else {
        sqlx::query(
            "SELECT u.eoa_address, r.name AS role, u.caller_info_json, u.created_at
             FROM users u JOIN roles r ON r.id = u.role_id
             ORDER BY u.created_at DESC
             LIMIT ? OFFSET ?",
        )
        .bind(limit)
        .bind(offset)
        .fetch_all(&state.pool)
        .await?
    };

    let out = rows
        .into_iter()
        .map(|r| {
            let info_str: String = r.get("caller_info_json");
            UserView {
                eoa_address: r.get("eoa_address"),
                role: r.get("role"),
                caller_info: serde_json::from_str(&info_str)
                    .unwrap_or(serde_json::Value::Null),
                created_at: r.get("created_at"),
            }
        })
        .collect();
    Ok(Json(out))
}

/// Capability 14: `GET /admin/users/:eoa`
pub async fn get_user(
    State(state): State<AppState>,
    Path(eoa): Path<String>,
) -> ApiResult<Json<UserView>> {
    let eoa = normalize_address(&eoa)?;
    let row = sqlx::query(
        "SELECT u.eoa_address, r.name AS role, u.caller_info_json, u.created_at
         FROM users u JOIN roles r ON r.id = u.role_id
         WHERE u.eoa_address = ?",
    )
    .bind(&eoa)
    .fetch_optional(&state.pool)
    .await?
    .ok_or_else(|| ApiError::not_found("user"))?;
    let info_str: String = row.get("caller_info_json");
    Ok(Json(UserView {
        eoa_address: row.get("eoa_address"),
        role: row.get("role"),
        caller_info: serde_json::from_str(&info_str).unwrap_or(serde_json::Value::Null),
        created_at: row.get("created_at"),
    }))
}

/// Capability 15 (and 16, via `role: "admin"`): `PUT /admin/users/:eoa`
pub async fn upsert_user(
    State(state): State<AppState>,
    Path(eoa): Path<String>,
    Json(req): Json<UpsertUserReq>,
) -> ApiResult<Json<UserView>> {
    let eoa = normalize_address(&eoa)?;
    let role_row = sqlx::query("SELECT id FROM roles WHERE name = ?")
        .bind(&req.role)
        .fetch_optional(&state.pool)
        .await?
        .ok_or_else(|| ApiError::bad_request("unknown role"))?;
    let role_id: i64 = role_row.get("id");
    let info = req.caller_info.unwrap_or_else(|| serde_json::json!({}));
    let info_str = info.to_string();
    let now = now_unix();

    sqlx::query(
        "INSERT INTO users (eoa_address, role_id, caller_info_json, created_at)
         VALUES (?, ?, ?, ?)
         ON CONFLICT(eoa_address) DO UPDATE SET
            role_id = excluded.role_id,
            caller_info_json = excluded.caller_info_json",
    )
    .bind(&eoa)
    .bind(role_id)
    .bind(&info_str)
    .bind(now)
    .execute(&state.pool)
    .await?;

    let row = sqlx::query(
        "SELECT u.eoa_address, r.name AS role, u.caller_info_json, u.created_at
         FROM users u JOIN roles r ON r.id = u.role_id
         WHERE u.eoa_address = ?",
    )
    .bind(&eoa)
    .fetch_one(&state.pool)
    .await?;
    let info_str: String = row.get("caller_info_json");
    Ok(Json(UserView {
        eoa_address: row.get("eoa_address"),
        role: row.get("role"),
        caller_info: serde_json::from_str(&info_str).unwrap_or(serde_json::Value::Null),
        created_at: row.get("created_at"),
    }))
}

/// Capability 17: `DELETE /admin/users/:eoa`
pub async fn delete_user(
    State(state): State<AppState>,
    Path(eoa): Path<String>,
) -> ApiResult<StatusCode> {
    let eoa = normalize_address(&eoa)?;
    let res = sqlx::query("DELETE FROM users WHERE eoa_address = ?")
        .bind(&eoa)
        .execute(&state.pool)
        .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("user"));
    }
    Ok(StatusCode::NO_CONTENT)
}

/// Capability 18: `DELETE /admin/users/:eoa/tokens`
pub async fn revoke_tokens(
    State(state): State<AppState>,
    Path(eoa): Path<String>,
) -> ApiResult<impl IntoResponse> {
    let eoa = normalize_address(&eoa)?;
    let res = sqlx::query("DELETE FROM auth_tokens WHERE eoa_address = ?")
        .bind(&eoa)
        .execute(&state.pool)
        .await?;
    Ok(Json(serde_json::json!({ "revoked": res.rows_affected() })))
}
