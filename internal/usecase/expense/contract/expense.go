package contract

import (
	"context"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
)

type ExpenseUseCase interface {
	Create(ctx context.Context, expense entity.Expense, multiplePayments bool, numInstallments int, paymentDay int) (entity.Expense, error)
	GetAllExpenses(ctx context.Context, userId int) ([]entity.Expense, error)
	GetExpense(ctx context.Context, id int) (entity.Expense, error)
	GetExpenseBySavingId(ctx context.Context, id int, savingId int) ([]entity.Expense, error)
	GetExpenses(ctx context.Context, filters map[string]any) ([]entity.Expense, error)
	Update(ctx context.Context, expense entity.Expense) error
	Delete(ctx context.Context, expense entity.Expense) error
}
