//! Address-parameterized read methods that the proxy gates separately from
//! contract calls. Each method gets a synthetic 4-byte selector in the
//! reserved `0xff______` range so admin rules can be stored in the same
//! `access_rules` table without a schema change.

use alloy::primitives::Address;
use serde_json::Value;

#[derive(Debug, Clone, Copy)]
pub struct GatedMethod {
    pub name: &'static str,
    pub selector: [u8; 4],
}

pub const ETH_GET_BALANCE: GatedMethod = GatedMethod {
    name: "eth_getBalance",
    selector: [0xff, 0x01, 0x00, 0x01],
};
pub const ETH_GET_TX_COUNT: GatedMethod = GatedMethod {
    name: "eth_getTransactionCount",
    selector: [0xff, 0x01, 0x00, 0x02],
};
pub const ETH_GET_CODE: GatedMethod = GatedMethod {
    name: "eth_getCode",
    selector: [0xff, 0x01, 0x00, 0x03],
};
pub const ETH_GET_STORAGE_AT: GatedMethod = GatedMethod {
    name: "eth_getStorageAt",
    selector: [0xff, 0x01, 0x00, 0x04],
};
pub const ETH_GET_PROOF: GatedMethod = GatedMethod {
    name: "eth_getProof",
    selector: [0xff, 0x01, 0x00, 0x05],
};

pub const ALL: &[GatedMethod] = &[
    ETH_GET_BALANCE,
    ETH_GET_TX_COUNT,
    ETH_GET_CODE,
    ETH_GET_STORAGE_AT,
    ETH_GET_PROOF,
];

pub fn lookup_by_method(method: &str) -> Option<GatedMethod> {
    ALL.iter().copied().find(|m| m.name == method)
}

pub fn lookup_by_selector(selector: [u8; 4]) -> Option<GatedMethod> {
    ALL.iter().copied().find(|m| m.selector == selector)
}

/// Extract the target address (`params[0]`) for a gated method.
pub fn extract_target(params: &Value) -> Option<Address> {
    let s = params.get(0)?.as_str()?;
    s.parse().ok()
}

/// Build a synthetic `call_data` payload — selector followed by an ABI-style
/// 32-byte left-padded address. Lambdas can decode `data[4..36]` to recover
/// the target. For `eth_getStorageAt`, the 32-byte slot is appended; future
/// lambdas can read it from `data[36..68]`.
pub fn encode_call_data(
    method: GatedMethod,
    target: Address,
    params: &Value,
) -> Result<Vec<u8>, &'static str> {
    let mut out = Vec::with_capacity(68);
    out.extend_from_slice(&method.selector);
    let mut padded = [0u8; 32];
    padded[12..].copy_from_slice(target.as_slice());
    out.extend_from_slice(&padded);

    if method.selector == ETH_GET_STORAGE_AT.selector {
        let s = params
            .get(1)
            .and_then(|v| v.as_str())
            .ok_or("eth_getStorageAt: missing slot parameter")?;
        let bytes = hex::decode(s.trim_start_matches("0x"))
            .map_err(|_| "eth_getStorageAt: slot is not hex")?;
        if bytes.len() > 32 {
            return Err("eth_getStorageAt: slot exceeds 32 bytes");
        }
        let mut slot = [0u8; 32];
        slot[32 - bytes.len()..].copy_from_slice(&bytes);
        out.extend_from_slice(&slot);
    }
    Ok(out)
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::json;

    #[test]
    fn lookup_round_trip() {
        for m in ALL {
            assert_eq!(lookup_by_method(m.name).unwrap().selector, m.selector);
            assert_eq!(lookup_by_selector(m.selector).unwrap().name, m.name);
        }
        assert!(lookup_by_method("eth_blockNumber").is_none());
        assert!(lookup_by_selector([0u8; 4]).is_none());
    }

    #[test]
    fn extract_target_parses_address() {
        let p = json!(["0x1111111111111111111111111111111111111111", "latest"]);
        let a = extract_target(&p).unwrap();
        assert_eq!(
            format!("0x{}", hex::encode(a.as_slice())),
            "0x1111111111111111111111111111111111111111"
        );
        assert!(extract_target(&json!([])).is_none());
        assert!(extract_target(&json!(["not_hex"])).is_none());
    }

    #[test]
    fn encode_call_data_layout() {
        let p = json!(["0x1111111111111111111111111111111111111111", "latest"]);
        let target = extract_target(&p).unwrap();
        let data = encode_call_data(ETH_GET_BALANCE, target, &p).unwrap();
        assert_eq!(&data[0..4], &ETH_GET_BALANCE.selector);
        assert_eq!(&data[4..16], &[0u8; 12]);
        assert_eq!(
            &data[16..36],
            &hex::decode("1111111111111111111111111111111111111111").unwrap()[..]
        );
        assert_eq!(data.len(), 36);
    }

    #[test]
    fn storage_at_appends_slot() {
        let p = json!([
            "0x1111111111111111111111111111111111111111",
            "0x000000000000000000000000000000000000000000000000000000000000002a",
            "latest"
        ]);
        let target = extract_target(&p).unwrap();
        let data = encode_call_data(ETH_GET_STORAGE_AT, target, &p).unwrap();
        assert_eq!(data.len(), 68);
        assert_eq!(data[67], 0x2a);
    }

    #[test]
    fn storage_at_rejects_oversized_slot() {
        let p = json!([
            "0x1111111111111111111111111111111111111111",
            // 33 bytes of hex
            "0x000000000000000000000000000000000000000000000000000000000000002a01",
            "latest"
        ]);
        let target = extract_target(&p).unwrap();
        let err = encode_call_data(ETH_GET_STORAGE_AT, target, &p).unwrap_err();
        assert!(err.contains("slot exceeds 32 bytes"));
    }

    #[test]
    fn storage_at_rejects_non_hex_slot() {
        let p = json!([
            "0x1111111111111111111111111111111111111111",
            "not-hex",
            "latest"
        ]);
        let target = extract_target(&p).unwrap();
        assert!(encode_call_data(ETH_GET_STORAGE_AT, target, &p).is_err());
    }
}
