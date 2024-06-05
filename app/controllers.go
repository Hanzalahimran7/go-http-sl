package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func (a *App) DeleteTaskById(w http.ResponseWriter, r *http.Request) {
	err := a.Store.DeleteByID(r.Context())
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
	w.WriteHeader(http.StatusAccepted)
}

func (a *App) UpdateTaskByID(w http.ResponseWriter, r *http.Request) {
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
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	if body.Status == "completed" {
		now := time.Now().UTC()
		res.CompletedAt = &now
		res.Status = "completed"
	} else if body.Status == "incomplete" {
		res.CompletedAt = nil
		res.Status = "incomplete"
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := a.Store.UpdateByID(r.Context(), res); err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(res)
}

func (a *App) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	log.Println(page, limit)
	if page == "0" {
		page = "1"
	}
	if limit == "0" {
		limit = "10"
	}
	page_number, err := strconv.Atoi(page)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	limit_number, err := strconv.Atoi(limit)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	res, err := a.Store.List(r.Context(), limit_number, page_number)
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
