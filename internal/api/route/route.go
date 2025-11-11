package route

import (
	"net/http"
	"os"

	"github.com/jonathanmoreiraa/2cents/internal/api/handler"
	"github.com/jonathanmoreiraa/2cents/internal/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

type HandlerGroup struct {
	UserHandler       *handler.UserHandler
	RevenueHandler    *handler.RevenueHandler
	ExpenseHandler    *handler.ExpenseHandler
	CategoryHandler   *handler.CategoryHandler
	SavingHandler     *handler.SavingHandler
	RendimentsHandler *handler.RendimentsHandler
}

func NewServerHTTP(Handlers HandlerGroup) *ServerHTTP {
	engine := gin.New()
	engine.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.AllowCredentials = true
	engine.Use(cors.New(config))

	engine.POST("/login", Handlers.UserHandler.Login)
	engine.POST("/user/new", Handlers.UserHandler.Create)

	api := engine.Group("/api", middleware.AuthorizationMiddleware)
	api.GET("/user/me", Handlers.UserHandler.GetCurrentUser)
	api.GET("/user/:id", Handlers.UserHandler.FindByID)
	api.PUT("/user/:id", Handlers.UserHandler.Update)
	// api.DELETE("/user/:id", Handlers.UserHandler.Delete)

	api.POST("/revenue/add", Handlers.RevenueHandler.Create)
	api.GET("/revenue/:id", Handlers.RevenueHandler.FindByID)
	api.GET("/revenues", Handlers.RevenueHandler.FindAll)
	api.POST("/revenue/filter", Handlers.RevenueHandler.FindByFilters)
	api.PUT("/revenue/:id", Handlers.RevenueHandler.Update)
	api.DELETE("/revenue/:id", Handlers.RevenueHandler.Delete)

	api.POST("/category/add", Handlers.CategoryHandler.Create)
	api.GET("/categories", Handlers.CategoryHandler.GetAllCategories)
	api.POST("/category", Handlers.CategoryHandler.FindCategory)
	api.DELETE("/category/:id", Handlers.CategoryHandler.Delete)

	api.POST("/expense/add", Handlers.ExpenseHandler.Create)
	api.GET("/expense/:id", Handlers.ExpenseHandler.FindByID)
	api.GET("/expenses", Handlers.ExpenseHandler.FindAll)
	api.POST("/expense/filter", Handlers.ExpenseHandler.FindByFilters)
	api.PUT("/expense/:id", Handlers.ExpenseHandler.Update)
	api.DELETE("/expense/:id", Handlers.ExpenseHandler.Delete)

	api.POST("/saving/add", Handlers.SavingHandler.Create)
	api.GET("/saving/:id", Handlers.SavingHandler.FindByID)
	api.GET("/savings", Handlers.SavingHandler.FindAll)
	api.PUT("/saving/:id", Handlers.SavingHandler.Update)
	api.DELETE("/saving/:id", Handlers.SavingHandler.Delete)

	api.POST("/rendiments", Handlers.RendimentsHandler.SimulateAllRendiments)
	api.POST("/month/rendiment", Handlers.RendimentsHandler.SimulateMonthlyRendiments)

	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"error":   "Rota não encontrada",
			"message": "Verifique a URL e tente novamente",
		})
	})

	return &ServerHTTP{engine: engine}
}

func (sh *ServerHTTP) Start() {
	sh.engine.Run(getEnv("PORT", ":8080"))
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return ":" + value
	}
	return defaultValue
}
