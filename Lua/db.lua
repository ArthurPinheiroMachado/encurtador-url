local pgmoon = require("pgmoon")
local config = require("config")
local json = require("dkjson")

local _db

function get_db()
  if not _db then
    _db = pgmoon.new({
      host = config.db.host,
      port = config.db.port,
      database = config.db.database,
      user = config.db.user,
      password = config.db.password,
    })
    assert(_db:connect(), "Failed to connect to PostgreSQL")
  end
  return _db
end

function migrate()
  local db = get_db()

  local ok, result = pcall(db.query, db, "SELECT COALESCE(MAX(id), -1) AS pos FROM migrations")
  local last_pos = -1
  if ok and result and #result > 0 then
    last_pos = tonumber(result[1].pos) or -1
  end

  local statements = {
    [[CREATE TABLE IF NOT EXISTS migrations(
      id INT NOT NULL,
      content TEXT NOT NULL,
      PRIMARY KEY(id)
    )]],
    [[CREATE TABLE IF NOT EXISTS url(
      id TEXT NOT NULL,
      original TEXT NOT NULL,
      accesses BIGINT DEFAULT 0,
      UNIQUE(original),
      PRIMARY KEY(id)
    )]],
  }

  for idx, stmt in ipairs(statements) do
    local i = idx - 1
    if i > last_pos then
      db:query(stmt)
      db:query("INSERT INTO migrations(id, content) VALUES($1, $2)", { i, stmt })
    end
  end
end

function load_urls_from_db()
  local db = get_db()
  local ok, result = pcall(db.query, db, "SELECT id, original, accesses FROM url")
  if ok and result then
    return result
  end
  return {}
end

function get_url_by_original(original)
  local db = get_db()
  local ok, result = pcall(db.query, db, "SELECT id FROM url WHERE original = $1", { original })
  if ok and result and #result > 0 then
    return result[1]
  end
  return nil
end

function insert_url(id, original)
  local db = get_db()
  return db:query("INSERT INTO url(id, original, accesses) VALUES($1, $2, 0)", { id, original })
end

function delete_url(id)
  local db = get_db()
  return db:query("DELETE FROM url WHERE id = $1", { id })
end

function update_accesses(id, accesses)
  local db = get_db()
  return db:query("UPDATE url SET accesses = $1 WHERE id = $2", { accesses, id })
end

return {
  migrate = migrate,
  load_urls_from_db = load_urls_from_db,
  get_url_by_original = get_url_by_original,
  insert_url = insert_url,
  delete_url = delete_url,
  update_accesses = update_accesses,
}
