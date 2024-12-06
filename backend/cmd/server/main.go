package main

import (
	"log"
	"net/http"
	"os"

	"greenroute/internal/database"
	"greenroute/internal/external"
	"greenroute/internal/routes"
	"greenroute/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize external clients
	mapsClient, err := external.NewMapsClient()
	if err != nil {
		log.Fatalf("Failed to create maps client: %v", err)
	}

	// Initialize databases
	postgres, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	mongodb, err := database.NewMongoDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongodb.Close()

	// Initialize services
	routeService := services.NewRouteService(mapsClient, postgres, mongodb)

	// Initialize handlers
	routeHandler := routes.NewRouteHandler(routeService)

	// Initialize router with CORS middleware
	router := gin.Default()
	router.Use(corsMiddleware())

	// Register routes
	routeHandler.RegisterRoutes(router)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
