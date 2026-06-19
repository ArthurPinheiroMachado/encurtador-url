const crypto = require("crypto");
const pool = require("../db/pool");
const cache = require("../cache");
const basicAuth = require("../middleware/auth");
const { Router } = require("express");

const router = Router();

router.use(basicAuth);

router.get("/urls", (req, res) => {
  res.json(cache.getAll());
});

router.post("/urls", async (req, res) => {
  const { url } = req.body;

  if (!url) {
    return res.status(400).json({ detail: "Invalid URL" });
  }

  let parsed;
  try {
    parsed = new URL(url);
    if (parsed.protocol !== "http:" && parsed.protocol !== "https:") throw new Error();
  } catch {
    return res.status(400).json({ detail: "Invalid URL" });
  }

  const { rows } = await pool.query("SELECT id FROM url WHERE original = $1", [url]);
  if (rows.length > 0) {
    return res.status(200).json({ id: rows[0].id, url });
  }

  const shortId = generateShortId(8, (id) => cache.exists(id));

  await pool.query(
    "INSERT INTO url(id, original, accesses) VALUES($1, $2, 0)",
    [shortId, url]
  );

  cache.set(shortId, { original: url, accesses: 0 });

  res.status(201).json({ id: shortId, url });
});

router.get("/urls/:id", (req, res) => {
  const info = cache.get(req.params.id);
  if (!info) {
    return res.status(400).json({ detail: "URL not found" });
  }
  res.json(info);
});

router.get("/:id", async (req, res) => {
  const info = cache.get(req.params.id);
  if (!info) {
    return res.status(404).json({ detail: "URL not found" });
  }

  const newAccesses = cache.incrementAccesses(req.params.id);
  await pool.query("UPDATE url SET accesses = $1 WHERE id = $2", [newAccesses, req.params.id]);

  res.redirect(302, info.original);
});

router.delete("/:id", async (req, res) => {
  if (!cache.exists(req.params.id)) {
    return res.status(400).json({ detail: "URL not found" });
  }

  await pool.query("DELETE FROM url WHERE id = $1", [req.params.id]);
  cache.delete(req.params.id);

  res.status(200).end();
});

function generateShortId(length, exists) {
  const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

  for (let i = 0; i < 100; i++) {
    const bytes = crypto.randomBytes(length);
    let id = "";
    for (let j = 0; j < length; j++) {
      id += charset[bytes[j] % charset.length];
    }
    if (!exists(id)) {
      return id;
    }
  }

  throw new Error("failed to generate unique ID after 100 attempts");
}

module.exports = router;
