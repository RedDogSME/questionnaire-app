package services

import (
	"context"
	"errors"
	"fmt"
	"questionnaire-app/internal/models"
	"questionnaire-app/internal/storage"
	"time"
	
	"github.com/google/uuid"
)

// AssessmentService handles the business logic for assessments
type AssessmentService struct {
	storage storage.Storage
}

// NewAssessmentService creates a new assessment service
func NewAssessmentService(storage storage.Storage) *AssessmentService {
	return &AssessmentService{
		storage: storage,
	}
}

// GetQuestions fetches all available questions
func (s *AssessmentService) GetQuestions(ctx context.Context) ([]*models.Question, error) {
	return s.storage.GetQuestions(ctx)
}

// StartAssessment creates a new assessment for an application
func (s *AssessmentService) StartAssessment(ctx context.Context, applicationID string) (*models.Assessment, error) {
	// Validate application exists
	app, err := s.storage.GetApplication(ctx, applicationID)
	if err != nil {
		return nil, fmt.Errorf("failed to find application: %w", err)
	}
	
	if app == nil {
		return nil, errors.New("application not found")
	}
	
	// Create new assessment
	assessment := &models.Assessment{
		ID:            uuid.NewString(),
		ApplicationID: applicationID,
		CreatedAt:     time.Now().Format(time.RFC3339),
		Answers:       make(map[string]string),
		Status:        "in_progress",
	}
	
	// Save assessment
	if err := s.storage.CreateAssessment(ctx, assessment); err != nil {
		return nil, fmt.Errorf("failed to create assessment: %w", err)
	}
	
	return assessment, nil
}

// GetAssessment retrieves an assessment by ID
func (s *AssessmentService) GetAssessment(ctx context.Context, id string) (*models.Assessment, error) {
	return s.storage.GetAssessment(ctx, id)
}

// SaveAnswer records an answer for a specific question
func (s *AssessmentService) SaveAnswer(ctx context.Context, assessmentID, questionID, optionID string) error {
	// Get assessment
	assessment, err := s.storage.GetAssessment(ctx, assessmentID)
	if err != nil {
		return fmt.Errorf("failed to get assessment: %w", err)
	}
	
	if assessment == nil {
		return errors.New("assessment not found")
	}
	
	// Validate question exists
	question, err := s.storage.GetQuestion(ctx, questionID)
	if err != nil {
		return fmt.Errorf("failed to get question: %w", err)
	}
	
	if question == nil {
		return errors.New("question not found")
	}
	
	// Validate option exists
	optionValid := false
	for _, option := range question.Options {
		if option.ID == optionID {
			optionValid = true
			break
		}
	}
	
	if !optionValid {
		return errors.New("option not found for question")
	}
	
	// Save answer
	assessment.Answers[questionID] = optionID
	
	// Update assessment
	if err := s.storage.UpdateAssessment(ctx, assessment); err != nil {
		return fmt.Errorf("failed to update assessment: %w", err)
	}
	
	return nil
}

// CompleteAssessment marks an assessment as complete and generates a report
func (s *AssessmentService) CompleteAssessment(ctx context.Context, assessmentID string) (*models.Report, error) {
	// Get assessment
	assessment, err := s.storage.GetAssessment(ctx, assessmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assessment: %w", err)
	}
	
	if assessment == nil {
		return nil, errors.New("assessment not found")
	}
	
	// Get all questions to calculate score
	questions, err := s.storage.GetQuestions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}
	
	// Mark assessment as complete
	assessment.Status = "completed"
	if err := s.storage.UpdateAssessment(ctx, assessment); err != nil {
		return nil, fmt.Errorf("failed to update assessment: %w", err)
	}
	
	// Generate report
	report, err := s.generateReport(ctx, assessment, questions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate report: %w", err)
	}
	
	// Save report
	if err := s.storage.SaveReport(ctx, report); err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}
	
	return report, nil
}

// GetReport retrieves a report by assessment ID
func (s *AssessmentService) GetReport(ctx context.Context, assessmentID string) (*models.Report, error) {
	return s.storage.GetReport(ctx, assessmentID)
}

// generateReport creates a suitability report based on assessment answers
func (s *AssessmentService) generateReport(ctx context.Context, 
                                          assessment *models.Assessment, 
                                          questions []*models.Question) (*models.Report, error) {
	// Initialize report
	report := &models.Report{
		AssessmentID:     assessment.ID,
		ApplicationID:    assessment.ApplicationID,
		GeneratedAt:      time.Now().Format(time.RFC3339),
		CategoryScores:   make(map[string]int),
		Recommendations:  []models.Recommendation{},
		Risks:            []models.Risk{},
		ModernizationPlan: []models.ModernizationStep{},
	}
	
	// Calculate scores
	totalScore := 0
	maxScore := 0
	categoryScores := make(map[string]int)
	categoryMaxScores := make(map[string]int)
	
	for _, question := range questions {
		optionID, answered := assessment.Answers[question.ID]
		
		// Add to max possible score
		maxScore += question.Weight * maxOptionPoints(question.Options)
		categoryMaxScores[question.Category] += question.Weight * maxOptionPoints(question.Options)
		
		if answered {
			// Find selected option
			for _, option := range question.Options {
				if option.ID == optionID {
					score := option.Points * question.Weight
					totalScore += score
					categoryScores[question.Category] += score
					break
				}
			}
		}
	}
	
	report.TotalScore = totalScore
	report.MaxPossibleScore = maxScore
	report.CategoryScores = categoryScores
	
	// Add recommendations based on scores (simplified)
	generateRecommendations(report, totalScore, maxScore, categoryScores, categoryMaxScores)
	
	// Add modernization plan
	report.ModernizationPlan = createModernizationPlan(totalScore, maxScore)
	
	return report, nil
}

// maxOptionPoints returns the maximum point value from options
func maxOptionPoints(options []models.Option) int {
	max := 0
	for _, option := range options {
		if option.Points > max {
			max = option.Points
		}
	}
	return max
}

// generateRecommendations adds recommendations and risks to the report
func generateRecommendations(report *models.Report, totalScore, maxScore int, categoryScores, categoryMaxScores map[string]int) {
	overallRatio := float64(totalScore) / float64(maxScore)
	
	// Overall recommendation
	if overallRatio < 0.5 {
		report.Recommendations = append(report.Recommendations, models.Recommendation{
			Category:    "General",
			Description: "Application requires significant modifications for Kubernetes deployment",
			Priority:    "High",
		})
		
		report.Risks = append(report.Risks, models.Risk{
			Category:    "Deployment",
			Description: "Application architecture not suitable for containerization",
			Severity:    "High",
		})
	} else if overallRatio < 0.7 {
		report.Recommendations = append(report.Recommendations, models.Recommendation{
			Category:    "General",
			Description: "Application needs moderate changes to be suitable for Kubernetes",
			Priority:    "Medium",
		})
	} else {
		report.Recommendations = append(report.Recommendations, models.Recommendation{
			Category:    "General",
			Description: "Application is a good candidate for Kubernetes deployment",
			Priority:    "Low",
		})
	}
	
	// Category-specific recommendations
	for category, score := range categoryScores {
		maxScore, ok := categoryMaxScores[category]
		if !ok || maxScore == 0 {
			continue
		}
		
		ratio := float64(score) / float64(maxScore)
		
		// Add category-specific recommendations
		if ratio < 0.5 && category == "Architecture" {
			report.Recommendations = append(report.Recommendations, models.Recommendation{
				Category:    category,
				Description: "Consider refactoring application architecture to be more containerization-friendly",
				Priority:    "High",
			})
			
			report.Risks = append(report.Risks, models.Risk{
				Category:    category,
				Description: "Complex architecture may lead to challenges in containerization",
				Severity:    "High",
			})
		} else if ratio < 0.6 && category == "Persistence" {
			report.Recommendations = append(report.Recommendations, models.Recommendation{
				Category:    category,
				Description: "Review database access patterns for compatibility with Kubernetes",
				Priority:    "Medium",
			})
			
			report.Risks = append(report.Risks, models.Risk{
				Category:    category,
				Description: "Data persistence implementation may cause issues in containerized environment",
				Severity:    "Medium",
			})
		}
	}
}

// createModernizationPlan creates a step-by-step plan based on scores
func createModernizationPlan(totalScore, maxScore int) []models.ModernizationStep {
	ratio := float64(totalScore) / float64(maxScore)
	plan := []models.ModernizationStep{}
	
	// Common steps for all applications
	plan = append(plan, models.ModernizationStep{
		Order:       1,
		Description: "Analyze application dependencies and external integrations",
		Effort:      "Low",
	})
	
	// Add different steps based on score
	if ratio < 0.5 {
		plan = append(plan, []models.ModernizationStep{
			{
				Order:       2,
				Description: "Refactor application architecture for microservices",
				Effort:      "High",
			},
			{
				Order:       3,
				Description: "Implement appropriate data persistence strategy",
				Effort:      "High",
			},
			{
				Order:       4,
				Description: "Create containerization strategy with multiple containers",
				Effort:      "Medium",
			},
		}...)
	} else if ratio < 0.7 {
		plan = append(plan, []models.ModernizationStep{
			{
				Order:       2,
				Description: "Refactor specific components for containerization",
				Effort:      "Medium",
			},
			{
				Order:       3,
				Description: "Adapt data persistence for cloud environment",
				Effort:      "Medium",
			},
		}...)
	}
	
	// Final common steps
	plan = append(plan, []models.ModernizationStep{
		{
			Order:       len(plan) + 1,
			Description: "Containerize application components",
			Effort:      "Medium",
		},
		{
			Order:       len(plan) + 2,
			Description: "Create Kubernetes deployment manifests",
			Effort:      "Medium",
		},
		{
			Order:       len(plan) + 3,
			Description: "Set up CI/CD pipeline for Kubernetes deployment",
			Effort:      "Medium",
		},
		{
			Order:       len(plan) + 4,
			Description: "Implement monitoring and observability",
			Effort:      "Medium",
		},
	}...)
	
	return plan
}
