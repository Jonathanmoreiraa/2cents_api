package expense

import (
	"context"

	error_message "github.com/jonathanmoreiraa/2cents/internal/domain/error"
	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	"github.com/jonathanmoreiraa/2cents/internal/domain/repository"
	services "github.com/jonathanmoreiraa/2cents/internal/usecase/expense/contract"
	"github.com/jonathanmoreiraa/2cents/pkg/util"
)

type expenseUseCase struct {
	expenseRepo repository.ExpenseRepository
}

func NewExpenseUseCase(repo repository.ExpenseRepository) services.ExpenseUseCase {
	return &expenseUseCase{
		expenseRepo: repo,
	}
}

func (useCase *expenseUseCase) Create(ctx context.Context, expense entity.Expense, multiplePayments bool, numInstallments int, paymentDay int) (entity.Expense, error) {
	expense, err := useCase.expenseRepo.Create(ctx, expense, multiplePayments, numInstallments, paymentDay)
	if err != nil {
		return entity.Expense{}, util.ErrorWithMessage(err, error_message.ErrCreateAccount)
	}

	return expense, nil
}

func (useCase *expenseUseCase) GetAllExpenses(ctx context.Context, userId int) ([]entity.Expense, error) {
	expenses, err := useCase.expenseRepo.FindAll(ctx, userId)
	if err != nil {
		return []entity.Expense{}, util.ErrorWithMessage(err, error_message.ErrFindExpense)
	}

	return expenses, nil
}

func (useCase *expenseUseCase) GetExpense(ctx context.Context, id int) (entity.Expense, error) {
	expense, err := useCase.expenseRepo.FindByID(ctx, id)
	if err != nil {
		return entity.Expense{}, util.ErrorWithMessage(err, error_message.ErrFindExpense)
	}

	return expense, nil
}

func (useCase *expenseUseCase) GetExpenses(ctx context.Context, filters map[string]any) ([]entity.Expense, error) {
	expense, err := useCase.expenseRepo.FindByFilter(ctx, filters)
	if err != nil {
		return []entity.Expense{}, util.ErrorWithMessage(err, error_message.ErrFindExpense)
	}

	return expense, nil
}

func (useCase *expenseUseCase) GetExpenseBySavingId(ctx context.Context, id int, savingId int) ([]entity.Expense, error) {
	expenses, err := useCase.expenseRepo.FindBySavingID(ctx, id, savingId)
	if err != nil {
		return []entity.Expense{}, util.ErrorWithMessage(err, error_message.ErrFindExpense)
	}

	return expenses, nil
}

func (useCase *expenseUseCase) Update(ctx context.Context, expense entity.Expense) error {
	err := useCase.expenseRepo.Update(ctx, expense)
	if err != nil {
		return util.ErrorWithMessage(err, error_message.ErrUpdateExpense)
	}

	return nil
}

func (useCase *expenseUseCase) Delete(ctx context.Context, expense entity.Expense) error {
	err := useCase.expenseRepo.Delete(ctx, expense)
	if err != nil {
		return util.ErrorWithMessage(err, error_message.ErrUpdateExpense)
	}

	return nil
}
