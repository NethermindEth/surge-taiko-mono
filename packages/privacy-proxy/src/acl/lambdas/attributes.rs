use crate::auth::CallerAttributes;
use crate::roles::{ROLE_ADMIN, ROLE_USER};

pub fn attribute_word(attrs: &CallerAttributes, name: &str) -> Option<[u8; 32]> {
    match attrs {
        CallerAttributes::Admin(a) => match name {
            "eoa" => Some(address_word(a.eoa.as_slice())),
            _ => None,
        },
        CallerAttributes::User(u) => match name {
            "eoa" => Some(address_word(&u.eoa.as_slice())),
            "kyc" => Some(bool_word(u.kyc)),
            "blacklisted" => Some(bool_word(u.blacklisted)),
            _ => None,
        },
    }
}

#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub struct AttributeSpec {
    pub name: &'static str,
    pub ty: &'static str,
}

pub fn known_attribute_specs_for_role(role: &str) -> &'static [AttributeSpec] {
    match role {
        ROLE_ADMIN => &[AttributeSpec { name: "eoa", ty: "address" }],
        ROLE_USER => &[
            AttributeSpec { name: "eoa", ty: "address" },
            AttributeSpec { name: "kyc", ty: "bool" },
            AttributeSpec { name: "blacklisted", ty: "bool" },
        ],
        _ => &[],
    }
}

fn address_word(addr: &[u8]) -> [u8; 32] {
    let mut w = [0u8; 32];
    let start = 32 - addr.len();
    w[start..].copy_from_slice(addr);
    w
}

fn bool_word(b: bool) -> [u8; 32] {
    let mut w = [0u8; 32];
    if b {
        w[31] = 1;
    }
    w
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::auth::{AdminCallerInfo, UserCallerInfo};
    use alloy::primitives::Address;

    #[test]
    fn user_attribute_words() {
        let eoa: Address = "0x1111111111111111111111111111111111111111".parse().unwrap();
        let attrs = CallerAttributes::User(UserCallerInfo {
            eoa,
            kyc: true,
            blacklisted: false,
        });
        let kyc = attribute_word(&attrs, "kyc").unwrap();
        assert_eq!(kyc[31], 1);
        assert!(kyc[..31].iter().all(|b| *b == 0));

        let bl = attribute_word(&attrs, "blacklisted").unwrap();
        assert!(bl.iter().all(|b| *b == 0));

        let eoa_w = attribute_word(&attrs, "eoa").unwrap();
        assert_eq!(&eoa_w[12..], eoa.as_slice());
        assert!(eoa_w[..12].iter().all(|b| *b == 0));

        assert!(attribute_word(&attrs, "unknown").is_none());
    }

    #[test]
    fn admin_attribute_words() {
        let eoa: Address = "0x2222222222222222222222222222222222222222".parse().unwrap();
        let attrs = CallerAttributes::Admin(AdminCallerInfo { eoa });
        let w = attribute_word(&attrs, "eoa").unwrap();
        assert_eq!(&w[12..], eoa.as_slice());
        assert!(attribute_word(&attrs, "kyc").is_none());
    }
}
