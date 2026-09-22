package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	controllers "github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/presenter"
	"github.com/mytheresa/go-hiring-challenge/app/usecase"
	"github.com/mytheresa/go-hiring-challenge/infrastructure/database"
	"github.com/mytheresa/go-hiring-challenge/infrastructure/logger"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	appLogger := &logger.ZapLogger{Logger: zapLogger}

	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		appLogger.Fatal("Error loading .env file", zap.Error(err))
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	db, close := database.NewConnection(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize handlers
	prodRepo := database.NewProductsRepository(db)
	uc := usecase.NewProductsUseCase(prodRepo, presenter.NewProductsPresenter())
	catalog := controllers.NewCatalogHandler(uc, appLogger)
	categoryRepo := database.NewCategoriesRepository(db)
	categoryUseCase := usecase.NewCategoriesUseCase(categoryRepo, presenter.NewCategoriesPresenter())
	categories := controllers.NewCategoriesHandler(categoryUseCase, appLogger)

	// Set up routing
	mux := http.NewServeMux()
	mux.HandleFunc("GET /catalog", catalog.GetAllProducts)
	mux.HandleFunc("GET /catalog/{code}", catalog.GetProductByCode)

	mux.HandleFunc("GET /categories", categories.GetAllCategories)
	mux.HandleFunc("GET /categories/{code}", categories.GetCategoryByCode)
	mux.HandleFunc("POST /categories", categories.CreateCategory)

	// Set up the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: mux,
	}

	// Start the server
	go func() {
		appLogger.Info("Starting server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("Server failed", zap.Error(err))
		}

		appLogger.Info("Server stopped gracefully")
	}()

	<-ctx.Done()
	appLogger.Info("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		appLogger.Error("Server shutdown failed", zap.Error(err))
	}
	stop()
}
