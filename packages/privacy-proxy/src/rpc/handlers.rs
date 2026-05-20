use alloy::primitives::Address;
use axum::body::Bytes;
use axum::extract::State;
use axum::response::{IntoResponse, Response};
use axum::Extension;
use axum::Json;
use serde::Deserialize;
use serde_json::{json, Value};

use crate::acl::evaluator::DenyReason;
use crate::acl::registry as acl_registry;
use crate::acl::{check_call, AccessDecision};
use crate::auth::CallerCtx;
use crate::rpc::gated_methods;
use crate::state::AppState;
use crate::tracer::{decode_raw_tx, trace_call, CallSite};

/// Whitelist of `eth_` methods every caller (including anonymous) is
/// allowed to invoke. These return only chain-global state with no
/// per-user data: chain id, latest block number, gas/fee oracles, sync
/// status. Everything else — block contents, logs, receipts, calls,
/// reads of address-parameterized state — requires authentication.
const ALWAYS_PUBLIC_METHODS: &[&str] = &[
    "eth_chainId",
    "eth_blockNumber",
    "eth_gasPrice",
    "eth_maxPriorityFeePerGas",
    "eth_feeHistory",
    "eth_syncing",
    "eth_protocolVersion",
];

#[derive(Debug, Deserialize)]
struct JsonRpcRequest {
    #[allow(dead_code)]
    #[serde(default)]
    jsonrpc: Option<String>,
    method: String,
    #[serde(default)]
    params: Value,
    #[serde(default)]
    id: Value,
}

pub async fn dispatch(
    State(state): State<AppState>,
    Extension(ctx): Extension<CallerCtx>,
    body: Bytes,
) -> Response {
    // Reject batch requests for the POC. Single requests only.
    if body.first() == Some(&b'[') {
        return rpc_error_resp(
            Value::Null,
            -32600,
            "batch JSON-RPC requests are not supported",
            None,
        );
    }

    let raw: Value = match serde_json::from_slice(&body) {
        Ok(v) => v,
        Err(e) => {
            return rpc_error_resp(
                Value::Null,
                -32700,
                &format!("parse error: {e}"),
                None,
            );
        }
    };
    let req: JsonRpcRequest = match serde_json::from_value(raw.clone()) {
        Ok(r) => r,
        Err(e) => {
            return rpc_error_resp(
                raw.get("id").cloned().unwrap_or(Value::Null),
                -32600,
                &format!("invalid request: {e}"),
                None,
            );
        }
    };

    let id = req.id.clone();

    if !req.method.starts_with("eth_") {
        return rpc_error_resp(
            id,
            -32601,
            "only eth_ namespace methods are accepted by this proxy",
            None,
        );
    }

    // Strict default: every method except a small whitelist of
    // chain-global status calls requires authentication. Anonymous
    // callers can fetch chain id / latest block / gas prices and
    // nothing else.
    if ctx.is_anonymous() && !ALWAYS_PUBLIC_METHODS.contains(&req.method.as_str()) {
        return access_denied_resp(
            id,
            Address::ZERO,
            [0; 4],
            &format!("{:?}", DenyReason::AnonymousAgainstGatedCall),
        );
    }

    // Gated address-parameterized read methods (eth_getBalance et al.).
    if let Some(method) = gated_methods::lookup_by_method(&req.method) {
        if let Err(resp) =
            check_gated_method(&state, &ctx, method, &req.params, &id).await
        {
            return resp;
        }
        return forward(&state, &raw, id).await;
    }

    match req.method.as_str() {
        "eth_call" | "eth_estimateGas" | "eth_createAccessList" | "eth_sendTransaction" => {
            let call_obj = req.params.get(0).cloned().unwrap_or(Value::Null);
            let block = req.params.get(1).cloned().unwrap_or_else(|| json!("latest"));
            handle_with_acl(&state, &ctx, &raw, id, call_obj, block).await
        }
        "eth_sendRawTransaction" => {
            let raw_tx = match req.params.get(0).and_then(|v| v.as_str()) {
                Some(s) => s,
                None => {
                    return rpc_error_resp(
                        id,
                        -32602,
                        "eth_sendRawTransaction requires a hex string parameter",
                        None,
                    );
                }
            };
            let decoded = match decode_raw_tx(raw_tx) {
                Ok(d) => d,
                Err(e) => {
                    return rpc_error_resp(
                        id,
                        -32602,
                        &format!("invalid raw transaction: {e}"),
                        None,
                    );
                }
            };
            let to_value = match decoded.to {
                Some(a) => json!(format!("0x{}", hex::encode(a.as_slice()))),
                None => Value::Null,
            };
            let call_obj = json!({
                "from": format!("0x{}", hex::encode(decoded.from.as_slice())),
                "to": to_value,
                "value": format!("0x{:x}", decoded.value),
                "data": format!("0x{}", hex::encode(&decoded.input)),
            });
            handle_with_acl(&state, &ctx, &raw, id, call_obj, json!("latest")).await
        }
        _ => forward(&state, &raw, id).await,
    }
}

/// Returns Ok(()) when the gated method should fall through to forwarding,
/// or Err(deny_response) when the request must be rejected.
async fn check_gated_method(
    state: &AppState,
    ctx: &CallerCtx,
    method: gated_methods::GatedMethod,
    params: &Value,
    id: &Value,
) -> Result<(), Response> {
    let Some(target) = gated_methods::extract_target(params) else {
        return Err(rpc_error_resp(
            id.clone(),
            -32602,
            &format!("{}: missing or invalid target address", method.name),
            None,
        ));
    };

    let target_hex = format!("0x{}", hex::encode(target.as_slice()));
    let selector_hex = format!("0x{}", hex::encode(method.selector));
    let rule_opt = match acl_registry::find_rule(&state.pool, &target_hex, &selector_hex).await
    {
        Ok(r) => r,
        Err(e) => {
            return Err(rpc_error_resp(
                id.clone(),
                -32000,
                &format!("registry lookup failed: {e}"),
                None,
            ));
        }
    };

    if rule_opt.is_some() {
        let call_data = match gated_methods::encode_call_data(method, target, params) {
            Ok(d) => d,
            Err(msg) => {
                return Err(rpc_error_resp(id.clone(), -32602, msg, None));
            }
        };
        let msg_sender = ctx.eoa.unwrap_or(Address::ZERO);
        return match check_call(&state.pool, ctx, &target, &call_data, msg_sender).await {
            Ok(AccessDecision::Allow) => Ok(()),
            Ok(AccessDecision::Deny {
                contract,
                selector,
                reason,
            }) => Err(access_denied_resp(
                id.clone(),
                contract,
                selector,
                &format!("{reason:?}"),
            )),
            Err(e) => Err(rpc_error_resp(
                id.clone(),
                -32000,
                &format!("acl error: {e}"),
                None,
            )),
        };
    }

    // No admin rule → apply default behavior.
    if ctx.is_anonymous() {
        return Err(access_denied_resp(
            id.clone(),
            target,
            method.selector,
            &format!("{:?}", DenyReason::AnonymousAgainstGatedCall),
        ));
    }

    match state.upstream.is_contract(target).await {
        Ok(true) => Ok(()),
        Ok(false) => {
            if ctx.eoa == Some(target) {
                Ok(())
            } else {
                Err(access_denied_resp(
                    id.clone(),
                    target,
                    method.selector,
                    &format!("{:?}", DenyReason::DefaultEoaSelfOnly),
                ))
            }
        }
        Err(e) => Err(rpc_error_resp(
            id.clone(),
            -32000,
            &format!("is_contract lookup failed: {e}"),
            None,
        )),
    }
}

async fn handle_with_acl(
    state: &AppState,
    ctx: &CallerCtx,
    raw: &Value,
    id: Value,
    call_obj: Value,
    block: Value,
) -> Response {
    // Fast-path: an empty call (no `to`) is a CREATE — currently allow.
    let to_opt = call_obj
        .get("to")
        .and_then(|v| v.as_str())
        .and_then(|s| s.parse::<Address>().ok());
    let input_hex = call_obj
        .get("data")
        .or_else(|| call_obj.get("input"))
        .and_then(|v| v.as_str())
        .unwrap_or("0x");
    let input_bytes = hex::decode(input_hex.trim_start_matches("0x")).unwrap_or_default();

    let tx_origin = ctx.eoa.unwrap_or(Address::ZERO);

    // Top-level ACL check (cheap; avoids tracer roundtrip when obviously denied).
    if let Some(to) = to_opt {
        match check_call(&state.pool, ctx, &to, &input_bytes, tx_origin).await {
            Ok(AccessDecision::Allow) => {}
            Ok(AccessDecision::Deny {
                contract,
                selector,
                reason,
            }) => {
                return access_denied_resp(id, contract, selector, &format!("{reason:?}"));
            }
            Err(e) => {
                return rpc_error_resp(id, -32000, &format!("acl error: {e}"), None);
            }
        }
    }

    // Trace and check every internal frame. We deliberately forward the
    // user-supplied `call_obj` (including its `from`) verbatim to
    // `debug_traceCall`: the top-level ACL above already evaluated with
    // `tx_origin = ctx.eoa` and the inner frames are populated by the
    // EVM, so the simulator's view of `from` doesn't affect authorization.
    let frame = match trace_call(&state.upstream, call_obj, block).await {
        Ok(f) => f,
        Err(e) => {
            tracing::warn!("debug_traceCall failed: {e}");
            return rpc_error_resp(
                id,
                -32000,
                &format!("internal trace failed: {e}"),
                None,
            );
        }
    };

    let mut sites: Vec<CallSite> = Vec::new();
    for inner in &frame.calls {
        inner.flatten(&mut sites);
    }

    for site in sites {
        match check_call(&state.pool, ctx, &site.contract, &site.input, site.from).await {
            Ok(AccessDecision::Allow) => {}
            Ok(AccessDecision::Deny {
                contract,
                selector,
                reason,
            }) => {
                return access_denied_resp(id, contract, selector, &format!("internal: {reason:?}"));
            }
            Err(e) => {
                return rpc_error_resp(id, -32000, &format!("acl error: {e}"), None);
            }
        }
    }

    forward(state, raw, id).await
}

async fn forward(state: &AppState, raw: &Value, id: Value) -> Response {
    match state.upstream.forward(raw).await {
        Ok(v) => Json(v).into_response(),
        Err(e) => rpc_error_resp(id, -32000, &format!("upstream error: {e}"), None),
    }
}

fn rpc_error_resp(id: Value, code: i64, message: &str, data: Option<Value>) -> Response {
    let mut err = json!({ "code": code, "message": message });
    if let Some(d) = data {
        err["data"] = d;
    }
    Json(json!({
        "jsonrpc": "2.0",
        "id": id,
        "error": err,
    }))
    .into_response()
}

fn access_denied_resp(id: Value, contract: Address, selector: [u8; 4], detail: &str) -> Response {
    rpc_error_resp(
        id,
        -32001,
        "access denied",
        Some(json!({
            "contract": format!("0x{}", hex::encode(contract.as_slice())),
            "selector": format!("0x{}", hex::encode(selector)),
            "detail": detail,
        })),
    )
}
