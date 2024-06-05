package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/hanzalahimran7/go-http-sl/model"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func ConnectToPostgresDb(host, port, user, password, dbname string) *Postgres {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to Database!")
	return &Postgres{db: db}
}

func (p *Postgres) Create(ctx context.Context, task model.Task) (model.Task, error) {
	return task, nil
}
func (p *Postgres) List(context.Context) ([]model.Task, error)     { return []model.Task{}, nil }
func (p *Postgres) GetByID(context.Context) (model.Task, error)    { return model.Task{}, nil }
func (p *Postgres) DeleteByID(context.Context) (string, error)     { return "", nil }
func (p *Postgres) UpdateByID(context.Context) (model.Task, error) { return model.Task{}, nil }
