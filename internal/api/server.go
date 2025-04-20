package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
	router *mux.Router
	server *http.Server
}

// NewServer creates a new API server
func NewServer(handler *Handler, port int) *Server {
	router := mux.NewRouter()
	
	// Register routes
	router.HandleFunc("/api/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/api/questions", handler.GetQuestions).Methods("GET")
	router.HandleFunc("/api/assessments", handler.StartAssessment).Methods("POST")
	router.HandleFunc("/api/assessments/{assessmentId}", handler.GetAssessment).Methods("GET")
	router.HandleFunc("/api/assessments/{assessmentId}/answers", handler.SaveAnswer).Methods("POST")
	router.HandleFunc("/api/assessments/{assessmentId}/complete", handler.CompleteAssessment).Methods("POST")
	router.HandleFunc("/api/assessments/{assessmentId}/report", handler.GetReport).Methods("GET")
	
	// Add middleware for logging, CORS, etc.
	router.Use(loggingMiddleware)
	router.Use(corsMiddleware)
	
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	return &Server{
		router: router,
		server: server,
	}
}

// Start begins the HTTP server
func (s *Server) Start() error {
	// Channel for server errors
	errChan := make(chan error, 1)
	
	// Start server in goroutine
	go func() {
		log.Printf("Starting server on %s", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()
	
	// Channel for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	
	// Wait for error or signal
	select {
	case err := <-errChan:
		return err
	case <-stop:
		log.Println("Shutting down server...")
		
		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		// Attempt graceful shutdown
		if err := s.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}
		
		log.Println("Server gracefully stopped")
	}
	
	return nil
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"up"}`))
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Process request
		next.ServeHTTP(w, r)
		
		// Log after request is processed
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Process request
		next.ServeHTTP(w, r)
	})
}
