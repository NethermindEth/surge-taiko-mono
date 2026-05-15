use alloy::primitives::{Address, U256};
use alloy::sol;
use alloy::sol_types::SolCall;

use super::LambdaCtx;

/// keccak256("transfer(address,uint256)")[..4]
pub const TRANSFER_SELECTOR: [u8; 4] = [0xa9, 0x05, 0x9c, 0xbb];

/// keccak256("balanceOf(address)")[..4]
pub const BALANCE_OF_SELECTOR: [u8; 4] = [0x70, 0xa0, 0x82, 0x31];
/// keccak256("allowance(address,address)")[..4]
pub const ALLOWANCE_SELECTOR: [u8; 4] = [0xdd, 0x62, 0xed, 0x3e];

sol! {
    function transfer(address to, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
    function allowance(address owner, address spender) external view returns (uint256);
}

pub fn require_kyc(ctx: &LambdaCtx) -> bool {
    ctx.caller_info
        .get("kyc")
        .and_then(|v| v.as_bool())
        .unwrap_or(false)
}

/// For synthetic gated-method rules: allow when the target address
/// (encoded in `call_data[4..36]` as a left-padded 32-byte word) is
/// present in `caller_info.balance_allowlist` (a JSON array of
/// lowercase hex address strings).
pub fn target_in_caller_allowlist(ctx: &LambdaCtx) -> bool {
    if ctx.call_data.len() < 36 {
        return false;
    }
    let target_bytes: &[u8] = &ctx.call_data[16..36];
    let target_hex = format!("0x{}", hex::encode(target_bytes));
    let Some(list) = ctx
        .caller_info
        .get("balance_allowlist")
        .and_then(|v| v.as_array())
    else {
        return false;
    };
    list.iter()
        .filter_map(|v| v.as_str())
        .any(|s| s.eq_ignore_ascii_case(&target_hex))
}

/// For ERC-20 `balanceOf(address)` and `allowance(address,address)`:
/// allow only when the queried `account` (balanceOf) or `owner`
/// (allowance) equals the caller's EOA, which the auth middleware
/// injects into `caller_info.eoa`. Spender on allowance is not
/// constrained. Any other selector → false.
pub fn erc20_self_only(ctx: &LambdaCtx) -> bool {
    let Some(eoa) = ctx
        .caller_info
        .get("eoa")
        .and_then(|v| v.as_str())
        .and_then(|s| s.parse::<Address>().ok())
    else {
        return false;
    };
    match ctx.selector {
        s if s == BALANCE_OF_SELECTOR => balanceOfCall::abi_decode(ctx.call_data)
            .map(|c| c.account == eoa)
            .unwrap_or(false),
        s if s == ALLOWANCE_SELECTOR => allowanceCall::abi_decode(ctx.call_data)
            .map(|c| c.owner == eoa)
            .unwrap_or(false),
        _ => false,
    }
}

pub fn transfer_under_limit(ctx: &LambdaCtx) -> bool {
    if ctx.selector != TRANSFER_SELECTOR {
        return false;
    }
    let decoded = match transferCall::abi_decode(ctx.call_data) {
        Ok(d) => d,
        Err(_) => return false,
    };
    let Some(limit_str) = ctx.caller_info.get("max_transfer").and_then(|v| v.as_str()) else {
        return false;
    };
    let Ok(limit) = U256::from_str_radix(limit_str, 10) else {
        return false;
    };
    decoded.amount <= limit
}

#[cfg(test)]
mod tests {
    use super::*;
    use alloy::primitives::{Address, U256};
    use alloy::sol_types::SolCall;
    use serde_json::json;

    fn ctx<'a>(info: &'a serde_json::Value, sel: [u8; 4], data: &'a [u8]) -> LambdaCtx<'a> {
        LambdaCtx {
            caller_info: info,
            selector: sel,
            call_data: data,
        }
    }

    #[test]
    fn require_kyc_pass_and_fail() {
        let yes = json!({ "kyc": true });
        let no = json!({ "kyc": false });
        let empty = json!({});
        assert!(require_kyc(&ctx(&yes, [0; 4], &[])));
        assert!(!require_kyc(&ctx(&no, [0; 4], &[])));
        assert!(!require_kyc(&ctx(&empty, [0; 4], &[])));
    }

    const SELF_EOA: &str = "0x1111111111111111111111111111111111111111";
    const OTHER_EOA: &str = "0x2222222222222222222222222222222222222222";

    fn self_info() -> serde_json::Value {
        json!({ "eoa": SELF_EOA })
    }

    fn encode_balance_of(account: &str) -> ([u8; 4], Vec<u8>) {
        let call = balanceOfCall {
            account: account.parse().unwrap(),
        };
        let data = call.abi_encode();
        let sel: [u8; 4] = data[0..4].try_into().unwrap();
        (sel, data)
    }

    fn encode_allowance(owner: &str, spender: &str) -> ([u8; 4], Vec<u8>) {
        let call = allowanceCall {
            owner: owner.parse().unwrap(),
            spender: spender.parse().unwrap(),
        };
        let data = call.abi_encode();
        let sel: [u8; 4] = data[0..4].try_into().unwrap();
        (sel, data)
    }

    #[test]
    fn erc20_self_only_selectors_match_constants() {
        let (sel, _) = encode_balance_of(SELF_EOA);
        assert_eq!(sel, BALANCE_OF_SELECTOR);
        let (sel, _) = encode_allowance(SELF_EOA, OTHER_EOA);
        assert_eq!(sel, ALLOWANCE_SELECTOR);
    }

    #[test]
    fn erc20_self_only_balance_of_self_allowed() {
        let info = self_info();
        let (sel, data) = encode_balance_of(SELF_EOA);
        assert!(erc20_self_only(&ctx(&info, sel, &data)));
    }

    #[test]
    fn erc20_self_only_balance_of_other_denied() {
        let info = self_info();
        let (sel, data) = encode_balance_of(OTHER_EOA);
        assert!(!erc20_self_only(&ctx(&info, sel, &data)));
    }

    #[test]
    fn erc20_self_only_allowance_self_owner_allowed() {
        let info = self_info();
        let (sel, data) = encode_allowance(SELF_EOA, OTHER_EOA);
        assert!(erc20_self_only(&ctx(&info, sel, &data)));
        // spender being self also passes (we don't constrain spender)
        let (sel2, data2) = encode_allowance(SELF_EOA, SELF_EOA);
        assert!(erc20_self_only(&ctx(&info, sel2, &data2)));
    }

    #[test]
    fn erc20_self_only_allowance_other_owner_denied() {
        let info = self_info();
        let (sel, data) = encode_allowance(OTHER_EOA, SELF_EOA);
        assert!(!erc20_self_only(&ctx(&info, sel, &data)));
    }

    #[test]
    fn erc20_self_only_missing_eoa_denies() {
        let info = json!({});
        let (sel, data) = encode_balance_of(SELF_EOA);
        assert!(!erc20_self_only(&ctx(&info, sel, &data)));
    }

    #[test]
    fn erc20_self_only_unparseable_eoa_denies() {
        let info = json!({ "eoa": "not-an-address" });
        let (sel, data) = encode_balance_of(SELF_EOA);
        assert!(!erc20_self_only(&ctx(&info, sel, &data)));
    }

    #[test]
    fn erc20_self_only_unknown_selector_denies() {
        let info = self_info();
        let (_, data) = encode_balance_of(SELF_EOA);
        // Attach to the transfer selector — not one the lambda knows.
        assert!(!erc20_self_only(&ctx(&info, TRANSFER_SELECTOR, &data)));
    }

    #[test]
    fn erc20_self_only_malformed_calldata_denies() {
        let info = self_info();
        assert!(!erc20_self_only(&ctx(&info, BALANCE_OF_SELECTOR, &[1, 2, 3])));
    }

    #[test]
    fn transfer_under_limit_decodes_and_compares() {
        let call = transferCall {
            to: Address::ZERO,
            amount: U256::from(100u64),
        };
        let data = call.abi_encode();
        let selector: [u8; 4] = data[0..4].try_into().unwrap();
        assert_eq!(selector, TRANSFER_SELECTOR);

        let allow_info = json!({ "max_transfer": "1000" });
        let deny_info = json!({ "max_transfer": "50" });
        assert!(transfer_under_limit(&ctx(&allow_info, selector, &data)));
        assert!(!transfer_under_limit(&ctx(&deny_info, selector, &data)));
    }
}
