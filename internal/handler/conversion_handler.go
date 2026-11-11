
package handler

import (
	"errors"
	"net/http"
	"os"
	"context"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/Sidi1901/mediamp3/internal/model"
	"github.com/Sidi1901/mediamp3/internal/service"
)

type ConversionHandler struct {
	conversionService service.ConversionService
}

func NewConversionHandler(
	conversionService service.ConversionService,
) *ConversionHandler {
	return &ConversionHandler{
		conversionService: conversionService,
	}
}

func (h *ConversionHandler) Convert(c *gin.Context) {
	ctx := c.Request.Context()

	var request model.ConversionRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid conversion request",
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

		case errors.Is(err, context.Canceled):
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": "Conversion cancelled",
			})

		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": "Conversion timed out",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to convert media",
			})
		}
		return
	}

	defer func() {
		_ = os.RemoveAll(filepath.Dir(result.FilePath))
	}()

	filename := "converted." + result.Format

	c.FileAttachment(result.FilePath, filename)
}
