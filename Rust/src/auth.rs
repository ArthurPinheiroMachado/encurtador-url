use base64::Engine as _;
use serde_json::json;

use crate::config::Config;

pub fn is_authorized(auth_header: Option<&str>, config: &Config) -> bool {
    match auth_header {
        Some(header) if header.starts_with("Basic ") => {
            let token = &header[6..];
            match base64::engine::general_purpose::STANDARD.decode(token) {
                Ok(decoded) => match String::from_utf8(decoded) {
                    Ok(creds) => {
                        let parts: Vec<&str> = creds.splitn(2, ':').collect();
                        parts.len() == 2 && parts[0] == config.auth_user && parts[1] == config.auth_pass
                    }
                    Err(_) => false,
                },
                Err(_) => false,
            }
        }
        _ => false,
    }
}

pub fn unauthorized_response() -> (axum::http::StatusCode, axum::Json<serde_json::Value>) {
    (axum::http::StatusCode::UNAUTHORIZED, axum::Json(json!({"error": "Unauthorized"})))
}
