package handler

import (
	"net/http"
	"strconv"

	error_message "github.com/jonathanmoreiraa/2cents/internal/domain/error"
	"github.com/jonathanmoreiraa/2cents/internal/domain/model"
	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	category_contract "github.com/jonathanmoreiraa/2cents/internal/usecase/category/contract"
	expense_contract "github.com/jonathanmoreiraa/2cents/internal/usecase/expense/contract"
	"github.com/jonathanmoreiraa/2cents/pkg/log"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryUseCase category_contract.CategoryUseCase
	expenseUseCase  expense_contract.ExpenseUseCase
}

func NewCategoryHandler(usecase category_contract.CategoryUseCase, expenseUseCase expense_contract.ExpenseUseCase) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: usecase,
		expenseUseCase:  expenseUseCase,
	}
}

func (cr *CategoryHandler) Create(ctx *gin.Context) {
	var category entity.Category

	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": error_message.ErrCreateCategory,
		})
		log.NewLogger().Error(err)
		return
	}

	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrCreateCategory,
			"more_info": "Verifique as informações do usuário logado!",
		})
		log.NewLogger().Error(err)
		return
	}

	category.UserID = &userId
	category, err = cr.categoryUseCase.Create(ctx.Request.Context(), category)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": error_message.ErrCreateCategory,
		})
		log.NewLogger().Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusOK,
		"message": "Categoria criada com sucesso!",
		"data": gin.H{
			"id": category.ID,
		},
	})
}

func (cr *CategoryHandler) GetAllCategories(ctx *gin.Context) {
	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrFindCategory,
			"more_info": "Verifique as informações do usuário logado!",
		})
		log.NewLogger().Error(err)
		return
	}

	categories, err := cr.categoryUseCase.GetAllCategories(ctx.Request.Context(), userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar as categorias",
		})
		return
	}

	if len(categories) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": error_message.ErrFindCategory,
		})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (cr *CategoryHandler) FindCategory(ctx *gin.Context) {
	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrCreateCategory,
			"more_info": "Verifique as informações do usuário logado!",
		})
		return
	}

	var category model.Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Erro ao localizar categoria, verifique os parâmetros e tente novamente!",
		})
		log.NewLogger().Error(err)
		return
	}

	categories, err := cr.categoryUseCase.GetCategory(ctx.Request.Context(), category.Name, &userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar a receita",
		})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

func (cr *CategoryHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Erro ao identificar o id da categoria",
		})
		return
	}

	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Erro ao apagar a categoria",
		})
		return
	}

	category, err := cr.categoryUseCase.GetCategoryById(ctx, id, userId)
	if err != nil || category.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Categoria não encontrada",
		})
		log.NewLogger().Error(err)
		return
	}

	filterExpense := make(map[string]any)
	filterExpense["category_id"] = id
	filterExpense["user_id"] = userId
	expense, err := cr.expenseUseCase.GetExpenses(ctx.Request.Context(), filterExpense)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Erro ao apagar a categoria",
		})
		log.NewLogger().Error(err)
		return
	}

	if len(expense) > 0 {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Erro ao apagar a categoria, pois existem despesas associadas a ela!",
		})
		log.NewLogger().Error("Erro ao apagar a categoria " + strconv.Itoa(id) + ", pois existem despesas associadas a ela")
		return
	}

	err = cr.categoryUseCase.Delete(ctx.Request.Context(), category)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Erro ao apagar a categoria",
		})
		log.NewLogger().Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Categoria deletada com sucesso!",
	})
}
