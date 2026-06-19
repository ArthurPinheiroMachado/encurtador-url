use axum::{
    extract::{Path, State},
    http::header::AUTHORIZATION,
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use serde_json::json;

use crate::auth::{is_authorized, unauthorized_response};
use crate::cache::UrlCache;
use crate::config::Config;
use crate::models::{CreateUrlPayload, Url, UrlInfo};
use crate::utils;

#[derive(Clone)]
pub struct AppState {
    pub cache: UrlCache,
    pub config: Config,
    pub pool: sqlx::PgPool,
}

fn check_auth(headers: &axum::http::HeaderMap, config: &Config) -> Result<(), Response> {
    let auth_header = headers.get(AUTHORIZATION).and_then(|v| v.to_str().ok());
    if !is_authorized(auth_header, config) {
        return Err(unauthorized_response().into_response());
    }
    Ok(())
}

pub async fn get_urls(
    State(state): State<AppState>,
    headers: axum::http::HeaderMap,
) -> Response {
    if let Err(resp) = check_auth(&headers, &state.config) {
        return resp;
    }

    Json(state.cache.get_all().await).into_response()
}

pub async fn create_url(
    State(state): State<AppState>,
    headers: axum::http::HeaderMap,
    Json(payload): Json<CreateUrlPayload>,
) -> Response {
    if let Err(resp) = check_auth(&headers, &state.config) {
        return resp;
    }

    if url::Url::parse(&payload.url).is_err() {
        return (StatusCode::BAD_REQUEST, Json(json!({"detail": "Invalid URL"}))).into_response();
    }

    match sqlx::query_as::<_, Url>("SELECT id, original, accesses FROM url WHERE original = $1")
        .bind(&payload.url)
        .fetch_optional(&state.pool)
        .await
    {
        Ok(Some(existing)) => {
            return (
                StatusCode::OK,
                Json(json!({"id": existing.id, "url": payload.url})),
            )
                .into_response();
        }
        Err(_) => {
            return (
                StatusCode::INTERNAL_SERVER_ERROR,
                Json(json!({"detail": "Database error"})),
            )
                .into_response();
        }
        Ok(None) => {}
    }

    let short_id = {
        let mut attempts = 0;
        loop {
            let id = utils::generate_short_id(8);
            if !state.cache.exists(&id).await {
                break id;
            }
            attempts += 1;
            if attempts >= 100 {
                return (StatusCode::INTERNAL_SERVER_ERROR,
                    Json(json!({"detail": "Failed to generate unique ID"}))).into_response();
            }
        }
    };

    match sqlx::query("INSERT INTO url(id, original, accesses) VALUES($1, $2, 0)")
        .bind(&short_id)
        .bind(&payload.url)
        .execute(&state.pool)
        .await
    {
        Ok(_) => {
            state
                .cache
                .set(
                    short_id.clone(),
                    UrlInfo {
                        original: payload.url.clone(),
                        accesses: 0,
                    },
                )
                .await;

            (StatusCode::CREATED, Json(json!({"id": short_id, "url": payload.url}))).into_response()
        }
        Err(_) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(json!({"detail": "Database error"})),
        )
            .into_response(),
    }
}

pub async fn get_url_info(
    State(state): State<AppState>,
    headers: axum::http::HeaderMap,
    Path(id): Path<String>,
) -> Response {
    if let Err(resp) = check_auth(&headers, &state.config) {
        return resp;
    }

    match state.cache.get(&id).await {
        Some(info) => Json(info).into_response(),
        None => (StatusCode::BAD_REQUEST, Json(json!({"detail": "URL not found"}))).into_response(),
    }
}

pub async fn redirect_url(
    State(state): State<AppState>,
    headers: axum::http::HeaderMap,
    Path(id): Path<String>,
) -> Response {
    if let Err(resp) = check_auth(&headers, &state.config) {
        return resp;
    }

    match state.cache.get(&id).await {
        Some(info) => {
            let new_accesses = state.cache.increment_accesses(&id).await;
            let _ = sqlx::query("UPDATE url SET accesses = $1 WHERE id = $2")
                .bind(new_accesses)
                .bind(&id)
                .execute(&state.pool)
                .await;

            let mut resp = Response::new(axum::body::Body::empty());
            *resp.status_mut() = StatusCode::FOUND;
            resp.headers_mut().insert(
                axum::http::header::LOCATION,
                axum::http::HeaderValue::from_str(&info.original).expect("valid URL"),
            );
            resp
        }
        None => (StatusCode::NOT_FOUND, Json(json!({"detail": "URL not found"}))).into_response(),
    }
}

pub async fn delete_url(
    State(state): State<AppState>,
    headers: axum::http::HeaderMap,
    Path(id): Path<String>,
) -> Response {
    if let Err(resp) = check_auth(&headers, &state.config) {
        return resp;
    }

    if !state.cache.exists(&id).await {
        return (StatusCode::BAD_REQUEST, Json(json!({"detail": "URL not found"}))).into_response();
    }

    match sqlx::query("DELETE FROM url WHERE id = $1")
        .bind(&id)
        .execute(&state.pool)
        .await
    {
        Ok(_) => {
            state.cache.delete(&id).await;
            StatusCode::OK.into_response()
        }
        Err(_) => (
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(json!({"detail": "Database error"})),
        )
            .into_response(),
    }
}
