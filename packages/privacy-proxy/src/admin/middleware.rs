use axum::extract::Request;
use axum::middleware::Next;
use axum::response::Response;

use crate::auth::CallerCtx;
use crate::error::ApiError;

/// Gate every `/admin/*` route. Requires the `caller_ctx_layer` to have
/// run before this on the same request so the `CallerCtx` extension is
/// present.
pub async fn admin_gate(req: Request, next: Next) -> Result<Response, ApiError> {
    let ctx = req
        .extensions()
        .get::<CallerCtx>()
        .cloned()
        .unwrap_or_else(CallerCtx::anonymous);
    if ctx.is_admin() {
        return Ok(next.run(req).await);
    }
    if ctx.is_anonymous() {
        Err(ApiError::unauthorized(
            "admin endpoints require an auth token",
        ))
    } else {
        Err(ApiError::forbidden("admin role required"))
    }
}
