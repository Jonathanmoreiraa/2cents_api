package repository

import (
	"context"
	"fmt"
	"time"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	interfaces "github.com/jonathanmoreiraa/2cents/internal/domain/repository"
	database "github.com/jonathanmoreiraa/2cents/internal/infra/database/interface"

	"gorm.io/gorm"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepository(Database database.DatabaseProvider) interfaces.UserRepository {
	return &userDatabase{DB: Database.GetDatabase()}
}

func (database *userDatabase) FindByID(ctx context.Context, id int) (entity.User, error) {
	var user entity.User

	err := database.DB.First(&user, id).Error
	return user, err
}

func (database *userDatabase) FindByColumn(ctx context.Context, column string, data any) (entity.User, error) {
	var user entity.User

	err := database.DB.First(&user, fmt.Sprintf("%s = ?", column), data).Error
	return user, err
}

func (database *userDatabase) Create(ctx context.Context, user entity.User) (entity.User, error) {
	err := database.DB.Create(&user).Error
	return user, err
}

func (database *userDatabase) Update(ctx context.Context, user entity.User) error {
	err := database.DB.Model(&user).Updates(map[string]interface{}{
		"name":       user.Name,
		"birth_date": user.BirthDate,
		"email":      user.Email,
		"updated_at": time.Now(),
	}).Error
	return err
}

func (database *userDatabase) Delete(ctx context.Context, user entity.User) error {
	err := database.DB.Delete(&user).Error
	return err
}
