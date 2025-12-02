package contract

import (
	"context"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
)

type RevenueUseCase interface {
	Create(ctx context.Context, revenue entity.Revenue) (entity.Revenue, error)
	GetAllRevenues(ctx context.Context, userId int) ([]entity.Revenue, error)
	GetRevenue(ctx context.Context, id int) (entity.Revenue, error)
	GetRevenues(ctx context.Context, filters map[string]any) ([]entity.Revenue, error)
	GetRevenuesByDates(ctx context.Context, userId int, dateStart string, dateEnd string) ([]entity.GraphicMonthTotal, error)
	Update(ctx context.Context, revenue entity.Revenue) error
	Delete(ctx context.Context, revenue entity.Revenue) error
}
