package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/robstave/gorag/docs"
	"github.com/robstave/gorag/internal/adapters/controller"
	"github.com/robstave/gorag/internal/adapters/repositories"
	"github.com/robstave/gorag/internal/domain"
	"github.com/robstave/gorag/internal/domain/types"
	"github.com/robstave/gorag/internal/logger"
	httpSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title gorag
// @version 1.0
// @description API documentation for gorag
// @BasePath /api
func main() {
	// Initialize Logger
	slogger := logger.InitializeLogger()
	logger.SetLogger(slogger)

	// Set database path
	dbPath := "./gorag.db"
	if path := os.Getenv("DB_PATH"); path != "" {
		dbPath = path
	}
	slogger.Info("DBPath set", "dbpath", dbPath)

	// Open SQLite database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		slogger.Error("Failed to connect to database", "path", dbPath, "error", err)
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate widget model
	if err = db.AutoMigrate(&types.Widget{}); err != nil {
		slogger.Error("Failed to migrate database", "error", err)
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Repository, Service, and Controller
	repo := repositories.NewRepositorySQLite(db)
	service := domain.NewService(slogger, repo)
	ctrl := controller.NewController(service, slogger)

	// Initialize Echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// API Routes
	api := e.Group("/api")
	widgetGroup := api.Group("/widgets")
	widgetGroup.POST("", ctrl.CreateWidget)
	widgetGroup.GET("", ctrl.GetAllWidgets)
	widgetGroup.GET("/:id", ctrl.GetWidget)
	widgetGroup.PUT("/:id", ctrl.UpdateWidget)
	widgetGroup.DELETE("/:id", ctrl.DeleteWidget)

	// Swagger endpoint
	e.GET("/swagger/*", httpSwagger.WrapHandler)

	// Start Server
	port := "8711"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	slogger.Info("Starting server", "port", port)
	if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
		slogger.Error("Shutting down the server", "error", err)
		log.Fatalf("Shutting down the server: %v", err)
	}
}
