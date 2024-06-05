package app

import (
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
	task, err = a.Store.Create(r.Context(), task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
