package handler

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	error_message "github.com/jonathanmoreiraa/2cents/internal/domain/error"
	"github.com/jonathanmoreiraa/2cents/internal/domain/repository"
	"github.com/jonathanmoreiraa/2cents/pkg/log"
)

type RendimentsHandler struct {
	metricRepository repository.MetricRepository
}

type SimulationRequest struct {
	InitialValue float64  `json:"initial_value"`
	Months       int      `json:"months"`
	Accumulated  *float64 `json:"accumulated"`
}

type SimulationResponse struct {
	CdbValue      float64 `json:"cdb"`
	PoupancaValue float64 `json:"poupanca"`
}

const (
	CDI_TYPE              = 1
	SELIC                 = 2
	CDI_PERCENT           = float64(100.00)
	BUSINESS_DAYS         = 21.0
	UNTIL_180_DAYS        = 22.5
	UNTIL_360_DAYS        = 20
	UNTIL_720_DAYS        = 17.5
	MORE_THAN_720_DAYS    = 15
	ONE_YEAR_DAYS         = 365
	MONTH_TAX_DEFAULT     = 0.005
	SELIC_LESS_THAN_EIGHT = 0.085
)

func NewRendimentsHandler(metricRepo repository.MetricRepository) *RendimentsHandler {
	return &RendimentsHandler{
		metricRepository: metricRepo,
	}
}

func (rh *RendimentsHandler) SimulateAllRendiments(ctx *gin.Context) {
	var request SimulationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Invalid simulation parameters",
		})
		log.NewLogger().Error(err)
		return
	}

	cdb, err := rh.SimulateCDB(ctx, request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": error_message.ErrSimulateCdb,
		})
		return
	}

	poupanca, err := rh.SimulatePoupanca(ctx, request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": error_message.ErrSimulatePoupanca,
		})
		return
	}

	var simulations SimulationResponse
	simulations.CdbValue = cdb
	simulations.PoupancaValue = poupanca

	ctx.JSON(http.StatusOK, simulations)
}

func (rh *RendimentsHandler) SimulateMonthlyRendiments(ctx *gin.Context) {
	var request SimulationRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": "Invalid simulation parameters",
		})
		log.NewLogger().Error(err)
		return
	}

	request.InitialValue = request.InitialValue / float64(request.Months)
	var cdbFinalValue float64
	monthsToCalculate := request.Months
	for i := 0; i < monthsToCalculate; i++ {
		request.Months = i
		cdbPerMonth, err := rh.SimulateCDB(ctx, request)
		if err != nil {
			continue
		}
		cdbFinalValue += cdbPerMonth
	}
	request.Months = monthsToCalculate

	poup, err := rh.SimulatePoupancaMonthly(ctx, request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"code":    http.StatusUnprocessableEntity,
			"message": error_message.ErrSimulatePoupanca,
		})
		log.NewLogger().Error(err)
	}

	if request.Accumulated != nil {
		request.InitialValue = *request.Accumulated
		accumulatedCdb, err := rh.SimulateCDB(ctx, request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": error_message.ErrSimulateCdb,
			})
			return
		}
		cdbFinalValue += accumulatedCdb
	}

	var simulations SimulationResponse
	simulations.CdbValue = math.Round(cdbFinalValue*100) / 100
	simulations.PoupancaValue = poup

	ctx.JSON(http.StatusOK, simulations)
}

func (rh *RendimentsHandler) SimulateCDB(ctx *gin.Context, request SimulationRequest) (float64, error) {
	metric, err := rh.metricRepository.GetLastMetric(ctx, CDI_TYPE)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": error_message.ErrFindMetric,
		})
		return 0.0, nil
	}

	dailyCdi := metric.Value.InexactFloat64() / 100.0

	initialValue := request.InitialValue
	months := request.Months

	businessDays := months * int(BUSINESS_DAYS)
	totalDays := (months / 12.0) * ONE_YEAR_DAYS

	adjustedDailyCdi := dailyCdi * (CDI_PERCENT / 100.0)

	finalAmount := math.Floor(initialValue) * math.Pow(1+adjustedDailyCdi, float64(businessDays))
	profit := finalAmount - initialValue

	var taxRate float64
	switch {
	case totalDays < 180:
		taxRate = UNTIL_180_DAYS
	case totalDays < 360:
		taxRate = UNTIL_360_DAYS
	case totalDays < 720:
		taxRate = UNTIL_720_DAYS
	default:
		taxRate = MORE_THAN_720_DAYS
	}

	taxDiscount := (profit * taxRate) / 100.0
	finalAmount -= taxDiscount

	return math.Round(finalAmount*100.0) / 100.0, nil
}

func (rh *RendimentsHandler) SimulatePoupanca(ctx *gin.Context, request SimulationRequest) (float64, error) {
	metric, err := rh.metricRepository.GetLastMetric(ctx, SELIC)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": error_message.ErrFindMetric,
		})
		return 0.0, nil
	}

	monthTax := MONTH_TAX_DEFAULT
	anualSelic := metric.Value.InexactFloat64() / 100
	if anualSelic < SELIC_LESS_THAN_EIGHT {
		monthTax = (0.7 * anualSelic) / 12
	}

	finalAmount := request.InitialValue * math.Pow(1+monthTax, float64(request.Months))

	return math.Round(finalAmount*100.0) / 100.0, nil
}

func (rh *RendimentsHandler) SimulatePoupancaMonthly(ctx *gin.Context, request SimulationRequest) (float64, error) {
	metric, err := rh.metricRepository.GetLastMetric(ctx, SELIC)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": error_message.ErrFindMetric,
		})
		return 0.0, nil
	}

	accumulatedValue := 0.0
	if request.Accumulated != nil {
		accumulatedValue = *request.Accumulated
	}

	monthTax := MONTH_TAX_DEFAULT
	anualSelic := metric.Value.InexactFloat64() / 100
	if anualSelic < SELIC_LESS_THAN_EIGHT {
		monthTax = (0.7 * anualSelic) / 12
	}

	futureAmount := request.InitialValue * (math.Pow(1+monthTax, float64(request.Months)) - 1) / monthTax
	totalAmount := accumulatedValue*math.Pow(1+monthTax, float64(request.Months)) + futureAmount

	return math.Round(totalAmount*100.0) / 100.0, nil
}
