package models

// Assessment represents a complete application assessment
type Assessment struct {
	ID            string            `json:"id"`
	ApplicationID string            `json:"applicationId"`
	CreatedAt     string            `json:"createdAt"`
	Answers       map[string]string `json:"answers"` // questionID -> optionID
	Status        string            `json:"status"`
}
