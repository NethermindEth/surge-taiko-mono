pub mod acl;
pub mod admin;
pub mod auth;
pub mod config;
pub mod db;
pub mod error;
pub mod roles;
pub mod rpc;
pub mod server;
pub mod state;
pub mod tracer;
pub mod upstream;

pub use config::Config;
pub use server::build_router;
pub use state::AppState;

pub async fn run() -> anyhow::Result<()> {
    tracing_subscriber::fmt()
        .with_env_filter(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "info,privacy_proxy=debug".into()),
        )
        .init();

    let cfg = Config::from_env()?;
    let pool = db::init_pool(&cfg.db_url).await?;
    roles::reconcile_roles(&pool).await?;
    admin::reconcile_seed_admins(&pool, &cfg.admin_eoas).await?;

    let bind = cfg.bind_addr.clone();
    let state = AppState::new(cfg, pool);
    let app = build_router(state);

    let listener = tokio::net::TcpListener::bind(&bind).await?;
    tracing::info!("privacy-proxy listening on {bind}");
    axum::serve(listener, app).await?;
    Ok(())
}
