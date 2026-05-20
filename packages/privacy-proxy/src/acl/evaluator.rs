use alloy::primitives::Address;
use anyhow::Result;

use crate::acl::lambdas::{eval as lambda_eval, loader as lambda_loader};
use crate::acl::registry;
use crate::auth::CallerCtx;
use crate::db::Pool;

#[derive(Clone, Debug, PartialEq, Eq)]
pub enum AccessDecision {
    Allow,
    Deny {
        contract: Address,
        selector: [u8; 4],
        reason: DenyReason,
    },
}

#[derive(Clone, Debug, PartialEq, Eq)]
pub enum DenyReason {
    NotInAllowList,
    InDenyList,
    LambdaRejected,
    UnknownLambda,
    LambdaRoleMismatch,
    AnonymousAgainstGatedCall,
    UnknownRuleMode,
    DefaultEoaSelfOnly,
}

pub async fn check_call(
    pool: &Pool,
    ctx: &CallerCtx,
    contract: &Address,
    call_data: &[u8],
    msg_sender: Address,
) -> Result<AccessDecision> {
    if call_data.len() < 4 {
        return Ok(AccessDecision::Allow);
    }
    let selector: [u8; 4] = call_data[0..4].try_into().expect("len checked above");
    let selector_hex = format!("0x{}", hex::encode(selector));
    let contract_hex = format!("0x{}", hex::encode(contract.as_slice()));

    let Some(rule) = registry::find_rule_for(pool, &contract_hex, &selector_hex).await? else {
        return Ok(AccessDecision::Allow);
    };

    let role_name = match ctx.role_name() {
        Some(r) => r,
        None => {
            return Ok(match rule.mode.as_str() {
                "allow" => AccessDecision::Deny {
                    contract: *contract,
                    selector,
                    reason: DenyReason::AnonymousAgainstGatedCall,
                },
                "deny" => AccessDecision::Allow,
                _ => AccessDecision::Deny {
                    contract: *contract,
                    selector,
                    reason: DenyReason::UnknownRuleMode,
                },
            });
        }
    };

    let entry = registry::entry_for_role(pool, rule.id, role_name).await?;

    let (entry_matches, lambda_outcome) = match entry {
        None => (false, true),
        Some(e) => {
            let outcome = match (&ctx.attributes, e.lambda_id) {
                (_, None) => true,
                (Some(attrs), Some(lambda_id)) => {
                    let Some(lambda) = lambda_loader::load_lambda_by_id(pool, lambda_id).await?
                    else {
                        return Ok(AccessDecision::Deny {
                            contract: *contract,
                            selector,
                            reason: DenyReason::UnknownLambda,
                        });
                    };
                    if lambda.role_id != e.role_id {
                        return Ok(AccessDecision::Deny {
                            contract: *contract,
                            selector,
                            reason: DenyReason::LambdaRoleMismatch,
                        });
                    }
                    let tx_origin = ctx.eoa.unwrap_or(Address::ZERO);
                    lambda_eval::evaluate(
                        &lambda,
                        attrs,
                        selector,
                        call_data,
                        tx_origin,
                        msg_sender,
                    )
                }
                (None, Some(_)) => unreachable!("anonymous handled earlier"),
            };
            (true, outcome)
        }
    };

    let decision = match rule.mode.as_str() {
        "allow" => {
            if !entry_matches {
                AccessDecision::Deny {
                    contract: *contract,
                    selector,
                    reason: DenyReason::NotInAllowList,
                }
            } else if !lambda_outcome {
                AccessDecision::Deny {
                    contract: *contract,
                    selector,
                    reason: DenyReason::LambdaRejected,
                }
            } else {
                AccessDecision::Allow
            }
        }
        "deny" => {
            if entry_matches && lambda_outcome {
                AccessDecision::Deny {
                    contract: *contract,
                    selector,
                    reason: DenyReason::InDenyList,
                }
            } else {
                AccessDecision::Allow
            }
        }
        _ => AccessDecision::Deny {
            contract: *contract,
            selector,
            reason: DenyReason::UnknownRuleMode,
        },
    };
    Ok(decision)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::auth::{AdminCallerInfo, CallerAttributes, UserCallerInfo};
    use crate::db;
    use crate::db::now_unix;
    use crate::roles::reconcile_roles;
    use alloy::primitives::Address;
    use sqlx::sqlite::SqlitePoolOptions;
    use sqlx::Row;

    async fn fresh_pool() -> Pool {
        let pool = SqlitePoolOptions::new()
            .max_connections(1)
            .connect("sqlite::memory:")
            .await
            .unwrap();
        sqlx::migrate!("./migrations").run(&pool).await.unwrap();
        reconcile_roles(&pool).await.unwrap();
        pool
    }

    fn user_ctx(kyc: bool) -> CallerCtx {
        CallerCtx {
            eoa: Some(Address::ZERO),
            attributes: Some(CallerAttributes::User(UserCallerInfo {
                eoa: Address::ZERO,
                kyc,
                blacklisted: false,
            })),
        }
    }

    fn _admin_ctx() -> CallerCtx {
        CallerCtx {
            eoa: Some(Address::ZERO),
            attributes: Some(CallerAttributes::Admin(AdminCallerInfo {
                eoa: Address::ZERO,
            })),
        }
    }

    fn anon_ctx() -> CallerCtx {
        CallerCtx::anonymous()
    }

    async fn insert_rule(pool: &Pool, mode: &str) -> i64 {
        let row = sqlx::query(
            "INSERT INTO access_rules (name, description, selector, mode)
             VALUES (?, NULL, '0x70a08231', ?)
             RETURNING id",
        )
        .bind(format!("rule-{mode}-{}", rand::random::<u32>()))
        .bind(mode)
        .fetch_one(pool)
        .await
        .unwrap();
        let rule_id: i64 = row.get("id");
        sqlx::query(
            "INSERT INTO contract_rules (contract_address, selector, rule_id)
             VALUES ('0x000000000000000000000000000000000000beef', '0x70a08231', ?)",
        )
        .bind(rule_id)
        .execute(pool)
        .await
        .unwrap();
        rule_id
    }

    async fn insert_kyc_lambda(pool: &Pool, role: &str) -> i64 {
        let role_id: i64 = sqlx::query("SELECT id FROM roles WHERE name = ?")
            .bind(role)
            .fetch_one(pool)
            .await
            .unwrap()
            .get("id");
        sqlx::query(
            "INSERT INTO lambdas (name, role_id, description, created_at)
             VALUES ('require_kyc', ?, NULL, ?)",
        )
        .bind(role_id)
        .bind(now_unix())
        .execute(pool)
        .await
        .unwrap();
        let lambda_id: i64 = sqlx::query("SELECT last_insert_rowid() AS id")
            .fetch_one(pool)
            .await
            .unwrap()
            .get("id");
        let one_word = format!("0x{}01", "00".repeat(31));
        sqlx::query(
            "INSERT INTO lambda_rules (lambda_id, selector, lhs_kind, lhs_offset, lhs_attribute,
                                       condition, rhs_kind, rhs_value)
             VALUES (?, '0x70a08231', 'attribute', NULL, 'kyc', 'eq', 'literal', ?)",
        )
        .bind(lambda_id)
        .bind(one_word)
        .execute(pool)
        .await
        .unwrap();
        lambda_id
    }

    async fn insert_entry(pool: &Pool, rule_id: i64, role: &str, lambda_id: Option<i64>) {
        let role_id: i64 = sqlx::query("SELECT id FROM roles WHERE name = ?")
            .bind(role)
            .fetch_one(pool)
            .await
            .unwrap()
            .get("id");
        sqlx::query(
            "INSERT INTO access_rule_entries (rule_id, role_id, lambda_id) VALUES (?, ?, ?)",
        )
        .bind(rule_id)
        .bind(role_id)
        .bind(lambda_id)
        .execute(pool)
        .await
        .unwrap();
    }

    fn contract() -> Address {
        "0x000000000000000000000000000000000000beef".parse().unwrap()
    }

    fn balance_of_call_data() -> Vec<u8> {
        let mut out = vec![0x70, 0xa0, 0x82, 0x31];
        out.extend_from_slice(&[0u8; 32]);
        out
    }

    #[tokio::test]
    async fn no_rule_means_free_access() {
        let pool = fresh_pool().await;
        let dec = check_call(&pool, &user_ctx(true), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert_eq!(dec, AccessDecision::Allow);
    }

    #[tokio::test]
    async fn allow_mode_no_entry_for_role_denies() {
        let pool = fresh_pool().await;
        insert_rule(&pool, "allow").await;
        let dec = check_call(&pool, &user_ctx(true), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert!(matches!(
            dec,
            AccessDecision::Deny {
                reason: DenyReason::NotInAllowList,
                ..
            }
        ));
    }

    #[tokio::test]
    async fn allow_mode_entry_no_lambda_allows() {
        let pool = fresh_pool().await;
        let r = insert_rule(&pool, "allow").await;
        insert_entry(&pool, r, "user", None).await;
        let dec = check_call(&pool, &user_ctx(true), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert_eq!(dec, AccessDecision::Allow);
    }

    #[tokio::test]
    async fn allow_mode_lambda_pass() {
        let pool = fresh_pool().await;
        let r = insert_rule(&pool, "allow").await;
        let lid = insert_kyc_lambda(&pool, "user").await;
        insert_entry(&pool, r, "user", Some(lid)).await;
        let dec = check_call(&pool, &user_ctx(true), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert_eq!(dec, AccessDecision::Allow);
    }

    #[tokio::test]
    async fn allow_mode_lambda_reject() {
        let pool = fresh_pool().await;
        let r = insert_rule(&pool, "allow").await;
        let lid = insert_kyc_lambda(&pool, "user").await;
        insert_entry(&pool, r, "user", Some(lid)).await;
        let dec = check_call(&pool, &user_ctx(false), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert!(matches!(
            dec,
            AccessDecision::Deny {
                reason: DenyReason::LambdaRejected,
                ..
            }
        ));
    }

    #[tokio::test]
    async fn deny_mode_no_entry_allows() {
        let pool = fresh_pool().await;
        insert_rule(&pool, "deny").await;
        let dec = check_call(&pool, &user_ctx(true), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert_eq!(dec, AccessDecision::Allow);
    }

    #[tokio::test]
    async fn deny_mode_entry_no_lambda_denies() {
        let pool = fresh_pool().await;
        let r = insert_rule(&pool, "deny").await;
        insert_entry(&pool, r, "user", None).await;
        let dec = check_call(&pool, &user_ctx(true), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert!(matches!(
            dec,
            AccessDecision::Deny {
                reason: DenyReason::InDenyList,
                ..
            }
        ));
    }

    #[tokio::test]
    async fn deny_mode_lambda_pass_denies() {
        let pool = fresh_pool().await;
        let r = insert_rule(&pool, "deny").await;
        let lid = insert_kyc_lambda(&pool, "user").await;
        insert_entry(&pool, r, "user", Some(lid)).await;
        let dec = check_call(&pool, &user_ctx(true), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert!(matches!(
            dec,
            AccessDecision::Deny {
                reason: DenyReason::InDenyList,
                ..
            }
        ));
    }

    #[tokio::test]
    async fn deny_mode_lambda_fail_allows() {
        let pool = fresh_pool().await;
        let r = insert_rule(&pool, "deny").await;
        let lid = insert_kyc_lambda(&pool, "user").await;
        insert_entry(&pool, r, "user", Some(lid)).await;
        let dec = check_call(&pool, &user_ctx(false), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert_eq!(dec, AccessDecision::Allow);
    }

    #[tokio::test]
    async fn anonymous_caller_denied_on_allow_rule() {
        let pool = fresh_pool().await;
        insert_rule(&pool, "allow").await;
        let dec = check_call(&pool, &anon_ctx(), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert!(matches!(
            dec,
            AccessDecision::Deny {
                reason: DenyReason::AnonymousAgainstGatedCall,
                ..
            }
        ));
    }

    #[tokio::test]
    async fn anonymous_caller_allowed_on_deny_rule() {
        let pool = fresh_pool().await;
        insert_rule(&pool, "deny").await;
        let dec = check_call(&pool, &anon_ctx(), &contract(), &balance_of_call_data(), Address::ZERO)
            .await
            .unwrap();
        assert_eq!(dec, AccessDecision::Allow);
    }

    #[allow(dead_code)]
    fn _typecheck(_: db::Pool) {}
}
