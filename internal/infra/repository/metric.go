package repository

import (
	"context"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	interfaces "github.com/jonathanmoreiraa/2cents/internal/domain/repository"
	database "github.com/jonathanmoreiraa/2cents/internal/infra/database/interface"

	"gorm.io/gorm"
)

type metricDatabase struct {
	DB *gorm.DB
}

func NewMetricRepository(Database database.DatabaseProvider) interfaces.MetricRepository {
	return &metricDatabase{DB: Database.GetDatabase()}
}

func (database *metricDatabase) Create(ctx context.Context, metric entity.Metric) (entity.Metric, error) {
	err := database.DB.Create(&metric).Error
	return metric, err
}

func (database *metricDatabase) GetInvestimentsTypes(ctx context.Context) ([]entity.InvestimentType, error) {
	var investimentTypes []entity.InvestimentType

	err := database.DB.
		Where("deleted_at IS NULL").
		Find(&investimentTypes).Error
	return investimentTypes, err
}

func (database *metricDatabase) GetLastMetric(ctx context.Context, investimentType int) (entity.Metric, error) {
	var metric entity.Metric

	err := database.DB.
		Where("investiment_type_id = ?", investimentType).
		Order("date DESC").
		First(&metric).Error
	return metric, err
}

func (database *metricDatabase) GetLastMetricGraphic(ctx context.Context) ([]entity.LastMetrics, error) {
	var metric []entity.LastMetrics

	subQuery := database.DB.
		Table("metrics").
		Select("investiment_type_id, MAX(created_at) as last_created").
		Group("investiment_type_id")

	err := database.DB.
		Table("metrics").
		Select("investiment_types.name, metrics.value").
		Joins("JOIN investiment_types ON investiment_types.id = metrics.investiment_type_id").
		Joins("JOIN (?) last ON last.investiment_type_id = metrics.investiment_type_id AND last.last_created = metrics.created_at", subQuery).
		Order("metrics.created_at DESC").
		Limit(2).
		Find(&metric).Error
	return metric, err
}

func (database *metricDatabase) Update(ctx context.Context, metric entity.Metric) error {
	err := database.DB.Model(&metric).Updates(map[string]interface{}{
		"value":      metric.Value,
		"updated_at": metric.UpdatedAt,
	}).Error
	return err
}

func (database *metricDatabase) Delete(ctx context.Context, metric entity.Metric) error {
	err := database.DB.Delete(&metric).Error
	return err
}
