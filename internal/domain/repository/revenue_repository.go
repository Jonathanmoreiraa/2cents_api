package repository

import (
	"context"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
)

type RevenueRepository interface {
	Create(ctx context.Context, revenue entity.Revenue) (entity.Revenue, error)
	FindAll(ctx context.Context, userId int) ([]entity.Revenue, error)
	FindByID(ctx context.Context, userId int) (entity.Revenue, error)
	FindByFilter(ctx context.Context, filters map[string]any) ([]entity.Revenue, error)
	FindByDates(ctx context.Context, userId int, dateStart string, dateEnd string) ([]entity.GraphicMonthTotal, error)
	Update(ctx context.Context, revenue entity.Revenue) error
	Delete(ctx context.Context, revenue entity.Revenue) error
}
