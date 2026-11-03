package handler

import (
	"errors"
	"github.com/Sidi1901/mediamp3/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strings"
)

var (
	NoURLError        = errors.New("No Query URL in request")
      	NoValidYoutubeURL = errors.New("YouTube URL is invalid")
	FormatBadGateway  = errors.New("Failed to retrieve video formats")
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
			"error": NoURLError,
		})

		return
	}

	if !isSupportedYoutubeURL(videoURL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": NoValidYoutubeURL,
		})
	}

	formatsResponse, err := h.formatService.GetFormats(ctx, videoURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": FormatBadGateway,
		})
	}

	c.JSON(http.StatusOK, formatsResponse)
}

func isSupportedYoutubeURL(rawURL string) bool {
	pasedURL, err := url.ParseRequestURI(rawURL)

	if err != nil {
		return false
	}

	if pasedURL.Scheme != "https" || pasedURL.User != nil {
		return false
	}

	host := strings.ToLower(pasedURL.Hostname())

	switch host {
	case "youtube.com",
		"www.youtube.com",
		"m.youtube.com",
		"youtu.be",
		"www.youtube-nocookie.com":
		return true
	default:
		return false
	}
}
