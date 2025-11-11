// wire.go
//go:build wireinject
// +build wireinject

package di

import (
	"github.com/jonathanmoreiraa/2cents/internal/api/handler"
	"github.com/jonathanmoreiraa/2cents/internal/api/route"
	"github.com/jonathanmoreiraa/2cents/internal/config"
	"github.com/jonathanmoreiraa/2cents/internal/infra/database"
	"github.com/jonathanmoreiraa/2cents/internal/infra/repository"
	"github.com/jonathanmoreiraa/2cents/internal/usecase/category"
	"github.com/jonathanmoreiraa/2cents/internal/usecase/expense"
	"github.com/jonathanmoreiraa/2cents/internal/usecase/revenue"
	"github.com/jonathanmoreiraa/2cents/internal/usecase/saving"
	"github.com/jonathanmoreiraa/2cents/internal/usecase/user"

	"github.com/google/wire"
)

func NewHandlerGroup(
	userHandler *handler.UserHandler,
	revenueHandler *handler.RevenueHandler,
	expenseHandler *handler.ExpenseHandler,
	categoryHandler *handler.CategoryHandler,
	savingHandler *handler.SavingHandler,
	rendimentHandler *handler.RendimentsHandler,
) route.HandlerGroup {
	return route.HandlerGroup{
		UserHandler:       userHandler,
		RevenueHandler:    revenueHandler,
		ExpenseHandler:    expenseHandler,
		CategoryHandler:   categoryHandler,
		SavingHandler:     savingHandler,
		RendimentsHandler: rendimentHandler,
	}
}

func InitializeAPI(cfg config.Config) (*route.ServerHTTP, error) {
	wire.Build(
		database.NewMySqlDatabase,

		repository.NewUserRepository,
		repository.NewRevenueRepository,
		repository.NewExpenseRepository,
		repository.NewCategoryRepository,
		repository.NewSavingRepository,
		repository.NewMetricRepository,

		user.NewUserUseCase,
		revenue.NewRevenueUseCase,
		expense.NewExpenseUseCase,
		category.NewCategoryUseCase,
		saving.NewSavingUseCase,

		handler.NewUserHandler,
		handler.NewRevenueHandler,
		handler.NewExpenseHandler,
		handler.NewCategoryHandler,
		handler.NewSavingHandler,
		handler.NewRendimentsHandler,

		NewHandlerGroup,
		route.NewServerHTTP,
	)
	return &route.ServerHTTP{}, nil
}
