use std::sync::Arc;

use crate::config::Config;
use crate::db::Pool;
use crate::upstream::UpstreamClient;

#[derive(Clone)]
pub struct AppState {
    pub config: Arc<Config>,
    pub pool: Pool,
    pub upstream: Arc<UpstreamClient>,
}

impl AppState {
    pub fn new(config: Config, pool: Pool) -> Self {
        let upstream = Arc::new(UpstreamClient::new(config.upstream_url.clone()));
        Self {
            config: Arc::new(config),
            pool,
            upstream,
        }
    }
}
