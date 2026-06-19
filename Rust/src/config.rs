use std::env;

#[derive(Clone)]
pub struct Config {
    pub db_host: String,
    pub db_port: String,
    pub db_name: String,
    pub db_user: String,
    pub db_pass: String,
    pub http_port: u16,
    pub http_base: String,
    pub auth_user: String,
    pub auth_pass: String,
}

impl Config {
    pub fn from_env() -> Self {
        Self {
            db_host: env::var("DB_HOST").unwrap_or_else(|_| "0.0.0.0".to_string()),
            db_port: env::var("DB_PORT").unwrap_or_else(|_| "5432".to_string()),
            db_name: env::var("DB_NAME").unwrap_or_else(|_| "encurtador".to_string()),
            db_user: env::var("DB_USER").unwrap_or_else(|_| "postgres".to_string()),
            db_pass: env::var("DB_PASS").unwrap_or_else(|_| "postgres".to_string()),
            http_port: env::var("HTTP_PORT")
                .unwrap_or_else(|_| "6060".to_string())
                .parse()
                .unwrap_or(6060),
            http_base: env::var("HTTP_BASE").unwrap_or_else(|_| "/api/".to_string()),
            auth_user: env::var("USER").unwrap_or_else(|_| "user".to_string()),
            auth_pass: env::var("PASS").unwrap_or_else(|_| "pass123".to_string()),
        }
    }

    pub fn http_base_trimmed(&self) -> String {
        self.http_base.trim_end_matches('/').to_string()
    }
}
