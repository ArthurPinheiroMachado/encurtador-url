class UrlCache {
  constructor() {
    this._urls = {};
  }

  load(urls) {
    this._urls = {};
    for (const u of urls) {
      this._urls[u.id] = { original: u.original, accesses: Number(u.accesses) };
    }
  }

  getAll() {
    const result = {};
    for (const [key, val] of Object.entries(this._urls)) {
      result[key] = { ...val };
    }
    return result;
  }

  get(id) {
    return this._urls[id] ? { ...this._urls[id] } : null;
  }

  exists(id) {
    return id in this._urls;
  }

  set(id, info) {
    this._urls[id] = info;
  }

  delete(id) {
    delete this._urls[id];
  }

  incrementAccesses(id) {
    if (this._urls[id]) {
      this._urls[id].accesses += 1;
      return this._urls[id].accesses;
    }
    return 0;
  }
}

module.exports = new UrlCache();
