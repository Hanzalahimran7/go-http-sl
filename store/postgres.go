package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to Database!")
	return &Postgres{db: db}
}

func (p *Postgres) RunMigrations(migrationPath string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	fmt.Println(currentDir)
	relativePath := "store/migrations/create_tasks_table.sql"
	sqlFile, err := os.ReadFile(filepath.Join("/", currentDir, relativePath))
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}
	_, err = p.db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}
	log.Println("Migration executed successfully!")
	return nil
}

func (p *Postgres) Create(ctx context.Context, task model.Task) (model.Task, error) {
	return task, nil
}
func (p *Postgres) List(context.Context) ([]model.Task, error)     { return []model.Task{}, nil }
func (p *Postgres) GetByID(context.Context) (model.Task, error)    { return model.Task{}, nil }
func (p *Postgres) DeleteByID(context.Context) (string, error)     { return "", nil }
func (p *Postgres) UpdateByID(context.Context) (model.Task, error) { return model.Task{}, nil }
