package handler 

import(
	"net/http"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/Sidi1901/mediamp3/internal/model"
	"github.com/Sidi1901/mediamp3/internal/service"
)


type ConversionHandler struct{
	conversionService service.ConversionService
}

func NewConversionHandler(conversionService service.ConversionService) *ConversionHandler{
	return &ConversionHandler{
		conversionService: conversionService,
	}
}

func (h *ConversionHandler) Convert(c *gin.Context){
	ctx := c.Request.Context()

	var request model.ConversionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Invalid Request",
		})

		return
	}

	response, err := h.conversionService.ConvertFormat(ctx, request)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidConversionRequest):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, service.ErrConversionNotImplemented):
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "Conversion functionality is not implemented yet",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create conversion",
			})
		}
		return
	}


	c.JSON(http.StatusAccepted, response)
}