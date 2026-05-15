#[tokio::main]
async fn main() -> anyhow::Result<()> {
    privacy_proxy::run().await
}
