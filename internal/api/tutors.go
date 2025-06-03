package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/wilgnert/webtutoria/internal/database"
)

func (c *Config) TutorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.listTutors(w, r)
	case http.MethodPost:
		c.createTutor(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
func (c *Config) listTutors(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if r.URL.Query().Has("q") {
		if len(q) < 3 {
			respondWithJSON(w, 200, nil)
			return
		}
		studs, err := c.DB.GetAllTutorsWithNameLike(r.Context(), q+"%")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		respondWithJSON(w, 200, studs)
	} else {
		studs, err := c.DB.GetAllTutors(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting tutors: %v", err), http.StatusInternalServerError)
			return
		}
		respondWithJSON(w, 200, studs)
	}
}
func (c *Config) createTutor(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Name string `json:"name"`
		RoleID string `json:"role_id"`
		ChannelID string `json:"channel_id"`
	}
	var newReqBody reqBody
	defer r.Body.Close()
	if err := DecodeJSON(r.Body, &newReqBody); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}
	result, err := c.DB.CreateTutor(r.Context(), database.CreateTutorParams{
		Name: newReqBody.Name,
		RoleID: newReqBody.RoleID,
		ChannelID: newReqBody.ChannelID,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating tutor: %v", err), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error retrieving last insert ID", http.StatusInternalServerError)
		return
	}
	stud, err := c.DB.GetTutorByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating tutor: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, 200, stud)
}

func (c *Config) TutorsByIdHandler(w http.ResponseWriter, r *http.Request) {
	id_str := r.PathValue("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		c.getTutorById(w, r, int32(id))
	case http.MethodPut:
		c.updateTutorById(w, r, int32(id))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
func (c *Config) getTutorById(w http.ResponseWriter, r *http.Request, id int32) {
	stud, err := c.DB.GetTutorByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, 200, stud)
}
func (c *Config) updateTutorById(w http.ResponseWriter, r *http.Request, id int32) {
	type reqBody struct {
		Name string `json:"name"`
		RoleID string `json:"role_id,omitempty"`
		ChannelID string `json:"channel_id,omitempty"`
	}
	var newReqBody reqBody
	defer r.Body.Close()
	if err := DecodeJSON(r.Body, &newReqBody); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
		return
	}
	_, err := c.DB.UpdateTutor(r.Context(), database.UpdateTutorParams{
		ID:   id,
		Name: newReqBody.Name,
		RoleID: newReqBody.RoleID,
		ChannelID: newReqBody.ChannelID,
	})
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	stud, err := c.DB.GetTutorByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating tutor: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, 200, stud)
}
