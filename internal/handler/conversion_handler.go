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

	var request model.ConversionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Invalid Request",
		})

		return
	}

	response, err := h.conversionService.ConvertFormat(request)

	if err != nil {
		if errors.Is(err, service.ErrConversionNotImplemented){
			c.JSON(http.StatusNotImplemented, gin.H{
				"error": "Conversion functionality is not implemented yet",
			})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"error" : "Failed to convert",
	})

	c.JSON(http.StatusAccepted, response)
}