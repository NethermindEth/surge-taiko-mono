use axum::extract::{Path, Query, State};
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use serde::{Deserialize, Serialize};
use sqlx::Row;

use crate::admin::normalize_address;
use crate::db::now_unix;
use crate::error::{ApiError, ApiResult};
use crate::roles::{ROLE_ADMIN, ROLE_USER, ROLES};
use crate::state::AppState;

#[derive(Serialize)]
pub struct UserAttributesView {
    pub kyc: bool,
    pub blacklisted: bool,
}

/// Typed member record. `attributes` is `null` for admin (identity-only)
/// and an object for user.
#[derive(Serialize)]
pub struct MemberView {
    pub eoa_address: String,
    pub role: &'static str,
    pub attributes: Option<UserAttributesView>,
    pub created_at: i64,
}

#[derive(Deserialize)]
pub struct ListMembersQuery {
    pub role: Option<String>,
    pub limit: Option<i64>,
    pub offset: Option<i64>,
}

#[derive(Deserialize)]
pub struct UpsertAttributes {
    #[serde(default)]
    pub kyc: Option<bool>,
    #[serde(default)]
    pub blacklisted: Option<bool>,
}

#[derive(Deserialize)]
pub struct UpsertMemberReq {
    pub role: String,
    #[serde(default)]
    pub attributes: Option<UpsertAttributes>,
}

async fn load_member(pool: &crate::db::Pool, eoa: &str) -> Result<MemberView, ApiError> {
    let row = sqlx::query(
        "SELECT m.eoa_address, r.name AS role, m.created_at
         FROM members m JOIN roles r ON r.id = m.role_id
         WHERE m.eoa_address = ?",
    )
    .bind(eoa)
    .fetch_optional(pool)
    .await?
    .ok_or_else(|| ApiError::not_found("member"))?;
    let role: String = row.get("role");
    let (role_static, attributes): (&'static str, Option<UserAttributesView>) = match role.as_str() {
        ROLE_ADMIN => (ROLE_ADMIN, None),
        ROLE_USER => {
            let attrs = sqlx::query(
                "SELECT kyc, blacklisted FROM user_attributes WHERE eoa_address = ?",
            )
            .bind(eoa)
            .fetch_optional(pool)
            .await?;
            let view = match attrs {
                Some(r) => {
                    let kyc: i64 = r.get("kyc");
                    let bl: i64 = r.get("blacklisted");
                    UserAttributesView {
                        kyc: kyc != 0,
                        blacklisted: bl != 0,
                    }
                }
                None => UserAttributesView {
                    kyc: false,
                    blacklisted: false,
                },
            };
            (ROLE_USER, Some(view))
        }
        other => {
            return Err(ApiError::Internal(anyhow::anyhow!(
                "member has unknown role `{other}`"
            )));
        }
    };
    Ok(MemberView {
        eoa_address: row.get("eoa_address"),
        role: role_static,
        attributes,
        created_at: row.get("created_at"),
    })
}

fn validate_role(role: &str) -> Result<&'static str, ApiError> {
    ROLES
        .iter()
        .copied()
        .find(|r| *r == role)
        .ok_or_else(|| ApiError::bad_request(format!("unknown role `{role}`")))
}

/// Capability 12: `GET /admin/members`
pub async fn list_members(
    State(state): State<AppState>,
    Query(q): Query<ListMembersQuery>,
) -> ApiResult<Json<Vec<MemberView>>> {
    let limit = q.limit.unwrap_or(100).clamp(1, 1000);
    let offset = q.offset.unwrap_or(0).max(0);

    let rows = if let Some(role) = &q.role {
        sqlx::query(
            "SELECT m.eoa_address
             FROM members m JOIN roles r ON r.id = m.role_id
             WHERE r.name = ?
             ORDER BY m.created_at DESC
             LIMIT ? OFFSET ?",
        )
        .bind(role)
        .bind(limit)
        .bind(offset)
        .fetch_all(&state.pool)
        .await?
    } else {
        sqlx::query(
            "SELECT m.eoa_address
             FROM members m
             ORDER BY m.created_at DESC
             LIMIT ? OFFSET ?",
        )
        .bind(limit)
        .bind(offset)
        .fetch_all(&state.pool)
        .await?
    };

    let mut out = Vec::with_capacity(rows.len());
    for r in rows {
        let eoa: String = r.get("eoa_address");
        out.push(load_member(&state.pool, &eoa).await?);
    }
    Ok(Json(out))
}

/// Capability 13: `GET /admin/members/:eoa`
pub async fn get_member(
    State(state): State<AppState>,
    Path(eoa): Path<String>,
) -> ApiResult<Json<MemberView>> {
    let eoa = normalize_address(&eoa)?;
    Ok(Json(load_member(&state.pool, &eoa).await?))
}

/// Capability 14: `PUT /admin/members/:eoa`
pub async fn upsert_member(
    State(state): State<AppState>,
    Path(eoa): Path<String>,
    Json(req): Json<UpsertMemberReq>,
) -> ApiResult<Json<MemberView>> {
    let eoa = normalize_address(&eoa)?;
    let role = validate_role(&req.role)?;
    if role == ROLE_ADMIN && req.attributes.is_some() {
        return Err(ApiError::bad_request(
            "admin role does not accept attributes",
        ));
    }
    let now = now_unix();

    let mut tx = state.pool.begin().await?;
    sqlx::query(
        "INSERT INTO members (eoa_address, role_id, created_at)
         VALUES (?, (SELECT id FROM roles WHERE name = ?), ?)
         ON CONFLICT(eoa_address) DO UPDATE SET
            role_id = (SELECT id FROM roles WHERE name = ?)",
    )
    .bind(&eoa)
    .bind(role)
    .bind(now)
    .bind(role)
    .execute(&mut *tx)
    .await?;

    match role {
        ROLE_ADMIN => {
            sqlx::query("DELETE FROM user_attributes WHERE eoa_address = ?")
                .bind(&eoa)
                .execute(&mut *tx)
                .await?;
        }
        ROLE_USER => {
            // Omitted attribute fields preserve existing values; explicit
            // values overwrite; defaults apply on first insert.
            let existing = sqlx::query(
                "SELECT kyc, blacklisted FROM user_attributes WHERE eoa_address = ?",
            )
            .bind(&eoa)
            .fetch_optional(&mut *tx)
            .await?;
            let (existing_kyc, existing_bl) = match existing {
                Some(r) => {
                    let k: i64 = r.get("kyc");
                    let b: i64 = r.get("blacklisted");
                    (k != 0, b != 0)
                }
                None => (false, false),
            };
            let attrs = req.attributes.unwrap_or(UpsertAttributes {
                kyc: None,
                blacklisted: None,
            });
            let final_kyc = attrs.kyc.unwrap_or(existing_kyc);
            let final_bl = attrs.blacklisted.unwrap_or(existing_bl);
            sqlx::query(
                "INSERT INTO user_attributes (eoa_address, kyc, blacklisted)
                 VALUES (?, ?, ?)
                 ON CONFLICT(eoa_address) DO UPDATE SET
                    kyc = excluded.kyc,
                    blacklisted = excluded.blacklisted",
            )
            .bind(&eoa)
            .bind(final_kyc as i64)
            .bind(final_bl as i64)
            .execute(&mut *tx)
            .await?;
        }
        _ => unreachable!("validate_role gates names"),
    }

    tx.commit().await?;
    Ok(Json(load_member(&state.pool, &eoa).await?))
}

/// Capability 15: `DELETE /admin/members/:eoa`
pub async fn delete_member(
    State(state): State<AppState>,
    Path(eoa): Path<String>,
) -> ApiResult<StatusCode> {
    let eoa = normalize_address(&eoa)?;
    let res = sqlx::query("DELETE FROM members WHERE eoa_address = ?")
        .bind(&eoa)
        .execute(&state.pool)
        .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("member"));
    }
    Ok(StatusCode::NO_CONTENT)
}

/// Capability 16: `DELETE /admin/members/:eoa/tokens`
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
