use axum::extract::{Path, Query, State};
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use serde::{Deserialize, Serialize};
use sqlx::Row;

use crate::acl::lambdas;
use crate::admin::{normalize_address, normalize_selector, validate_mode};
use crate::error::{ApiError, ApiResult};
use crate::state::AppState;

#[derive(Serialize)]
pub struct RuleView {
    pub id: i64,
    pub contract_address: String,
    pub function_selector: String,
    pub mode: String,
    pub entries: Vec<EntryView>,
}

#[derive(Serialize)]
pub struct EntryView {
    pub id: i64,
    pub role: String,
    pub lambda_name: Option<String>,
}

#[derive(Deserialize)]
pub struct ListRulesQuery {
    pub contract: Option<String>,
    pub limit: Option<i64>,
    pub offset: Option<i64>,
}

#[derive(Deserialize)]
pub struct EntryInput {
    pub role: String,
    #[serde(default)]
    pub lambda_name: Option<String>,
}

#[derive(Deserialize)]
pub struct CreateRuleReq {
    pub contract_address: String,
    pub function_selector: String,
    pub mode: String,
    #[serde(default)]
    pub entries: Vec<EntryInput>,
}

#[derive(Deserialize)]
pub struct ReplaceRuleReq {
    pub mode: String,
    #[serde(default)]
    pub entries: Vec<EntryInput>,
}

#[derive(Deserialize)]
pub struct UpdateEntryReq {
    #[serde(default)]
    pub lambda_name: Option<String>,
}

fn ensure_lambda_known(name: &str) -> Result<(), ApiError> {
    if lambdas::lookup(name).is_none() {
        return Err(ApiError::bad_request(format!(
            "unknown lambda `{name}`"
        )));
    }
    Ok(())
}

async fn resolve_role_id(pool: &crate::db::Pool, role: &str) -> Result<i64, ApiError> {
    sqlx::query("SELECT id FROM roles WHERE name = ?")
        .bind(role)
        .fetch_optional(pool)
        .await?
        .map(|r| r.get::<i64, _>("id"))
        .ok_or_else(|| ApiError::bad_request(format!("unknown role `{role}`")))
}

async fn load_rule(pool: &crate::db::Pool, id: i64) -> Result<RuleView, ApiError> {
    let rule_row = sqlx::query(
        "SELECT id, contract_address, function_selector, mode FROM access_rules WHERE id = ?",
    )
    .bind(id)
    .fetch_optional(pool)
    .await?
    .ok_or_else(|| ApiError::not_found("rule"))?;
    let entry_rows = sqlx::query(
        "SELECT e.id, r.name AS role, e.lambda_name
         FROM access_rule_entries e
         JOIN roles r ON r.id = e.role_id
         WHERE e.rule_id = ?
         ORDER BY r.name",
    )
    .bind(id)
    .fetch_all(pool)
    .await?;
    let entries = entry_rows
        .into_iter()
        .map(|r| EntryView {
            id: r.get("id"),
            role: r.get("role"),
            lambda_name: r.get("lambda_name"),
        })
        .collect();
    Ok(RuleView {
        id: rule_row.get("id"),
        contract_address: rule_row.get("contract_address"),
        function_selector: rule_row.get("function_selector"),
        mode: rule_row.get("mode"),
        entries,
    })
}

/// Capability 1: `GET /admin/registry/rules`
pub async fn list_rules(
    State(state): State<AppState>,
    Query(q): Query<ListRulesQuery>,
) -> ApiResult<Json<Vec<RuleView>>> {
    let limit = q.limit.unwrap_or(100).clamp(1, 1000);
    let offset = q.offset.unwrap_or(0).max(0);
    let rows = if let Some(contract) = q.contract.as_deref() {
        let norm = normalize_address(contract)?;
        sqlx::query(
            "SELECT id FROM access_rules WHERE contract_address = ? ORDER BY id LIMIT ? OFFSET ?",
        )
        .bind(&norm)
        .bind(limit)
        .bind(offset)
        .fetch_all(&state.pool)
        .await?
    } else {
        sqlx::query("SELECT id FROM access_rules ORDER BY id LIMIT ? OFFSET ?")
            .bind(limit)
            .bind(offset)
            .fetch_all(&state.pool)
            .await?
    };
    let mut out = Vec::with_capacity(rows.len());
    for r in rows {
        let id: i64 = r.get("id");
        out.push(load_rule(&state.pool, id).await?);
    }
    Ok(Json(out))
}

/// Capability 2: `GET /admin/registry/rules/:id`
pub async fn get_rule(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> ApiResult<Json<RuleView>> {
    Ok(Json(load_rule(&state.pool, id).await?))
}

/// Capability 3: `POST /admin/registry/rules`
pub async fn create_rule(
    State(state): State<AppState>,
    Json(req): Json<CreateRuleReq>,
) -> ApiResult<impl IntoResponse> {
    let contract = normalize_address(&req.contract_address)?;
    let selector = normalize_selector(&req.function_selector)?;
    let mode = validate_mode(&req.mode)?;
    // Resolve every entry's role + validate lambdas *before* opening the
    // transaction. Otherwise resolve_role_id holds a second pool connection
    // while the tx holds the first, which deadlocks under pools sized to 1.
    let mut entries_resolved: Vec<(i64, Option<String>)> = Vec::with_capacity(req.entries.len());
    for e in &req.entries {
        if let Some(name) = &e.lambda_name {
            ensure_lambda_known(name)?;
        }
        let role_id = resolve_role_id(&state.pool, &e.role).await?;
        entries_resolved.push((role_id, e.lambda_name.clone()));
    }

    let mut tx = state.pool.begin().await?;
    sqlx::query(
        "INSERT INTO access_rules (contract_address, function_selector, mode) VALUES (?, ?, ?)",
    )
    .bind(&contract)
    .bind(&selector)
    .bind(mode)
    .execute(&mut *tx)
    .await?;
    let rule_id: i64 = sqlx::query("SELECT last_insert_rowid() AS id")
        .fetch_one(&mut *tx)
        .await?
        .get("id");

    for (role_id, lambda_name) in &entries_resolved {
        sqlx::query(
            "INSERT INTO access_rule_entries (rule_id, role_id, lambda_name) VALUES (?, ?, ?)",
        )
        .bind(rule_id)
        .bind(role_id)
        .bind(lambda_name.as_deref())
        .execute(&mut *tx)
        .await?;
    }
    tx.commit().await?;

    let view = load_rule(&state.pool, rule_id).await?;
    Ok((StatusCode::CREATED, Json(view)))
}

/// Capability 4: `PUT /admin/registry/rules/:id`
pub async fn replace_rule(
    State(state): State<AppState>,
    Path(id): Path<i64>,
    Json(req): Json<ReplaceRuleReq>,
) -> ApiResult<Json<RuleView>> {
    let mode = validate_mode(&req.mode)?;
    // Same deadlock avoidance as create_rule: resolve roles before tx.
    let mut entries_resolved: Vec<(i64, Option<String>)> = Vec::with_capacity(req.entries.len());
    for e in &req.entries {
        if let Some(name) = &e.lambda_name {
            ensure_lambda_known(name)?;
        }
        let role_id = resolve_role_id(&state.pool, &e.role).await?;
        entries_resolved.push((role_id, e.lambda_name.clone()));
    }

    let mut tx = state.pool.begin().await?;
    let res = sqlx::query("UPDATE access_rules SET mode = ? WHERE id = ?")
        .bind(mode)
        .bind(id)
        .execute(&mut *tx)
        .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("rule"));
    }
    sqlx::query("DELETE FROM access_rule_entries WHERE rule_id = ?")
        .bind(id)
        .execute(&mut *tx)
        .await?;
    for (role_id, lambda_name) in &entries_resolved {
        sqlx::query(
            "INSERT INTO access_rule_entries (rule_id, role_id, lambda_name) VALUES (?, ?, ?)",
        )
        .bind(id)
        .bind(role_id)
        .bind(lambda_name.as_deref())
        .execute(&mut *tx)
        .await?;
    }
    tx.commit().await?;

    Ok(Json(load_rule(&state.pool, id).await?))
}

/// Capability 5: `DELETE /admin/registry/rules/:id`
pub async fn delete_rule(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> ApiResult<StatusCode> {
    let res = sqlx::query("DELETE FROM access_rules WHERE id = ?")
        .bind(id)
        .execute(&state.pool)
        .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("rule"));
    }
    Ok(StatusCode::NO_CONTENT)
}

/// Capability 6: `POST /admin/registry/rules/:id/entries`
pub async fn add_entry(
    State(state): State<AppState>,
    Path(rule_id): Path<i64>,
    Json(req): Json<EntryInput>,
) -> ApiResult<impl IntoResponse> {
    if let Some(name) = &req.lambda_name {
        ensure_lambda_known(name)?;
    }
    let _ = sqlx::query("SELECT 1 FROM access_rules WHERE id = ?")
        .bind(rule_id)
        .fetch_optional(&state.pool)
        .await?
        .ok_or_else(|| ApiError::not_found("rule"))?;

    let role_id = resolve_role_id(&state.pool, &req.role).await?;
    sqlx::query(
        "INSERT INTO access_rule_entries (rule_id, role_id, lambda_name) VALUES (?, ?, ?)",
    )
    .bind(rule_id)
    .bind(role_id)
    .bind(req.lambda_name.as_deref())
    .execute(&state.pool)
    .await?;
    let entry_id: i64 = sqlx::query("SELECT last_insert_rowid() AS id")
        .fetch_one(&state.pool)
        .await?
        .get("id");
    Ok((
        StatusCode::CREATED,
        Json(EntryView {
            id: entry_id,
            role: req.role,
            lambda_name: req.lambda_name,
        }),
    ))
}

/// Capability 7: `PUT /admin/registry/rules/:id/entries/:entry_id`
pub async fn update_entry(
    State(state): State<AppState>,
    Path((rule_id, entry_id)): Path<(i64, i64)>,
    Json(req): Json<UpdateEntryReq>,
) -> ApiResult<Json<EntryView>> {
    if let Some(name) = &req.lambda_name {
        ensure_lambda_known(name)?;
    }
    let res = sqlx::query(
        "UPDATE access_rule_entries SET lambda_name = ? WHERE id = ? AND rule_id = ?",
    )
    .bind(req.lambda_name.as_deref())
    .bind(entry_id)
    .bind(rule_id)
    .execute(&state.pool)
    .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("entry"));
    }
    let row = sqlx::query(
        "SELECT e.id, r.name AS role, e.lambda_name FROM access_rule_entries e
         JOIN roles r ON r.id = e.role_id WHERE e.id = ?",
    )
    .bind(entry_id)
    .fetch_one(&state.pool)
    .await?;
    Ok(Json(EntryView {
        id: row.get("id"),
        role: row.get("role"),
        lambda_name: row.get("lambda_name"),
    }))
}

/// Capability 8: `DELETE /admin/registry/rules/:id/entries/:entry_id`
pub async fn delete_entry(
    State(state): State<AppState>,
    Path((rule_id, entry_id)): Path<(i64, i64)>,
) -> ApiResult<StatusCode> {
    let res = sqlx::query("DELETE FROM access_rule_entries WHERE id = ? AND rule_id = ?")
        .bind(entry_id)
        .bind(rule_id)
        .execute(&state.pool)
        .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("entry"));
    }
    Ok(StatusCode::NO_CONTENT)
}
