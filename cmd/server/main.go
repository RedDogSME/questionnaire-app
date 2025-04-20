package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"questionnaire-app/internal/api"
	"questionnaire-app/internal/models"
	"questionnaire-app/internal/services"
	"questionnaire-app/internal/storage"
	"strconv"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", getEnvInt("PORT", 8080), "Server port")
	dataDir := flag.String("data", getEnvStr("DATA_DIR", "./data"), "Data directory")
	flag.Parse()
	
	log.Printf("Starting questionnaire application on port %d with data directory %s", *port, *dataDir)
	
	// Ensure data directory exists
	if err := os.MkdirAll(*dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}
	
	// Initialize storage
	store, err := storage.NewFileStorage(*dataDir)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	
	// Add sample data if needed
	if err := ensureSampleData(store); err != nil {
		log.Fatalf("Failed to add sample data: %v", err)
	}
	
	// Initialize services
	assessmentService := services.NewAssessmentService(store)
	
	// Initialize HTTP handlers
	handler := api.NewHandler(assessmentService)
	
	// Initialize and start server
	server := api.NewServer(handler, *port)
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// getEnvStr gets a string environment variable with a fallback
func getEnvStr(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvInt gets an integer environment variable with a fallback
func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// ensureSampleData adds sample questions and applications if none exist
func ensureSampleData(store storage.Storage) error {
	// Check if we already have questions
	questions, err := store.GetQuestions(nil)
	if err != nil {
		return err
	}
	
	// If we already have questions, don't add sample data
	if len(questions) > 0 {
		return nil
	}
	
	// Add sample questions
	sampleQuestions := []*models.Question{
		{
			ID:       "q1",
			Text:     "Is the application stateless?",
			Category: "Architecture",
			Weight:   5,
			Options: []models.Option{
				{ID: "q1_a1", Text: "Yes, completely stateless", Points: 10},
				{ID: "q1_a2", Text: "Mostly stateless with minimal state", Points: 7},
				{ID: "q1_a3", Text: "Partially stateless", Points: 4},
				{ID: "q1_a4", Text: "Heavily stateful", Points: 1},
			},
		},
		{
			ID:       "q2",
			Text:     "Does the application use external configuration?",
			Category: "Configuration",
			Weight:   3,
			Options: []models.Option{
				{ID: "q2_a1", Text: "Yes, all configuration is external", Points: 10},
				{ID: "q2_a2", Text: "Most configuration is external", Points: 7},
				{ID: "q2_a3", Text: "Some configuration is external", Points: 4},
				{ID: "q2_a4", Text: "No, all configuration is internal", Points: 1},
			},
		},
		{
			ID:       "q3",
			Text:     "How is application logging handled?",
			Category: "Observability",
			Weight:   2,
			Options: []models.Option{
				{ID: "q3_a1", Text: "Logs to stdout/stderr", Points: 10},
				{ID: "q3_a2", Text: "Logs to configurable location", Points: 7},
				{ID: "q3_a3", Text: "Logs to fixed file location", Points: 3},
				{ID: "q3_a4", Text: "No logging capability", Points: 0},
			},
		},
		{
			ID:       "q4",
			Text:     "How does the application store persistent data?",
			Category: "Persistence",
			Weight:   4,
			Options: []models.Option{
				{ID: "q4_a1", Text: "Uses external databases with connection strings", Points: 10},
				{ID: "q4_a2", Text: "Uses external storage with configurable location", Points: 7},
				{ID: "q4_a3", Text: "Uses local filesystem with fixed paths", Points: 3},
				{ID: "q4_a4", Text: "Embedded database or storage", Points: 1},
			},
		},
		{
			ID:       "q5",
			Text:     "Does the application support horizontal scaling?",
			Category: "Scalability",
			Weight:   5,
			Options: []models.Option{
				{ID: "q5_a1", Text: "Designed for horizontal scaling", Points: 10},
				{ID: "q5_a2", Text: "Can scale horizontally with minor changes", Points: 7},
				{ID: "q5_a3", Text: "Requires significant changes to scale horizontally", Points: 3},
				{ID: "q5_a4", Text: "Cannot scale horizontally", Points: 0},
			},
		},
	}
	
	// Save sample questions
	for _, question := range sampleQuestions {
		path := filepath.Join(store.(*storage.FileStorage).BasePath, "questions", question.ID+".json")
		content, _ := json.Marshal(question)
		if err := os.WriteFile(path, content, 0644); err != nil {
			return err
		}
	}
	
	// Add sample application
	sampleApp := &models.Application{
		ID:          "app1",
		Name:        "Sample Application",
		Description: "A sample application for testing the assessment tool",
		Tags: map[string]string{
			"language": "Java",
			"type":     "Web Application",
		},
	}
	
	// Save sample application
	return store.SaveApplication(nil, sampleApp)
}
