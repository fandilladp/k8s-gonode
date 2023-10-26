package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "host=127.0.0.1 user=postgres password=mysecretpassword dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	router := gin.Default()

	router.GET("/cats", getCats)
	router.POST("/cats", addCat)

	initializeDatabase()
	router.Run(":8080")
}

func initializeDatabase() {
	var tableName string
	err := db.QueryRow("SELECT to_regclass('public.cats')").Scan(&tableName)
	if err != nil {
		log.Fatalf("Failed to check table existence: %v", err)
	}

	if tableName == "" {
		_, err := db.Exec(`
            CREATE TABLE cats (
                name VARCHAR(255)
            )
        `)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
		log.Println("Table cats created")
	} else {
		log.Println("Table cats already exists")
	}
}

func getCats(c *gin.Context) {
	rows, err := db.Query("SELECT name FROM cats")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()

	var cats []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cats = append(cats, name)
	}

	c.JSON(http.StatusOK, gin.H{"cats": cats})
}

func addCat(c *gin.Context) {
	var input struct {
		Name string `json: "name"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO cats (name) VALUES ($1)", input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Cat added!"})
}
