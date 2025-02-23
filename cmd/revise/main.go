package main 

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"revise/internal/config"
	"github.com/gin-gonic/gin"
)


func main() {
	config := config.MustLoad()
	db, err := New(config.StorageConfig)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router.GET("/documents", getDocuments)
	router.POST("/documents", saveDocument)

	go func() {
		router.Run(":8080")
	}()

	// Graceful shutdown
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	fmt.Println("Shutting down server...")
}
