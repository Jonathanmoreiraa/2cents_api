package handler

import (
	"fmt"
	"net/http"
	"strconv"

	error_message "github.com/jonathanmoreiraa/2cents/internal/domain/error"
	"github.com/jonathanmoreiraa/2cents/internal/domain/model"
	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	revenue_contract "github.com/jonathanmoreiraa/2cents/internal/usecase/revenue/contract"
	"github.com/jonathanmoreiraa/2cents/pkg/log"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

type RevenueHandler struct {
	revenueUseCase revenue_contract.RevenueUseCase
}

type RevenueFilters struct {
	Description string          `json:"description"`
	Min         decimal.Decimal `json:"min"`
	Max         decimal.Decimal `json:"max"`
	Status      struct {
		Received bool `json:"received"`
		Pending  bool `json:"pending"`
		Overdue  bool `json:"overdue"`
	} `json:"status"`
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
}

func NewRevenueHandler(usecase revenue_contract.RevenueUseCase) *RevenueHandler {
	return &RevenueHandler{
		revenueUseCase: usecase,
	}
}

func (cr *RevenueHandler) Create(ctx *gin.Context) {
	var revenue entity.Revenue

	if err := ctx.ShouldBindJSON(&revenue); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": error_message.ErrCreateRevenue,
		})
		log.NewLogger().Error(err)
		return
	}

	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrCreateRevenue,
			"more_info": "Verifique as informações do usuário logado!",
		})
		return
	}

	revenue.UserID = userId
	_, err = cr.revenueUseCase.Create(ctx.Request.Context(), revenue)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": error_message.ErrCreateRevenue,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusOK,
		"message": "Receita criada com sucesso!",
	})
}

func (cr *RevenueHandler) FindAll(ctx *gin.Context) {
	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrCreateRevenue,
			"more_info": "Verifique as informações do usuário logado!",
		})
		return
	}

	revenues, err := cr.revenueUseCase.GetAllRevenues(ctx.Request.Context(), userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar as receitas",
		})
		return
	}

	if len(revenues) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": error_message.ErrFindRevenue,
		})
		return
	}

	ctx.JSON(http.StatusOK, revenues)
}

func (cr *RevenueHandler) FindByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Erro ao identificar o id da receita",
		})
		return
	}

	revenue, err := cr.revenueUseCase.GetRevenue(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar a receita",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          revenue.ID,
		"description": revenue.Description,
		"due_date":    revenue.DueDate,
		"received":    revenue.Received,
		"value":       revenue.Value,
	})
}

func (cr *RevenueHandler) FindByFilters(ctx *gin.Context) {
	var revenuesFilter RevenueFilters
	if err := ctx.ShouldBindJSON(&revenuesFilter); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": error_message.ErrFindRevenue,
		})
		log.NewLogger().Error(err)
		return
	}

	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrCreateRevenue,
			"more_info": "Verifique as informações do usuário logado!",
		})
		return
	}

	filters := make(map[string]any)
	filters["description"] = revenuesFilter.Description
	filters["date_start"] = revenuesFilter.DateStart
	filters["date_end"] = revenuesFilter.DateEnd
	filters["min"] = revenuesFilter.Min
	filters["max"] = revenuesFilter.Max
	filters["status"] = revenuesFilter.Status
	filters["user_id"] = userId

	revenues, err := cr.revenueUseCase.GetRevenues(ctx.Request.Context(), filters)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": error_message.ErrFindRevenue,
		})
		return
	}

	if len(revenues) <= 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Nenhuma receita encontrada!",
		})
		return
	}

	ctx.JSON(http.StatusOK, revenues)
}

func (cr *RevenueHandler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Erro ao identificar o id da receita",
		})
		return
	}

	var revenue model.Revenue

	if err := ctx.ShouldBindJSON(&revenue); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Erro ao realizar o login!",
		})
		log.NewLogger().Error(err)
		return
	}

	revenue.ID = id

	err = cr.revenueUseCase.Update(ctx.Request.Context(), revenue)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Erro ao criar a conta!",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Receita editada com sucesso!",
		"data": gin.H{
			"id": revenue.ID,
		},
	})
}

func (cr *RevenueHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Erro ao identificar o id da receita",
		})
		return
	}

	requestCtx := ctx.Request.Context()
	revenue, err := cr.revenueUseCase.GetRevenue(requestCtx, id)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
	}

	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Erro ao apagar a receita",
		})
		return
	}
	if revenue.UserID != userId {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Erro ao apagar a receita com esse usuário",
		})
		log.NewLogger().Error(fmt.Errorf("As receita com id %d não está relacionado com o usuário logado com id %d", revenue.UserID, userId))
		return
	}

	err = cr.revenueUseCase.Delete(ctx, revenue)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Erro ao apagar a receita",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Receita deletada com sucesso!",
	})
}
