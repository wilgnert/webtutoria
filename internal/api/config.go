package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/wilgnert/webtutoria/internal/database"
)

type Config struct {
	DB *database.Queries
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	res, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(res)
	return nil
}

func EncodeJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

// Generic function to decode JSON from a Reader into a provided data structure
func DecodeJSON(r io.Reader, v any) error {
	dec := json.NewDecoder(r)
	return dec.Decode(v)
}

func (c *Config) HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset: utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func (c *Config) Init() error {
	if c.DB != nil {
		return nil
	}
	data, err := os.ReadFile("config.json")
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Unmarshal the JSON data into a map[string]any
	var result map[string]any
	err = json.Unmarshal(data, &result)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	db_url := result["db_url"].(string)
	db, err := sql.Open("mysql", db_url)
	if err != nil {
		return fmt.Errorf("error opening mysql with conn string %v: %w", db_url, err) 
	}
	db.Exec("USE webtutoria") // Ensure the database is selected
	if err := db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}
	c.DB = database.New(db)
	time.Sleep(1 * time.Second)

	return nil
}

func (c *Config) ResetHandler(w http.ResponseWriter, r *http.Request) {
	// reset categories, students, tutors and subjects
	err := c.DB.ResetCategories(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reset categories: %v", err), http.StatusInternalServerError)
		return
	}
	err = c.DB.ResetStudents(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reset students: %v", err), http.StatusInternalServerError)
		return
	}
	err = c.DB.ResetTutors(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reset tutors: %v", err), http.StatusInternalServerError)
		return
	}
	err = c.DB.ResetSubjects(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to reset subjects: %v", err), http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
