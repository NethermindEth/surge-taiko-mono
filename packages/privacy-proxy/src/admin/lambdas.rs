use axum::Json;
use serde::Serialize;

use crate::acl::lambdas::user as user_lambdas;
use crate::error::ApiResult;
use crate::roles::{ROLE_ADMIN, ROLE_USER};
use crate::rpc::gated_methods;

#[derive(Serialize)]
pub struct LambdaDescriptor {
    pub name: &'static str,
    pub description: &'static str,
    /// 4-byte selectors the lambda is built to evaluate. Empty array means
    /// the lambda is selector-agnostic.
    pub expected_selectors: Vec<String>,
}

#[derive(Serialize)]
pub struct RoleLambdas {
    pub role: &'static str,
    pub lambdas: Vec<LambdaDescriptor>,
}

#[derive(Serialize)]
pub struct SyntheticSelector {
    pub method: &'static str,
    pub selector: String,
}

/// Capability 9: `GET /admin/registry/lambdas`. Returns per-role groups
/// in the order roles are declared in `src/roles.rs::ROLES`. Roles
/// without lambdas (e.g. `admin`) return an empty `lambdas` array so
/// the response shape is stable.
pub async fn list_lambdas() -> ApiResult<Json<Vec<RoleLambdas>>> {
    let out = vec![
        RoleLambdas {
            role: ROLE_ADMIN,
            lambdas: Vec::new(),
        },
        RoleLambdas {
            role: ROLE_USER,
            lambdas: user_lambdas::list_specs()
                .into_iter()
                .map(|s| LambdaDescriptor {
                    name: s.name,
                    description: s.description,
                    expected_selectors: s
                        .expected_selectors
                        .iter()
                        .map(|b| format!("0x{}", hex::encode(b)))
                        .collect(),
                })
                .collect(),
        },
    ];
    Ok(Json(out))
}

/// Capability 10: `GET /admin/registry/synthetic-selectors`
pub async fn list_synthetic_selectors() -> ApiResult<Json<Vec<SyntheticSelector>>> {
    let out = gated_methods::ALL
        .iter()
        .map(|m| SyntheticSelector {
            method: m.name,
            selector: format!("0x{}", hex::encode(m.selector)),
        })
        .collect();
    Ok(Json(out))
}
