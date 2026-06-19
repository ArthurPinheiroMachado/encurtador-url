const express = require("express");
const config = require("./config");
const urlsRouter = require("./routes/urls");

const app = express();

app.use(express.json());
app.use(config.http.base, urlsRouter);

module.exports = app;
