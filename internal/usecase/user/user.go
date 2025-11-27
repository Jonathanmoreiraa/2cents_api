package user

import (
	"context"

	"github.com/jonathanmoreiraa/2cents/internal/api/middleware"
	error_message "github.com/jonathanmoreiraa/2cents/internal/domain/error"
	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	"github.com/jonathanmoreiraa/2cents/internal/domain/repository"
	services "github.com/jonathanmoreiraa/2cents/internal/usecase/user/contract"
	"github.com/jonathanmoreiraa/2cents/pkg/util"
)

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) services.UserUseCase {
	return &userUseCase{
		userRepo: repo,
	}
}

func (useCase *userUseCase) Create(ctx context.Context, user entity.User) (entity.User, error) {
	passwordHashed, err := middleware.HashPassword(user.Password)
	if err != nil {
		return entity.User{}, util.ErrorWithMessage(err, error_message.ErrProccessingPassword)
	}
	user.Password = passwordHashed

	user, err = useCase.userRepo.Create(ctx, user)
	if err != nil {
		return entity.User{}, util.ErrorWithMessage(err, error_message.ErrCreateAccount)
	}

	return user, nil
}

func (useCase *userUseCase) Login(ctx context.Context, email string, password string) (map[string]any, error) {
	user, err := useCase.userRepo.FindByColumn(ctx, "email", email)
	if err != nil {
		return nil, util.ErrorWithMessage(err, error_message.ErrUserNotFound)
	}

	if !middleware.CheckPasswordHash(user.Password, password) {
		return nil, util.ErrorWithMessage(nil, error_message.ErrWrongPassword)
	}

	token, err := middleware.CreateToken(user.ID)
	if err != nil {
		return nil, util.ErrorWithMessage(err, error_message.ErrWrongPassword)
	}

	return map[string]any{
		"token":  token,
		"userId": user.ID,
	}, nil

}

func (useCase *userUseCase) Update(ctx context.Context, user entity.User) error {
	err := useCase.userRepo.Update(ctx, user)
	if err != nil {
		return util.ErrorWithMessage(err, error_message.ErrCreateAccount)
	}

	return nil
}

func (useCase *userUseCase) GetUser(ctx context.Context, id int) (entity.User, error) {
	user, err := useCase.userRepo.FindByID(ctx, id)
	if err != nil {
		return entity.User{}, util.ErrorWithMessage(err, error_message.ErrFindUser)
	}

	return user, nil
}

func (useCase *userUseCase) GetUserByColumn(ctx context.Context, column string, data string) (entity.User, error) {
	user, err := useCase.userRepo.FindByColumn(ctx, column, data)
	if err != nil {
		return entity.User{}, util.ErrorWithMessage(err, error_message.ErrFindUser)
	}

	return user, nil
}

// func (useCase *userUseCase) Delete(ctx context.Context, user entity.User) error {
// 	err := useCase.userRepo.Delete(ctx, user)

// 	return err
// }
