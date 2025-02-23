package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"

	"revise/internal/config"
	"revise/internal/domain/models"
)

type Storage struct {
	db *pgx.Conn
}

func New(config config.StorageConfig) (*Storage, error) {
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

func (s *Storage) GetDocuments(ctx context.Context) ([]models.Document, error) {
	const op = "storage.postgres.GetDocuments"

	stmt, err := s.db.Prepare(ctx, "get_documents", "SELECT id, title, content FROM documents")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.db.Query(ctx, stmt.SQL)	
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var documents []models.Document
	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(&doc.ID, &doc.Title, &doc.Content); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func (s *Storage) SaveDocument(ctx context.Context, title string, content string) (int64, error) {
	const op = "storage.postgres.SaveDocument"

	// TODO: add validation fields and use prepared statement
	stmt := `INSERT INTO documents (title, content) VALUES ($1, $2) RETURNING id`

	var id int64
	err := s.db.QueryRow(ctx, stmt, title, content).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
