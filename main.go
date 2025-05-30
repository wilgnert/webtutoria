package main

import (
	"net/http"
	"os"

	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wilgnert/webtutoria/internal/api"
)

func main() {
	cfg := api.Config{}
	err := cfg.Init()
	if err != nil {
		fmt.Printf("error reading config: %v\n", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz",	cfg.HealthzHandler)
	mux.HandleFunc("/reset", cfg.ResetHandler)
	mux.HandleFunc("/subjects", cfg.SubjectHandler)
	mux.HandleFunc("/subjects/{id}", cfg.SubjectByIDHandler)
	mux.HandleFunc("/tutors", cfg.TutorsHandler)
	mux.HandleFunc("/tutors/{id}", cfg.TutorsByIdHandler)
	mux.HandleFunc("/students", cfg.StudentsHandler)
	mux.HandleFunc("/students/{id}", cfg.StudentsByIdHandler)
	mux.HandleFunc("/students-tutors", cfg.StudentTutorHandler)
	mux.HandleFunc("/students-tutors/{id}", cfg.StudentTutorByIDHandler)
	mux.HandleFunc("/students-subjects", cfg.StudentSubjectsHandler)
	mux.HandleFunc("/students-subjects/{id}", cfg.StudentSubjectsByIDHandler)
	
	mux.HandleFunc("/student-discords", cfg.StudentDiscordsHandler)
	mux.HandleFunc("/student-discords/{id}", cfg.StudentDiscordByIDHandler)
	mux.HandleFunc("/tutor-discords", cfg.TutorDiscordsHandler)
	mux.HandleFunc("/tutor-discords/{id}", cfg.TutorDiscordByIDHandler)

	fmt.Printf("listening on http://localhost:8080\n")
	http.ListenAndServe(":8080", mux)
	
}

