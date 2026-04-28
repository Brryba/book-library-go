package author_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"book-library-go/internal/author"
)

type mockRepository struct {
	authors []author.Author
	err     error
}

func (m *mockRepository) FindAll(ctx context.Context) ([]author.Author, error) {
	return m.authors, m.err
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (*author.Author, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.authors) == 0 {
		return nil, errors.New("author not found")
	}
	return &m.authors[0], nil
}

func (m *mockRepository) Create(ctx context.Context, req author.CreateRequest) (*author.Author, error) {
	if m.err != nil {
		return nil, m.err
	}
	a := &author.Author{Name: req.Name, Bio: req.Bio}
	return a, nil
}

func (m *mockRepository) Update(ctx context.Context, id uuid.UUID, req author.UpdateRequest) (*author.Author, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &m.authors[0], nil
}

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.err
}

func TestCreate_EmptyName(t *testing.T) {
	svc := author.NewService(&mockRepository{})

	_, err := svc.Create(context.Background(), author.CreateRequest{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestCreate_Success(t *testing.T) {
	svc := author.NewService(&mockRepository{})

	a, err := svc.Create(context.Background(), author.CreateRequest{Name: "Orwell", Bio: "Writer"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name != "Orwell" {
		t.Errorf("expected name Orwell, got %s", a.Name)
	}
}

func TestGetAll_ReturnsError(t *testing.T) {
	svc := author.NewService(&mockRepository{err: errors.New("db error")})

	_, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc := author.NewService(&mockRepository{})

	_, err := svc.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetByID_Success(t *testing.T) {
	svc := author.NewService(&mockRepository{
		authors: []author.Author{{Name: "Orwell"}},
	})

	a, err := svc.GetByID(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name != "Orwell" {
		t.Errorf("expected Orwell, got %s", a.Name)
	}
}

func TestUpdate_EmptyName(t *testing.T) {
	svc := author.NewService(&mockRepository{
		authors: []author.Author{{Name: "Orwell"}},
	})

	_, err := svc.Update(context.Background(), uuid.New(), author.UpdateRequest{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	svc := author.NewService(&mockRepository{})

	_, err := svc.Update(context.Background(), uuid.New(), author.UpdateRequest{Name: "Orwell"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDelete_NotFound(t *testing.T) {
	svc := author.NewService(&mockRepository{})

	err := svc.Delete(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDelete_Success(t *testing.T) {
	svc := author.NewService(&mockRepository{
		authors: []author.Author{{Name: "Orwell"}},
	})

	err := svc.Delete(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
