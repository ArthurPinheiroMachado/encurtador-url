use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::RwLock;

use crate::models::{Url, UrlInfo};

#[derive(Clone)]
pub struct UrlCache {
    inner: Arc<RwLock<HashMap<String, UrlInfo>>>,
}

impl UrlCache {
    pub fn new() -> Self {
        Self {
            inner: Arc::new(RwLock::new(HashMap::new())),
        }
    }

    pub async fn load(&self, urls: Vec<Url>) {
        let mut map = self.inner.write().await;
        map.clear();
        for u in urls {
            map.insert(u.id, UrlInfo { original: u.original, accesses: u.accesses });
        }
    }

    pub async fn get_all(&self) -> HashMap<String, UrlInfo> {
        self.inner.read().await.clone()
    }

    pub async fn get(&self, id: &str) -> Option<UrlInfo> {
        self.inner.read().await.get(id).cloned()
    }

    pub async fn exists(&self, id: &str) -> bool {
        self.inner.read().await.contains_key(id)
    }

    pub async fn set(&self, id: String, info: UrlInfo) {
        self.inner.write().await.insert(id, info);
    }

    pub async fn delete(&self, id: &str) {
        self.inner.write().await.remove(id);
    }

    pub async fn increment_accesses(&self, id: &str) -> i64 {
        let mut map = self.inner.write().await;
        if let Some(info) = map.get_mut(id) {
            info.accesses += 1;
            info.accesses
        } else {
            0
        }
    }
}
