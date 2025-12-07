package contract

import (
	"context"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
)

type SavingUseCase interface {
	Create(ctx context.Context, saving entity.Saving) (entity.Saving, error)
	GetAllSavings(ctx context.Context, userId int) ([]entity.Saving, error)
	GetSaving(ctx context.Context, id int) (entity.Saving, error)
	GetSavings(ctx context.Context, filters map[string]any) ([]entity.Saving, error)
	Update(ctx context.Context, saving entity.Saving) error
	Delete(ctx context.Context, saving entity.Saving) error
}
