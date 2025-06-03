package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/wilgnert/webtutoria/internal/database"
)

func (c *Config) SubjectHandler(w http.ResponseWriter, r * http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.listAllSubjects(w, r)
	case http.MethodPost:
		c.createSubject(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func (c *Config) createSubject(w http.ResponseWriter, r * http.Request) {
	var body struct{
		Code        string         `json:"code"`
		Name        string         `json:"name"`
		Description string         `json:"description"`
		Class       string         `json:"class"`
		Categories  []string 			 `json:"categories"`
	}
	if err := DecodeJSON(r.Body, &body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	result, err := c.DB.CreateSubject(r.Context(), database.CreateSubjectParams{
		Code: body.Code,
		Name: body.Name,
		Description: sql.NullString{
			String: body.Description,
			Valid: len(body.Description) > 0,
		},
		Class: body.Class,
	})
	if err != nil {
		if strings.Contains(err.Error(), "subjects_code_key") {
			http.Error(w, "code already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Internal server error while creating subject: %v", err), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error while getting last insert id: %v", err), http.StatusInternalServerError)
		return
	}
	sub, err := c.DB.GetSubjectByID(r.Context(), int32(id))
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error while getting subject by id: %v", err), http.StatusInternalServerError)
		return
	}
	
	for _, cat := range body.Categories {
		var some database.Category
		some, err := c.DB.GetCategoryByName(r.Context(), cat)
		if err != nil {
			_, err := c.DB.CreateCategory(r.Context(), cat)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			some, err = c.DB.GetCategoryByName(r.Context(), cat)	
			if err != nil {
				http.Error(w, fmt.Sprintf("Internal server error while getting category by name: %v", err), http.StatusInternalServerError)
				return
			}
		}
		_, err = c.DB.CreateSubjectCategory(r.Context(), database.CreateSubjectCategoryParams{
			SubjectID: sub.ID,
			CategoryID: some.ID,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	respondWithJSON(w, http.StatusOK, c.populateCategoriesOne(r.Context(), sub))

}
func (c *Config) listAllSubjects(w http.ResponseWriter, r * http.Request) {
	var subjects []database.Subject
	var err error
	if r.URL.Query().Has("class") {
		subjects, err = c.DB.ListSubjectsByClass(r.Context(), r.URL.Query().Get("class"))
	} else {
		subjects, err = c.DB.ListSubjects(r.Context())
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Not found: %v", err), http.StatusNotFound)
	}

	respondWithJSON(w, http.StatusOK, c.populateCategoriesSlice(r.Context(), subjects))
}

type subjectsWithCategories struct{
	ID 					int32          `json:"id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Class       string         `json:"class"`
	Categories  []string 			 `json:"categories"`
}

func (c *Config) populateCategoriesSlice(ctx context.Context, subjects []database.Subject) []subjectsWithCategories {
	var swc []subjectsWithCategories
	for _, s := range subjects {
		swc = append(swc, c.populateCategoriesOne(ctx, s))
	}
	return swc
}
func (c *Config) populateCategoriesOne(ctx context.Context, subject database.Subject) subjectsWithCategories {
	cats, err := c.DB.ListCategoriesBySubjectID(ctx, subject.ID)
	if err != nil {
		fmt.Println(err.Error())
		return subjectsWithCategories{}
	}
	return subjectsWithCategories{
		ID: subject.ID,
		Code: subject.Code,
		Name: subject.Name,
		Description: subject.Description.String,
		Class: subject.Class,
		Categories: cats,
	}
}

func (c *Config) SubjectByIDHandler(w http.ResponseWriter, r * http.Request) {
	id_str := r.PathValue("id")
	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
	}
	switch r.Method {
	case http.MethodGet:
		c.getSubjectById(w, r, int32(id))
	case http.MethodPut:
		c.updateSubjectById(w, r, int32(id))
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (c *Config) getSubjectById(w http.ResponseWriter, r * http.Request, id int32) {
	sub, err := c.DB.GetSubjectByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	respondWithJSON(w, 200, c.populateCategoriesOne(r.Context(), sub))
}

func (c *Config) updateSubjectById(w http.ResponseWriter, r * http.Request, id int32) {
	sub, err := c.DB.GetSubjectByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	var body struct{
		Code        string         `json:"code"`
		Name        string         `json:"name"`
		Description string         `json:"description"`
		Class       string         `json:"class"`
		Categories  []string 			 `json:"categories"`
	}
	if err := DecodeJSON(r.Body, &body); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}
	_, err = c.DB.UpdateSubject(r.Context(), database.UpdateSubjectParams{
		ID: sub.ID,
		Code: body.Code,
		Name: body.Name,
		Description: sql.NullString{
			String: body.Description,
			Valid: len(body.Description) > 0,
		},
		Class: body.Class,
	})
	if err != nil {
		if strings.Contains(err.Error(), "subjects_code_key") {
			http.Error(w, "code already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Internal server error while updating subject: %v", err), http.StatusInternalServerError)
		return
	}
	// delete all subject categories
	err = c.DB.DeleteSubjectCategoriesBySubjectID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error while deleting subject categories: %v", err), http.StatusInternalServerError)
		return
	}
	// create new subject categories
	for _, cat := range body.Categories {
		var some database.Category
		some, err := c.DB.GetCategoryByName(r.Context(), cat)
		if err != nil {
			_, err = c.DB.CreateCategory(r.Context(), cat)
			if err != nil {
				fmt.Println(err.Error())
			}
			some, err = c.DB.GetCategoryByName(r.Context(), cat)
			if err != nil {
				http.Error(w, fmt.Sprintf("Internal server error while getting category by name: %v", err), http.StatusInternalServerError)
				return
			}
		}
		_, err = c.DB.CreateSubjectCategory(r.Context(), database.CreateSubjectCategoryParams{
			SubjectID: sub.ID,
			CategoryID: some.ID,
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	sub, err = c.DB.GetSubjectByID(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Internal server error while getting subject by id: %v", err), http.StatusInternalServerError)
		return
	}
	// get the updated subject
	respondWithJSON(w, 200, c.populateCategoriesOne(r.Context(), sub))
}