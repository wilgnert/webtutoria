package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wilgnert/webtutoria/internal/database"
)

// --- Handler for /students-subjects (List and Create) ---
func (c *Config) StudentSubjectsHandler (w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.listStudentSubjectCompletions(w, r)
	case http.MethodPost:
		c.createStudentSubjectCompletion(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// --- Handler for /students-subjects/{id} (Get, Update, Delete) ---
func (c *Config) StudentSubjectsByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	// Example: /students-subjects/123 -> "123"
	idStr := r.URL.Path[len("/students-subjects/"):]
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 32) // Parse as int32 for PostgreSQL SERIAL
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		c.getStudentSubjectCompletionByID(w, r, int32(id))
	case http.MethodDelete:
		c.deleteStudentSubjectCompletion(w, r, int32(id))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// --- CRUD Implementation Functions ---

// listStudentSubjectCompletions handles GET requests to /students-subjects
func (c *Config) listStudentSubjectCompletions(w http.ResponseWriter, r *http.Request) {
	// You might want to add query parameters for filtering (e.g., by student_id, subject_id)
	// For simplicity, this example lists all or by student_id if provided.
	studentIDStr := r.URL.Query().Get("student_id")
	subjectIDStr := r.URL.Query().Get("subject_id")

	var completions []database.Studentsubjectcompletion // Use sqlc generated type
	var err error

	if studentIDStr != "" {
		studentID, err := strconv.ParseInt(studentIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid student_id format", http.StatusBadRequest)
			return
		}
		completions, err = c.DB.ListStudentSubjectCompletionsByStudent(r.Context(), int32(studentID))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list completions for student: %v", err), http.StatusInternalServerError)
			return
		}
	} else if subjectIDStr != "" {
		subjectID, err := strconv.ParseInt(subjectIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid subject_id format", http.StatusBadRequest)
			return
		}
		completions, err = c.DB.ListStudentSubjectCompletionsBySubject(r.Context(), int32(subjectID))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list completions for subject: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		// List all if no specific filter
		completions, err = c.DB.ListStudentSubjectCompletions(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list all completions: %v", err), http.StatusInternalServerError)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, completions)
}

// createStudentSubjectCompletion handles POST requests to /students-subjects
func (c *Config) createStudentSubjectCompletion(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		StudentID  int32      `json:"student_id"`
		SubjectID  int32      `json:"subject_id"`
	}

	if err := DecodeJSON(r.Body, &reqPayload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var newCompletion database.Studentsubjectcompletion
	var err error
	// Let the database set the timestamp (NOW())
	result, err := c.DB.CreateStudentSubjectCompletion(r.Context(), database.CreateStudentSubjectCompletionParams{
		StudentID: reqPayload.StudentID,
		SubjectID: reqPayload.SubjectID,
	})
	if err != nil {
		// Handle unique constraint violation specifically
		if err.Error() == "pq: duplicate key value violates unique constraint \"studentsubjectcompletion_student_id_subject_id_key\"" {
			http.Error(w, "Student has already completed this subject", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to create student subject completion: %v", err), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve new completion ID: %v", err), http.StatusInternalServerError)
		return
	}
	newCompletion, err = c.DB.GetStudentSubjectCompletionByID(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Newly created Student Subject Completion not found", http.StatusInternalServerError)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to retrieve newly created completion: %v", err), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, newCompletion)
}

// getStudentSubjectCompletionByID handles GET requests to /students-subjects/{id}
func (c *Config) getStudentSubjectCompletionByID(w http.ResponseWriter, r *http.Request, id int32) {
	completion, err :=  c.DB.GetStudentSubjectCompletionByID(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student Subject Completion not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get student subject completion: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, completion)
}

// deleteStudentSubjectCompletion handles DELETE requests to /students-subjects/{id}
func (c *Config) deleteStudentSubjectCompletion(w http.ResponseWriter, r *http.Request, id int32) {
	err := c.DB.DeleteStudentSubjectCompletion(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete student subject completion: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content for successful deletion
}