package contract

import (
	"context"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
)

type MetricUseCase interface {
	Create(ctx context.Context, metric entity.Metric) (entity.Metric, error)
	GetLastMetric(ctx context.Context, investimentType int) (entity.Metric, error)
	GetLastMetricGraphic(ctx context.Context) ([]entity.LastMetrics, error)
}
