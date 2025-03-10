package main

import (
	"dklautomationgo/handlers"
	"dklautomationgo/services/email"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize services
	emailService, err := email.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	// Initialize handlers
	emailHandler := handlers.NewEmailHandler(emailService)

	// Setup Gin
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://localhost:3000",
		"http://127.0.0.1:5173",
		"http://127.0.0.1:3000",
		"https://dekoninklijkeloop.nl",
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Content-Length",
		"Accept",
		"Authorization",
		"X-Requested-With",
	}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Content-Length"}
	config.MaxAge = 12 * 60 * 60 // 12 hours
	r.Use(cors.New(config))

	// API routes
	api := r.Group("/api")
	{
		// Email routes
		emails := api.Group("/emails")
		{
			emails.GET("", emailHandler.GetEmails)
			emails.GET("/stats", emailHandler.GetEmailStats)
			emails.PUT("/:id/read", emailHandler.MarkEmailAsRead)
		}

		// Contact form routes
		api.POST("/contact", emailHandler.HandleContactEmail)

		// Aanmelding routes
		api.POST("/aanmelding", emailHandler.HandleAanmeldingEmail)
	}

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
