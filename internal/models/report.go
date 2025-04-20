package models

// Report represents the generated suitability report
type Report struct {
	AssessmentID      string             `json:"assessmentId"`
	ApplicationID     string             `json:"applicationId"`
	GeneratedAt       string             `json:"generatedAt"`
	TotalScore        int                `json:"totalScore"`
	MaxPossibleScore  int                `json:"maxPossibleScore"`
	CategoryScores    map[string]int     `json:"categoryScores"`
	Recommendations   []Recommendation   `json:"recommendations"`
	Risks             []Risk             `json:"risks"`
	ModernizationPlan []ModernizationStep `json:"modernizationPlan"`
}

// Recommendation provides guidance based on assessment answers
type Recommendation struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

// Risk represents potential migration challenges
type Risk struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
}

// ModernizationStep defines a step in the adoption plan
type ModernizationStep struct {
	Order       int    `json:"order"`
	Description string `json:"description"`
	Effort      string `json:"effort"`
}
