package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/wilgnert/webtutoria/internal/database"
)

// --- Student Discords Handlers ---

// StudentDiscordsHandler handles requests to /student-discords (GET for list, POST for create).
func (c *Config) StudentDiscordsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.listStudentDiscords(w, r)
	case http.MethodPost:
		c.createStudentDiscord(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// StudentDiscordByIDHandler handles requests to /student-discords/{id} (GET, PUT, DELETE).
func (c *Config) StudentDiscordByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path.
	// Example: /student-discords/123 -> "123"
	idStr := r.URL.Path[len("/student-discords/"):]
	if idStr == "" {
		http.Error(w, "Student ID is required in the URL path", http.StatusBadRequest)
		return
	}

	studentID, err := strconv.ParseInt(idStr, 10, 32) // Parse as int32 for PostgreSQL INT
	if err != nil {
		http.Error(w, "Invalid Student ID format", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		c.getStudentDiscordByStudentID(w, r, int32(studentID))
	case http.MethodDelete:
		c.deleteStudentDiscordByStudentID(w, r, int32(studentID))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// listStudentDiscords handles GET requests to /student-discords.
func (c *Config) listStudentDiscords(w http.ResponseWriter, r *http.Request) {
	// Optional: Allow filtering by discord_id as a query parameter
	discordIDQuery := r.URL.Query().Get("discord_id")

	var studentDiscords []database.Studentdiscord
	var err error

	if discordIDQuery != "" {
		temp, err := c.DB.GetStudentDiscordByDiscordID(r.Context(), discordIDQuery)
		// GetStudentDiscordByDiscordID returns a single item, wrap it in a slice for consistency
		if err == nil {
			studentDiscords = []database.Studentdiscord{temp} // Assuming it returns a slice of 1
		} else if err == sql.ErrNoRows {
			studentDiscords = []database.Studentdiscord{} // No rows found
		} else {
			http.Error(w, fmt.Sprintf("Failed to get student discord by Discord ID: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		studentDiscords, err = c.DB.ListStudentDiscords(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list student discords: %v", err), http.StatusInternalServerError)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, studentDiscords)
}

// createStudentDiscord handles POST requests to /student-discords.
func (c *Config) createStudentDiscord(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		StudentID int32  `json:"student_id"`
		DiscordID string `json:"discord_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := c.DB.CreateStudentDiscord(r.Context(), database.CreateStudentDiscordParams{
		StudentID: reqPayload.StudentID,
		DiscordID: reqPayload.DiscordID,
	})
	if err != nil {
		// Handle unique constraint violation specifically (e.g., if student_id or discord_id already exists)
		// Error message might vary by DB driver, check for specific error codes or strings.
		if err.Error() == "pq: duplicate key value violates unique constraint \"studentdiscords_pkey\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"studentdiscords_discord_id_key\"" {
			http.Error(w, "Student ID or Discord ID already associated", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to create student discord: %v", err), http.StatusInternalServerError)
		return
	}
	newStudentDiscord, err := c.DB.GetStudentDiscordByDiscordID(r.Context(), reqPayload.DiscordID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve newly created student discord: %v", err), http.StatusInternalServerError)
		return
	} 

	respondWithJSON(w, http.StatusCreated, newStudentDiscord)
}

// getStudentDiscordByStudentID handles GET requests to /student-discords/{id}.
func (c *Config) getStudentDiscordByStudentID(w http.ResponseWriter, r *http.Request, studentID int32) {
	studentDiscord, err := c.DB.GetStudentDiscordByStudentID(r.Context(), studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student Discord association not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get student discord: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, studentDiscord)
}

// deleteStudentDiscordByStudentID handles DELETE requests to /student-discords/{id}.
func (c *Config) deleteStudentDiscordByStudentID(w http.ResponseWriter, r *http.Request, studentID int32) {
	err := c.DB.DeleteStudentDiscordByStudentID(r.Context(), studentID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student Discord association not found for deletion", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete student discord: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content for successful deletion
}

// --- Tutor Discords Handlers ---

// TutorDiscordsHandler handles requests to /tutor-discords (GET for list, POST for create).
func (c *Config) TutorDiscordsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.listTutorDiscords(w, r)
	case http.MethodPost:
		c.createTutorDiscord(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// TutorDiscordByIDHandler handles requests to /tutor-discords/{id} (GET, PUT, DELETE).
func (c *Config) TutorDiscordByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path.
	// Example: /tutor-discords/123 -> "123"
	idStr := r.URL.Path[len("/tutor-discords/"):]
	if idStr == "" {
		http.Error(w, "Tutor ID is required in the URL path", http.StatusBadRequest)
		return
	}

	tutorID, err := strconv.ParseInt(idStr, 10, 32) // Parse as int32 for PostgreSQL INT
	if err != nil {
		http.Error(w, "Invalid Tutor ID format", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		c.getTutorDiscordByTutorID(w, r, int32(tutorID))
	case http.MethodDelete:
		c.deleteTutorDiscordByTutorID(w, r, int32(tutorID))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// listTutorDiscords handles GET requests to /tutor-discords.
func (c *Config) listTutorDiscords(w http.ResponseWriter, r *http.Request) {
	// Optional: Allow filtering by discord_id as a query parameter
	discordIDQuery := r.URL.Query().Get("discord_id")

	var tutorDiscords []database.Tutordiscord
	var err error

	if discordIDQuery != "" {
		temp, err := c.DB.GetTutorDiscordByDiscordID(r.Context(), discordIDQuery)
		// GetTutorDiscordByDiscordID returns a single item, wrap it in a slice for consistency
		if err == nil {
			tutorDiscords = []database.Tutordiscord{temp} // Assuming it returns a slice of 1
		} else if err == sql.ErrNoRows {
			tutorDiscords = []database.Tutordiscord{} // No rows found
		} else {
			http.Error(w, fmt.Sprintf("Failed to get tutor discord by Discord ID: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		tutorDiscords, err = c.DB.ListTutorDiscords(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list tutor discords: %v", err), http.StatusInternalServerError)
			return
		}
	}

	respondWithJSON(w, http.StatusOK, tutorDiscords)
}

// createTutorDiscord handles POST requests to /tutor-discords.
func (c *Config) createTutorDiscord(w http.ResponseWriter, r *http.Request) {
	var reqPayload struct {
		TutorID   int32  `json:"tutor_id"`
		DiscordID string `json:"discord_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqPayload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := c.DB.CreateTutorDiscord(r.Context(), database.CreateTutorDiscordParams{
		TutorID:   reqPayload.TutorID,
		DiscordID: reqPayload.DiscordID,
	})
	if err != nil {
		// Handle unique constraint violation specifically
		if err.Error() == "pq: duplicate key value violates unique constraint \"tutordiscords_pkey\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"tutordiscords_discord_id_key\"" {
			http.Error(w, "Tutor ID or Discord ID already associated", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to create tutor discord: %v", err), http.StatusInternalServerError)
		return
	}

	newTutorDiscord, err := c.DB.GetTutorDiscordByDiscordID(r.Context(), reqPayload.DiscordID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve newly created tutor discord: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, newTutorDiscord)
}

// getTutorDiscordByTutorID handles GET requests to /tutor-discords/{id}.
func (c *Config) getTutorDiscordByTutorID(w http.ResponseWriter, r *http.Request, tutorID int32) {
	tutorDiscord, err := c.DB.GetTutorDiscordByTutorID(r.Context(), tutorID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Tutor Discord association not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get tutor discord: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, tutorDiscord)
}

// deleteTutorDiscordByTutorID handles DELETE requests to /tutor-discords/{id}.
func (c *Config) deleteTutorDiscordByTutorID(w http.ResponseWriter, r *http.Request, tutorID int32) {
	err := c.DB.DeleteTutorDiscordByTutorID(r.Context(), tutorID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Tutor Discord association not found for deletion", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete tutor discord: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204 No Content for successful deletion
}