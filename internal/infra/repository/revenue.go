package repository

import (
	"context"
	"strings"
	"time"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	interfaces "github.com/jonathanmoreiraa/2cents/internal/domain/repository"
	database "github.com/jonathanmoreiraa/2cents/internal/infra/database/interface"
	"github.com/shopspring/decimal"

	"gorm.io/gorm"
)

type revenueDatabase struct {
	DB *gorm.DB
}

func NewRevenueRepository(Database database.DatabaseProvider) interfaces.RevenueRepository {
	return &revenueDatabase{DB: Database.GetDatabase()}
}

func (database *revenueDatabase) Create(ctx context.Context, revenue entity.Revenue) (entity.Revenue, error) {
	err := database.DB.Create(&revenue).Error
	return revenue, err
}

func (database *revenueDatabase) FindAll(ctx context.Context, userId int) ([]entity.Revenue, error) {
	var revenues []entity.Revenue

	err := database.DB.
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Order("due_date ASC").
		Find(&revenues).Error
	return revenues, err
}

func (database *revenueDatabase) FindByID(ctx context.Context, id int) (entity.Revenue, error) {
	var revenue entity.Revenue

	err := database.DB.First(&revenue, id).Error
	return revenue, err
}

func (database *revenueDatabase) FindByFilter(ctx context.Context, filters map[string]any) (revenues []entity.Revenue, err error) {
	query := database.DB.
		Where("user_id = ?", filters["user_id"]).
		Where("deleted_at IS NULL")

	if description, ok := filters["description"]; ok && description != "" {
		query = query.Where("description LIKE ?", "%"+description.(string)+"%")
	}
	if dateStart, ok := filters["date_start"]; ok && dateStart != "" {
		query = query.Where("due_date >= ?", dateStart)
	}
	if dateEnd, ok := filters["date_end"]; ok && dateEnd != "" {
		query = query.Where("due_date <= ?", dateEnd)
	}
	if min, ok := filters["min"]; ok && !min.(decimal.Decimal).IsZero() {
		query = query.Where("value >= ?", min)
	}
	if max, ok := filters["max"]; ok && !max.(decimal.Decimal).IsZero() {
		query = query.Where("value <= ?", max)
	}
	if status, ok := filters["status"]; ok && status != nil {
		statusStruct := status.(struct {
			Received bool `json:"received"`
			Pending  bool `json:"pending"`
			Overdue  bool `json:"overdue"`
		})
		conditions := []string{}

		if statusStruct.Received {
			conditions = append(conditions, "received = 1")
		}
		if statusStruct.Pending {
			conditions = append(conditions, "(received = 0 AND due_date > '"+time.Now().Format("2006-01-02")+"')")
		}
		if statusStruct.Overdue {
			conditions = append(conditions, "(received = 0 AND due_date < '"+time.Now().Format("2006-01-02")+"')")
		}

		if len(conditions) > 0 {
			query = query.Where(strings.Join(conditions, " OR "))
		}
	}

	query = query.Order("due_date ASC")
	err = query.Find(&revenues).Error

	return revenues, err
}

// TODO: verificar se o gráfico deveria ser retornado com base nas receitas recebidas ou as que ainda não recebi?
func (database *revenueDatabase) FindByDates(ctx context.Context, userId int, dateStart string, dateEnd string, received *bool) ([]entity.GraphicMonthTotal, error) {
	var revenues []entity.GraphicMonthTotal

	query := database.DB.
		Table("revenues").
		Select(`
			CONCAT(
				UPPER(LEFT(DATE_FORMAT(due_date, '%M'), 1)),
				SUBSTRING(DATE_FORMAT(due_date, '%M'), 2)
			) AS month,
			SUM(value) AS total
		`).
		Where("deleted_at IS NULL AND user_id = ? AND due_date BETWEEN ? AND ?", userId, dateStart, dateEnd).
		Group("month").
		Order("MIN(due_date)").
		Exec("SET lc_time_names = 'pt_BR'")

	receivedQuery := "received = 1"
	if received != nil && *received == false {
		receivedQuery = ""
	}

	query = query.Where(receivedQuery)

	err := query.Find(&revenues).Error

	return revenues, err
}

func (database *revenueDatabase) Update(ctx context.Context, revenue entity.Revenue) error {
	err := database.DB.Model(&revenue).Updates(map[string]interface{}{
		"description": revenue.Description,
		"due_date":    revenue.DueDate,
		"received":    revenue.Received,
		"value":       revenue.Value,
		"updated_at":  time.Now(),
	}).Error
	return err
}

func (database *revenueDatabase) Delete(ctx context.Context, revenue entity.Revenue) error {
	err := database.DB.Delete(&revenue).Error
	return err
}
