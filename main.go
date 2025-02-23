package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	""
)

// https://github.com/GolangLessons
// https://gin-gonic.com/docs/quickstart/
/**
* Config
*/


/**
* Config end
*/


/**
* Database 
*/

/**
* Database end
*/


/**
* handlers
*/

func getDocuments(c *gin.Context) {
	
}

func saveDocument(c *gin.Context) {
	var doc Document

	if networkErr := c.ShouldBindJSON(&doc); networkErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": networkErr.Error()})
		return
	}

	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching database connection"})
		return
	}

	s, ok := db.(*Storage)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching database connection"})
		return
	}
	
	id, err := s.SaveDocument(context.Background(), doc.Title, doc.Content)	

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document", "error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Document saved successfully", "id": id})
}
