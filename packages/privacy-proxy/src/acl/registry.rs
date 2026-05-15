use anyhow::Result;
use sqlx::Row;

use crate::db::Pool;

#[derive(Clone, Debug)]
pub struct AccessRule {
    pub id: i64,
    pub contract_address: String,
    pub function_selector: String,
    pub mode: String,
}

#[derive(Clone, Debug)]
pub struct AccessRuleEntry {
    pub id: i64,
    pub rule_id: i64,
    pub role_id: i64,
    pub role_name: String,
    pub lambda_name: Option<String>,
}

pub async fn find_rule(
    pool: &Pool,
    contract_hex: &str,
    selector_hex: &str,
) -> Result<Option<AccessRule>> {
    let row = sqlx::query(
        "SELECT id, contract_address, function_selector, mode
         FROM access_rules WHERE contract_address = ? AND function_selector = ?",
    )
    .bind(contract_hex)
    .bind(selector_hex)
    .fetch_optional(pool)
    .await?;
    Ok(row.map(|r| AccessRule {
        id: r.get("id"),
        contract_address: r.get("contract_address"),
        function_selector: r.get("function_selector"),
        mode: r.get("mode"),
    }))
}

pub async fn entry_for_role(
    pool: &Pool,
    rule_id: i64,
    role_name: &str,
) -> Result<Option<AccessRuleEntry>> {
    let row = sqlx::query(
        "SELECT e.id, e.rule_id, e.role_id, r.name AS role_name, e.lambda_name
         FROM access_rule_entries e
         JOIN roles r ON r.id = e.role_id
         WHERE e.rule_id = ? AND r.name = ?",
    )
    .bind(rule_id)
    .bind(role_name)
    .fetch_optional(pool)
    .await?;
    Ok(row.map(|r| AccessRuleEntry {
        id: r.get("id"),
        rule_id: r.get("rule_id"),
        role_id: r.get("role_id"),
        role_name: r.get("role_name"),
        lambda_name: r.get("lambda_name"),
    }))
}

pub async fn list_entries(pool: &Pool, rule_id: i64) -> Result<Vec<AccessRuleEntry>> {
    let rows = sqlx::query(
        "SELECT e.id, e.rule_id, e.role_id, r.name AS role_name, e.lambda_name
         FROM access_rule_entries e
         JOIN roles r ON r.id = e.role_id
         WHERE e.rule_id = ?
         ORDER BY r.name",
    )
    .bind(rule_id)
    .fetch_all(pool)
    .await?;
    Ok(rows
        .into_iter()
        .map(|r| AccessRuleEntry {
            id: r.get("id"),
            rule_id: r.get("rule_id"),
            role_id: r.get("role_id"),
            role_name: r.get("role_name"),
            lambda_name: r.get("lambda_name"),
        })
        .collect())
}
