package metric

import (
	"context"

	error_message "github.com/jonathanmoreiraa/2cents/internal/domain/error"
	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	"github.com/jonathanmoreiraa/2cents/internal/domain/repository"
	services "github.com/jonathanmoreiraa/2cents/internal/usecase/metric/contract"
	"github.com/jonathanmoreiraa/2cents/pkg/util"
)

type metricUseCase struct {
	metricRepo repository.MetricRepository
}

func NewMetricUseCase(repo repository.MetricRepository) services.MetricUseCase {
	return &metricUseCase{
		metricRepo: repo,
	}
}

func (useCase *metricUseCase) Create(ctx context.Context, metric entity.Metric) (entity.Metric, error) {
	metric, err := useCase.metricRepo.Create(ctx, metric)
	if err != nil {
		return entity.Metric{}, util.ErrorWithMessage(err, error_message.ErrCreateAccount)
	}

	return metric, nil
}

func (useCase *metricUseCase) GetLastMetric(ctx context.Context, investimentType int) (entity.Metric, error) {
	metric, err := useCase.metricRepo.GetLastMetric(ctx, investimentType)
	if err != nil {
		return entity.Metric{}, util.ErrorWithMessage(err, error_message.ErrFindMetric)
	}

	return metric, nil
}

func (useCase *metricUseCase) GetLastMetricGraphic(ctx context.Context) ([]entity.LastMetrics, error) {
	metric, err := useCase.metricRepo.GetLastMetricGraphic(ctx)
	if err != nil {
		return []entity.LastMetrics{}, util.ErrorWithMessage(err, error_message.ErrFindMetric)
	}

	return metric, nil
}
