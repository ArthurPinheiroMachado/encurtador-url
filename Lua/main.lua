local socket = require("socket")
local json = require("dkjson")
local config = require("config")
local db = require("db")
local cache = require("cache")
local utils = require("utils")

math.randomseed(os.time() + os.clock() * 1e6)

db.migrate()
local urls = db.load_urls_from_db()
cache:load(urls)

local prefix = config.http.base

local function send_json(client, status, status_text, data, extra_headers)
  local body = json.encode(data)
  local resp = "HTTP/1.1 " .. status .. " " .. status_text .. "\r\n"
    .. "Content-Type: application/json\r\n"
    .. "Content-Length: " .. #body .. "\r\n"
  if extra_headers then
    for k, v in pairs(extra_headers) do
      resp = resp .. k .. ": " .. v .. "\r\n"
    end
  end
  resp = resp .. "Connection: close\r\n"
    .. "\r\n"
    .. body
  client:send(resp)
end

local function send_redirect(client, location)
  local resp = "HTTP/1.1 302 Found\r\n"
    .. "Location: " .. location .. "\r\n"
    .. "Content-Length: 0\r\n"
    .. "Connection: close\r\n"
    .. "\r\n"
  client:send(resp)
end

local function parse_request(client)
  local line, err = client:receive("*l")
  if not line then return nil end

  local method, path = line:match("^(%S+)%s+(%S+)")
  if not method then return nil end

  local headers = {}
  while true do
    local h, err = client:receive("*l")
    if not h or h == "" then break end
    local key, val = h:match("^([%w%-]+):%s*(.+)$")
    if key then
      headers[key:lower()] = val
    end
  end

  local body = ""
  local cl = headers["content-length"]
  if cl then
    local b, err = client:receive(tonumber(cl))
    if b then body = b end
  end

  return method, path, headers, body
end

local function strip_prefix(path)
  if path:sub(1, #prefix) == prefix then
    path = path:sub(#prefix + 1)
  end
  if path == "" then path = "/" end
  return path
end

local function handle_request(client, method, path, headers, body)
  path = strip_prefix(path)

  if not utils.check_auth(headers["authorization"]) then
    send_json(client, 401, "Unauthorized", { error = "Unauthorized" }, { ["WWW-Authenticate"] = 'Basic realm="Protected"' })
    return
  end

  if method == "GET" and path == "/urls" then
    send_json(client, 200, "OK", cache:get_all())

  elseif method == "POST" and path == "/urls" then
    local ok, payload = pcall(json.decode, body or "{}")
    if not ok or not payload.url or type(payload.url) ~= "string" then
      send_json(client, 400, "Bad Request", { detail = "Invalid URL" })
      return
    end
    if not payload.url:match("^https?://") then
      send_json(client, 400, "Bad Request", { detail = "Invalid URL" })
      return
    end

    local existing = db.get_url_by_original(payload.url)
    if existing then
      send_json(client, 200, "OK", { id = existing.id, url = payload.url })
      return
    end

    local short_id = utils.generate_short_id(8, function(id) return cache:exists(id) end)
    local ok, err = pcall(db.insert_url, short_id, payload.url)
    if not ok then
      send_json(client, 500, "Internal Server Error", { detail = tostring(err) })
      return
    end

    cache:set(short_id, { original = payload.url, accesses = 0 })
    send_json(client, 201, "Created", { id = short_id, url = payload.url })

  elseif method == "GET" and path:match("^/urls/(.+)$") then
    local id = path:match("^/urls/(.+)$")
    local info = cache:get(id)
    if not info then
      send_json(client, 400, "Bad Request", { detail = "URL not found" })
      return
    end
    send_json(client, 200, "OK", info)

  elseif method == "GET" then
    local id = path:match("^/(.+)$")
    if not id then
      send_json(client, 404, "Not Found", { detail = "URL not found" })
      return
    end
    local info = cache:get(id)
    if not info then
      send_json(client, 404, "Not Found", { detail = "URL not found" })
      return
    end
    local new_accesses = cache:increment_accesses(id)
    pcall(db.update_accesses, id, new_accesses)
    send_redirect(client, info.original)

  elseif method == "DELETE" then
    local id = path:match("^/(.+)$")
    if not id or not cache:exists(id) then
      send_json(client, 400, "Bad Request", { detail = "URL not found" })
      return
    end
    pcall(db.delete_url, id)
    cache:delete(id)
    send_json(client, 200, "OK", {})

  else
    send_json(client, 404, "Not Found", { detail = "Not found" })
  end
end

local server = socket.tcp()
server:setoption("reuseaddr", true)
server:bind("*", config.http.port)
server:listen(128)
print("Starting ENCURTADOR at port " .. config.http.port)

while true do
  local client, err = server:accept()
  if client then
    client:settimeout(5)
    local method, path, headers, body = parse_request(client)
    if method then
      handle_request(client, method, path, headers, body)
    end
    client:close()
  end
end
