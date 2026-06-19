mod auth;
mod cache;
mod config;
mod db;
mod handlers;
mod models;
mod utils;

use axum::routing::get;
use axum::Router;
use handlers::AppState;

#[tokio::main]
async fn main() {
    let config = config::Config::from_env();
    let pool = db::create_pool(&config).await.expect("Failed to create database pool");
    db::migrate(&pool).await.expect("Failed to run migrations");

    let cache = cache::UrlCache::new();

    let urls: Vec<models::Url> = sqlx::query_as("SELECT id, original, accesses FROM url")
        .fetch_all(&pool)
        .await
        .expect("Failed to fetch URLs");
    cache.load(urls).await;

    let state = AppState {
        cache,
        config: config.clone(),
        pool,
    };

    let prefix = config.http_base_trimmed();

    let app = Router::new()
        .route(&format!("{prefix}/urls"), get(handlers::get_urls).post(handlers::create_url))
        .route(
            &format!("{prefix}/urls/{{id}}"),
            get(handlers::get_url_info),
        )
        .route(
            &format!("{prefix}/{{id}}"),
            get(handlers::redirect_url).delete(handlers::delete_url),
        )
        .with_state(state);

    let addr = format!("0.0.0.0:{}", config.http_port);
    println!("Starting ENCURTADOR at port {}", config.http_port);

    let listener = tokio::net::TcpListener::bind(&addr)
        .await
        .expect("Failed to bind address");
    axum::serve(listener, app)
        .await
        .expect("Server error");
}
