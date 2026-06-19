use serde::{Deserialize, Serialize};



#[derive(Debug, Clone, Serialize, Deserialize, sqlx::FromRow)]
pub struct Url {
    pub id: String,
    pub original: String,
    pub accesses: i64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UrlInfo {
    pub original: String,
    pub accesses: i64,
}

#[derive(Debug, Deserialize)]
pub struct CreateUrlPayload {
    pub url: String,
}

