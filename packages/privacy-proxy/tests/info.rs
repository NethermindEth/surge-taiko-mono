use std::time::Duration;

use axum::body::{to_bytes, Body};
use axum::http::{Method, Request, StatusCode};
use axum::routing::post;
use axum::{Json, Router};
use privacy_proxy::config::Config;
use privacy_proxy::roles::reconcile_roles;
use privacy_proxy::{build_router, AppState};
use serde_json::{json, Value};
use sqlx::sqlite::SqlitePoolOptions;
use tokio::net::TcpListener;
use tower::ServiceExt;

const UPSTREAM_CHAIN_ID_HEX: &str = "0xba5ed";
const UPSTREAM_CHAIN_ID_DEC: u64 = 0xba5ed;

async fn mock_handler(Json(req): Json<Value>) -> Json<Value> {
    let id = req.get("id").cloned().unwrap_or(Value::Null);
    let method = req["method"].as_str().unwrap_or("");
    let result = match method {
        "eth_chainId" => Value::from(UPSTREAM_CHAIN_ID_HEX),
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

async fn build_app() -> Router {
    let upstream_url = spawn_mock_upstream().await;
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
        domain: "test-domain".to_string(),
    };
    let state = AppState::new(config, pool);
    build_router(state)
}

#[tokio::test]
async fn info_returns_expected_shape_unauthenticated() {
    let app = build_app().await;

    let response = app
        .oneshot(
            Request::builder()
                .uri("/info")
                .body(Body::empty())
                .unwrap(),
        )
        .await
        .unwrap();

    assert_eq!(response.status(), StatusCode::OK);
    let body = to_bytes(response.into_body(), usize::MAX).await.unwrap();
    let v: Value = serde_json::from_slice(&body).unwrap();

    assert_eq!(v["name"], "privacy-proxy");
    assert_eq!(v["chain_id"], UPSTREAM_CHAIN_ID_DEC);
    assert_eq!(v["domain"], "test-domain");
    assert_eq!(v["auth"]["scheme"], "bearer");
    assert_eq!(v["auth"]["challenge_path"], "/auth/challenge");
    assert_eq!(v["auth"]["verify_path"], "/auth/verify");
    assert!(v["version"].is_string());
}

#[tokio::test]
async fn cors_preflight_for_info_succeeds() {
    let app = build_app().await;

    let response = app
        .oneshot(
            Request::builder()
                .method(Method::OPTIONS)
                .uri("/info")
                .header("origin", "https://wallet.example")
                .header("access-control-request-method", "GET")
                .body(Body::empty())
                .unwrap(),
        )
        .await
        .unwrap();

    // tower-http CorsLayer responds 200 OK with Access-Control-Allow-* headers.
    assert!(response.status().is_success());
    assert!(response
        .headers()
        .get("access-control-allow-origin")
        .is_some());
}
