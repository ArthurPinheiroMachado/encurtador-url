use sqlx::postgres::PgPool;

use crate::config::Config;

pub async fn create_pool(config: &Config) -> Result<PgPool, sqlx::Error> {
    let conn_string = format!(
        "postgres://{}:{}@{}:{}/{}",
        config.db_user, config.db_pass, config.db_host, config.db_port, config.db_name
    );
    PgPool::connect(&conn_string).await
}

pub async fn migrate(pool: &PgPool) -> Result<(), sqlx::Error> {
    let last_pos: i32 = sqlx::query_scalar("SELECT COALESCE(MAX(id), -1) FROM migrations")
        .fetch_one(pool)
        .await
        .unwrap_or(-1);

    let statements = vec![
        "CREATE TABLE IF NOT EXISTS migrations(
            id INT NOT NULL,
            content TEXT NOT NULL,
            PRIMARY KEY(id)
        )",
        "CREATE TABLE IF NOT EXISTS url(
            id TEXT NOT NULL,
            original TEXT NOT NULL,
            accesses BIGINT DEFAULT 0,
            UNIQUE(original),
            PRIMARY KEY(id)
        )",
    ];

    for (idx, stmt) in statements.iter().enumerate() {
        let idx = idx as i32;
        if idx <= last_pos {
            continue;
        }
        sqlx::query(stmt).execute(pool).await?;
        sqlx::query("INSERT INTO migrations(id, content) VALUES($1, $2)")
            .bind(idx)
            .bind(stmt)
            .execute(pool)
            .await?;
    }

    Ok(())
}
