use axum::extract::Request;
use axum::middleware::Next;
use axum::response::Response;
use axum::routing::{get, post};
use axum::{Extension, Router};
use tower_http::trace::TraceLayer;

use crate::admin;
use crate::auth::{caller_ctx_layer, challenge, verify};
use crate::rpc;
use crate::state::AppState;

pub fn build_router(state: AppState) -> Router {
    let admin_router = admin::router();
    let auth_router = Router::new()
        .route("/auth/challenge", get(challenge::handler))
        .route("/auth/verify", post(verify::handler));

    // Wrap caller_ctx_layer in a closure with explicit extractor types so
    // axum's FromFn can resolve the extractor tuple at .layer() time.
    let ctx_mw =
        |ext: Extension<AppState>, req: Request, next: Next| -> std::pin::Pin<
            Box<dyn std::future::Future<Output = Response> + Send>,
        > { Box::pin(caller_ctx_layer(ext, req, next)) };

    Router::new()
        .route("/", post(rpc::dispatch))
        .merge(auth_router)
        .merge(admin_router)
        .layer(axum::middleware::from_fn(ctx_mw))
        .layer(Extension(state.clone()))
        .layer(TraceLayer::new_for_http())
        .with_state(state)
}
