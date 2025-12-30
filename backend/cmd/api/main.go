package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/username/url-shortener/internal/handler"
	"github.com/username/url-shortener/internal/repository"
	"github.com/username/url-shortener/internal/service"
)

func main() {
	// Configuration from environment
	config := loadConfig()

	// Initialize PostgreSQL
	db, err := initPostgres(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	rdb := initRedis(config.RedisURL)
	defer rdb.Close()

	// Initialize repositories (Hexagonal Architecture - Adapters)
	postgresRepo := repository.NewPostgresRepository(db)
	redisRepo := repository.NewRedisRepository(rdb)

	// Initialize service (Hexagonal Architecture - Application/Domain)
	urlService := service.NewURLService(postgresRepo, redisRepo)

	// Initialize handlers (Hexagonal Architecture - Ports)
	urlHandler := handler.NewURLHandler(urlService)

	// Setup Gin router
	router := setupRouter(urlHandler)

	// CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      corsHandler.Handler(router),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server starting on port %s", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited gracefully")
}

// Config holds application configuration
type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
	BaseURL     string
}

func loadConfig() Config {
	return Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/urlshortener?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		BaseURL:     getEnv("BASE_URL", "http://localhost"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func initPostgres(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	log.Println("✅ Connected to PostgreSQL")
	return db, nil
}

func runMigrations(db *sql.DB) error {
	migration := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(10) UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		clicks INTEGER DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);
	`
	_, err := db.Exec(migration)
	return err
}

func initRedis(redisURL string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         redisURL,
		Password:     getEnv("REDIS_PASSWORD", ""),
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Printf("⚠️ Redis connection failed: %v (continuing without cache)", err)
	} else {
		log.Println("✅ Connected to Redis")
	}

	return rdb
}

func setupRouter(h *handler.URLHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
		})
	})

	// API routes
	api := router.Group("/api")
	{
		api.POST("/shorten", h.CreateShortURL)
		api.GET("/stats/:code", h.GetURLStats)
	}

	// Redirect route (at root level)
	router.GET("/:code", h.RedirectURL)

	return router
}
