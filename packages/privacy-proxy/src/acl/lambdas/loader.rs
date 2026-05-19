use anyhow::{anyhow, Result};
use sqlx::Row;

use super::{Condition, Lambda, LambdaRule, Lhs, Rhs};
use crate::db::Pool;

pub async fn load_lambda_by_id(pool: &Pool, id: i64) -> Result<Option<Lambda>> {
    let head = sqlx::query(
        "SELECT l.id, l.name, l.role_id, r.name AS role_name, l.description
         FROM lambdas l
         JOIN roles r ON r.id = l.role_id
         WHERE l.id = ?",
    )
    .bind(id)
    .fetch_optional(pool)
    .await?;
    let Some(row) = head else {
        return Ok(None);
    };
    let id: i64 = row.get("id");
    let mut lambda = Lambda {
        id,
        name: row.get("name"),
        role_id: row.get("role_id"),
        role: row.get("role_name"),
        description: row.get("description"),
        rules: Vec::new(),
    };
    lambda.rules = load_rules_for(pool, id).await?;
    Ok(Some(lambda))
}

pub async fn load_lambda_by_name_role(
    pool: &Pool,
    name: &str,
    role_id: i64,
) -> Result<Option<Lambda>> {
    let head = sqlx::query(
        "SELECT l.id, l.name, l.role_id, r.name AS role_name, l.description
         FROM lambdas l
         JOIN roles r ON r.id = l.role_id
         WHERE l.name = ? AND l.role_id = ?",
    )
    .bind(name)
    .bind(role_id)
    .fetch_optional(pool)
    .await?;
    let Some(row) = head else {
        return Ok(None);
    };
    let id: i64 = row.get("id");
    let mut lambda = Lambda {
        id,
        name: row.get("name"),
        role_id: row.get("role_id"),
        role: row.get("role_name"),
        description: row.get("description"),
        rules: Vec::new(),
    };
    lambda.rules = load_rules_for(pool, id).await?;
    Ok(Some(lambda))
}

pub async fn list_lambdas(pool: &Pool) -> Result<Vec<Lambda>> {
    let rows = sqlx::query(
        "SELECT l.id, l.name, l.role_id, r.name AS role_name, l.description
         FROM lambdas l
         JOIN roles r ON r.id = l.role_id
         ORDER BY r.name, l.name",
    )
    .fetch_all(pool)
    .await?;
    let mut out = Vec::with_capacity(rows.len());
    for row in rows {
        let id: i64 = row.get("id");
        let rules = load_rules_for(pool, id).await?;
        out.push(Lambda {
            id,
            name: row.get("name"),
            role_id: row.get("role_id"),
            role: row.get("role_name"),
            description: row.get("description"),
            rules,
        });
    }
    Ok(out)
}

pub async fn count_references(pool: &Pool, lambda_id: i64) -> Result<i64> {
    let row = sqlx::query("SELECT COUNT(*) AS c FROM access_rule_entries WHERE lambda_id = ?")
        .bind(lambda_id)
        .fetch_one(pool)
        .await?;
    Ok(row.get("c"))
}

async fn load_rules_for(pool: &Pool, lambda_id: i64) -> Result<Vec<LambdaRule>> {
    let rows = sqlx::query(
        "SELECT id, selector, lhs_kind, lhs_offset, lhs_attribute,
                condition, rhs_kind, rhs_value
         FROM lambda_rules
         WHERE lambda_id = ?
         ORDER BY id",
    )
    .bind(lambda_id)
    .fetch_all(pool)
    .await?;
    let mut out = Vec::with_capacity(rows.len());
    for r in rows {
        let id: i64 = r.get("id");
        let selector_hex: String = r.get("selector");
        let selector = parse_selector(&selector_hex)?;
        let lhs_kind: String = r.get("lhs_kind");
        let lhs = match lhs_kind.as_str() {
            "calldata" => {
                let off: i64 = r.get("lhs_offset");
                Lhs::Calldata { offset: off as u32 }
            }
            "attribute" => Lhs::Attribute {
                name: r.get("lhs_attribute"),
            },
            other => return Err(anyhow!("unknown lhs_kind in DB: {other}")),
        };
        let condition_db: String = r.get("condition");
        let condition = Condition::from_db(&condition_db)
            .ok_or_else(|| anyhow!("unknown condition in DB: {condition_db}"))?;
        let rhs_kind: String = r.get("rhs_kind");
        let rhs = match rhs_kind.as_str() {
            "tx_origin" => Rhs::TxOrigin,
            "msg_sender" => Rhs::MsgSender,
            "literal" => Rhs::Literal {
                value_hex: r.get("rhs_value"),
            },
            other => return Err(anyhow!("unknown rhs_kind in DB: {other}")),
        };
        out.push(LambdaRule {
            id,
            selector,
            lhs,
            condition,
            rhs,
        });
    }
    Ok(out)
}

fn parse_selector(hex_str: &str) -> Result<[u8; 4]> {
    let trimmed = hex_str.trim_start_matches("0x");
    let bytes = hex::decode(trimmed)?;
    if bytes.len() != 4 {
        return Err(anyhow!("selector must be 4 bytes, got {}", bytes.len()));
    }
    let mut out = [0u8; 4];
    out.copy_from_slice(&bytes);
    Ok(out)
}
