use axum::extract::{Path, State};
use axum::http::StatusCode;
use axum::response::IntoResponse;
use axum::Json;
use serde::{Deserialize, Serialize};
use sqlx::Row;

use crate::acl::lambdas::attributes::{known_attribute_specs_for_role, AttributeSpec};
use crate::acl::lambdas::loader as lambda_loader;
use crate::acl::lambdas::{Condition, Lhs, Rhs};
use crate::admin::normalize_selector;
use crate::db::now_unix;
use crate::error::{ApiError, ApiResult};
use crate::roles::ROLES;
use crate::rpc::gated_methods;
use crate::state::AppState;

const MAX_RULES_PER_LAMBDA: usize = 64;

#[derive(Serialize)]
pub struct LambdaRuleView {
    pub id: i64,
    pub selector: String,
    pub lhs_kind: &'static str,
    pub lhs_offset: Option<u32>,
    pub lhs_attribute: Option<String>,
    pub condition: &'static str,
    pub rhs_kind: &'static str,
    pub rhs_value: Option<String>,
}

#[derive(Serialize)]
pub struct LambdaView {
    pub id: i64,
    pub name: String,
    pub role: String,
    pub description: Option<String>,
    pub rules: Vec<LambdaRuleView>,
    pub in_use: bool,
}

#[derive(Serialize)]
pub struct LambdaGroup {
    pub role: &'static str,
    pub lambdas: Vec<LambdaView>,
}

#[derive(Serialize)]
pub struct SyntheticSelector {
    pub method: &'static str,
    pub selector: String,
}

#[derive(Serialize)]
pub struct RoleAttributeView {
    pub name: &'static str,
    #[serde(rename = "type")]
    pub ty: &'static str,
}

#[derive(Serialize)]
pub struct RoleAttributes {
    pub role: &'static str,
    pub attributes: Vec<RoleAttributeView>,
}

#[derive(Deserialize)]
pub struct CreateLambdaRuleInput {
    pub selector: String,
    pub lhs_kind: String,
    #[serde(default)]
    pub lhs_offset: Option<u32>,
    #[serde(default)]
    pub lhs_attribute: Option<String>,
    pub condition: String,
    pub rhs_kind: String,
    #[serde(default)]
    pub rhs_value: Option<String>,
}

#[derive(Deserialize)]
pub struct CreateLambdaReq {
    pub name: String,
    pub role: String,
    #[serde(default)]
    pub description: Option<String>,
    pub rules: Vec<CreateLambdaRuleInput>,
}

pub async fn list_lambdas(State(state): State<AppState>) -> ApiResult<Json<Vec<LambdaGroup>>> {
    let lambdas = lambda_loader::list_lambdas(&state.pool).await?;
    let mut groups: Vec<LambdaGroup> = ROLES
        .iter()
        .map(|r| LambdaGroup {
            role: *r,
            lambdas: Vec::new(),
        })
        .collect();
    for l in lambdas {
        let in_use = lambda_loader::count_references(&state.pool, l.id).await? > 0;
        let view = to_view(&l, in_use);
        if let Some(g) = groups.iter_mut().find(|g| g.role == l.role) {
            g.lambdas.push(view);
        }
    }
    Ok(Json(groups))
}

pub async fn get_lambda(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> ApiResult<Json<LambdaView>> {
    let lambda = lambda_loader::load_lambda_by_id(&state.pool, id)
        .await?
        .ok_or_else(|| ApiError::not_found("lambda"))?;
    let in_use = lambda_loader::count_references(&state.pool, lambda.id).await? > 0;
    Ok(Json(to_view(&lambda, in_use)))
}

pub async fn create_lambda(
    State(state): State<AppState>,
    Json(req): Json<CreateLambdaReq>,
) -> ApiResult<impl IntoResponse> {
    let name = req.name.trim();
    if name.is_empty() {
        return Err(ApiError::bad_request("name must not be empty"));
    }
    if !ROLES.contains(&req.role.as_str()) {
        return Err(ApiError::bad_request(format!("unknown role `{}`", req.role)));
    }
    if req.rules.is_empty() {
        return Err(ApiError::bad_request("at least one rule is required"));
    }
    if req.rules.len() > MAX_RULES_PER_LAMBDA {
        return Err(ApiError::bad_request(format!(
            "too many rules (limit {MAX_RULES_PER_LAMBDA})"
        )));
    }

    let attributes = known_attribute_specs_for_role(&req.role);
    let mut parsed: Vec<NormalizedRule> = Vec::with_capacity(req.rules.len());
    for rule in &req.rules {
        parsed.push(normalize_rule(rule, attributes)?);
    }

    let role_id: i64 = sqlx::query("SELECT id FROM roles WHERE name = ?")
        .bind(&req.role)
        .fetch_one(&state.pool)
        .await?
        .get("id");

    let exists = sqlx::query("SELECT 1 FROM lambdas WHERE name = ? AND role_id = ?")
        .bind(name)
        .bind(role_id)
        .fetch_optional(&state.pool)
        .await?;
    if exists.is_some() {
        return Err(ApiError::conflict(format!(
            "lambda `{name}` already exists for role `{}`",
            req.role
        )));
    }

    let mut tx = state.pool.begin().await?;
    sqlx::query(
        "INSERT INTO lambdas (name, role_id, description, created_at) VALUES (?, ?, ?, ?)",
    )
    .bind(name)
    .bind(role_id)
    .bind(req.description.as_deref())
    .bind(now_unix())
    .execute(&mut *tx)
    .await?;
    let lambda_id: i64 = sqlx::query("SELECT last_insert_rowid() AS id")
        .fetch_one(&mut *tx)
        .await?
        .get("id");

    for r in &parsed {
        sqlx::query(
            "INSERT INTO lambda_rules
                (lambda_id, selector, lhs_kind, lhs_offset, lhs_attribute,
                 condition, rhs_kind, rhs_value)
             VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
        )
        .bind(lambda_id)
        .bind(&r.selector)
        .bind(r.lhs_kind)
        .bind(r.lhs_offset)
        .bind(r.lhs_attribute.as_deref())
        .bind(r.condition)
        .bind(r.rhs_kind)
        .bind(r.rhs_value.as_deref())
        .execute(&mut *tx)
        .await?;
    }
    tx.commit().await?;

    let lambda = lambda_loader::load_lambda_by_id(&state.pool, lambda_id)
        .await?
        .ok_or_else(|| ApiError::from(anyhow::anyhow!("lambda vanished after insert")))?;
    Ok((StatusCode::CREATED, Json(to_view(&lambda, false))))
}

pub async fn delete_lambda(
    State(state): State<AppState>,
    Path(id): Path<i64>,
) -> ApiResult<StatusCode> {
    let refs = lambda_loader::count_references(&state.pool, id).await?;
    if refs > 0 {
        return Err(ApiError::conflict(format!(
            "lambda is referenced by {refs} rule entry(ies); detach before deleting"
        )));
    }
    let res = sqlx::query("DELETE FROM lambdas WHERE id = ?")
        .bind(id)
        .execute(&state.pool)
        .await?;
    if res.rows_affected() == 0 {
        return Err(ApiError::not_found("lambda"));
    }
    Ok(StatusCode::NO_CONTENT)
}

pub async fn list_role_attributes() -> ApiResult<Json<Vec<RoleAttributes>>> {
    let out = ROLES
        .iter()
        .map(|r| RoleAttributes {
            role: *r,
            attributes: known_attribute_specs_for_role(r)
                .iter()
                .map(|s| RoleAttributeView { name: s.name, ty: s.ty })
                .collect(),
        })
        .collect();
    Ok(Json(out))
}

pub async fn list_synthetic_selectors() -> ApiResult<Json<Vec<SyntheticSelector>>> {
    let out = gated_methods::ALL
        .iter()
        .map(|m| SyntheticSelector {
            method: m.name,
            selector: format!("0x{}", hex::encode(m.selector)),
        })
        .collect();
    Ok(Json(out))
}

struct NormalizedRule {
    selector: String,
    lhs_kind: &'static str,
    lhs_offset: Option<u32>,
    lhs_attribute: Option<String>,
    condition: &'static str,
    rhs_kind: &'static str,
    rhs_value: Option<String>,
}

fn normalize_rule(
    input: &CreateLambdaRuleInput,
    attributes: &'static [AttributeSpec],
) -> Result<NormalizedRule, ApiError> {
    let selector = normalize_selector(&input.selector)?;
    let (lhs_kind, lhs_offset, lhs_attribute) = match input.lhs_kind.as_str() {
        "calldata" => {
            let off = input
                .lhs_offset
                .ok_or_else(|| ApiError::bad_request("lhs_offset required for calldata lhs"))?;
            if input.lhs_attribute.is_some() {
                return Err(ApiError::bad_request(
                    "lhs_attribute must be null when lhs_kind=calldata",
                ));
            }
            ("calldata", Some(off), None)
        }
        "attribute" => {
            let name = input.lhs_attribute.as_deref().ok_or_else(|| {
                ApiError::bad_request("lhs_attribute required for attribute lhs")
            })?;
            if !attributes.iter().any(|a| a.name == name) {
                return Err(ApiError::bad_request(format!(
                    "unknown attribute `{name}` for this role"
                )));
            }
            if input.lhs_offset.is_some() {
                return Err(ApiError::bad_request(
                    "lhs_offset must be null when lhs_kind=attribute",
                ));
            }
            ("attribute", None, Some(name.to_string()))
        }
        other => {
            return Err(ApiError::bad_request(format!(
                "lhs_kind must be calldata|attribute, got `{other}`"
            )))
        }
    };
    let condition = Condition::from_db(&input.condition)
        .ok_or_else(|| ApiError::bad_request(format!("invalid condition `{}`", input.condition)))?;

    let (rhs_kind, rhs_value) = match input.rhs_kind.as_str() {
        "tx_origin" => {
            if input.rhs_value.is_some() {
                return Err(ApiError::bad_request("rhs_value must be null for tx_origin"));
            }
            ("tx_origin", None)
        }
        "msg_sender" => {
            if input.rhs_value.is_some() {
                return Err(ApiError::bad_request("rhs_value must be null for msg_sender"));
            }
            ("msg_sender", None)
        }
        "literal" => {
            let v = input
                .rhs_value
                .as_deref()
                .ok_or_else(|| ApiError::bad_request("rhs_value required for literal"))?;
            let trimmed = v.trim_start_matches("0x");
            if trimmed.len() != 64 {
                return Err(ApiError::bad_request(
                    "rhs_value must be a 32-byte hex (0x + 64 chars)",
                ));
            }
            hex::decode(trimmed)
                .map_err(|_| ApiError::bad_request("rhs_value is not hex"))?;
            ("literal", Some(format!("0x{}", trimmed.to_ascii_lowercase())))
        }
        other => {
            return Err(ApiError::bad_request(format!(
                "rhs_kind must be tx_origin|msg_sender|literal, got `{other}`"
            )))
        }
    };

    Ok(NormalizedRule {
        selector,
        lhs_kind,
        lhs_offset,
        lhs_attribute,
        condition: condition.as_db(),
        rhs_kind,
        rhs_value,
    })
}

fn to_view(l: &crate::acl::lambdas::Lambda, in_use: bool) -> LambdaView {
    LambdaView {
        id: l.id,
        name: l.name.clone(),
        role: l.role.clone(),
        description: l.description.clone(),
        rules: l.rules.iter().map(rule_to_view).collect(),
        in_use,
    }
}

fn rule_to_view(r: &crate::acl::lambdas::LambdaRule) -> LambdaRuleView {
    let (lhs_kind, lhs_offset, lhs_attribute) = match &r.lhs {
        Lhs::Calldata { offset } => ("calldata", Some(*offset), None),
        Lhs::Attribute { name } => ("attribute", None, Some(name.clone())),
    };
    let (rhs_kind, rhs_value) = match &r.rhs {
        Rhs::TxOrigin => ("tx_origin", None),
        Rhs::MsgSender => ("msg_sender", None),
        Rhs::Literal { value_hex } => ("literal", Some(value_hex.clone())),
    };
    LambdaRuleView {
        id: r.id,
        selector: format!("0x{}", hex::encode(r.selector)),
        lhs_kind,
        lhs_offset,
        lhs_attribute,
        condition: r.condition.as_db(),
        rhs_kind,
        rhs_value,
    }
}
