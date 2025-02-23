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
	"github.com/jackc/pgx/v5"
)

// https://github.com/GolangLessons
// https://gin-gonic.com/docs/quickstart/
/**
* Config
*/
type Config struct {
	Env            string     `yaml:"env" env-default:"local"`
	StorageConfig  StorageConfig     `yaml:"storage" env-required:"true"`
}

type StorageConfig struct {
	Host     string `yml:"host"`
	Port     int    `yml:"port"`
	Username string `yml:"username"`
	Password string `yml:"password"`
	Database string `yml:"database"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	fmt.Println("config path: ", res)
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}


/**
* Config end
*/


/**
* Database 
*/

type Storage struct {
	db *pgx.Conn
}

func New(config StorageConfig) (*Storage, error) {
	const op = "storage.postgres.New"

	var err error
	
	dns := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	conn, err := pgx.Connect(context.Background(), dns)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: add migration
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS documents (
		id SERIAL PRIMARY KEY,
		title TEXT,
		content TEXT
	)`)

	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	fmt.Println("Successfully connected to the database")

	return &Storage{
		db: conn,
	}, err
}

func (s *Storage) Close() error {
	return s.db.Close(context.Background())
}

func (s *Storage) SaveDocument(ctx context.Context, title string, content string) (int64, error) {
	const op = "storage.postgres.SaveDocument"

	stmt, err := s.db.Prepare(ctx, "save_document", "INSERT INTO documents (title, content) VALUES ($1, $2)")
	// TODO: add more detiled error handling https://github.com/GolangLessons/sso/blob/main/internal/storage/sqlite/sqlite.go#L46
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := s.db.Query(ctx, stmt.SQL, title, content)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	var id int64
	if err := res.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

var db *pgx.Conn

func connectToDB(config StorageConfig) {
	
}

type Document struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

/**
* Database end
*/


/**
* handlers
*/

func getDocuments(c *gin.Context) {
	var documents []Document

	err := db.QueryRow(context.Background(), "SELECT id, title, content FROM documents").Scan(&documents)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching documents", "error": err})
		return
	}

	c.JSON(http.StatusOK, documents)
}

func saveDocument(c *gin.Context) {
	var doc Document

	if networkErr := c.ShouldBindJSON(&doc); networkErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input", "error": networkErr.Error()})
		return
	}
	
	err := db.QueryRow(context.Background(), "INSERT INTO documents (title, content) VALUES ($1, $2)", doc.Title, doc.Content)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document", "error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Document saved successfully"})
}

func main() {
	config := MustLoad()
	connectToDB(config.StorageConfig)
	defer db.Close(context.Background())

	router := gin.Default()

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
