
package handler

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Sidi1901/mediamp3/internal/service"
	"github.com/Sidi1901/mediamp3/internal/validation"
)

var (
	NoURLError       = errors.New("no URL provided in request")
	NoValidYoutubeURL = errors.New("YouTube URL is invalid")
	FormatBadGateway = errors.New("failed to retrieve video formats")
)

type FormatHandler struct {
	formatService service.FormatService
}

func NewFormatHandler(formatService service.FormatService) *FormatHandler {
	return &FormatHandler{
		formatService: formatService,
	}
}

func (h *FormatHandler) GetFormats(c *gin.Context) {
	ctx := c.Request.Context()

	videoURL := strings.TrimSpace(c.Query("url"))

	if videoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": NoURLError.Error(),
		})
		return
	}

	if !validation.IsValidYoutubeURL(videoURL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": NoValidYoutubeURL.Error(),
		})
		return
	}

	formatsResponse, err := h.formatService.GetFormats(ctx, videoURL)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": FormatBadGateway.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, formatsResponse)
}

