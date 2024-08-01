package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/harshkhangarot07/backend/handlers"
	"github.com/harshkhangarot07/backend/middleware"
	"github.com/harshkhangarot07/backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{}, &models.Post{})

	// Set up router
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/api/register", handlers.Register(db)).Methods("POST")
	router.HandleFunc("/api/login", handlers.Login(db)).Methods("POST")

	// Protected routes
	protected := router.PathPrefix("/api").Subrouter()
	protected.Use(middleware.JwtVerify)
	protected.HandleFunc("/posts", handlers.CreatePost(db)).Methods("POST")
	protected.HandleFunc("/posts/{id:[0-9]+}", handlers.GetPost(db)).Methods("GET")
	protected.HandleFunc("/posts", handlers.GetPosts(db)).Methods("GET")

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", router))
}
