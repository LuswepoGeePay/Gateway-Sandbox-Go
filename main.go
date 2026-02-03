package main

import (
	"log"
	"pg_sandbox/config"
	"pg_sandbox/routes"
	"pg_sandbox/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	err := utils.InitLogger("pg-sandbox.log")

	if err != nil {
		panic("Failed to initialize logger:" + err.Error())
	}

	config.InitDB()

	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true                                  // Allow all origins, or specify specific origins
	config.AllowMethods = []string{"GET", "POST", "DELETE", "PUT"} // Allow specific HTTP methods
	config.AllowHeaders = []string{"*"}                            // Allow specific headers

	r.Use(cors.New(config))

	routes.SetupRoutes(r)

	// Fall back to HTTP if SSL is not configured
	log.Println("Starting HTTP server on port 2000")
	if err := r.Run(":2000"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
