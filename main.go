package main

import (
	authHandlers "dklautomationgo/auth/handlers"
	"dklautomationgo/auth/middleware"
	"dklautomationgo/auth/service"
	"dklautomationgo/database"
	"dklautomationgo/database/repository"
	"dklautomationgo/handlers"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"dklautomationgo/services/email"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database connection
	dbConfig := database.NewConfig()
	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run auto migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize repositories
	contactRepo := repository.NewContactRepository(db)
	aanmeldingRepo := repository.NewAanmeldingRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Load email templates
	templatesDir := "templates"
	templates := make(map[string]*template.Template)
	templateFiles, err := filepath.Glob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		log.Fatalf("Failed to find template files: %v", err)
	}
	for _, file := range templateFiles {
		tmpl, err := template.ParseFiles(file)
		if err != nil {
			log.Fatalf("Failed to parse template %s: %v", file, err)
		}
		templates[filepath.Base(file)] = tmpl
	}

	// Initialize services
	emailService, err := email.NewEmailService()
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(userRepo, tokenService)
	aanmeldingService := services.NewAanmeldingService(aanmeldingRepo, emailService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService, userRepo)

	// Initialize handlers
	emailHandler := handlers.NewEmailHandler(emailService)
	contactHandler := handlers.NewContactHandler(emailService, contactRepo)
	aanmeldingHandler := handlers.NewAanmeldingHandler(aanmeldingService)
	authHandler := authHandlers.NewAuthHandler(authService, authMiddleware)

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

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		// Check database connection
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database connection error"})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database ping failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Service is healthy"})
	})

	// Register auth routes
	authHandler.RegisterRoutes(r)

	// API routes
	api := r.Group("/api")
	{
		// Email routes - beschermd met auth
		emails := api.Group("/emails")
		emails.Use(authMiddleware.RequireAuth())
		emails.Use(authMiddleware.RequireRole(models.RoleBeheerder, models.RoleAdmin))
		{
			emails.GET("", emailHandler.GetEmails)
			emails.GET("/stats", emailHandler.GetEmailStats)
			emails.PUT("/:id/read", emailHandler.MarkEmailAsRead)
		}

		// Contact form routes - gedeeltelijk beschermd
		contacts := api.Group("/contacts")
		{
			contacts.POST("", contactHandler.HandleContactEmail) // Publiek

			// Beschermde routes
			contactsAdmin := contacts.Group("")
			contactsAdmin.Use(authMiddleware.RequireAuth())
			contactsAdmin.Use(authMiddleware.RequireRole(models.RoleBeheerder, models.RoleAdmin))
			{
				contactsAdmin.GET("", contactHandler.GetContacts)
				contactsAdmin.PUT("/:id/status", contactHandler.UpdateContactStatus)
			}
		}

		// Backwards compatibility
		api.POST("/contact", contactHandler.HandleContactEmail)

		// Aanmelding routes - gedeeltelijk beschermd
		aanmeldingen := api.Group("/aanmeldingen")
		{
			aanmeldingen.POST("", aanmeldingHandler.CreateAanmelding) // Publiek

			// Beschermde routes
			aanmeldingenAdmin := aanmeldingen.Group("")
			aanmeldingenAdmin.Use(authMiddleware.RequireAuth())
			aanmeldingenAdmin.Use(authMiddleware.RequireRole(models.RoleBeheerder, models.RoleAdmin))
			{
				aanmeldingenAdmin.GET("", aanmeldingHandler.GetAanmeldingen)
				aanmeldingenAdmin.GET("/:id", aanmeldingHandler.GetAanmeldingByID)
				aanmeldingenAdmin.PUT("/:id", aanmeldingHandler.UpdateAanmelding)
				aanmeldingenAdmin.DELETE("/:id", aanmeldingHandler.DeleteAanmelding)
			}
		}

		// Backwards compatibility
		api.POST("/aanmelding", aanmeldingHandler.CreateAanmelding)
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
