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
	posgresDb := store.ConnectToPostgresDb(os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DB"))
	err := posgresDb.RunMigrations("../store/migrations/create_tasks_table.sql")
	if err != nil {
		log.Fatal(err)
	}
	return &App{
		Router: r,
		Store:  posgresDb,
	}
}

func (a *App) LoadRoutes() {
	a.Router.Use(middleware.Timeout(2 * time.Second))
	a.Router.Use(TimeOutMiddleware)
	a.Router.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {

	})
	a.Router.Post("/tasks", a.CreateTask)
	a.Router.Route("/tasks/{id}", func(r chi.Router) {
		r.Use(TaskIdMiddleWare)
		r.Get("/", a.FindTaskById)
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) })
		r.Put("/", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusCreated) })
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
