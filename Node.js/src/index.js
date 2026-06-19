const pool = require("./db/pool");
const { migrate } = require("./db/migrate");
const cache = require("./cache");
const config = require("./config");
const app = require("./app");

async function main() {
  await migrate();

  const { rows } = await pool.query("SELECT id, original, accesses FROM url");
  cache.load(rows);

  app.listen(config.http.port, () => {
    console.log(`Starting ENCURTADOR at port ${config.http.port}`);
  });
}

main().catch((err) => {
  console.error("Failed to start server:", err);
  process.exit(1);
});
