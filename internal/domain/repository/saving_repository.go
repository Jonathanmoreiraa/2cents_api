package repository

import (
	"context"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
)

type SavingRepository interface {
	Create(ctx context.Context, saving entity.Saving) (entity.Saving, error)
	FindAll(ctx context.Context, userId int) ([]entity.Saving, error)
	FindByID(ctx context.Context, userId int) (entity.Saving, error)
	FindByFilter(ctx context.Context, filters map[string]any) ([]entity.Saving, error)
	Update(ctx context.Context, saving entity.Saving) error
	Delete(ctx context.Context, saving entity.Saving) error
}
