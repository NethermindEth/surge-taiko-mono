use std::collections::HashMap;
use std::sync::{Mutex, OnceLock};
use std::time::{Duration, Instant};

use alloy::primitives::Address;
use anyhow::{Context, Result};
use serde_json::Value;

const IS_CONTRACT_TTL: Duration = Duration::from_secs(60);

#[derive(Debug)]
pub struct UpstreamClient {
    url: String,
    client: reqwest::Client,
    is_contract_cache: Mutex<HashMap<Address, (bool, Instant)>>,
    chain_id: OnceLock<u64>,
}

impl UpstreamClient {
    pub fn new(url: String) -> Self {
        let client = reqwest::Client::builder()
            .timeout(std::time::Duration::from_secs(30))
            .build()
            .expect("reqwest client should build");
        Self {
            url,
            client,
            is_contract_cache: Mutex::new(HashMap::new()),
            chain_id: OnceLock::new(),
        }
    }

    /// Lazily fetch and cache the upstream chain_id. Two concurrent first
    /// callers may both issue the RPC; whichever sets first wins and the
    /// other is discarded — benign since chain_id is immutable per chain.
    pub async fn chain_id(&self) -> Result<u64> {
        if let Some(&id) = self.chain_id.get() {
            return Ok(id);
        }
        let result = self.call("eth_chainId", serde_json::json!([])).await?;
        let id_hex = result
            .as_str()
            .context("eth_chainId returned non-string result")?;
        let id = u64::from_str_radix(id_hex.trim_start_matches("0x"), 16)
            .with_context(|| format!("invalid chain_id hex: {id_hex}"))?;
        let _ = self.chain_id.set(id);
        Ok(id)
    }

    /// Forward a raw JSON-RPC payload verbatim and return the response body.
    ///
    /// Non-2xx responses are still parsed and returned when their body is
    /// valid JSON: many JSON-RPC servers return structured `{ error: ... }`
    /// bodies with HTTP 4xx/5xx, and clients need the JSON-RPC error to
    /// surface to the caller. Only treat the response as a transport
    /// failure when the body isn't parseable JSON.
    pub async fn forward(&self, body: &Value) -> Result<Value> {
        // CodeQL: ssrf - false positive. `self.url` is set once from the
        // UPSTREAM_URL env var at boot (see config.rs); it is never derived
        // from request input.
        let resp = self
            .client
            .post(&self.url)
            .json(body)
            .send()
            .await?;
        let status = resp.status();
        let bytes = resp.bytes().await?;
        match serde_json::from_slice::<Value>(&bytes) {
            Ok(v) => Ok(v),
            Err(e) => {
                if status.is_success() {
                    Err(anyhow::anyhow!("decode upstream response: {e}"))
                } else {
                    let body_preview = String::from_utf8_lossy(&bytes);
                    let trimmed: String = body_preview.chars().take(256).collect();
                    Err(anyhow::anyhow!(
                        "upstream returned HTTP {status} with non-JSON body: {trimmed}"
                    ))
                }
            }
        }
    }

    /// Issue a one-off JSON-RPC call (used internally for debug_traceCall, eth_getCode).
    pub async fn call(&self, method: &str, params: Value) -> Result<Value> {
        let body = serde_json::json!({
            "jsonrpc": "2.0",
            "id": 1,
            "method": method,
            "params": params,
        });
        let resp = self.forward(&body).await?;
        if let Some(err) = resp.get("error") {
            anyhow::bail!("upstream error: {err}");
        }
        Ok(resp
            .get("result")
            .cloned()
            .unwrap_or(Value::Null))
    }

    /// True if `addr` has non-empty bytecode at latest. Result is cached for
    /// `IS_CONTRACT_TTL` to avoid a per-request roundtrip. The cache is
    /// safe to be stale for short windows: contracts don't become EOAs;
    /// EOAs that later deploy as contracts will be misclassified for at
    /// most one TTL window.
    pub async fn is_contract(&self, addr: Address) -> Result<bool> {
        if let Some(cached) = self.cache_get(&addr) {
            return Ok(cached);
        }
        let addr_hex = format!("0x{}", hex::encode(addr.as_slice()));
        let result = self
            .call("eth_getCode", serde_json::json!([addr_hex, "latest"]))
            .await?;
        let code = result.as_str().unwrap_or("0x");
        let trimmed = code.trim_start_matches("0x");
        let has_code = !trimmed.is_empty() && trimmed != "0";
        self.cache_put(addr, has_code);
        Ok(has_code)
    }

    fn cache_get(&self, addr: &Address) -> Option<bool> {
        let map = self.is_contract_cache.lock().ok()?;
        let (value, at) = map.get(addr).copied()?;
        if at.elapsed() < IS_CONTRACT_TTL {
            Some(value)
        } else {
            None
        }
    }

    fn cache_put(&self, addr: Address, value: bool) {
        if let Ok(mut map) = self.is_contract_cache.lock() {
            map.insert(addr, (value, Instant::now()));
        }
    }

    #[cfg(test)]
    pub fn cache_seed(&self, addr: Address, value: bool) {
        self.cache_put(addr, value);
    }
}
