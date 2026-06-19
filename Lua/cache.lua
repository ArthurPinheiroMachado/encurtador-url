local cache = {}

function cache:load(urls)
  self._urls = {}
  for _, u in ipairs(urls) do
    self._urls[u.id] = { original = u.original, accesses = tonumber(u.accesses) or 0 }
  end
end

function cache:get_all()
  local result = {}
  for k, v in pairs(self._urls) do
    result[k] = { original = v.original, accesses = v.accesses }
  end
  return result
end

function cache:get(id)
  local v = self._urls[id]
  if v then
    return { original = v.original, accesses = v.accesses }
  end
  return nil
end

function cache:exists(id)
  return self._urls[id] ~= nil
end

function cache:set(id, info)
  self._urls[id] = info
end

function cache:delete(id)
  self._urls[id] = nil
end

function cache:increment_accesses(id)
  local v = self._urls[id]
  if v then
    v.accesses = v.accesses + 1
    return v.accesses
  end
  return 0
end

return cache
