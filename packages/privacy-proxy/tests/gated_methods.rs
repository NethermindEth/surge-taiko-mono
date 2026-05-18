use std::time::Duration;

use alloy::primitives::Address;
use axum::body::{to_bytes, Body};
use axum::http::{Request, StatusCode};
use axum::routing::post;
use axum::{Json, Router};
use privacy_proxy::config::Config;
use privacy_proxy::roles::reconcile_roles;
use privacy_proxy::{build_router, db, AppState};
use serde_json::{json, Value};
use sha2::{Digest, Sha256};
use sqlx::sqlite::SqlitePoolOptions;
use tokio::net::TcpListener;
use tower::ServiceExt;

const USER_EOA: &str = "0x1111111111111111111111111111111111111111";
const OTHER_EOA: &str = "0x2222222222222222222222222222222222222222";
const CONTRACT_ADDR: &str = "0xcccccccccccccccccccccccccccccccccccccccc";

async fn mock_handler(Json(req): Json<Value>) -> Json<Value> {
    let id = req.get("id").cloned().unwrap_or(Value::Null);
    let method = req["method"].as_str().unwrap_or("");
    let result = match method {
        "eth_getCode" => {
            let target = req["params"][0].as_str().unwrap_or("").to_lowercase();
            if target == CONTRACT_ADDR {
                Value::from("0xdeadbeef")
            } else {
                Value::from("0x")
            }
        }
        "eth_getBalance" => Value::from("0xabcd"),
        _ => Value::Null,
    };
    Json(json!({ "jsonrpc": "2.0", "id": id, "result": result }))
}

async fn spawn_mock_upstream() -> String {
    let listener = TcpListener::bind("127.0.0.1:0").await.unwrap();
    let addr = listener.local_addr().unwrap();
    let app = Router::new().route("/", post(mock_handler));
    tokio::spawn(async move {
        axum::serve(listener, app).await.unwrap();
    });
    format!("http://{addr}")
}

async fn build_app_with_upstream(upstream_url: String) -> (AppState, axum::Router) {
    let pool = SqlitePoolOptions::new()
        .max_connections(1)
        .connect("sqlite::memory:")
        .await
        .unwrap();
    sqlx::migrate!("./migrations").run(&pool).await.unwrap();
    reconcile_roles(&pool).await.unwrap();

    let config = Config {
        bind_addr: "127.0.0.1:0".to_string(),
        upstream_url,
        db_url: "sqlite::memory:".to_string(),
        admin_eoas: vec![],
        challenge_ttl: Duration::from_secs(300),
        token_ttl: Duration::from_secs(3600),
        domain: "test".to_string(),
    };
    let state = AppState::new(config, pool);
    let app = build_router(state.clone());
    (state, app)
}

async fn issue_token(state: &AppState, role: &str, eoa: Address) -> String {
    let addr_hex = format!("0x{}", hex::encode(eoa.as_slice()));
    let now = db::now_unix();
    sqlx::query(
        "INSERT INTO members (eoa_address, role_id, created_at)
         VALUES (?, (SELECT id FROM roles WHERE name = ?), ?)
         ON CONFLICT(eoa_address) DO UPDATE SET role_id = excluded.role_id",
    )
    .bind(&addr_hex)
    .bind(role)
    .bind(now)
    .execute(&state.pool)
    .await
    .unwrap();
    if role == "user" {
        sqlx::query(
            "INSERT OR IGNORE INTO user_attributes (eoa_address, kyc, blacklisted) VALUES (?, 0, 0)",
        )
        .bind(&addr_hex)
        .execute(&state.pool)
        .await
        .unwrap();
    }
    use rand::RngCore;
    let mut buf = [0u8; 32];
    rand::thread_rng().fill_bytes(&mut buf);
    let token = hex::encode(buf);
    let token_hash = hex::encode(Sha256::digest(token.as_bytes()));
    sqlx::query(
        "INSERT INTO auth_tokens (token_hash, eoa_address, issued_at, expires_at)
         VALUES (?, ?, ?, ?)",
    )
    .bind(&token_hash)
    .bind(&addr_hex)
    .bind(now)
    .bind(now + 3600)
    .execute(&state.pool)
    .await
    .unwrap();
    token
}

fn rpc_req(token: Option<&str>, method: &str, params: Value) -> Request<Body> {
    let body = json!({ "jsonrpc": "2.0", "id": 1, "method": method, "params": params });
    let mut b = Request::builder()
        .method("POST")
        .uri("/")
        .header("content-type", "application/json");
    if let Some(t) = token {
        b = b.header("authorization", format!("Bearer {t}"));
    }
    b.body(Body::from(serde_json::to_vec(&body).unwrap())).unwrap()
}

async fn body_json(res: axum::response::Response) -> Value {
    let bytes = to_bytes(res.into_body(), 8192).await.unwrap();
    serde_json::from_slice(&bytes).unwrap()
}

#[tokio::test]
async fn self_balance_forwarded() {
    let upstream = spawn_mock_upstream().await;
    let (state, app) = build_app_with_upstream(upstream).await;
    let user: Address = USER_EOA.parse().unwrap();
    let token = issue_token(&state, "user", user).await;

    let res = app
        .clone()
        .oneshot(rpc_req(Some(&token), "eth_getBalance", json!([USER_EOA, "latest"])))
        .await
        .unwrap();
    assert_eq!(res.status(), StatusCode::OK);
    let v = body_json(res).await;
    assert_eq!(v["result"], "0xabcd");
    assert!(v.get("error").is_none());
}

#[tokio::test]
async fn other_eoa_balance_denied() {
    let upstream = spawn_mock_upstream().await;
    let (state, app) = build_app_with_upstream(upstream).await;
    let user: Address = USER_EOA.parse().unwrap();
    let token = issue_token(&state, "user", user).await;

    let res = app
        .clone()
        .oneshot(rpc_req(
            Some(&token),
            "eth_getBalance",
            json!([OTHER_EOA, "latest"]),
        ))
        .await
        .unwrap();
    let v = body_json(res).await;
    assert_eq!(v["error"]["code"], -32001);
    assert_eq!(
        v["error"]["data"]["contract"].as_str().unwrap().to_lowercase(),
        OTHER_EOA
    );
    assert_eq!(v["error"]["data"]["detail"], "DefaultEoaSelfOnly");
}

#[tokio::test]
async fn contract_balance_forwarded_by_default() {
    let upstream = spawn_mock_upstream().await;
    let (state, app) = build_app_with_upstream(upstream).await;
    let user: Address = USER_EOA.parse().unwrap();
    let token = issue_token(&state, "user", user).await;

    let res = app
        .clone()
        .oneshot(rpc_req(
            Some(&token),
            "eth_getBalance",
            json!([CONTRACT_ADDR, "latest"]),
        ))
        .await
        .unwrap();
    let v = body_json(res).await;
    assert!(v.get("error").is_none(), "got error: {v}");
    assert_eq!(v["result"], "0xabcd");
}

#[tokio::test]
async fn anonymous_balance_denied() {
    let upstream = spawn_mock_upstream().await;
    let (_state, app) = build_app_with_upstream(upstream).await;

    let res = app
        .clone()
        .oneshot(rpc_req(None, "eth_getBalance", json!([CONTRACT_ADDR, "latest"])))
        .await
        .unwrap();
    let v = body_json(res).await;
    assert_eq!(v["error"]["code"], -32001);
    assert_eq!(v["error"]["data"]["detail"], "AnonymousAgainstGatedCall");
}

#[tokio::test]
async fn admin_deny_rule_on_contract_blocks_user() {
    let upstream = spawn_mock_upstream().await;
    let (state, app) = build_app_with_upstream(upstream).await;
    let user: Address = USER_EOA.parse().unwrap();
    let token = issue_token(&state, "user", user).await;

    sqlx::query(
        "INSERT INTO access_rules (contract_address, function_selector, mode) VALUES (?, ?, 'deny')",
    )
    .bind(CONTRACT_ADDR)
    .bind("0xff010001") // synthetic selector for eth_getBalance
    .execute(&state.pool)
    .await
    .unwrap();
    let rule_id: i64 = sqlx::query_scalar("SELECT last_insert_rowid()")
        .fetch_one(&state.pool)
        .await
        .unwrap();
    sqlx::query(
        "INSERT INTO access_rule_entries (rule_id, role_id, lambda_name)
         VALUES (?, (SELECT id FROM roles WHERE name = 'user'), NULL)",
    )
    .bind(rule_id)
    .execute(&state.pool)
    .await
    .unwrap();

    let res = app
        .clone()
        .oneshot(rpc_req(
            Some(&token),
            "eth_getBalance",
            json!([CONTRACT_ADDR, "latest"]),
        ))
        .await
        .unwrap();
    let v = body_json(res).await;
    assert_eq!(v["error"]["code"], -32001);
    assert_eq!(v["error"]["data"]["detail"], "InDenyList");
}

#[tokio::test]
async fn admin_can_create_rule_with_method_name_selector() {
    let upstream = spawn_mock_upstream().await;
    let (state, app) = build_app_with_upstream(upstream).await;
    let admin: Address = "0xa1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1".parse().unwrap();
    let token = issue_token(&state, "admin", admin).await;

    let body = json!({
        "contract_address": CONTRACT_ADDR,
        "function_selector": "eth_getBalance",
        "mode": "deny",
        "entries": [ { "role": "user" } ]
    });
    let req = Request::builder()
        .method("POST")
        .uri("/admin/registry/rules")
        .header("content-type", "application/json")
        .header("authorization", format!("Bearer {token}"))
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::CREATED);
    let v = body_json(res).await;
    assert_eq!(v["function_selector"], "0xff010001");
    assert_eq!(v["contract_address"], CONTRACT_ADDR);
    assert_eq!(v["mode"], "deny");
}

#[tokio::test]
async fn anonymous_eth_call_denied() {
    let upstream = spawn_mock_upstream().await;
    let (_state, app) = build_app_with_upstream(upstream).await;

    let res = app
        .clone()
        .oneshot(rpc_req(
            None,
            "eth_call",
            json!([{ "to": CONTRACT_ADDR, "data": "0x95d89b41" }, "latest"]),
        ))
        .await
        .unwrap();
    let v = body_json(res).await;
    assert_eq!(v["error"]["code"], -32001);
    assert_eq!(v["error"]["data"]["detail"], "AnonymousAgainstGatedCall");
}

#[tokio::test]
async fn anonymous_eth_chain_id_allowed() {
    let upstream = spawn_mock_upstream().await;
    let (_state, app) = build_app_with_upstream(upstream).await;

    let res = app
        .clone()
        .oneshot(rpc_req(None, "eth_chainId", json!([])))
        .await
        .unwrap();
    let v = body_json(res).await;
    assert!(v.get("error").is_none(), "got error: {v}");
}

#[tokio::test]
async fn anonymous_eth_block_number_allowed() {
    let upstream = spawn_mock_upstream().await;
    let (_state, app) = build_app_with_upstream(upstream).await;

    let res = app
        .clone()
        .oneshot(rpc_req(None, "eth_blockNumber", json!([])))
        .await
        .unwrap();
    let v = body_json(res).await;
    assert!(v.get("error").is_none(), "got error: {v}");
}

#[tokio::test]
async fn anonymous_eth_get_logs_denied() {
    let upstream = spawn_mock_upstream().await;
    let (_state, app) = build_app_with_upstream(upstream).await;

    let res = app
        .clone()
        .oneshot(rpc_req(None, "eth_getLogs", json!([{ "fromBlock": "latest" }])))
        .await
        .unwrap();
    let v = body_json(res).await;
    assert_eq!(v["error"]["code"], -32001);
    assert_eq!(v["error"]["data"]["detail"], "AnonymousAgainstGatedCall");
}

#[tokio::test]
async fn synthetic_selectors_endpoint_lists_methods() {
    let upstream = spawn_mock_upstream().await;
    let (state, app) = build_app_with_upstream(upstream).await;
    let admin: Address = "0xa2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2".parse().unwrap();
    let token = issue_token(&state, "admin", admin).await;

    let req = Request::builder()
        .method("GET")
        .uri("/admin/registry/synthetic-selectors")
        .header("authorization", format!("Bearer {token}"))
        .body(Body::empty())
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::OK);
    let v = body_json(res).await;
    let methods: Vec<&str> = v
        .as_array()
        .unwrap()
        .iter()
        .map(|r| r["method"].as_str().unwrap())
        .collect();
    assert!(methods.contains(&"eth_getBalance"));
    assert!(methods.contains(&"eth_getProof"));
    assert!(methods.contains(&"eth_getStorageAt"));
    assert!(methods.contains(&"eth_getCode"));
    assert!(methods.contains(&"eth_getTransactionCount"));
}
