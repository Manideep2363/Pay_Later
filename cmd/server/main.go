package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"paylater/internal/config"
	"paylater/internal/routes"
)

func main() {

	// Load environment variables
	cfg := config.LoadConfig()

	// Connect to MySQL
	dbConn, err := config.NewMySQL(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Create Gin router
	router := gin.Default()

	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	// Register all routes
	routes.SetupRoutes(router, dbConn)

	// Start server
	log.Printf("Server running on port %s", cfg.ServerPort)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}