package main

import (
	"database/sql"
	"example/web-service-gin/app/controller"
	"example/web-service-gin/app/middleware"
	"example/web-service-gin/app/repository"
	"example/web-service-gin/app/service"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

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

	execFile(db, "app/db/create-tables.sql")
	execFile(db, "app/db/insert-tables.sql")
	fmt.Println("Database initialized successfully!")

	albumRepo := repository.NewAlbumRepository(db)
	albumService := service.NewAlbumService(albumRepo)
	albumController := controller.NewAlbumController(albumService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	router := gin.Default()

	// Public routes
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)

	// Protected routes
	authorized := router.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/albums", albumController.GetAlbums)
		authorized.POST("/add", albumController.AddAlbum)
		authorized.GET("/get/:id", albumController.GetAlbumByID)
	}

	router.Run("localhost:8080")

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
