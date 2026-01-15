package main

import (
	"database/sql"
	"example/web-service-gin/models"
	"example/web-service-gin/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var albumRepo *repository.AlbumRepository

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to the database
	// Use environment variables for credentials to avoid hardcoding secrets
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")

	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/web_service_recordings?multiStatements=true", dbUser, dbPass)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	execFile(db, "db/create-tables.sql")
	execFile(db, "db/insert-tables.sql")
	fmt.Println("Database initialized successfully!")

	// Initialize repository
	albumRepo = repository.NewAlbumRepository(db)

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/add", addAlbum)
	router.GET("/get/:id", getAlbumByID)

	router.Run("localhost:8081")

}
func execFile(db *sql.DB, filepath string) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", filepath, err)
	}

	_, err = db.Exec(string(content))
	if err != nil {
		log.Fatalf("Error executing SQL in %s: %v", filepath, err)
	}
	fmt.Printf("Successfully executed %s\n", filepath)
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(context *gin.Context) {
	albums, err := albumRepo.GetAll()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, albums)
}

func addAlbum(context *gin.Context) {
	var newAlbum models.Album
	if err := context.BindJSON(&newAlbum); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := albumRepo.Add(newAlbum)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newAlbum.ID = id
	context.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	alb, err := albumRepo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, alb)
}
