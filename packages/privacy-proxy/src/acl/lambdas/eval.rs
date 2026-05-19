use alloy::primitives::{Address, U256};

use super::attributes::attribute_word;
use super::{Condition, Lambda, Lhs, Rhs};
use crate::auth::CallerAttributes;

pub fn evaluate(
    lambda: &Lambda,
    attrs: &CallerAttributes,
    selector: [u8; 4],
    call_data: &[u8],
    tx_origin: Address,
    msg_sender: Address,
) -> bool {
    let matching = lambda.rules.iter().filter(|r| r.selector == selector);
    for rule in matching {
        let Some(lhs) = resolve_lhs(&rule.lhs, attrs, call_data) else {
            return false;
        };
        let Some(rhs) = resolve_rhs(&rule.rhs, tx_origin, msg_sender) else {
            return false;
        };
        if !compare(rule.condition, lhs, rhs) {
            return false;
        }
    }
    true
}

fn resolve_lhs(lhs: &Lhs, attrs: &CallerAttributes, call_data: &[u8]) -> Option<[u8; 32]> {
    match lhs {
        Lhs::Calldata { offset } => {
            let start = *offset as usize;
            let end = start.checked_add(32)?;
            if end > call_data.len() {
                return None;
            }
            let mut w = [0u8; 32];
            w.copy_from_slice(&call_data[start..end]);
            Some(w)
        }
        Lhs::Attribute { name } => attribute_word(attrs, name),
    }
}

fn resolve_rhs(rhs: &Rhs, tx_origin: Address, msg_sender: Address) -> Option<[u8; 32]> {
    match rhs {
        Rhs::TxOrigin => Some(address_word(tx_origin)),
        Rhs::MsgSender => Some(address_word(msg_sender)),
        Rhs::Literal { value_hex } => parse_word(value_hex),
    }
}

fn address_word(addr: Address) -> [u8; 32] {
    let mut w = [0u8; 32];
    w[12..].copy_from_slice(addr.as_slice());
    w
}

fn parse_word(s: &str) -> Option<[u8; 32]> {
    let trimmed = s.trim_start_matches("0x");
    if trimmed.len() != 64 {
        return None;
    }
    let bytes = hex::decode(trimmed).ok()?;
    let mut w = [0u8; 32];
    w.copy_from_slice(&bytes);
    Some(w)
}

fn compare(cond: Condition, lhs: [u8; 32], rhs: [u8; 32]) -> bool {
    let a = U256::from_be_bytes(lhs);
    let b = U256::from_be_bytes(rhs);
    match cond {
        Condition::Eq => a == b,
        Condition::Neq => a != b,
        Condition::Gt => a > b,
        Condition::Lt => a < b,
        Condition::Gte => a >= b,
        Condition::Lte => a <= b,
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::acl::lambdas::{Lambda, LambdaRule};
    use crate::auth::{AdminCallerInfo, UserCallerInfo};
    use alloy::primitives::Address;

    fn user(eoa: Address, kyc: bool) -> CallerAttributes {
        CallerAttributes::User(UserCallerInfo {
            eoa,
            kyc,
            blacklisted: false,
        })
    }

    fn admin(eoa: Address) -> CallerAttributes {
        CallerAttributes::Admin(AdminCallerInfo { eoa })
    }

    fn lambda_with(rules: Vec<LambdaRule>) -> Lambda {
        Lambda {
            id: 1,
            name: "test".into(),
            role_id: 1,
            role: "user".into(),
            description: None,
            rules,
        }
    }

    fn rule(selector: [u8; 4], lhs: Lhs, condition: Condition, rhs: Rhs) -> LambdaRule {
        LambdaRule {
            id: 0,
            selector,
            lhs,
            condition,
            rhs,
        }
    }

    fn balance_of_calldata(account: Address) -> Vec<u8> {
        let mut out = vec![0x70, 0xa0, 0x82, 0x31];
        out.extend_from_slice(&[0u8; 12]);
        out.extend_from_slice(account.as_slice());
        out
    }

    const BALANCE_OF: [u8; 4] = [0x70, 0xa0, 0x82, 0x31];

    #[test]
    fn vacuous_true_when_no_matching_rules() {
        let l = lambda_with(vec![rule(
            [0xaa, 0xbb, 0xcc, 0xdd],
            Lhs::Attribute { name: "kyc".into() },
            Condition::Eq,
            Rhs::Literal {
                value_hex: format!("0x{}", "00".repeat(31)) + "01",
            },
        )]);
        let eoa: Address = "0x1111111111111111111111111111111111111111".parse().unwrap();
        assert!(evaluate(&l, &user(eoa, false), BALANCE_OF, &[], eoa, eoa));
    }

    #[test]
    fn calldata_eq_tx_origin_pass() {
        let l = lambda_with(vec![rule(
            BALANCE_OF,
            Lhs::Calldata { offset: 4 },
            Condition::Eq,
            Rhs::TxOrigin,
        )]);
        let eoa: Address = "0x1111111111111111111111111111111111111111".parse().unwrap();
        let data = balance_of_calldata(eoa);
        assert!(evaluate(&l, &user(eoa, false), BALANCE_OF, &data, eoa, eoa));
    }

    #[test]
    fn calldata_eq_tx_origin_reject() {
        let l = lambda_with(vec![rule(
            BALANCE_OF,
            Lhs::Calldata { offset: 4 },
            Condition::Eq,
            Rhs::TxOrigin,
        )]);
        let caller: Address = "0x1111111111111111111111111111111111111111".parse().unwrap();
        let other: Address = "0x3333333333333333333333333333333333333333".parse().unwrap();
        let data = balance_of_calldata(other);
        assert!(!evaluate(&l, &user(caller, false), BALANCE_OF, &data, caller, caller));
    }

    #[test]
    fn attribute_kyc_eq_literal_one() {
        let one_word = format!("0x{}01", "00".repeat(31));
        let l = lambda_with(vec![rule(
            BALANCE_OF,
            Lhs::Attribute { name: "kyc".into() },
            Condition::Eq,
            Rhs::Literal { value_hex: one_word },
        )]);
        let eoa: Address = Address::ZERO;
        let data = balance_of_calldata(eoa);
        assert!(evaluate(&l, &user(eoa, true), BALANCE_OF, &data, eoa, eoa));
        assert!(!evaluate(&l, &user(eoa, false), BALANCE_OF, &data, eoa, eoa));
    }

    #[test]
    fn attribute_blacklisted_neq_literal_one() {
        let one_word = format!("0x{}01", "00".repeat(31));
        let l = lambda_with(vec![rule(
            BALANCE_OF,
            Lhs::Attribute { name: "blacklisted".into() },
            Condition::Neq,
            Rhs::Literal { value_hex: one_word },
        )]);
        let eoa: Address = Address::ZERO;
        // blacklisted defaults to false → 0; 0 != 1 → true → rule passes
        assert!(evaluate(&l, &user(eoa, true), BALANCE_OF, &[], eoa, eoa));
    }

    #[test]
    fn unknown_attribute_returns_false() {
        let one = format!("0x{}01", "00".repeat(31));
        let l = lambda_with(vec![rule(
            BALANCE_OF,
            Lhs::Attribute { name: "nonexistent".into() },
            Condition::Eq,
            Rhs::Literal { value_hex: one },
        )]);
        let eoa: Address = Address::ZERO;
        assert!(!evaluate(&l, &user(eoa, true), BALANCE_OF, &[], eoa, eoa));
    }

    #[test]
    fn calldata_out_of_bounds_returns_false() {
        let l = lambda_with(vec![rule(
            BALANCE_OF,
            Lhs::Calldata { offset: 4 },
            Condition::Eq,
            Rhs::TxOrigin,
        )]);
        let eoa: Address = Address::ZERO;
        assert!(!evaluate(&l, &user(eoa, false), BALANCE_OF, &[0, 1, 2], eoa, eoa));
    }

    #[test]
    fn all_conditions() {
        let two = format!("0x{}02", "00".repeat(31));
        let three = format!("0x{}03", "00".repeat(31));
        let mk = |cond: Condition, v: &str| {
            lambda_with(vec![rule(
                BALANCE_OF,
                Lhs::Attribute { name: "kyc".into() },
                cond,
                Rhs::Literal {
                    value_hex: v.to_string(),
                },
            )])
        };
        let eoa: Address = Address::ZERO;
        let attrs = user(eoa, true); // kyc=true → 1
        let call = balance_of_calldata(eoa);
        let one = format!("0x{}01", "00".repeat(31));
        assert!(evaluate(&mk(Condition::Eq, &one), &attrs, BALANCE_OF, &call, eoa, eoa));
        assert!(evaluate(&mk(Condition::Neq, &two), &attrs, BALANCE_OF, &call, eoa, eoa));
        assert!(evaluate(&mk(Condition::Lt, &two), &attrs, BALANCE_OF, &call, eoa, eoa));
        assert!(evaluate(&mk(Condition::Lte, &one), &attrs, BALANCE_OF, &call, eoa, eoa));
        assert!(!evaluate(&mk(Condition::Gt, &two), &attrs, BALANCE_OF, &call, eoa, eoa));
        assert!(!evaluate(&mk(Condition::Gte, &three), &attrs, BALANCE_OF, &call, eoa, eoa));
    }

    #[test]
    fn msg_sender_resolves_independent_of_tx_origin() {
        // Calldata arg0 must equal msg_sender, even when tx_origin differs.
        let l = lambda_with(vec![rule(
            BALANCE_OF,
            Lhs::Calldata { offset: 4 },
            Condition::Eq,
            Rhs::MsgSender,
        )]);
        let origin: Address = "0x1111111111111111111111111111111111111111".parse().unwrap();
        let sender: Address = "0x2222222222222222222222222222222222222222".parse().unwrap();
        let data = balance_of_calldata(sender);
        assert!(evaluate(&l, &user(origin, false), BALANCE_OF, &data, origin, sender));
        // Same calldata, sender = origin → still passes only when arg0 = sender.
        let data2 = balance_of_calldata(origin);
        assert!(!evaluate(&l, &user(origin, false), BALANCE_OF, &data2, origin, sender));
    }

    #[test]
    fn admin_attributes_available() {
        let l = lambda_with(vec![rule(
            BALANCE_OF,
            Lhs::Attribute { name: "eoa".into() },
            Condition::Eq,
            Rhs::TxOrigin,
        )]);
        let a: Address = "0x4444444444444444444444444444444444444444".parse().unwrap();
        assert!(evaluate(&l, &admin(a), BALANCE_OF, &[], a, a));
    }

    #[test]
    fn malformed_literal_returns_false() {
        let l = lambda_with(vec![rule(
            BALANCE_OF,
            Lhs::Attribute { name: "kyc".into() },
            Condition::Eq,
            Rhs::Literal { value_hex: "0xdead".into() },
        )]);
        let eoa: Address = Address::ZERO;
        assert!(!evaluate(&l, &user(eoa, true), BALANCE_OF, &[], eoa, eoa));
    }
}
