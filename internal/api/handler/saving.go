package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	error_message "github.com/jonathanmoreiraa/2cents/internal/domain/error"
	"github.com/jonathanmoreiraa/2cents/internal/domain/model"
	entity "github.com/jonathanmoreiraa/2cents/internal/domain/model"
	saving_contract "github.com/jonathanmoreiraa/2cents/internal/usecase/saving/contract"
	"github.com/jonathanmoreiraa/2cents/pkg/log"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
)

type SavingHandler struct {
	savingUseCase  saving_contract.SavingUseCase
	expenseHandler *ExpenseHandler
}

type SavingAddRequest struct {
	Description     string          `json:"description"`
	Goal            decimal.Decimal `json:"goal"`
	Accumulated     decimal.Decimal `json:"accumulated"`
	IsEmergencyFund int             `json:"is_emergency_fund"`
	ShouldBeExpense int             `json:"should_be_expense"`
	MonthsToGoal    int             `json:"months_to_goal"`
}

func NewSavingHandler(usecase saving_contract.SavingUseCase, expenseHandler *ExpenseHandler) *SavingHandler {
	return &SavingHandler{
		savingUseCase:  usecase,
		expenseHandler: expenseHandler,
	}
}

func (cr *SavingHandler) Create(ctx *gin.Context) {
	var savingRequest SavingAddRequest

	if err := ctx.ShouldBindJSON(&savingRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": error_message.ErrCreateSaving,
		})
		log.NewLogger().Error(err)
		return
	}

	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrCreateSaving,
			"more_info": "Verifique as informações do usuário logado!",
		})
		return
	}

	saving := entity.Saving{
		Description:     savingRequest.Description,
		Goal:            savingRequest.Goal,
		Accumulated:     savingRequest.Accumulated,
		IsEmergencyFund: savingRequest.IsEmergencyFund,
	}

	saving.UserID = userId
	saving.Priority = 1
	if list, err := cr.savingUseCase.GetAllSavings(ctx.Request.Context(), userId); err == nil && len(list) > 0 {
		maxPriority := 0
		existEmergencyFund := false
		for _, s := range list {
			if s.Priority > maxPriority {
				maxPriority = s.Priority
			}

			if s.IsEmergencyFund > 0 {
				existEmergencyFund = true
			}
		}
		saving.Priority = maxPriority + 1
		if saving.IsEmergencyFund == 1 && existEmergencyFund {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    http.StatusUnprocessableEntity,
				"message": error_message.ErrCreateDuplicateEmergencyFund,
			})
			return
		}
	}

	_, err = cr.savingUseCase.Create(ctx.Request.Context(), saving)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": error_message.ErrCreateSaving,
		})
		return
	}

	if savingRequest.ShouldBeExpense == 1 && savingRequest.MonthsToGoal > 0 {
		var expenseInput ExpenseInput
		monthValue := savingRequest.Goal.Div(decimal.NewFromInt(int64(savingRequest.MonthsToGoal)))

		expenseInput.Description = savingRequest.Description
		expenseInput.Value = monthValue
		expenseInput.MultiplePayments = true
		expenseInput.NumInstallments = savingRequest.MonthsToGoal
		expenseInput.PaymentDay = 1

		t := time.Now()
		expenseInput.DueDate = &t

		categoryId, err := cr.expenseHandler.categoryUseCase.GetCategory(ctx, "Caixinha", nil)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    http.StatusUnprocessableEntity,
				"message": error_message.ErrCreateExpenseFromSaving,
			})
			return
		}
		expenseInput.CategoryID = categoryId[0].ID

		err = cr.expenseHandler.createExpenseInternal(ctx, expenseInput)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"code":    http.StatusUnprocessableEntity,
				"message": error_message.ErrCreateExpenseFromSaving,
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusOK,
		"message": "Caixinha criada com sucesso!",
	})
}

func (cr *SavingHandler) FindAll(ctx *gin.Context) {
	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrCreateSaving,
			"more_info": "Verifique as informações do usuário logado!",
		})
		return
	}

	savings, err := cr.savingUseCase.GetAllSavings(ctx.Request.Context(), userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar as caixinhas",
		})
		return
	}

	if len(savings) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": error_message.ErrFindSaving,
		})
		return
	}

	ctx.JSON(http.StatusOK, savings)
}

func (cr *SavingHandler) FindByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Erro ao identificar o id da caixinha",
		})
		return
	}

	saving, err := cr.savingUseCase.GetSaving(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar a caixinha",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": saving.ID,
	})
}

func (cr *SavingHandler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Erro ao identificar o id da caixinha",
		})
		return
	}

	var saving model.Saving

	if err := ctx.ShouldBindJSON(&saving); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": error_message.ErrUpdateSaving,
		})
		log.NewLogger().Error(err)
		return
	}

	saving.ID = id

	err = cr.savingUseCase.Update(ctx.Request.Context(), saving)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": error_message.ErrUpdateSaving,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Caixinha editada com sucesso!",
		"data": gin.H{
			"id": saving.ID,
		},
	})
}

func (cr *SavingHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Erro ao identificar o id da caixinha",
		})
		return
	}

	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrDeleteSaving,
			"more_info": "Verifique as informações do usuário logado!",
		})
		return
	}

	savingToDelete, err := cr.savingUseCase.GetSaving(ctx.Request.Context(), id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": error_message.ErrFindSaving,
		})
		return
	}

	if savingToDelete.UserID != userId {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Erro ao apagar a caixinha com esse usuário",
		})
		log.NewLogger().Error(
			fmt.Errorf("As caixinha com id %d não está relacionado com o usuário logado com id %d", savingToDelete.UserID, userId),
		)
		return
	}

	deletedPriority := savingToDelete.Priority

	allSavings, err := cr.savingUseCase.GetAllSavings(ctx.Request.Context(), userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": error_message.ErrDeleteSaving,
		})
		return
	}

	err = cr.savingUseCase.Delete(ctx, savingToDelete)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": error_message.ErrDeleteSaving,
		})
		return
	}

	if len(allSavings) > 0 {
		for _, saving := range allSavings {
			if saving.Priority > deletedPriority {
				saving.Priority = saving.Priority - 1
				err = cr.savingUseCase.Update(ctx.Request.Context(), savingToDelete)
				if err != nil {
					log.NewLogger().Error("Erro ao atualizar prioridade do saving:", err)
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Caixinha deletada com sucesso!",
	})
}
