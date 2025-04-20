package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"questionnaire-app/internal/models"
)

// Storage defines the interface for persistence operations
type Storage interface {
	// Application operations
	GetApplication(ctx context.Context, id string) (*models.Application, error)
	ListApplications(ctx context.Context) ([]*models.Application, error)
	SaveApplication(ctx context.Context, app *models.Application) error
	
	// Question operations
	GetQuestions(ctx context.Context) ([]*models.Question, error)
	GetQuestion(ctx context.Context, id string) (*models.Question, error)
	
	// Assessment operations
	CreateAssessment(ctx context.Context, assessment *models.Assessment) error
	GetAssessment(ctx context.Context, id string) (*models.Assessment, error)
	UpdateAssessment(ctx context.Context, assessment *models.Assessment) error
	ListAssessments(ctx context.Context, applicationID string) ([]*models.Assessment, error)
	
	// Report operations
	SaveReport(ctx context.Context, report *models.Report) error
	GetReport(ctx context.Context, assessmentID string) (*models.Report, error)
}

// FileStorage implements Storage interface using local file system
type FileStorage struct {
	basePath string
}

// NewFileStorage creates a new file-based storage
func NewFileStorage(basePath string) (*FileStorage, error) {
	// Create necessary directories
	dirs := []string{
		filepath.Join(basePath, "applications"),
		filepath.Join(basePath, "questions"),
		filepath.Join(basePath, "assessments"),
		filepath.Join(basePath, "reports"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	return &FileStorage{basePath: basePath}, nil
}

// GetApplication retrieves an application by ID
func (s *FileStorage) GetApplication(ctx context.Context, id string) (*models.Application, error) {
	path := filepath.Join(s.basePath, "applications", id+".json")
	
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read application file: %w", err)
	}
	
	var app models.Application
	if err := json.Unmarshal(data, &app); err != nil {
		return nil, fmt.Errorf("failed to unmarshal application: %w", err)
	}
	
	return &app, nil
}

// ListApplications returns all applications
func (s *FileStorage) ListApplications(ctx context.Context) ([]*models.Application, error) {
	dir := filepath.Join(s.basePath, "applications")
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read applications directory: %w", err)
	}
	
	var apps []*models.Application
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		path := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read application file %s: %w", file.Name(), err)
		}
		
		var app models.Application
		if err := json.Unmarshal(data, &app); err != nil {
			return nil, fmt.Errorf("failed to unmarshal application %s: %w", file.Name(), err)
		}
		
		apps = append(apps, &app)
	}
	
	return apps, nil
}

// SaveApplication stores an application
func (s *FileStorage) SaveApplication(ctx context.Context, app *models.Application) error {
	data, err := json.Marshal(app)
	if err != nil {
		return fmt.Errorf("failed to marshal application: %w", err)
	}
	
	path := filepath.Join(s.basePath, "applications", app.ID+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write application file: %w", err)
	}
	
	return nil
}

// GetQuestions returns all questions
func (s *FileStorage) GetQuestions(ctx context.Context) ([]*models.Question, error) {
	dir := filepath.Join(s.basePath, "questions")
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read questions directory: %w", err)
	}
	
	var questions []*models.Question
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		path := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read question file %s: %w", file.Name(), err)
		}
		
		var question models.Question
		if err := json.Unmarshal(data, &question); err != nil {
			return nil, fmt.Errorf("failed to unmarshal question %s: %w", file.Name(), err)
		}
		
		questions = append(questions, &question)
	}
	
	return questions, nil
}

// GetQuestion retrieves a question by ID
func (s *FileStorage) GetQuestion(ctx context.Context, id string) (*models.Question, error) {
	path := filepath.Join(s.basePath, "questions", id+".json")
	
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read question file: %w", err)
	}
	
	var question models.Question
	if err := json.Unmarshal(data, &question); err != nil {
		return nil, fmt.Errorf("failed to unmarshal question: %w", err)
	}
	
	return &question, nil
}

// CreateAssessment creates a new assessment
func (s *FileStorage) CreateAssessment(ctx context.Context, assessment *models.Assessment) error {
	data, err := json.Marshal(assessment)
	if err != nil {
		return fmt.Errorf("failed to marshal assessment: %w", err)
	}
	
	path := filepath.Join(s.basePath, "assessments", assessment.ID+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write assessment file: %w", err)
	}
	
	return nil
}

// GetAssessment retrieves an assessment by ID
func (s *FileStorage) GetAssessment(ctx context.Context, id string) (*models.Assessment, error) {
	path := filepath.Join(s.basePath, "assessments", id+".json")
	
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read assessment file: %w", err)
	}
	
	var assessment models.Assessment
	if err := json.Unmarshal(data, &assessment); err != nil {
		return nil, fmt.Errorf("failed to unmarshal assessment: %w", err)
	}
	
	return &assessment, nil
}

// UpdateAssessment updates an existing assessment
func (s *FileStorage) UpdateAssessment(ctx context.Context, assessment *models.Assessment) error {
	// Check if assessment exists
	path := filepath.Join(s.basePath, "assessments", assessment.ID+".json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("assessment not found: %s", assessment.ID)
	}
	
	// Update assessment
	data, err := json.Marshal(assessment)
	if err != nil {
		return fmt.Errorf("failed to marshal assessment: %w", err)
	}
	
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write assessment file: %w", err)
	}
	
	return nil
}

// ListAssessments returns all assessments for an application
func (s *FileStorage) ListAssessments(ctx context.Context, applicationID string) ([]*models.Assessment, error) {
	dir := filepath.Join(s.basePath, "assessments")
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read assessments directory: %w", err)
	}
	
	var assessments []*models.Assessment
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		
		path := filepath.Join(dir, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read assessment file %s: %w", file.Name(), err)
		}
		
		var assessment models.Assessment
		if err := json.Unmarshal(data, &assessment); err != nil {
			return nil, fmt.Errorf("failed to unmarshal assessment %s: %w", file.Name(), err)
		}
		
		// Filter by applicationID if provided
		if applicationID == "" || assessment.ApplicationID == applicationID {
			assessments = append(assessments, &assessment)
		}
	}
	
	return assessments, nil
}

// SaveReport stores a report
func (s *FileStorage) SaveReport(ctx context.Context, report *models.Report) error {
	data, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}
	
	path := filepath.Join(s.basePath, "reports", report.AssessmentID+".json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}
	
	return nil
}

// GetReport retrieves a report by assessment ID
func (s *FileStorage) GetReport(ctx context.Context, assessmentID string) (*models.Report, error) {
	path := filepath.Join(s.basePath, "reports", assessmentID+".json")
	
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read report file: %w", err)
	}
	
	var report models.Report
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("failed to unmarshal report: %w", err)
	}
	
	return &report, nil
}
