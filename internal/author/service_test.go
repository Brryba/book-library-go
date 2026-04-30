package author_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"book-library-go/internal/author"
)

type mockRepository struct {
	authors      []author.Author
	err          error
	deleteCalled bool
	createCalled bool
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
	m.createCalled = true
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
	m.deleteCalled = true
	return m.err
}

func newTestService(repo *mockRepository) *author.Service {
	return author.NewService(repo, nil)
}

func TestCreate_EmptyName(t *testing.T) {
	repo := &mockRepository{}
	svc := newTestService(repo)

	_, err := svc.Create(context.Background(), author.CreateRequest{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
	if repo.createCalled {
		t.Fatal("expected repository not to be called")
	}
}

func TestCreate_Success(t *testing.T) {
	repo := &mockRepository{}
	svc := newTestService(repo)

	a, err := svc.Create(context.Background(), author.CreateRequest{Name: "Orwell", Bio: "Writer"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Name != "Orwell" {
		t.Errorf("expected name Orwell, got %s", a.Name)
	}
	if !repo.createCalled {
		t.Fatal("expected repository Create to be called")
	}
}

func TestGetAll_ReturnsError(t *testing.T) {
	svc := newTestService(&mockRepository{err: errors.New("db error")})

	_, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetByID_NotFound(t *testing.T) {
	svc := newTestService(&mockRepository{})

	_, err := svc.GetByID(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetByID_Success(t *testing.T) {
	svc := newTestService(&mockRepository{
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
	svc := newTestService(&mockRepository{
		authors: []author.Author{{Name: "Orwell"}},
	})

	_, err := svc.Update(context.Background(), uuid.New(), author.UpdateRequest{Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	svc := newTestService(&mockRepository{})

	_, err := svc.Update(context.Background(), uuid.New(), author.UpdateRequest{Name: "Orwell"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDelete_NotFound(t *testing.T) {
	repo := &mockRepository{}
	svc := newTestService(repo)

	err := svc.Delete(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if repo.deleteCalled {
		t.Fatal("expected repository Delete not to be called")
	}
}

func TestDelete_Success(t *testing.T) {
	repo := &mockRepository{
		authors: []author.Author{{Name: "Orwell"}},
	}
	svc := newTestService(repo)

	err := svc.Delete(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repo.deleteCalled {
		t.Fatal("expected repository Delete to be called")
	}
}
