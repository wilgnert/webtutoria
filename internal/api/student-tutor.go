package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/wilgnert/webtutoria/internal/database"
)

func (c *Config) StudentTutorHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.listStudentTutors(w, r)
	case http.MethodPost:
		c.createStudentTutor(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (c *Config) StudentTutorByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path (e.g., /student-tutors/123)
	idStr := r.URL.Path[len("/student-tutors/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		c.getStudentTutorByID(w, r, id)
	case http.MethodDelete:
		c.deleteStudentTutor(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (c *Config) listStudentTutors(w http.ResponseWriter, r *http.Request) {
	studentIDStr := r.URL.Query().Get("student_id")
	tutorIDStr := r.URL.Query().Get("tutor_id")
	// Example: Get all student tutors (you might want to add pagination or filtering)
	var studentTutors []database.Studenttutor
	var err error

	if tutorIDStr != "" {
		tutor_id, err := strconv.ParseInt(tutorIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid student_id format", http.StatusBadRequest)
			return
		}
		studentTutors, err = c.DB.ListStudentTutorsByTutor(r.Context(), int32(tutor_id))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list student tutors: %v", err), http.StatusInternalServerError)
			return
		}
	} else if studentIDStr != "" {
		studentID, err := strconv.ParseInt(studentIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid student_id format", http.StatusBadRequest)
			return
		}
		studentTutors, err = c.DB.ListStudentTutorsByStudent(r.Context(), int32(studentID))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list student tutors: %v", err), http.StatusInternalServerError)
			return
		}
	} else {
		studentTutors, err = c.DB.ListStudentTutors(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list student tutors: %v", err), http.StatusInternalServerError)
			return
		}
	}
	respondWithJSON(w, http.StatusOK, studentTutors)
}

func (c *Config) createStudentTutor(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StudentID int32 `json:"student_id"`
		TutorID   int32 `json:"tutor_id"`
	}
	if err := DecodeJSON(r.Body, &req); err != nil {
		fmt.Println("Error decoding JSON:", err, r.Body)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	result, err := c.DB.CreateStudentTutor(r.Context(), database.CreateStudentTutorParams{
		StudentID: req.StudentID,
		TutorID:   req.TutorID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create student tutor: %v", err), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve last insert ID: %v", err), http.StatusInternalServerError)
		return
	}
	res, err := c.DB.GetStudentTutorByID(r.Context(), int32(id))
	if err != nil {
		if err.Error() == "sql: no rows in result" {
			http.Error(w, "Student Tutor not found after creation", http.StatusInternalServerError)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get student tutor after creation: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusCreated, res)
}

func (c *Config) getStudentTutorByID(w http.ResponseWriter, r *http.Request, id int64) {
	studentTutor, err := c.DB.GetStudentTutorByID(r.Context(), int32(id))
	if err != nil {
		if err.Error() == "sql: no rows in result" {
			http.Error(w, "Student Tutor not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get student tutor: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, studentTutor)
}

func (c *Config) deleteStudentTutor(w http.ResponseWriter, r *http.Request, id int64) {
	err := c.DB.DeleteStudentTutorByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete student tutor: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent) // Successful deletion, no content to return
}