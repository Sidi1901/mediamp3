package handler 

import(
	"net/http"
	"os"
	"errors"
	"path/filepath"
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

	result, err := h.conversionService.ConvertFormat(ctx, request)

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

	// remove temporary files after the download finishes

	defer func(){
		_ = os.RemoveAll(filepath.Dir(result.FilePath))
	}()

	c.FileAttachment(result.FilePath, "converted."+ result.Format)


	c.JSON(http.StatusAccepted, result.Response)
}