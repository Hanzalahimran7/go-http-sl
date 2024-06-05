package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hanzalahimran7/go-http-sl/model"
)

func (a *App) CreateTask(w http.ResponseWriter, r *http.Request) {
	taskRequest := model.TaskRequest{}
	err := json.NewDecoder(r.Body).Decode(&taskRequest)
	if err != nil {
		log.Printf("Bad request: %v+\n", err)
	}
	now := time.Now().UTC()
	if taskRequest.Title == "" || taskRequest.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	task := model.Task{
		ID:          uuid.New(),
		Title:       taskRequest.Title,
		Description: taskRequest.Description,
		CreatedAt:   &now,
		Status:      "incomplete",
		CompletedAt: nil,
	}
	if err := a.Store.Create(r.Context(), task); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func (a *App) FindTaskById(w http.ResponseWriter, r *http.Request) {
	res, err := a.Store.GetByID(r.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(res)
}
