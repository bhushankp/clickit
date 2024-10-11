package main

import (
	"clickit/internal/config"
	"clickit/internal/handlers"
	"log"

	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize MySQL and Redis
	config.ConnectDatabase()
	config.ConnectRedis()

	// Setup the router
	router := mux.NewRouter()

	// Routes for Record operations
	router.HandleFunc("/upload", handlers.UploadExcel).Methods("POST")
	router.HandleFunc("/records", handlers.GetPaginatedRecords).Methods("GET")
	router.HandleFunc("/records/{id}", handlers.UpdateRecord).Methods("PUT")
	router.HandleFunc("/records/{id}", handlers.DeleteRecord).Methods("DELETE")

	// Start server
	log.Printf("Server running on port %s", os.Getenv("API_PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("API_PORT"), router))
}
