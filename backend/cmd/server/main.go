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

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/mattn/go-sqlite3"

	"github.com/ivikasavnish/istio-ui/backend/internal/api"
	"github.com/ivikasavnish/istio-ui/backend/internal/kube"
	"github.com/ivikasavnish/istio-ui/backend/internal/scheduler"
	"github.com/ivikasavnish/istio-ui/backend/internal/storage"
)

func main() {
	log.Println("🚀 Starting MeshControl Center Backend...")

	// Initialize database
	dbPath := getEnv("DB_PATH", "./meshcontrol.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize storage
	store := storage.NewStore(db)
	if err := store.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}
	log.Println("✓ Database initialized")

	// Initialize Kubernetes client
	kubeClient, err := kube.NewClient()
	if err != nil {
		log.Printf("⚠ Warning: Kubernetes client initialization failed: %v", err)
		log.Println("  Running in standalone mode (some features will be limited)")
	} else {
		log.Println("✓ Kubernetes client connected")
	}

	// Initialize scheduler
	sched := scheduler.NewScheduler(store, kubeClient)
	sched.Start()
	defer sched.Stop()
	log.Println("✓ Scheduler started")

	// Initialize API server
	apiServer := api.NewServer(store, kubeClient)
	router := mux.NewRouter()

	// API v1 routes
	v1 := router.PathPrefix("/api/v1").Subrouter()
	apiServer.RegisterRoutes(v1)

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods("GET")

	// WebSocket endpoint
	router.HandleFunc("/ws", apiServer.HandleWebSocket)

	// Static file serving for frontend (in production)
	staticDir := getEnv("STATIC_DIR", "../frontend/dist")
	if _, err := os.Stat(staticDir); err == nil {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDir)))
		log.Printf("✓ Serving static files from %s\n", staticDir)
	}

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// HTTP Server
	port := getEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Printf("✓ Server listening on port %s\n", port)
		log.Printf("🌐 API available at http://localhost:%s/api/v1\n", port)
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

	log.Println("✓ Server stopped gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
