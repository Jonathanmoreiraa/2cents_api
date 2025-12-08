package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	error_message "github.com/jonathanmoreiraa/2cents/internal/domain/error"
	"github.com/jonathanmoreiraa/2cents/internal/domain/model"
	expense_contract "github.com/jonathanmoreiraa/2cents/internal/usecase/expense/contract"
	metric_contract "github.com/jonathanmoreiraa/2cents/internal/usecase/metric/contract"
	revenue_contract "github.com/jonathanmoreiraa/2cents/internal/usecase/revenue/contract"
	saving_contract "github.com/jonathanmoreiraa/2cents/internal/usecase/saving/contract"
)

type DashboardHandler struct {
	expenseUseCase expense_contract.ExpenseUseCase
	revenueUseCase revenue_contract.RevenueUseCase
	savingUseCase  saving_contract.SavingUseCase
	metricUseCase  metric_contract.MetricUseCase
}

type FormatedDate struct {
	dateStart string
	dateEnd   string
}

type GraphicMonthsBar struct {
	Revenues []model.GraphicMonthTotal `json:"revenues"`
	Expenses []model.GraphicMonthTotal `json:"expenses"`
}

func NewDashboardHandler(
	expenseUseCase expense_contract.ExpenseUseCase,
	revenueUseCase revenue_contract.RevenueUseCase,
	savingUseCase saving_contract.SavingUseCase,
	metricUseCase metric_contract.MetricUseCase,
) *DashboardHandler {
	return &DashboardHandler{
		expenseUseCase: expenseUseCase,
		revenueUseCase: revenueUseCase,
		savingUseCase:  savingUseCase,
		metricUseCase:  metricUseCase,
	}
}

func (cr *DashboardHandler) LastExpensesRevenueGraphic(ctx *gin.Context) {
	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrCreateExpense,
			"more_info": "Verifique as informações do usuário logado!",
		})
		return
	}

	end := time.Now()
	sixMonthsAgo := end.AddDate(0, -6, 0)
	dates := GetMonthsFormatedDateStartEnd(sixMonthsAgo, end)

	expenses, err := cr.expenseUseCase.GetExpensesByDates(ctx, userId, dates.dateStart, dates.dateEnd)
	if err != nil || len(expenses) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar as despesas",
		})
		return
	}

	var received bool = true
	revenues, err := cr.revenueUseCase.GetRevenuesByDates(ctx, userId, dates.dateStart, dates.dateEnd, &received)
	if err != nil || len(revenues) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar as despesas",
		})
		return
	}

	graphicBarResponse := GraphicMonthsBar{
		Revenues: revenues,
		Expenses: expenses,
	}

	ctx.JSON(http.StatusOK, graphicBarResponse)
}

func (cr *DashboardHandler) LastMonthSituation(ctx *gin.Context) {
	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrCreateExpense,
			"more_info": "Verifique as informações do usuário logado!",
		})
		return
	}

	y, m, _ := time.Now().Date()
	end := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC)
	oneMonthAgo := time.Now().AddDate(0, 0, 0)
	dates := GetMonthsFormatedDateStartEnd(oneMonthAgo, end)

	expenses, err := cr.expenseUseCase.GetExpensesByDates(ctx, userId, dates.dateStart, dates.dateEnd)
	if err != nil || len(expenses) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar as despesas",
		})
		return
	}

	received := true
	if ctx.Request.URL.Query().Get("received") == "false" {
		received = false
	}
	revenues, err := cr.revenueUseCase.GetRevenuesByDates(ctx, userId, dates.dateStart, dates.dateEnd, &received)
	if err != nil || len(revenues) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar as despesas",
		})
		return
	}

	graphicBarResponse := GraphicMonthsBar{
		Revenues: revenues,
		Expenses: expenses,
	}

	ctx.JSON(http.StatusOK, graphicBarResponse)
}

func (cr *DashboardHandler) ExpensesByCategory(ctx *gin.Context) {
	userId, err := GetUserIdByToken(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":      http.StatusUnprocessableEntity,
			"message":   error_message.ErrCreateExpense,
			"more_info": "Verifique as informações do usuário logado!",
		})
		return
	}

	y, m, _ := time.Now().Date()
	end := time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC)
	oneMonthAgo := time.Now().AddDate(0, 0, 0)
	dates := GetMonthsFormatedDateStartEnd(oneMonthAgo, end)

	expenses, err := cr.expenseUseCase.GetCountCategoryForExpenses(ctx, userId, dates.dateStart, dates.dateEnd)
	if err != nil || len(expenses) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar as despesas por categorias",
		})
		return
	}

	ctx.JSON(http.StatusOK, expenses)
}

func (cr *DashboardHandler) GetTaxes(ctx *gin.Context) {
	lastTaxes, err := cr.metricUseCase.GetLastMetricGraphic(ctx)
	if err != nil || len(lastTaxes) == 0 {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Erro ao encontrar as taxas",
		})
		return
	}

	ctx.JSON(http.StatusOK, lastTaxes)
}

func GetMonthsFormatedDateStartEnd(start time.Time, end time.Time) FormatedDate {
	firstDate := time.Date(
		start.Year(),
		start.Month(),
		1,
		0, 0, 0, 0,
		time.Local,
	)

	layout := "2006-01-02"

	dateStart := firstDate.Format(layout)
	dateEnd := end.Format(layout)

	return FormatedDate{dateStart, dateEnd}
}
