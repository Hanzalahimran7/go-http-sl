package store

import (
	"context"

	"github.com/hanzalahimran7/go-http-sl/model"
)

type Store interface {
	Create(context.Context, model.Task) error
	List(context.Context) ([]model.Task, error)
	GetByID(context.Context) (model.Task, error)
	DeleteByID(context.Context) error
	UpdateByID(context.Context) error
}
