local config = require("config")

local b64chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

local function b64_value(c)
  if c == "=" or c == "" then return 0 end
  local pos = b64chars:find(c, 1, true)
  return (pos or 1) - 1
end

function base64_decode(str)
  if not str then return "" end
  local cleaned = str:gsub("[^%w%+/=]", "")
  if #cleaned == 0 then return "" end

  local result = {}
  local padding = cleaned:match("=+$") or ""
  local main = cleaned:sub(1, #cleaned - #padding)

  for i = 1, #main, 4 do
    local chunk = main:sub(i, i + 3)
    local pad = 4 - #chunk
    for _ = 1, pad do chunk = chunk .. "=" end

    local c1, c2, c3, c4 = chunk:sub(1, 1), chunk:sub(2, 2), chunk:sub(3, 3), chunk:sub(4, 4)
    local a, bb, c, d = b64_value(c1), b64_value(c2), b64_value(c3), b64_value(c4)

    result[#result + 1] = string.char(a * 4 + math.floor(bb / 16))
    if c3 ~= "=" and c ~= 0 then
      result[#result + 1] = string.char((bb % 16) * 16 + math.floor(c / 4))
    end
    if c4 ~= "=" and d ~= 0 then
      result[#result + 1] = string.char((c % 4) * 64 + d)
    end
  end

  return table.concat(result)
end

function check_auth(auth_header)
  if not auth_header or not auth_header:match("^Basic ") then
    return false
  end
  local token = auth_header:sub(7)
  local ok, decoded = pcall(base64_decode, token)
  if not ok or not decoded then return false end
  local colon = decoded:find(":")
  if not colon then return false end
  local user = decoded:sub(1, colon - 1)
  local pass = decoded:sub(colon + 1)
  return user == config.auth.user and pass == config.auth.pass
end

local CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

function generate_short_id(length, exists)
  length = length or 8
  for _ = 1, 100 do
    local id_chars = {}
    for i = 1, length do
      local idx = math.random(1, #CHARSET)
      id_chars[i] = CHARSET:sub(idx, idx)
    end
    local id = table.concat(id_chars)
    if not exists or not exists(id) then
      return id
    end
  end
  error("failed to generate unique ID after 100 attempts")
end

return {
  check_auth = check_auth,
  generate_short_id = generate_short_id,
}
