use alloy::sol;
use alloy::sol_types::SolCall;

use super::UserCallerInfo;
use crate::acl::lambdas::LambdaCtx;

/// keccak256("balanceOf(address)")[..4]
pub const BALANCE_OF_SELECTOR: [u8; 4] = [0x70, 0xa0, 0x82, 0x31];
/// keccak256("allowance(address,address)")[..4]
pub const ALLOWANCE_SELECTOR: [u8; 4] = [0xdd, 0x62, 0xed, 0x3e];

sol! {
    function balanceOf(address account) external view returns (uint256);
    function allowance(address owner, address spender) external view returns (uint256);
}

/// For ERC-20 `balanceOf(address)` and `allowance(address,address)`:
/// allow only when the queried `account` (balanceOf) or `owner`
/// (allowance) equals the caller's EOA. Spender on allowance is not
/// constrained. Any other selector → false.
pub fn run(ctx: &LambdaCtx<UserCallerInfo>) -> bool {
    let eoa = ctx.caller_info.eoa;
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

#[cfg(test)]
mod tests {
    use super::*;

    const SELF_EOA: &str = "0x1111111111111111111111111111111111111111";
    const OTHER_EOA: &str = "0x2222222222222222222222222222222222222222";

    fn info(eoa_hex: &str) -> UserCallerInfo {
        UserCallerInfo {
            eoa: eoa_hex.parse().unwrap(),
            kyc: false,
            blacklisted: false,
        }
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

    fn ctx<'a>(info: &'a UserCallerInfo, sel: [u8; 4], data: &'a [u8]) -> LambdaCtx<'a, UserCallerInfo> {
        LambdaCtx {
            caller_info: info,
            selector: sel,
            call_data: data,
        }
    }

    #[test]
    fn selectors_match_constants() {
        let (sel, _) = encode_balance_of(SELF_EOA);
        assert_eq!(sel, BALANCE_OF_SELECTOR);
        let (sel, _) = encode_allowance(SELF_EOA, OTHER_EOA);
        assert_eq!(sel, ALLOWANCE_SELECTOR);
    }

    #[test]
    fn balance_of_self_allowed() {
        let i = info(SELF_EOA);
        let (sel, data) = encode_balance_of(SELF_EOA);
        assert!(run(&ctx(&i, sel, &data)));
    }

    #[test]
    fn balance_of_other_denied() {
        let i = info(SELF_EOA);
        let (sel, data) = encode_balance_of(OTHER_EOA);
        assert!(!run(&ctx(&i, sel, &data)));
    }

    #[test]
    fn allowance_self_owner_allowed() {
        let i = info(SELF_EOA);
        let (sel, data) = encode_allowance(SELF_EOA, OTHER_EOA);
        assert!(run(&ctx(&i, sel, &data)));
        // spender being self also passes (spender not constrained)
        let (sel2, data2) = encode_allowance(SELF_EOA, SELF_EOA);
        assert!(run(&ctx(&i, sel2, &data2)));
    }

    #[test]
    fn allowance_other_owner_denied() {
        let i = info(SELF_EOA);
        let (sel, data) = encode_allowance(OTHER_EOA, SELF_EOA);
        assert!(!run(&ctx(&i, sel, &data)));
    }

    #[test]
    fn unknown_selector_denies() {
        let i = info(SELF_EOA);
        let (_, data) = encode_balance_of(SELF_EOA);
        // Selector outside the whitelist.
        assert!(!run(&ctx(&i, [0xaa, 0xbb, 0xcc, 0xdd], &data)));
    }

    #[test]
    fn malformed_calldata_denies() {
        let i = info(SELF_EOA);
        assert!(!run(&ctx(&i, BALANCE_OF_SELECTOR, &[1, 2, 3])));
    }
}
