package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/wilgnert/webtutoria/internal/database"
)

func (c *Config) StudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.listStudents(w, r)
	case http.MethodPost:
		c.createStudent(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
func (c *Config) listStudents(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if r.URL.Query().Has("q") {
		if len(q) < 3 {
			respondWithJSON(w, 200, nil)
			return
		}
		studs, err := c.DB.GetAllStudentsWithNameLike(r.Context(), q+"%")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		respondWithJSON(w, 200, studs)
	} else {
		studs, err := c.DB.GetAllStudents(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting students: %v", err), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, 200, studs)
	}
}
func (c *Config) createStudent(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Name string `json:"name"`
	}
	var newReqBody reqBody
	defer r.Body.Close()
	if err := DecodeJSON(r.Body, &newReqBody); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}
	result, err := c.DB.CreateStudent(r.Context(), newReqBody.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating student: %v", err), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error retrieving last insert ID", http.StatusInternalServerError)
		return
	}
	stud, err := c.DB.GetStudentByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating student: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, 200, stud)
}

func (c *Config) StudentsByIdHandler(w http.ResponseWriter, r *http.Request) {
	id_str := r.PathValue("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		c.getStudentById(w, r, int32(id))
	case http.MethodPost:
		c.updateStudentById(w, r, int32(id))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
func (c *Config) getStudentById(w http.ResponseWriter, r *http.Request, id int32) {
	stud, err := c.DB.GetStudentByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, 200, stud)
}
func (c *Config) updateStudentById(w http.ResponseWriter, r *http.Request, id int32) {
	type reqBody struct {
		Name string `json:"name"`
	}
	var newReqBody reqBody
	defer r.Body.Close()
	if err := DecodeJSON(r.Body, &newReqBody); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}
	_, err := c.DB.UpdateStudent(r.Context(), database.UpdateStudentParams{
		ID:   id,
		Name: newReqBody.Name,
	})
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	stud, err := c.DB.GetStudentByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating student: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, 200, stud)
}
