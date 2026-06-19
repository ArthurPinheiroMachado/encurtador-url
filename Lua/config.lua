local M = {}

M.db = {
  host = os.getenv("DB_HOST") or "0.0.0.0",
  port = tonumber(os.getenv("DB_PORT")) or 5432,
  database = os.getenv("DB_NAME") or "encurtador",
  user = os.getenv("DB_USER") or "postgres",
  password = os.getenv("DB_PASS") or "postgres",
}

M.http = {
  port = tonumber(os.getenv("HTTP_PORT")) or 6060,
  base = (os.getenv("HTTP_BASE") or "/api/"):gsub("/+$", ""),
}

M.auth = {
  user = os.getenv("USER") or "user",
  pass = os.getenv("PASS") or "pass123",
}

return M
