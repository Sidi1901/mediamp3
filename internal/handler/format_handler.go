package handler 

import(
	"net/http"
	"net/url"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/Sidi1901/mediamp3/internal/service"
)

var (
	NoURLError := errors.New("No Query URL in request")
	NoValidYoutubeURL := errors.New("YouTube URL is invalid")
	FormatBadGateway := errors.New("Failed to retrieve video formats")
)

type FormatHandler struct {
	formatService service.FormatService
}

func NewFormatHandler(formatService service.FormatService) *FormatHandler{
	return &FormatHandler{
		formatService : formatService,
	}
}

func (h *FormatHandler) GetFormats(c *gin.Context){

	ctx := c.request.Context()

	videoURL := strings.Trimpspace(c.Query("url"))

	if videoURL == "" {
		c.SJON(http.StatusbadRequest, gin.H{
			"error" : NoURLError,
		})

		return
	}

	if !isSupportedYoutubeURL(videoURL) {
		c.JSON(http.StatusbadRequest, gin.H{
			"error" : NoValidYoutubeURL,
		})
	}

	formatsResponse, err := h.formatService.GetFormats(ctx, videoURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error" : FormatBadGateway,
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