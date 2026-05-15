use std::time::Duration;

use alloy::primitives::Address;
use axum::body::{to_bytes, Body};
use axum::http::{Request, StatusCode};
use privacy_proxy::config::Config;
use privacy_proxy::{admin, build_router, db, AppState};
use sha2::{Digest, Sha256};
use sqlx::sqlite::SqlitePoolOptions;
use tower::ServiceExt;

async fn build_test_app() -> (AppState, axum::Router) {
    let pool = SqlitePoolOptions::new()
        .max_connections(1)
        .connect("sqlite::memory:")
        .await
        .unwrap();
    sqlx::migrate!("./migrations").run(&pool).await.unwrap();

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

async fn insert_user_with_token(state: &AppState, role_name: &str, eoa: Address) -> String {
    let addr_hex = format!("0x{}", hex::encode(eoa.as_slice()));
    let now = db::now_unix();
    sqlx::query(
        "INSERT INTO users (eoa_address, role_id, caller_info_json, created_at)
         VALUES (?, (SELECT id FROM roles WHERE name = ?), '{}', ?)
         ON CONFLICT(eoa_address) DO UPDATE SET role_id = excluded.role_id",
    )
    .bind(&addr_hex)
    .bind(role_name)
    .bind(now)
    .execute(&state.pool)
    .await
    .unwrap();

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
    let token = insert_user_with_token(&state, "user", eoa).await;
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
    let token = insert_user_with_token(&state, "admin", eoa).await;
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
async fn lambdas_endpoint_lists_in_build_lambdas() {
    let (state, app) = build_test_app().await;
    let eoa: Address = "0x3333333333333333333333333333333333333333".parse().unwrap();
    let token = insert_user_with_token(&state, "admin", eoa).await;
    let res = app
        .clone()
        .oneshot(make_req("GET", "/admin/registry/lambdas", Some(&token)))
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
    assert!(names.contains(&"require_kyc"));
    assert!(names.contains(&"transfer_under_limit"));
}

#[tokio::test]
async fn reconcile_promotes_seed_eoa_to_admin() {
    let pool = SqlitePoolOptions::new()
        .max_connections(1)
        .connect("sqlite::memory:")
        .await
        .unwrap();
    sqlx::migrate!("./migrations").run(&pool).await.unwrap();

    let seed: Address = "0x4444444444444444444444444444444444444444".parse().unwrap();
    let addr_hex = format!("0x{}", hex::encode(seed.as_slice()));

    // Pre-seed the user at role 'user'.
    let now = db::now_unix();
    sqlx::query(
        "INSERT INTO users (eoa_address, role_id, caller_info_json, created_at)
         VALUES (?, (SELECT id FROM roles WHERE name = 'user'), '{}', ?)",
    )
    .bind(&addr_hex)
    .bind(now)
    .execute(&pool)
    .await
    .unwrap();

    // Run reconciliation with the EOA in ADMIN_EOAS.
    admin::reconcile_seed_admins(&pool, &[seed]).await.unwrap();

    let row = sqlx::query(
        "SELECT r.name FROM users u JOIN roles r ON r.id = u.role_id WHERE u.eoa_address = ?",
    )
    .bind(&addr_hex)
    .fetch_one(&pool)
    .await
    .unwrap();
    use sqlx::Row;
    let role: String = row.get("name");
    assert_eq!(role, "admin");
}

#[tokio::test]
async fn reconcile_creates_missing_seed_eoa() {
    let pool = SqlitePoolOptions::new()
        .max_connections(1)
        .connect("sqlite::memory:")
        .await
        .unwrap();
    sqlx::migrate!("./migrations").run(&pool).await.unwrap();

    let seed: Address = "0x5555555555555555555555555555555555555555".parse().unwrap();
    let addr_hex = format!("0x{}", hex::encode(seed.as_slice()));
    admin::reconcile_seed_admins(&pool, &[seed]).await.unwrap();

    let row = sqlx::query(
        "SELECT r.name FROM users u JOIN roles r ON r.id = u.role_id WHERE u.eoa_address = ?",
    )
    .bind(&addr_hex)
    .fetch_one(&pool)
    .await
    .unwrap();
    use sqlx::Row;
    let role: String = row.get("name");
    assert_eq!(role, "admin");
}

#[tokio::test]
async fn non_eth_namespace_method_rejected() {
    let (_state, app) = build_test_app().await;
    let body = serde_json::json!({
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
