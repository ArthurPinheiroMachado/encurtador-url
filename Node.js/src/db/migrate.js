const pool = require("./pool");

const queries = [
  `CREATE TABLE IF NOT EXISTS migrations(
    id INT NOT NULL,
    content TEXT NOT NULL,
    PRIMARY KEY(id)
  )`,
  `CREATE TABLE IF NOT EXISTS url(
    id TEXT NOT NULL,
    original TEXT NOT NULL,
    accesses BIGINT DEFAULT 0,
    UNIQUE(original),
    PRIMARY KEY(id)
  )`,
];

async function migrate() {
  const client = await pool.connect();
  try {
    const { rows } = await client.query("SELECT COALESCE(MAX(id), -1) AS pos FROM migrations");
    const lastPos = rows[0].pos;

    for (let i = 0; i < queries.length; i++) {
      if (i <= lastPos) continue;

      await client.query(queries[i]);
      await client.query("INSERT INTO migrations(id, content) VALUES($1, $2)", [i, queries[i]]);
    }

    console.log("Migrations applied successfully");
  } finally {
    client.release();
  }
}

module.exports = { migrate };
