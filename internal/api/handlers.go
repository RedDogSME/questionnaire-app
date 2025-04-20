package api

import (
	"encoding/json"
	"net/http"
	"questionnaire-app/internal/models"
	"questionnaire-app/internal/services"
	
	"github.com/gorilla/mux"
)

// Handler manages HTTP requests
type Handler struct {
	assessmentService *services.AssessmentService
}

// NewHandler creates a new API handler
func NewHandler(assessmentService *services.AssessmentService) *Handler {
	return &Handler{
		assessmentService: assessmentService,
	}
}

// GetQuestions returns all assessment questions
func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.assessmentService.GetQuestions(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get questions: "+err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, questions)
}

// StartAssessment creates a new assessment
func (h *Handler) StartAssessment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ApplicationID string `json:"applicationId"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	
	if req.ApplicationID == "" {
		respondWithError(w, http.StatusBadRequest, "Application ID is required")
		return
	}
	
	assessment, err := h.assessmentService.StartAssessment(r.Context(), req.ApplicationID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to start assessment: "+err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusCreated, assessment)
}

// GetAssessment returns an assessment by ID
func (h *Handler) GetAssessment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assessmentID := vars["assessmentId"]
	
	assessment, err := h.assessmentService.GetAssessment(r.Context(), assessmentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get assessment: "+err.Error())
		return
	}
	
	if assessment == nil {
		respondWithError(w, http.StatusNotFound, "Assessment not found")
		return
	}
	
	respondWithJSON(w, http.StatusOK, assessment)
}

// SaveAnswer saves an answer for a question
func (h *Handler) SaveAnswer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assessmentID := vars["assessmentId"]
	
	var req struct {
		QuestionID string `json:"questionId"`
		OptionID   string `json:"optionId"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}
	
	if req.QuestionID == "" || req.OptionID == "" {
		respondWithError(w, http.StatusBadRequest, "Question ID and Option ID are required")
		return
	}
	
	if err := h.assessmentService.SaveAnswer(r.Context(), assessmentID, req.QuestionID, req.OptionID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save answer: "+err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// CompleteAssessment finishes an assessment and generates a report
func (h *Handler) CompleteAssessment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assessmentID := vars["assessmentId"]
	
	report, err := h.assessmentService.CompleteAssessment(r.Context(), assessmentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to complete assessment: "+err.Error())
		return
	}
	
	respondWithJSON(w, http.StatusOK, report)
}

// GetReport returns a generated report
func (h *Handler) GetReport(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assessmentID := vars["assessmentId"]
	
	report, err := h.assessmentService.GetReport(r.Context(), assessmentID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get report: "+err.Error())
		return
	}
	
	if report == nil {
		respondWithError(w, http.StatusNotFound, "Report not found")
		return
	}
	
	respondWithJSON(w, http.StatusOK, report)
}

// Helper functions for HTTP responses

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	if payload != nil {
		response, _ := json.Marshal(payload)
		w.Write(response)
	}
}
