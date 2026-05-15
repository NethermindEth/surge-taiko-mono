use axum::Json;
use serde::Serialize;

use crate::acl::lambdas;
use crate::error::ApiResult;
use crate::rpc::gated_methods;

#[derive(Serialize)]
pub struct LambdaDescriptor {
    pub name: &'static str,
    pub description: &'static str,
    pub expected_selector: Option<String>,
}

#[derive(Serialize)]
pub struct SyntheticSelector {
    pub method: &'static str,
    pub selector: String,
}

/// Capability 9: `GET /admin/registry/lambdas`
pub async fn list_lambdas() -> ApiResult<Json<Vec<LambdaDescriptor>>> {
    let specs = lambdas::list_specs();
    let out = specs
        .into_iter()
        .map(|s| LambdaDescriptor {
            name: s.name,
            description: s.description,
            expected_selector: s
                .expected_selector
                .map(|b| format!("0x{}", hex::encode(b))),
        })
        .collect();
    Ok(Json(out))
}

/// Capability 19: `GET /admin/registry/synthetic-selectors`
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
