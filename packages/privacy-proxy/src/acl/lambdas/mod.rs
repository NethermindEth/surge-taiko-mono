pub mod examples;

use std::collections::HashMap;
use std::sync::LazyLock;

use serde_json::Value;

pub struct LambdaCtx<'a> {
    pub caller_info: &'a Value,
    pub selector: [u8; 4],
    pub call_data: &'a [u8],
}

pub type LambdaFn = fn(&LambdaCtx) -> bool;

pub struct LambdaSpec {
    pub name: &'static str,
    pub description: &'static str,
    pub expected_selector: Option<[u8; 4]>,
    pub run: LambdaFn,
}

/// Every lambda available in this build. Adding a lambda = adding an entry
/// here and a function in `examples.rs` (or a sibling module). Descriptions
/// are the source of truth shown to admins via `GET /admin/registry/lambdas`.
pub fn registry() -> &'static HashMap<&'static str, &'static LambdaSpec> {
    static REGISTRY: LazyLock<HashMap<&'static str, &'static LambdaSpec>> = LazyLock::new(|| {
        let specs: Box<[LambdaSpec]> = Box::new([
            LambdaSpec {
                name: "require_kyc",
                description: "Allow only callers whose stored caller_info has `kyc: true`.",
                expected_selector: None,
                run: examples::require_kyc,
            },
            LambdaSpec {
                name: "transfer_under_limit",
                description: "For ERC-20 transfer(address,uint256): require amount <= caller_info.max_transfer (decimal string in wei).",
                expected_selector: Some(examples::TRANSFER_SELECTOR),
                run: examples::transfer_under_limit,
            },
            LambdaSpec {
                name: "target_in_caller_allowlist",
                description: "For gated address-parameterized reads (eth_getBalance et al.): allow if params[0] is present in caller_info.balance_allowlist (array of lowercase hex addresses).",
                expected_selector: None,
                run: examples::target_in_caller_allowlist,
            },
            LambdaSpec {
                name: "erc20_self_only",
                description: "For ERC-20 balanceOf(address) and allowance(address,address): allow only when the queried account (balanceOf) or owner (allowance) equals the caller's EOA injected as caller_info.eoa.",
                expected_selector: None,
                run: examples::erc20_self_only,
            },
        ]);
        let specs: &'static [LambdaSpec] = Box::leak(specs);
        let mut m = HashMap::new();
        for spec in specs {
            m.insert(spec.name, spec);
        }
        m
    });
    &REGISTRY
}

pub fn list_specs() -> Vec<&'static LambdaSpec> {
    let mut v: Vec<&'static LambdaSpec> = registry().values().copied().collect();
    v.sort_by_key(|s| s.name);
    v
}

pub fn lookup(name: &str) -> Option<&'static LambdaSpec> {
    registry().get(name).copied()
}
