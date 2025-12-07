package saving

import (
	"context"

	error_message "github.com/jonathanmoreiraa/2cents/internal/domain/error"
	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	"github.com/jonathanmoreiraa/2cents/internal/domain/repository"
	services "github.com/jonathanmoreiraa/2cents/internal/usecase/saving/contract"
	"github.com/jonathanmoreiraa/2cents/pkg/util"
)

type savingUseCase struct {
	savingRepo repository.SavingRepository
}

func NewSavingUseCase(repo repository.SavingRepository) services.SavingUseCase {
	return &savingUseCase{
		savingRepo: repo,
	}
}

func (useCase *savingUseCase) Create(ctx context.Context, saving entity.Saving) (entity.Saving, error) {
	saving, err := useCase.savingRepo.Create(ctx, saving)
	if err != nil {
		return entity.Saving{}, util.ErrorWithMessage(err, error_message.ErrCreateAccount)
	}

	return saving, nil
}

func (useCase *savingUseCase) GetAllSavings(ctx context.Context, userId int) ([]entity.Saving, error) {
	savings, err := useCase.savingRepo.FindAll(ctx, userId)
	if err != nil {
		return []entity.Saving{}, util.ErrorWithMessage(err, error_message.ErrFindSaving)
	}

	return savings, nil
}

func (useCase *savingUseCase) GetSaving(ctx context.Context, id int) (entity.Saving, error) {
	saving, err := useCase.savingRepo.FindByID(ctx, id)
	if err != nil {
		return entity.Saving{}, util.ErrorWithMessage(err, error_message.ErrFindSaving)
	}

	return saving, nil
}

func (useCase *savingUseCase) GetSavings(ctx context.Context, filters map[string]any) ([]entity.Saving, error) {
	revenue, err := useCase.savingRepo.FindByFilter(ctx, filters)
	if err != nil {
		return []entity.Saving{}, util.ErrorWithMessage(err, error_message.ErrFindSaving)
	}

	return revenue, nil
}

func (useCase *savingUseCase) Update(ctx context.Context, saving entity.Saving) error {
	err := useCase.savingRepo.Update(ctx, saving)
	if err != nil {
		return util.ErrorWithMessage(err, error_message.ErrUpdateSaving)
	}

	return nil
}

func (useCase *savingUseCase) Delete(ctx context.Context, saving entity.Saving) error {
	err := useCase.savingRepo.Delete(ctx, saving)
	if err != nil {
		return util.ErrorWithMessage(err, error_message.ErrUpdateSaving)
	}

	return nil
}
