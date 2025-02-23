package document

import (
	"net/http"

	"revise/internal/service"

	"github.com/gin-gonic/gin"
)

func getDocuments(c *gin.Context) {
	ctx := c.Request.Context()

	// documents, err := service.DocumentService.New()
	// TODO: педавать через context параметры для инициализации или передавать service через контекст?
	if err != nil {
		if err == service.ErrDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
	})
}
