pub mod challenge;
pub mod middleware;
pub mod verify;

use alloy::primitives::Address;
use serde::Serialize;

use crate::roles::{ROLE_ADMIN, ROLE_USER};

pub use middleware::caller_ctx_layer;

#[derive(Clone, Debug, Serialize)]
pub struct AdminCallerInfo {
    pub eoa: Address,
}

#[derive(Clone, Debug, Serialize)]
pub struct UserCallerInfo {
    pub eoa: Address,
    pub kyc: bool,
    pub blacklisted: bool,
}

#[derive(Clone, Debug, Serialize)]
#[serde(tag = "role", rename_all = "lowercase")]
pub enum CallerAttributes {
    Admin(AdminCallerInfo),
    User(UserCallerInfo),
}

impl CallerAttributes {
    pub fn role_name(&self) -> &'static str {
        match self {
            CallerAttributes::Admin(_) => ROLE_ADMIN,
            CallerAttributes::User(_) => ROLE_USER,
        }
    }

    pub fn eoa(&self) -> Address {
        match self {
            CallerAttributes::Admin(a) => a.eoa,
            CallerAttributes::User(u) => u.eoa,
        }
    }
}

#[derive(Clone, Debug, Default, Serialize)]
pub struct CallerCtx {
    pub eoa: Option<Address>,
    pub attributes: Option<CallerAttributes>,
}

impl CallerCtx {
    pub fn anonymous() -> Self {
        Self {
            eoa: None,
            attributes: None,
        }
    }

    pub fn is_admin(&self) -> bool {
        matches!(self.attributes, Some(CallerAttributes::Admin(_)))
    }

    pub fn is_anonymous(&self) -> bool {
        self.attributes.is_none()
    }

    pub fn role_name(&self) -> Option<&'static str> {
        self.attributes.as_ref().map(|a| a.role_name())
    }
}

pub fn build_signin_message(domain: &str, address: &Address, nonce: &str) -> String {
    format!(
        "{domain} sign-in\nAddress: {addr}\nNonce: {nonce}",
        domain = domain,
        addr = format_address(address),
        nonce = nonce,
    )
}

pub fn format_address(address: &Address) -> String {
    format!("0x{}", hex::encode(address.as_slice()))
}
