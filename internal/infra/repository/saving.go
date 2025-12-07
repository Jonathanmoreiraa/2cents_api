package repository

import (
	"context"
	"strings"
	"time"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	interfaces "github.com/jonathanmoreiraa/2cents/internal/domain/repository"
	database "github.com/jonathanmoreiraa/2cents/internal/infra/database/interface"

	"gorm.io/gorm"
)

type savingDatabase struct {
	DB *gorm.DB
}

func NewSavingRepository(Database database.DatabaseProvider) interfaces.SavingRepository {
	return &savingDatabase{DB: Database.GetDatabase()}
}

func (database *savingDatabase) Create(ctx context.Context, saving entity.Saving) (entity.Saving, error) {
	err := database.DB.Create(&saving).Error
	return saving, err
}

func (database *savingDatabase) FindAll(ctx context.Context, userId int) ([]entity.Saving, error) {
	var savings []entity.Saving

	query := database.DB.
		Where("user_id = ?", userId).
		Where("deleted_at IS NULL").
		Where("priority > 0").
		Order("priority ASC").
		Find(&savings)

	err := query.Error
	return savings, err
}

func (database *savingDatabase) FindByID(ctx context.Context, id int) (entity.Saving, error) {
	var saving entity.Saving

	err := database.DB.First(&saving, id).Error
	return saving, err
}

func (database *savingDatabase) FindByFilter(ctx context.Context, filters map[string]any) (savings []entity.Saving, err error) {
	query := database.DB.
		Where("user_id = ?", filters["user_id"]).
		Where("deleted_at IS NULL")

	if description, ok := filters["description"]; ok && description != "" {
		query = query.Where("description LIKE ?", "%"+description.(string)+"%")
	}
	if status, ok := filters["status"]; ok && status != nil {
		statusStruct := status.(struct {
			Completed bool `json:"completed"`
			Pending   bool `json:"pending"`
		})
		conditions := []string{}

		if statusStruct.Completed {
			conditions = append(conditions, "priority = 0")
		}
		if statusStruct.Pending {
			conditions = append(conditions, "priority > 0")
		}

		if len(conditions) > 0 {
			query = query.Where(strings.Join(conditions, " OR "))
		}
	}

	query = query.Order("priority ASC")
	err = query.Find(&savings).Error

	return savings, err
}

func (database *savingDatabase) Update(ctx context.Context, saving entity.Saving) error {
	err := database.DB.Model(&saving).Updates(map[string]interface{}{
		"description":       saving.Description,
		"goal":              saving.Goal,
		"accumulated":       saving.Accumulated,
		"is_emergency_fund": saving.IsEmergencyFund,
		"priority":          saving.Priority,
		"updated_at":        time.Now(),
	}).Error
	return err
}

func (database *savingDatabase) Delete(ctx context.Context, saving entity.Saving) error {
	err := database.DB.Delete(&saving).Error
	return err
}
