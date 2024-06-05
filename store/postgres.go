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

func (p *Postgres) Create(ctx context.Context, task model.Task) error {
	query := `
        INSERT INTO tasks (id, title, description, status, created_at, completed_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at
    `
	_, err := p.db.ExecContext(ctx, query, task.ID, task.Title, task.Description, task.Status, task.CreatedAt, task.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to insert task: %w", err)
	}
	return nil
}

func (p *Postgres) List(context.Context) ([]model.Task, error) {
	return []model.Task{}, nil
}
func (p *Postgres) GetByID(ctx context.Context) (model.Task, error) {
	id := ctx.Value("taskID")
	task := model.Task{}
	if err := p.db.QueryRow("SELECT * from tasks where id = $1", id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.CompletedAt,
	); err != nil {
		return model.Task{}, err
	}
	return task, nil
}

func (p *Postgres) DeleteByID(ctx context.Context) error {
	id := ctx.Value("taskID")
	res, err := p.db.Exec("DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Not deleted")
	}
	return nil
}

func (p *Postgres) UpdateByID(ctx context.Context, task model.Task) error {
	id := ctx.Value("taskID")
	query := `
        UPDATE tasks
        SET title = $1, description = $2, status = $3, completed_at = $4
        WHERE id = $5
    `
	_, err := p.db.Exec(query, task.Title, task.Description, task.Status, task.CompletedAt, id)
	if err != nil {
		return err
	}
	return nil
}
