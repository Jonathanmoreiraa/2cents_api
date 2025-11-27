package contract

import (
	"context"

	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
)

// TODO: Criar UseCases faltantes
// 1. Recuperar senha
// 2. Visualizar métricas

type UserUseCase interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Login(ctx context.Context, email string, password string) (map[string]any, error)
	Update(ctx context.Context, user entity.User) error
	GetUser(ctx context.Context, id int) (entity.User, error)
	GetUserByColumn(ctx context.Context, column string, data string) (entity.User, error)
}
