package book_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"book-library-go/internal/book"
)

type mockRepository struct {
	books []book.Book
	err   error
}

func (m *mockRepository) FindAll(ctx context.Context) ([]book.Book, error) {
	return m.books, m.err
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (*book.Book, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.books) == 0 {
		return nil, errors.New("book not found")
	}
	return &m.books[0], nil
}

func (m *mockRepository) FindByAuthorIDs(ctx context.Context, authorIDs []uuid.UUID) ([]book.Book, error) {
	return m.books, m.err
}

func (m *mockRepository) Create(ctx context.Context, authorID uuid.UUID, req book.CreateRequest) (*book.Book, error) {
	if m.err != nil {
		return nil, m.err
	}
	b := &book.Book{Title: req.Title, Year: req.Year}
	return b, nil
}

func (m *mockRepository) Update(ctx context.Context, id uuid.UUID, req book.UpdateRequest) (*book.Book, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &m.books[0], nil
}

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.err
}

func TestCreate_EmptyTitle(t *testing.T) {
	svc := book.NewService(&mockRepository{})

	_, err := svc.Create(context.Background(), uuid.New(), book.CreateRequest{Title: ""})
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}

func TestCreate_InvalidYear(t *testing.T) {
	svc := book.NewService(&mockRepository{})

	_, err := svc.Create(context.Background(), uuid.New(), book.CreateRequest{Title: "1984", Year: 0})
	if err == nil {
		t.Fatal("expected error for invalid year, got nil")
	}
}

func TestCreate_Success(t *testing.T) {
	svc := book.NewService(&mockRepository{})

	b, err := svc.Create(context.Background(), uuid.New(), book.CreateRequest{Title: "1984", Year: 1949})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Title != "1984" {
		t.Errorf("expected title 1984, got %s", b.Title)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc := book.NewService(&mockRepository{})

	_, err := svc.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDelete_Success(t *testing.T) {
	svc := book.NewService(&mockRepository{
		books: []book.Book{{Title: "1984"}},
	})

	err := svc.Delete(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
