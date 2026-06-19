const crypto = require("crypto");
const config = require("../config");

function basicAuth(req, res, next) {
  const authHeader = req.headers.authorization;

  if (!authHeader || !authHeader.startsWith("Basic ")) {
    res.set("WWW-Authenticate", 'Basic realm="Protected"');
    return res.status(401).json({ error: "Unauthorized" });
  }

  const token = authHeader.slice(6);
  let decoded;
  try {
    decoded = Buffer.from(token, "base64").toString("utf-8");
  } catch {
    res.set("WWW-Authenticate", 'Basic realm="Protected"');
    return res.status(401).json({ error: "Unauthorized" });
  }

  const parts = decoded.split(":");
  if (parts.length !== 2) {
    res.set("WWW-Authenticate", 'Basic realm="Protected"');
    return res.status(401).json({ error: "Unauthorized" });
  }

  const [username, password] = parts;

  const userBuf = Buffer.from(username, "utf-8");
  const expectedUserBuf = Buffer.from(config.auth.user, "utf-8");
  const passBuf = Buffer.from(password, "utf-8");
  const expectedPassBuf = Buffer.from(config.auth.pass, "utf-8");

  const userOk =
    userBuf.length === expectedUserBuf.length &&
    crypto.timingSafeEqual(userBuf, expectedUserBuf);
  const passOk =
    passBuf.length === expectedPassBuf.length &&
    crypto.timingSafeEqual(passBuf, expectedPassBuf);

  if (!userOk || !passOk) {
    res.set("WWW-Authenticate", 'Basic realm="Protected"');
    return res.status(401).json({ error: "Unauthorized" });
  }

  next();
}

module.exports = basicAuth;
