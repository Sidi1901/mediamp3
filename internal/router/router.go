package router 

import(
	"github.com/gin-gonic/gin"
	"github.com/Sidi1901/mediamp3/internal/config"
	"github.com/Sidi1901/mediamp3/internal/handler"
	"github.com/Sidi1901/mediamp3/internal/service"
)

func Setup(cfg config.Config) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Services
	formatService := service.NewFormatService(cfg)
	conversionService := service.NewConversionService(cfg)

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