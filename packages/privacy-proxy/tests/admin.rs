use std::time::Duration;

use alloy::primitives::Address;
use axum::body::{to_bytes, Body};
use axum::http::{Request, StatusCode};
use privacy_proxy::config::Config;
use privacy_proxy::roles::reconcile_roles;
use privacy_proxy::{admin, build_router, db, AppState};
use serde_json::json;
use sha2::{Digest, Sha256};
use sqlx::sqlite::SqlitePoolOptions;
use sqlx::Row;
use tower::ServiceExt;

async fn build_test_app() -> (AppState, axum::Router) {
    let pool = SqlitePoolOptions::new()
        .max_connections(1)
        .connect("sqlite::memory:")
        .await
        .unwrap();
    sqlx::migrate!("./migrations").run(&pool).await.unwrap();
    reconcile_roles(&pool).await.unwrap();

    let config = Config {
        bind_addr: "127.0.0.1:0".to_string(),
        upstream_url: "http://127.0.0.1:1".to_string(),
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

fn random_token() -> String {
    use rand::RngCore;
    let mut buf = [0u8; 32];
    rand::thread_rng().fill_bytes(&mut buf);
    hex::encode(buf)
}

async fn insert_member_with_token(state: &AppState, role_name: &str, eoa: Address) -> String {
    let addr_hex = format!("0x{}", hex::encode(eoa.as_slice()));
    let now = db::now_unix();
    sqlx::query(
        "INSERT INTO members (eoa_address, role_id, created_at)
         VALUES (?, (SELECT id FROM roles WHERE name = ?), ?)
         ON CONFLICT(eoa_address) DO UPDATE SET role_id = excluded.role_id",
    )
    .bind(&addr_hex)
    .bind(role_name)
    .bind(now)
    .execute(&state.pool)
    .await
    .unwrap();
    if role_name == "user" {
        sqlx::query(
            "INSERT OR IGNORE INTO user_attributes (eoa_address, kyc, blacklisted) VALUES (?, 0, 0)",
        )
        .bind(&addr_hex)
        .execute(&state.pool)
        .await
        .unwrap();
    }

    let token = random_token();
    let token_hash = hex::encode(Sha256::digest(token.as_bytes()));
    let expires = now + 3600;
    sqlx::query(
        "INSERT INTO auth_tokens (token_hash, eoa_address, issued_at, expires_at)
         VALUES (?, ?, ?, ?)",
    )
    .bind(&token_hash)
    .bind(&addr_hex)
    .bind(now)
    .bind(expires)
    .execute(&state.pool)
    .await
    .unwrap();
    token
}

fn make_req(method: &str, path: &str, token: Option<&str>) -> Request<Body> {
    let mut b = Request::builder().method(method).uri(path);
    if let Some(t) = token {
        b = b.header("authorization", format!("Bearer {t}"));
    }
    b.body(Body::empty()).unwrap()
}

#[tokio::test]
async fn admin_endpoint_without_token_returns_401() {
    let (_state, app) = build_test_app().await;
    let res = app
        .clone()
        .oneshot(make_req("GET", "/admin/roles", None))
        .await
        .unwrap();
    assert_eq!(res.status(), StatusCode::UNAUTHORIZED);
}

#[tokio::test]
async fn admin_endpoint_with_user_token_returns_403() {
    let (state, app) = build_test_app().await;
    let eoa: Address = "0x1111111111111111111111111111111111111111".parse().unwrap();
    let token = insert_member_with_token(&state, "user", eoa).await;
    let res = app
        .clone()
        .oneshot(make_req("GET", "/admin/roles", Some(&token)))
        .await
        .unwrap();
    assert_eq!(res.status(), StatusCode::FORBIDDEN);
}

#[tokio::test]
async fn admin_endpoint_with_admin_token_returns_200() {
    let (state, app) = build_test_app().await;
    let eoa: Address = "0x2222222222222222222222222222222222222222".parse().unwrap();
    let token = insert_member_with_token(&state, "admin", eoa).await;
    let res = app
        .clone()
        .oneshot(make_req("GET", "/admin/roles", Some(&token)))
        .await
        .unwrap();
    assert_eq!(res.status(), StatusCode::OK);
    let body = to_bytes(res.into_body(), 4096).await.unwrap();
    let v: serde_json::Value = serde_json::from_slice(&body).unwrap();
    let names: Vec<&str> = v
        .as_array()
        .unwrap()
        .iter()
        .map(|r| r["name"].as_str().unwrap())
        .collect();
    assert!(names.contains(&"admin"));
    assert!(names.contains(&"user"));
}

#[tokio::test]
async fn lambdas_endpoint_returns_empty_groups_on_fresh_db() {
    let (state, app) = build_test_app().await;
    let eoa: Address = "0x3333333333333333333333333333333333333333".parse().unwrap();
    let token = insert_member_with_token(&state, "admin", eoa).await;
    let res = app
        .clone()
        .oneshot(make_req("GET", "/admin/registry/lambdas", Some(&token)))
        .await
        .unwrap();
    assert_eq!(res.status(), StatusCode::OK);
    let body = to_bytes(res.into_body(), 4096).await.unwrap();
    let v: serde_json::Value = serde_json::from_slice(&body).unwrap();
    let groups = v.as_array().unwrap();
    for g in groups {
        assert!(g["lambdas"].as_array().unwrap().is_empty());
    }
}

#[tokio::test]
async fn create_lambda_then_list_returns_it() {
    let (state, app) = build_test_app().await;
    let eoa: Address = "0xb1b1b1b1b1b1b1b1b1b1b1b1b1b1b1b1b1b1b1b1".parse().unwrap();
    let token = insert_member_with_token(&state, "admin", eoa).await;

    let one_word = format!("0x{}01", "00".repeat(31));
    let body = json!({
        "name": "require_kyc",
        "role": "user",
        "description": "kyc must be true",
        "rules": [{
            "selector": "0x70a08231",
            "lhs_kind": "attribute",
            "lhs_attribute": "kyc",
            "condition": "eq",
            "rhs_kind": "literal",
            "rhs_value": one_word
        }]
    });
    let req = Request::builder()
        .method("POST")
        .uri("/admin/registry/lambdas")
        .header("authorization", format!("Bearer {token}"))
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::CREATED);

    let res = app
        .clone()
        .oneshot(make_req("GET", "/admin/registry/lambdas", Some(&token)))
        .await
        .unwrap();
    let v: serde_json::Value =
        serde_json::from_slice(&to_bytes(res.into_body(), 4096).await.unwrap()).unwrap();
    let user_group = v.as_array().unwrap().iter().find(|g| g["role"] == "user").unwrap();
    let names: Vec<&str> = user_group["lambdas"]
        .as_array()
        .unwrap()
        .iter()
        .map(|r| r["name"].as_str().unwrap())
        .collect();
    assert!(names.contains(&"require_kyc"));
}

#[tokio::test]
async fn create_lambda_with_tx_origin_and_msg_sender_rhs() {
    let (state, app) = build_test_app().await;
    let eoa: Address = "0xb3b3b3b3b3b3b3b3b3b3b3b3b3b3b3b3b3b3b3b3".parse().unwrap();
    let token = insert_member_with_token(&state, "admin", eoa).await;

    let body = json!({
        "name": "self_only",
        "role": "user",
        "rules": [
            {
                "selector": "0x70a08231",
                "lhs_kind": "calldata",
                "lhs_offset": 4,
                "condition": "eq",
                "rhs_kind": "tx_origin"
            },
            {
                "selector": "0xa9059cbb",
                "lhs_kind": "calldata",
                "lhs_offset": 4,
                "condition": "neq",
                "rhs_kind": "msg_sender"
            }
        ]
    });
    let req = Request::builder()
        .method("POST")
        .uri("/admin/registry/lambdas")
        .header("authorization", format!("Bearer {token}"))
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::CREATED);
    let v: serde_json::Value =
        serde_json::from_slice(&to_bytes(res.into_body(), 8192).await.unwrap()).unwrap();
    let rules = v["rules"].as_array().unwrap();
    assert_eq!(rules.len(), 2);
    let kinds: Vec<&str> = rules.iter().map(|r| r["rhs_kind"].as_str().unwrap()).collect();
    assert!(kinds.contains(&"tx_origin"));
    assert!(kinds.contains(&"msg_sender"));
}

#[tokio::test]
async fn member_upsert_writes_typed_user_attributes() {
    let (state, app) = build_test_app().await;
    let admin_eoa: Address = "0xa0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0a0".parse().unwrap();
    let admin_token = insert_member_with_token(&state, "admin", admin_eoa).await;

    let target = "0x1010101010101010101010101010101010101010";
    let body = json!({
        "role": "user",
        "attributes": { "kyc": true, "blacklisted": false }
    });
    let req = Request::builder()
        .method("PUT")
        .uri(format!("/admin/members/{target}"))
        .header("authorization", format!("Bearer {admin_token}"))
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::OK);
    let v: serde_json::Value =
        serde_json::from_slice(&to_bytes(res.into_body(), 4096).await.unwrap()).unwrap();
    assert_eq!(v["role"], "user");
    assert_eq!(v["attributes"]["kyc"], true);
    assert_eq!(v["attributes"]["blacklisted"], false);

    // Verify the row.
    let row = sqlx::query("SELECT kyc, blacklisted FROM user_attributes WHERE eoa_address = ?")
        .bind(target)
        .fetch_one(&state.pool)
        .await
        .unwrap();
    let kyc: i64 = row.get("kyc");
    let bl: i64 = row.get("blacklisted");
    assert_eq!(kyc, 1);
    assert_eq!(bl, 0);
}

#[tokio::test]
async fn member_upsert_admin_has_null_attributes() {
    let (state, app) = build_test_app().await;
    let admin_eoa: Address = "0xa1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1".parse().unwrap();
    let admin_token = insert_member_with_token(&state, "admin", admin_eoa).await;

    let target = "0x2020202020202020202020202020202020202020";
    let body = json!({ "role": "admin" });
    let req = Request::builder()
        .method("PUT")
        .uri(format!("/admin/members/{target}"))
        .header("authorization", format!("Bearer {admin_token}"))
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::OK);
    let v: serde_json::Value =
        serde_json::from_slice(&to_bytes(res.into_body(), 4096).await.unwrap()).unwrap();
    assert_eq!(v["role"], "admin");
    assert!(v["attributes"].is_null());

    // No row in user_attributes.
    let row_count: i64 =
        sqlx::query("SELECT COUNT(*) AS c FROM user_attributes WHERE eoa_address = ?")
            .bind(target)
            .fetch_one(&state.pool)
            .await
            .unwrap()
            .get("c");
    assert_eq!(row_count, 0);
}

#[tokio::test]
async fn member_upsert_admin_with_attributes_rejected() {
    let (state, app) = build_test_app().await;
    let admin_eoa: Address = "0xa2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2a2".parse().unwrap();
    let admin_token = insert_member_with_token(&state, "admin", admin_eoa).await;

    let target = "0x3030303030303030303030303030303030303030";
    let body = json!({ "role": "admin", "attributes": { "kyc": true, "blacklisted": false } });
    let req = Request::builder()
        .method("PUT")
        .uri(format!("/admin/members/{target}"))
        .header("authorization", format!("Bearer {admin_token}"))
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::BAD_REQUEST);
}

async fn create_user_lambda(app: &axum::Router, token: &str, name: &str) -> i64 {
    let one_word = format!("0x{}01", "00".repeat(31));
    let body = json!({
        "name": name,
        "role": "user",
        "rules": [{
            "selector": "0x70a08231",
            "lhs_kind": "attribute",
            "lhs_attribute": "kyc",
            "condition": "eq",
            "rhs_kind": "literal",
            "rhs_value": one_word
        }]
    });
    let req = Request::builder()
        .method("POST")
        .uri("/admin/registry/lambdas")
        .header("authorization", format!("Bearer {token}"))
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::CREATED);
    let v: serde_json::Value =
        serde_json::from_slice(&to_bytes(res.into_body(), 8192).await.unwrap()).unwrap();
    v["id"].as_i64().unwrap()
}

#[tokio::test]
async fn rule_with_role_mismatched_lambda_rejected() {
    let (state, app) = build_test_app().await;
    let admin_eoa: Address = "0xa3a3a3a3a3a3a3a3a3a3a3a3a3a3a3a3a3a3a3a3".parse().unwrap();
    let admin_token = insert_member_with_token(&state, "admin", admin_eoa).await;
    let user_lambda_id = create_user_lambda(&app, &admin_token, "kyc_check").await;

    let body = json!({
        "name": "transfer-admin-with-user-lambda",
        "selector": "0xa9059cbb",
        "mode": "deny",
        "entries": [ { "role": "admin", "lambda_id": user_lambda_id } ]
    });
    let req = Request::builder()
        .method("POST")
        .uri("/admin/registry/rules")
        .header("authorization", format!("Bearer {admin_token}"))
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::BAD_REQUEST);
}

#[tokio::test]
async fn rule_with_user_lambda_accepted() {
    let (state, app) = build_test_app().await;
    let admin_eoa: Address = "0xa4a4a4a4a4a4a4a4a4a4a4a4a4a4a4a4a4a4a4a4".parse().unwrap();
    let admin_token = insert_member_with_token(&state, "admin", admin_eoa).await;
    let lambda_id = create_user_lambda(&app, &admin_token, "kyc_v2").await;

    let body = json!({
        "name": "balanceof-user-allowlist",
        "selector": "0x70a08231",
        "mode": "allow",
        "entries": [ { "role": "user", "lambda_id": lambda_id } ]
    });
    let req = Request::builder()
        .method("POST")
        .uri("/admin/registry/rules")
        .header("authorization", format!("Bearer {admin_token}"))
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::CREATED);
}

#[tokio::test]
async fn delete_lambda_blocked_when_referenced() {
    let (state, app) = build_test_app().await;
    let admin_eoa: Address = "0xb2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2".parse().unwrap();
    let token = insert_member_with_token(&state, "admin", admin_eoa).await;
    let lambda_id = create_user_lambda(&app, &token, "kyc_locked").await;

    let body = json!({
        "name": "balanceof-locked-rule",
        "selector": "0x70a08231",
        "mode": "allow",
        "entries": [ { "role": "user", "lambda_id": lambda_id } ]
    });
    let req = Request::builder()
        .method("POST")
        .uri("/admin/registry/rules")
        .header("authorization", format!("Bearer {token}"))
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::CREATED);

    let res = app
        .clone()
        .oneshot(make_req(
            "DELETE",
            &format!("/admin/registry/lambdas/{lambda_id}"),
            Some(&token),
        ))
        .await
        .unwrap();
    assert_eq!(res.status(), StatusCode::CONFLICT);
}

#[tokio::test]
async fn reconcile_promotes_seed_eoa_to_admin() {
    let pool = SqlitePoolOptions::new()
        .max_connections(1)
        .connect("sqlite::memory:")
        .await
        .unwrap();
    sqlx::migrate!("./migrations").run(&pool).await.unwrap();
    reconcile_roles(&pool).await.unwrap();

    let seed: Address = "0x4444444444444444444444444444444444444444".parse().unwrap();
    let addr_hex = format!("0x{}", hex::encode(seed.as_slice()));

    // Pre-seed the member at role 'user' with an attribute row.
    let now = db::now_unix();
    sqlx::query(
        "INSERT INTO members (eoa_address, role_id, created_at)
         VALUES (?, (SELECT id FROM roles WHERE name = 'user'), ?)",
    )
    .bind(&addr_hex)
    .bind(now)
    .execute(&pool)
    .await
    .unwrap();
    sqlx::query(
        "INSERT INTO user_attributes (eoa_address, kyc, blacklisted) VALUES (?, 1, 0)",
    )
    .bind(&addr_hex)
    .execute(&pool)
    .await
    .unwrap();

    admin::reconcile_seed_admins(&pool, &[seed]).await.unwrap();

    // Role flipped to admin.
    let role: String = sqlx::query(
        "SELECT r.name FROM members m JOIN roles r ON r.id = m.role_id WHERE m.eoa_address = ?",
    )
    .bind(&addr_hex)
    .fetch_one(&pool)
    .await
    .unwrap()
    .get("name");
    assert_eq!(role, "admin");

    // Attribute row dropped.
    let row_count: i64 =
        sqlx::query("SELECT COUNT(*) AS c FROM user_attributes WHERE eoa_address = ?")
            .bind(&addr_hex)
            .fetch_one(&pool)
            .await
            .unwrap()
            .get("c");
    assert_eq!(row_count, 0);
}

#[tokio::test]
async fn reconcile_creates_missing_seed_eoa() {
    let pool = SqlitePoolOptions::new()
        .max_connections(1)
        .connect("sqlite::memory:")
        .await
        .unwrap();
    sqlx::migrate!("./migrations").run(&pool).await.unwrap();
    reconcile_roles(&pool).await.unwrap();

    let seed: Address = "0x5555555555555555555555555555555555555555".parse().unwrap();
    let addr_hex = format!("0x{}", hex::encode(seed.as_slice()));
    admin::reconcile_seed_admins(&pool, &[seed]).await.unwrap();

    let role: String = sqlx::query(
        "SELECT r.name FROM members m JOIN roles r ON r.id = m.role_id WHERE m.eoa_address = ?",
    )
    .bind(&addr_hex)
    .fetch_one(&pool)
    .await
    .unwrap()
    .get("name");
    assert_eq!(role, "admin");
}

#[tokio::test]
async fn non_eth_namespace_method_rejected() {
    let (_state, app) = build_test_app().await;
    let body = json!({
        "jsonrpc": "2.0", "id": 1, "method": "net_version", "params": []
    });
    let req = Request::builder()
        .method("POST")
        .uri("/")
        .header("content-type", "application/json")
        .body(Body::from(serde_json::to_vec(&body).unwrap()))
        .unwrap();
    let res = app.clone().oneshot(req).await.unwrap();
    assert_eq!(res.status(), StatusCode::OK);
    let bytes = to_bytes(res.into_body(), 4096).await.unwrap();
    let v: serde_json::Value = serde_json::from_slice(&bytes).unwrap();
    assert_eq!(v["error"]["code"], -32601);
}
