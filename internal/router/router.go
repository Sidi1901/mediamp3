package router 

import(
	"github.com/gin-gonic/gin"
	"github.com/Sidi1901/mediamp3/internal/handler"
	"github.com/Sidi1901/mediamp3/internal/service"
)

func Setup() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Services
	formatService := service.NewFormatService()
	conversionService := service.NewConversionService()

	//handler
	healthHandler := handler.NewHealthHandler()
	formatHandler := handler.NewFormatHandler(formatService)
	conversionHandler := handler.NewConversionHandler(conversionService)

	// Endpoints

	router.GET("/health",healthHandler.CheckHealth)

	api := router.Group("/api/v1")

	{
		api.GET("/formats", formatHandler.GetFormats)
		api.POST("/conversions", conversionHandler.Convert)
	}

	return router
}