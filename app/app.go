package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hanzalahimran7/go-http-sl/store"
)

type App struct {
	Router *chi.Mux
	Store  store.Store
}

func Initialise() *App {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	db := store.ConnectToPostgresDb(os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DB"))
	err := db.RunMigrations("../store/migrations/create_tasks_table.sql")
	if err != nil {
		log.Fatal(err)
	}
	return &App{
		Router: r,
		Store:  db,
	}
}

func (a *App) LoadRoutes() {
	a.Router.Use(middleware.Timeout(2 * time.Second))
	a.Router.Use(TimeOutMiddleware)
	a.Router.Get("/tasks", a.GetAllTasks)
	a.Router.Post("/tasks", a.CreateTask)
	a.Router.Route("/tasks/{id}", func(r chi.Router) {
		r.Use(TaskIdMiddleWare)
		r.Get("/", a.FindTaskById)
		r.Delete("/", a.DeleteTaskById)
		r.Patch("/", a.UpdateTaskByID)
	})
}

func TimeOutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respch := make(chan int)
		go func() {
			next.ServeHTTP(w, r)
			respch <- 1
		}()
		select {
		case <-respch:
			return
		case <-r.Context().Done():
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}
	})
}

func TaskIdMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "id")
		ctx := context.WithValue(r.Context(), "taskID", taskID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
