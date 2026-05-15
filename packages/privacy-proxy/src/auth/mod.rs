pub mod challenge;
pub mod middleware;
pub mod verify;

use alloy::primitives::Address;
use serde::Serialize;
use serde_json::Value;

pub use middleware::caller_ctx_layer;

/// Resolved caller for a single request. `eoa` and `role` are `None` for
/// anonymous (no token / expired token / unknown token) requests.
#[derive(Clone, Debug, Default, Serialize)]
pub struct CallerCtx {
    pub eoa: Option<Address>,
    pub role: Option<String>,
    pub caller_info: Value,
}

impl CallerCtx {
    pub fn anonymous() -> Self {
        Self {
            eoa: None,
            role: None,
            caller_info: Value::Null,
        }
    }

    pub fn is_admin(&self) -> bool {
        self.role.as_deref() == Some("admin")
    }

    pub fn is_anonymous(&self) -> bool {
        self.role.is_none()
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
