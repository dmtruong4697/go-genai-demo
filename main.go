package main

import (
	"go-genai-demo/src/routes"
	"log"
	"os"

	_ "go-genai-demo/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title go-genai-demo
// @version 1.0
// @description go-genai-demo
// @BasePath /api
func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	envFile := ".env." + env

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("⚠️  Warning: Could not load %s. Using system environment variables.", envFile)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8910"
	}

	r := gin.Default()

	// r.Static("/static", "./static")
	r.Static("/", "./publics")

	routes.SetupRoutes(r)

	log.Printf("🚀 Server running on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server: ", err)
	}
}
