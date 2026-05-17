use std::collections::HashMap;
use std::sync::LazyLock;

use alloy::primitives::Address;
use serde::Serialize;

use super::LambdaSpec;

pub mod erc20_self_only;
pub mod require_kyc;

/// Typed attribute set for the `user` role.
///
/// Persisted as a row in `user_attributes`. The `eoa` field comes from
/// `members.eoa_address` (the table's primary key) and is set by the auth
/// middleware; the other fields are admin-managed (`PUT /admin/members/:eoa`).
#[derive(Clone, Debug, Serialize)]
pub struct UserCallerInfo {
    pub eoa: Address,
    pub kyc: bool,
    pub blacklisted: bool,
}

/// Lambdas available for `user`-role entries. Adding a lambda = adding
/// a sibling module under `user/` and an entry here.
pub fn registry() -> &'static HashMap<&'static str, &'static LambdaSpec<UserCallerInfo>> {
    static REGISTRY: LazyLock<HashMap<&'static str, &'static LambdaSpec<UserCallerInfo>>> =
        LazyLock::new(|| {
            let specs: Box<[LambdaSpec<UserCallerInfo>]> = Box::new([
                LambdaSpec {
                    name: "require_kyc",
                    description: "Allow only callers whose stored attributes have kyc=true.",
                    expected_selector: None,
                    run: require_kyc::run,
                },
                LambdaSpec {
                    name: "erc20_self_only",
                    description: "For ERC-20 balanceOf(address) and allowance(address,address): allow only when the queried account (balanceOf) or owner (allowance) equals the caller's EOA.",
                    expected_selector: None,
                    run: erc20_self_only::run,
                },
            ]);
            let specs: &'static [LambdaSpec<UserCallerInfo>] = Box::leak(specs);
            let mut m = HashMap::new();
            for spec in specs {
                m.insert(spec.name, spec);
            }
            m
        });
    &REGISTRY
}

pub fn lookup(name: &str) -> Option<&'static LambdaSpec<UserCallerInfo>> {
    registry().get(name).copied()
}

pub fn list_specs() -> Vec<&'static LambdaSpec<UserCallerInfo>> {
    let mut v: Vec<&'static LambdaSpec<UserCallerInfo>> = registry().values().copied().collect();
    v.sort_by_key(|s| s.name);
    v
}
