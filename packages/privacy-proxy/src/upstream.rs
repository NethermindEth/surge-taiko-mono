use std::collections::HashMap;
use std::sync::Mutex;
use std::time::{Duration, Instant};

use alloy::primitives::Address;
use anyhow::Result;
use serde_json::Value;

const IS_CONTRACT_TTL: Duration = Duration::from_secs(60);

#[derive(Debug)]
pub struct UpstreamClient {
    url: String,
    client: reqwest::Client,
    is_contract_cache: Mutex<HashMap<Address, (bool, Instant)>>,
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
        }
    }

    /// Forward a raw JSON-RPC payload verbatim and return the response body.
    pub async fn forward(&self, body: &Value) -> Result<Value> {
        let resp = self
            .client
            .post(&self.url)
            .json(body)
            .send()
            .await?
            .error_for_status()?
            .json::<Value>()
            .await?;
        Ok(resp)
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
