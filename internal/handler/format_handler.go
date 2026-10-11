package handler 

import(
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Sidi1901/mediamp3/internal/service"
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
	response := h.formatService.GetFormats()

	c.JSON(http.StatusOK, response)
}