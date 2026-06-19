const config = {
  db: {
    host: process.env.DB_HOST || "0.0.0.0",
    port: parseInt(process.env.DB_PORT, 10) || 5432,
    database: process.env.DB_NAME || "encurtador",
    user: process.env.DB_USER || "postgres",
    password: process.env.DB_PASS || "postgres",
  },
  http: {
    port: parseInt(process.env.HTTP_PORT, 10) || 6060,
    base: (process.env.HTTP_BASE || "/api/").replace(/\/+$/, ""),
    timeout: parseInt(process.env.TIMEOUT_TIME, 10) || 3,
  },
  auth: {
    user: process.env.USER || "user",
    pass: process.env.PASS || "pass123",
  },
};

module.exports = config;
