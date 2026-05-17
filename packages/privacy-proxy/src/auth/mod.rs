pub mod challenge;
pub mod middleware;
pub mod verify;

use alloy::primitives::Address;
use serde::Serialize;

use crate::acl::lambdas::user::UserCallerInfo;
use crate::roles::{ROLE_ADMIN, ROLE_USER};

pub use middleware::caller_ctx_layer;

/// Typed attributes for the `admin` role. Identity-only: an admin has no
/// extra state beyond the EOA used to authenticate.
#[derive(Clone, Debug, Serialize)]
pub struct AdminCallerInfo {
    pub eoa: Address,
}

/// Tagged union over the per-role attribute structs. The evaluator
/// unwraps this to dispatch to the matching role's lambda registry.
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

/// Resolved caller for a single request. `eoa` and `attributes` are
/// `None` for anonymous (no token / expired token / unknown token)
/// requests.
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

/// EIP-191 personal_sign message the wallet is asked to sign.
/// Format is stable across the binary; any change is a wire-breaking change.
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
