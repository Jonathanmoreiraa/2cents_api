package repository

import (
	"context"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
)

type ExpenseRepository interface {
	Create(ctx context.Context, expense entity.Expense, multiplePayments bool, numInstallments int, paymentDay int) (entity.Expense, error)
	FindAll(ctx context.Context, userId int) ([]entity.Expense, error)
	FindByID(ctx context.Context, userId int) (entity.Expense, error)
	FindBySavingID(ctx context.Context, userId int, savingID int) ([]entity.Expense, error)
	FindByFilter(ctx context.Context, filters map[string]any) ([]entity.Expense, error)
	Update(ctx context.Context, revenue entity.Expense) error
	Delete(ctx context.Context, revenue entity.Expense) error
}
