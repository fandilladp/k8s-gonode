const express = require("express");
const { Pool } = require("pg");
const app = express();
const bodyParser = require("body-parser");
app.use(bodyParser.json());

const pool = new Pool({
  user: "myuser",
  host: "127.0.0.1",
  database: "mydb",
  password: "mypassword",
  port: 5432,
});

app.get("/cats", async (req, res) => {
  try {
    const result = await pool.query("SELECT name FROM  cats");
    res.json(result.rows);
  } catch (err) {
    res.status(500).send(err.message);
  }
});

app.post("/cats", async (req, res) => {
  try {
    const { name } = req.body;
    await pool.query("INSERT INTO cats (name) VALUES ($1)", [name]);
    res.status(201).send("Cat added!");
  } catch (err) {
    res.status(500).send(err.message);
  }
});
async function initializeDatabase() {
  try {
    const result = await pool.query("SELECT to_regclass('public.cats')");
    if (result.rows[0].to_regclass === null) {
      await pool.query(`
                CREATE TABLE cats (
                    name VARCHAR(255)
                )
            `);
      console.log("Table cats created");
    } else {
      console.log("Table cats already exists");
    }
  } catch (err) {
    console.error("Failed to initialize database:", err.message);
  }
}

initializeDatabase().then(() => {
  app.listen(3000, () => {
    console.log("Node app listening on port 3000");
  });
});
