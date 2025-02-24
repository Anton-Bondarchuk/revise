package handlers

import (
	"net/http"

	"revise/internal/domain/models"
	"revise/internal/service"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	DocumentService service.DocumentService
}

func New(documentService service.DocumentService) *DocumentHandler {
	return &DocumentHandler{
		DocumentService: documentService,
	}
}

func (h *DocumentHandler) GetDocuments(c *gin.Context) {
	ctx := c.Request.Context()

	documents, err := h.DocumentService.GetDocuments(ctx)
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

func (h *DocumentHandler) SaveDocument(c *gin.Context) {
	var doc models.Document

	if networkErr := c.ShouldBindJSON(&doc); networkErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": networkErr.Error()})
		return
	}
	
	id, err := h.DocumentService.SaveDocument(c.Request.Context(), doc.Title, doc.Content)	

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document", "error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}