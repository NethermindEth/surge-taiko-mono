use axum::extract::{Path, Query, State};
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use serde::{Deserialize, Serialize};
use sqlx::Row;

use crate::admin::{normalize_address, normalize_selector, validate_mode};
use crate::error::{ApiError, ApiResult};
use crate::state::AppState;

const MAX_ENTRIES_PER_RULE: usize = 64;

#[derive(Serialize)]
pub struct RuleView {
    pub id: i64,
    pub name: String,
    pub description: Option<String>,
    pub selector: String,
    pub mode: String,
    pub entries: Vec<EntryView>,
    pub binding_count: i64,
}

#[derive(Serialize)]
pub struct EntryView {
    pub id: i64,
    pub role: String,
    pub lambda_id: Option<i64>,
    pub lambda_name: Option<String>,
}

#[derive(Serialize)]
pub struct BindingView {
    pub id: i64,
    pub contract_address: String,
    pub selector: String,
    pub rule_id: i64,
    pub rule_name: String,
    pub mode: String,
}

#[derive(Deserialize)]
pub struct ListRulesQuery {
    pub limit: Option<i64>,
    pub offset: Option<i64>,
}

#[derive(Deserialize)]
pub struct ListBindingsQuery {
    pub contract: Option<String>,
    pub rule_id: Option<i64>,
    pub limit: Option<i64>,
    pub offset: Option<i64>,
}

#[derive(Deserialize)]
pub struct EntryInput {
    pub role: String,
    #[serde(default)]
    pub lambda_id: Option<i64>,
}

#[derive(Deserialize)]
pub struct CreateRuleReq {
    pub name: String,
    #[serde(default)]
    pub description: Option<String>,
    pub selector: String,
    pub mode: String,
    #[serde(default)]
    pub entries: Vec<EntryInput>,
}

#[derive(Deserialize)]
pub struct ReplaceRuleReq {
    pub name: String,
    #[serde(default)]
    pub description: Option<String>,
    pub mode: String,
    #[serde(default)]
    pub entries: Vec<EntryInput>,
}

#[derive(Deserialize)]
pub struct UpdateEntryReq {
    #[serde(default)]
    pub lambda_id: Option<i64>,
}

#[derive(Deserialize)]
pub struct CreateBindingReq {
    pub contract_address: String,
    pub rule_id: i64,
}

async fn ensure_lambda_attachable(
    pool: &crate::db::Pool,
    role_id: i64,
    lambda_id: i64,
) -> Result<(), ApiError> {
    let row = sqlx::query("SELECT role_id FROM lambdas WHERE id = ?")
        .bind(lambda_id)
        .fetch_optional(pool)
        .await?
        .ok_or_else(|| ApiError::bad_request(format!("unknown lambda id {lambda_id}")))?;
    let lambda_role_id: i64 = row.get("role_id");
    if lambda_role_id != role_id {
        return Err(ApiError::bad_request(format!(
            "lambda {lambda_id} belongs to a different role than this entry"
        )));
    }
    Ok(())
}

async fn lookup_entry_role_id(pool: &crate::db::Pool, entry_id: i64) -> Result<i64, ApiError> {
    sqlx::query("SELECT role_id FROM access_rule_entries WHERE id = ?")
        .bind(entry_id)
        .fetch_optional(pool)
        .await?
        .map(|r| r.get::<i64, _>("role_id"))
        .ok_or_else(|| ApiError::not_found("entry"))
}

async fn resolve_role_id(pool: &crate::db::Pool, role: &str) -> Result<i64, ApiError> {
    sqlx::query("SELECT id FROM roles WHERE name = ?")
        .bind(role)
        .fetch_optional(pool)
        .await?
        .map(|r| r.get::<i64, _>("id"))
        .ok_or_else(|| ApiError::bad_request(format!("unknown role `{role}`")))
}

async fn binding_count_for(pool: &crate::db::Pool, rule_id: i64) -> Result<i64, ApiError> {
    let row = sqlx::query("SELECT COUNT(*) AS c FROM contract_rules WHERE rule_id = ?")
        .bind(rule_id)
        .fetch_one(pool)
        .await?;
    Ok(row.get("c"))
}

async fn load_rule(pool: &crate::db::Pool, id: i64) -> Result<RuleView, ApiError> {
    let rule_row = sqlx::query(
        "SELECT id, name, description, selector, mode FROM access_rules WHERE id = ?",
    )
    .bind(id)
    .fetch_optional(pool)
    .await?
    .ok_or_else(|| ApiError::not_found("rule"))?;
    let entry_rows = sqlx::query(
        "SELECT e.id, r.name AS role, e.lambda_id, l.name AS lambda_name
         FROM access_rule_entries e
         JOIN roles r ON r.id = e.role_id
         LEFT JOIN lambdas l ON l.id = e.lambda_id
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
            lambda_id: r.get("lambda_id"),
            lambda_name: r.get("lambda_name"),
        })
        .collect();
    let binding_count = binding_count_for(pool, id).await?;
    Ok(RuleView {
        id: rule_row.get("id"),
        name: rule_row.get("name"),
        description: rule_row.get("description"),
        selector: rule_row.get("selector"),
        mode: rule_row.get("mode"),
        entries,
        binding_count,
    })
}

pub async fn list_rules(
    State(state): State<AppState>,
    Query(q): Query<ListRulesQuery>,
) -> ApiResult<Json<Vec<RuleView>>> {
    let limit = q.limit.unwrap_or(100).clamp(1, 1000);
    let offset = q.offset.unwrap_or(0).max(0);
    let rows = sqlx::query("SELECT id FROM access_rules ORDER BY id LIMIT ? OFFSET ?")
        .bind(limit)
        .bind(offset)
        .fetch_all(&state.pool)
        .await?;
    let mut out = Vec::with_capacity(rows.len());
    for r in rows {
        let id: i64 = r.get("id");
        out.push(load_rule(&state.pool, id).await?);
    }
    Ok(Json(out))
}

pub async fn get_rule(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> ApiResult<Json<RuleView>> {
    Ok(Json(load_rule(&state.pool, id).await?))
}

pub async fn create_rule(
    State(state): State<AppState>,
    Json(req): Json<CreateRuleReq>,
) -> ApiResult<impl IntoResponse> {
    let name = req.name.trim();
    if name.is_empty() {
        return Err(ApiError::bad_request("name must not be empty"));
    }
    let selector = normalize_selector(&req.selector)?;
    let mode = validate_mode(&req.mode)?;
    if req.entries.len() > MAX_ENTRIES_PER_RULE {
        return Err(ApiError::bad_request(format!(
            "too many entries (limit {MAX_ENTRIES_PER_RULE})"
        )));
    }
    let mut entries_resolved: Vec<(i64, Option<i64>)> = Vec::with_capacity(req.entries.len());
    for e in &req.entries {
        let role_id = resolve_role_id(&state.pool, &e.role).await?;
        if let Some(lid) = e.lambda_id {
            ensure_lambda_attachable(&state.pool, role_id, lid).await?;
        }
        entries_resolved.push((role_id, e.lambda_id));
    }

    let mut tx = state.pool.begin().await?;
    sqlx::query(
        "INSERT INTO access_rules (name, description, selector, mode) VALUES (?, ?, ?, ?)",
    )
    .bind(name)
    .bind(req.description.as_deref())
    .bind(&selector)
    .bind(mode)
    .execute(&mut *tx)
    .await?;
    let rule_id: i64 = sqlx::query("SELECT last_insert_rowid() AS id")
        .fetch_one(&mut *tx)
        .await?
        .get("id");

    for (role_id, lambda_id) in &entries_resolved {
        sqlx::query(
            "INSERT INTO access_rule_entries (rule_id, role_id, lambda_id) VALUES (?, ?, ?)",
        )
        .bind(rule_id)
        .bind(role_id)
        .bind(lambda_id)
        .execute(&mut *tx)
        .await?;
    }
    tx.commit().await?;

    let view = load_rule(&state.pool, rule_id).await?;
    Ok((StatusCode::CREATED, Json(view)))
}

/// Replace a rule's metadata (name, description, mode) and its entries in
/// one shot. The rule's `selector` is immutable post-creation because
/// existing `contract_rules` bindings denormalize it.
pub async fn replace_rule(
    State(state): State<AppState>,
    Path(id): Path<i64>,
    Json(req): Json<ReplaceRuleReq>,
) -> ApiResult<Json<RuleView>> {
    let name = req.name.trim();
    if name.is_empty() {
        return Err(ApiError::bad_request("name must not be empty"));
    }
    let mode = validate_mode(&req.mode)?;
    if req.entries.len() > MAX_ENTRIES_PER_RULE {
        return Err(ApiError::bad_request(format!(
            "too many entries (limit {MAX_ENTRIES_PER_RULE})"
        )));
    }
    let mut entries_resolved: Vec<(i64, Option<i64>)> = Vec::with_capacity(req.entries.len());
    for e in &req.entries {
        let role_id = resolve_role_id(&state.pool, &e.role).await?;
        if let Some(lid) = e.lambda_id {
            ensure_lambda_attachable(&state.pool, role_id, lid).await?;
        }
        entries_resolved.push((role_id, e.lambda_id));
    }

    let mut tx = state.pool.begin().await?;
    let res = sqlx::query(
        "UPDATE access_rules SET name = ?, description = ?, mode = ? WHERE id = ?",
    )
    .bind(name)
    .bind(req.description.as_deref())
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
    for (role_id, lambda_id) in &entries_resolved {
        sqlx::query(
            "INSERT INTO access_rule_entries (rule_id, role_id, lambda_id) VALUES (?, ?, ?)",
        )
        .bind(id)
        .bind(role_id)
        .bind(lambda_id)
        .execute(&mut *tx)
        .await?;
    }
    tx.commit().await?;

    Ok(Json(load_rule(&state.pool, id).await?))
}

pub async fn delete_rule(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> ApiResult<StatusCode> {
    let bindings = binding_count_for(&state.pool, id).await?;
    if bindings > 0 {
        return Err(ApiError::conflict(format!(
            "rule is bound to {bindings} contract selector(s); unbind before deleting"
        )));
    }
    let res = sqlx::query("DELETE FROM access_rules WHERE id = ?")
        .bind(id)
        .execute(&state.pool)
        .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("rule"));
    }
    Ok(StatusCode::NO_CONTENT)
}

pub async fn add_entry(
    State(state): State<AppState>,
    Path(rule_id): Path<i64>,
    Json(req): Json<EntryInput>,
) -> ApiResult<impl IntoResponse> {
    let _ = sqlx::query("SELECT 1 FROM access_rules WHERE id = ?")
        .bind(rule_id)
        .fetch_optional(&state.pool)
        .await?
        .ok_or_else(|| ApiError::not_found("rule"))?;

    let role_id = resolve_role_id(&state.pool, &req.role).await?;
    if let Some(lid) = req.lambda_id {
        ensure_lambda_attachable(&state.pool, role_id, lid).await?;
    }

    sqlx::query(
        "INSERT INTO access_rule_entries (rule_id, role_id, lambda_id) VALUES (?, ?, ?)",
    )
    .bind(rule_id)
    .bind(role_id)
    .bind(req.lambda_id)
    .execute(&state.pool)
    .await?;
    let entry_id: i64 = sqlx::query("SELECT last_insert_rowid() AS id")
        .fetch_one(&state.pool)
        .await?
        .get("id");

    let lambda_name = if let Some(lid) = req.lambda_id {
        sqlx::query("SELECT name FROM lambdas WHERE id = ?")
            .bind(lid)
            .fetch_optional(&state.pool)
            .await?
            .map(|r| r.get::<String, _>("name"))
    } else {
        None
    };

    Ok((
        StatusCode::CREATED,
        Json(EntryView {
            id: entry_id,
            role: req.role,
            lambda_id: req.lambda_id,
            lambda_name,
        }),
    ))
}

pub async fn update_entry(
    State(state): State<AppState>,
    Path((rule_id, entry_id)): Path<(i64, i64)>,
    Json(req): Json<UpdateEntryReq>,
) -> ApiResult<Json<EntryView>> {
    if let Some(lid) = req.lambda_id {
        let role_id = lookup_entry_role_id(&state.pool, entry_id).await?;
        ensure_lambda_attachable(&state.pool, role_id, lid).await?;
    }
    let res = sqlx::query(
        "UPDATE access_rule_entries SET lambda_id = ? WHERE id = ? AND rule_id = ?",
    )
    .bind(req.lambda_id)
    .bind(entry_id)
    .bind(rule_id)
    .execute(&state.pool)
    .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("entry"));
    }
    let row = sqlx::query(
        "SELECT e.id, r.name AS role, e.lambda_id, l.name AS lambda_name
         FROM access_rule_entries e
         JOIN roles r ON r.id = e.role_id
         LEFT JOIN lambdas l ON l.id = e.lambda_id
         WHERE e.id = ?",
    )
    .bind(entry_id)
    .fetch_one(&state.pool)
    .await?;
    Ok(Json(EntryView {
        id: row.get("id"),
        role: row.get("role"),
        lambda_id: row.get("lambda_id"),
        lambda_name: row.get("lambda_name"),
    }))
}

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

// ----- Contract <-> rule bindings -----

pub async fn list_bindings(
    State(state): State<AppState>,
    Query(q): Query<ListBindingsQuery>,
) -> ApiResult<Json<Vec<BindingView>>> {
    let limit = q.limit.unwrap_or(500).clamp(1, 1000);
    let offset = q.offset.unwrap_or(0).max(0);
    let contract_norm = match q.contract.as_deref() {
        Some(c) => Some(normalize_address(c)?),
        None => None,
    };

    let rows = match (contract_norm.as_deref(), q.rule_id) {
        (Some(c), Some(r)) => {
            sqlx::query(
                "SELECT cr.id, cr.contract_address, cr.selector, cr.rule_id,
                        ar.name AS rule_name, ar.mode
                 FROM contract_rules cr
                 JOIN access_rules ar ON ar.id = cr.rule_id
                 WHERE cr.contract_address = ? AND cr.rule_id = ?
                 ORDER BY cr.id LIMIT ? OFFSET ?",
            )
            .bind(c)
            .bind(r)
            .bind(limit)
            .bind(offset)
            .fetch_all(&state.pool)
            .await?
        }
        (Some(c), None) => {
            sqlx::query(
                "SELECT cr.id, cr.contract_address, cr.selector, cr.rule_id,
                        ar.name AS rule_name, ar.mode
                 FROM contract_rules cr
                 JOIN access_rules ar ON ar.id = cr.rule_id
                 WHERE cr.contract_address = ?
                 ORDER BY cr.id LIMIT ? OFFSET ?",
            )
            .bind(c)
            .bind(limit)
            .bind(offset)
            .fetch_all(&state.pool)
            .await?
        }
        (None, Some(r)) => {
            sqlx::query(
                "SELECT cr.id, cr.contract_address, cr.selector, cr.rule_id,
                        ar.name AS rule_name, ar.mode
                 FROM contract_rules cr
                 JOIN access_rules ar ON ar.id = cr.rule_id
                 WHERE cr.rule_id = ?
                 ORDER BY cr.id LIMIT ? OFFSET ?",
            )
            .bind(r)
            .bind(limit)
            .bind(offset)
            .fetch_all(&state.pool)
            .await?
        }
        (None, None) => {
            sqlx::query(
                "SELECT cr.id, cr.contract_address, cr.selector, cr.rule_id,
                        ar.name AS rule_name, ar.mode
                 FROM contract_rules cr
                 JOIN access_rules ar ON ar.id = cr.rule_id
                 ORDER BY cr.id LIMIT ? OFFSET ?",
            )
            .bind(limit)
            .bind(offset)
            .fetch_all(&state.pool)
            .await?
        }
    };

    Ok(Json(
        rows.into_iter()
            .map(|r| BindingView {
                id: r.get("id"),
                contract_address: r.get("contract_address"),
                selector: r.get("selector"),
                rule_id: r.get("rule_id"),
                rule_name: r.get("rule_name"),
                mode: r.get("mode"),
            })
            .collect(),
    ))
}

pub async fn create_binding(
    State(state): State<AppState>,
    Json(req): Json<CreateBindingReq>,
) -> ApiResult<impl IntoResponse> {
    let contract = normalize_address(&req.contract_address)?;
    let rule = sqlx::query("SELECT id, selector, name, mode FROM access_rules WHERE id = ?")
        .bind(req.rule_id)
        .fetch_optional(&state.pool)
        .await?
        .ok_or_else(|| ApiError::bad_request(format!("unknown rule id {}", req.rule_id)))?;
    let selector: String = rule.get("selector");

    let res = sqlx::query(
        "INSERT INTO contract_rules (contract_address, selector, rule_id) VALUES (?, ?, ?)",
    )
    .bind(&contract)
    .bind(&selector)
    .bind(req.rule_id)
    .execute(&state.pool)
    .await?;
    let binding_id = res.last_insert_rowid();

    Ok((
        StatusCode::CREATED,
        Json(BindingView {
            id: binding_id,
            contract_address: contract,
            selector,
            rule_id: req.rule_id,
            rule_name: rule.get("name"),
            mode: rule.get("mode"),
        }),
    ))
}

pub async fn delete_binding(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> ApiResult<StatusCode> {
    let res = sqlx::query("DELETE FROM contract_rules WHERE id = ?")
        .bind(id)
        .execute(&state.pool)
        .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("binding"));
    }
    Ok(StatusCode::NO_CONTENT)
}
