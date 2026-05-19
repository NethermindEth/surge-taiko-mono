use alloy::consensus::transaction::SignerRecoverable;
use alloy::consensus::{Transaction, TxEnvelope};
use alloy::eips::eip2718::Decodable2718;
use alloy::primitives::{Address, TxKind, U256};
use anyhow::{anyhow, Context, Result};
use serde::Deserialize;
use serde_json::{json, Value};

use crate::upstream::UpstreamClient;

/// One frame in a `callTracer` result, recursively containing inner calls.
#[derive(Debug, Deserialize, Clone)]
pub struct CallFrame {
    #[serde(rename = "type")]
    pub typ: String,
    pub from: Option<String>,
    pub to: Option<String>,
    pub input: Option<String>,
    pub value: Option<String>,
    #[serde(default)]
    pub calls: Vec<CallFrame>,
    #[serde(default)]
    pub error: Option<String>,
}

/// A single (caller, contract, calldata) triple extracted from a call frame.
/// `from` is the per-frame msg.sender — for the top-level frame this equals
/// the original tx.from, for internal frames it's the parent contract's
/// address.
#[derive(Debug, Clone)]
pub struct CallSite {
    pub from: Address,
    pub contract: Address,
    pub input: Vec<u8>,
}

impl CallFrame {
    pub fn flatten(&self, out: &mut Vec<CallSite>) {
        // Only CALL / STATICCALL / DELEGATECALL / CALLCODE produce a callee
        // that should be subject to ACL. CREATE / CREATE2 / SELFDESTRUCT are
        // ignored (no `to` selector to gate).
        let gateable = matches!(
            self.typ.as_str(),
            "CALL" | "STATICCALL" | "DELEGATECALL" | "CALLCODE"
        );
        if gateable {
            if let (Some(from), Some(to), Some(input)) = (&self.from, &self.to, &self.input) {
                let trimmed = input.trim_start_matches("0x");
                if let (Ok(from_addr), Ok(to_addr), Ok(bytes)) = (
                    from.parse::<Address>(),
                    to.parse::<Address>(),
                    hex::decode(trimmed),
                ) {
                    out.push(CallSite {
                        from: from_addr,
                        contract: to_addr,
                        input: bytes,
                    });
                }
            }
        }
        for c in &self.calls {
            c.flatten(out);
        }
    }
}

/// Ask the upstream node for the full call tree of a hypothetical call.
pub async fn trace_call(
    upstream: &UpstreamClient,
    call_object: Value,
    block: Value,
) -> Result<CallFrame> {
    let tracer_cfg = json!({ "tracer": "callTracer" });
    let params = json!([call_object, block, tracer_cfg]);
    let result = upstream
        .call("debug_traceCall", params)
        .await
        .context("debug_traceCall request failed")?;
    serde_json::from_value(result).context("decode callTracer response")
}

/// Decoded outline of a signed raw transaction. Enough to issue a
/// `debug_traceCall` simulation against the upstream.
#[derive(Debug, Clone)]
pub struct DecodedRawTx {
    pub from: Address,
    pub to: Option<Address>,
    pub value: U256,
    pub input: Vec<u8>,
}

pub fn decode_raw_tx(raw_hex: &str) -> Result<DecodedRawTx> {
    let bytes = hex::decode(raw_hex.trim_start_matches("0x"))
        .map_err(|e| anyhow!("invalid raw tx hex: {e}"))?;
    let envelope = TxEnvelope::decode_2718(&mut bytes.as_slice())
        .map_err(|e| anyhow!("decode raw tx: {e}"))?;
    let from = envelope
        .recover_signer()
        .map_err(|e| anyhow!("recover signer: {e}"))?;
    let to = match envelope.kind() {
        TxKind::Call(a) => Some(a),
        TxKind::Create => None,
    };
    let value = envelope.value();
    let input = envelope.input().to_vec();
    Ok(DecodedRawTx {
        from,
        to,
        value,
        input,
    })
}
