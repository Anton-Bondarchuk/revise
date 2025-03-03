package service

import (
	"context"
	"errors"
	"log"

	"revise/internal/domain/models"
)

var (
	ErrDocumentNotFound = errors.New("document not found")
	ErrFailedToSaveDocument = errors.New("failed to save document")
)

type DocumentProvider interface {
	GetDocuments(ctx context.Context) ([]models.Document, error)
	SaveDocument(ctx context.Context, title, content string) (int64, error)
}

type DocumentService struct {
	log               *log.Logger
	documentProvider  DocumentProvider
}

func New(log *log.Logger, documentProvider DocumentProvider) *DocumentService {
	return &DocumentService{
		log:              log,
		documentProvider: documentProvider,
	}
}

func (s *DocumentService) GetDocuments(ctx context.Context) ([]models.Document, error) {
	documents, err := s.documentProvider.GetDocuments(ctx)
	if err != nil {
		s.log.Printf("Error retrieving documents: %v", err)
		
		if errors.Is(err, ErrDocumentNotFound) {
			return nil, ErrDocumentNotFound
		}

		return nil, err
	}
	return documents, nil
}

func (s *DocumentService) SaveDocument(ctx context.Context, title, content string) (int64, error) {
	id, err := s.documentProvider.SaveDocument(ctx, title, content)
	if err != nil {
		s.log.Printf("Error saving document: %v", err)
		return 0, ErrFailedToSaveDocument
	}
	return id, nil
}
