package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"revise/internal/config"
	"revise/internal/service"
	"revise/internal/storage"
	"revise/internal/handlers"

	"github.com/gin-gonic/gin"
)


func main() {
	config := config.MustLoad()
	db, err := storage.New(config.StorageConfig)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()
	service := service.New(log.Default(), db)
	handlers := handlers.New(*service)

	router := gin.Default()

	router.GET("/documents", handlers.GetDocuments)
	router.POST("/documents", handlers.SaveDocument)

	go func() {
		router.Run(":8080")
	}()

	// Graceful shutdown
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	fmt.Println("Shutting down server...")
}
