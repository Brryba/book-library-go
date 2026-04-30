package book_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"book-library-go/internal/book"
)

type mockRepository struct {
	books         []book.Book
	err           error
	findAllCalled bool
	createCalled  bool
	deleteCalled  bool
	updateCalled  bool
}

func (m *mockRepository) FindAll(ctx context.Context) ([]book.Book, error) {
	m.findAllCalled = true
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
	m.createCalled = true
	if m.err != nil {
		return nil, m.err
	}
	return &book.Book{Title: req.Title, Year: req.Year}, nil
}

func (m *mockRepository) Update(ctx context.Context, id uuid.UUID, req book.UpdateRequest) (*book.Book, error) {
	m.updateCalled = true
	if m.err != nil {
		return nil, m.err
	}
	return &m.books[0], nil
}

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.deleteCalled = true
	return m.err
}

type mockCache struct {
	data             []book.Book
	addCalled        bool
	invalidateCalled bool
}

func (m *mockCache) Get(key string) ([]book.Book, bool) {
	if m.data != nil {
		return m.data, true
	}
	return nil, false
}

func (m *mockCache) Add(key string, value []book.Book) {
	m.addCalled = true
	m.data = value
}

func (m *mockCache) Invalidate(key string) {
	m.invalidateCalled = true
	m.data = nil
}

func (m *mockCache) Close() {}

func newTestService(repo *mockRepository, c *mockCache) *book.Service {
	return book.NewService(repo, c)
}

func TestGetAll_CacheHit(t *testing.T) {
	repo := &mockRepository{}
	c := &mockCache{data: []book.Book{{Title: "1984"}}}

	books, err := newTestService(repo, c).GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(books) != 1 {
		t.Errorf("expected 1 book, got %d", len(books))
	}
	if repo.findAllCalled {
		t.Fatal("repository should not be called on cache hit")
	}
}

func TestGetAll_CacheMiss(t *testing.T) {
	repo := &mockRepository{books: []book.Book{{Title: "1984"}}}
	c := &mockCache{}

	books, err := newTestService(repo, c).GetAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(books) != 1 {
		t.Errorf("expected 1 book, got %d", len(books))
	}
	if !repo.findAllCalled {
		t.Fatal("repository should be called on cache miss")
	}
	if !c.addCalled {
		t.Fatal("cache Add should be called after repo fetch")
	}
}

func TestCreate_EmptyTitle(t *testing.T) {
	repo := &mockRepository{}
	_, err := newTestService(repo, &mockCache{}).Create(context.Background(), uuid.New(), book.CreateRequest{Title: ""})
	if err == nil {
		t.Fatal("expected error for empty title")
	}
	if repo.createCalled {
		t.Fatal("repository should not be called")
	}
}

func TestCreate_InvalidYear(t *testing.T) {
	repo := &mockRepository{}
	_, err := newTestService(repo, &mockCache{}).Create(context.Background(), uuid.New(), book.CreateRequest{Title: "1984", Year: 0})
	if err == nil {
		t.Fatal("expected error for invalid year")
	}
	if repo.createCalled {
		t.Fatal("repository should not be called")
	}
}

func TestCreate_InvalidatesCache(t *testing.T) {
	c := &mockCache{}
	_, err := newTestService(&mockRepository{}, c).Create(context.Background(), uuid.New(), book.CreateRequest{Title: "1984", Year: 1949})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.invalidateCalled {
		t.Fatal("cache Invalidate should be called after create")
	}
}

func TestUpdate_EmptyTitle(t *testing.T) {
	repo := &mockRepository{books: []book.Book{{Title: "1984"}}}
	_, err := newTestService(repo, &mockCache{}).Update(context.Background(), uuid.New(), book.UpdateRequest{Title: ""})
	if err == nil {
		t.Fatal("expected error for empty title")
	}
	if repo.updateCalled {
		t.Fatal("repository should not be called")
	}
}

func TestUpdate_InvalidatesCache(t *testing.T) {
	repo := &mockRepository{books: []book.Book{{Title: "1984", Year: 1949}}}
	c := &mockCache{}
	_, err := newTestService(repo, c).Update(context.Background(), uuid.New(), book.UpdateRequest{Title: "Animal Farm", Year: 1945})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.invalidateCalled {
		t.Fatal("cache Invalidate should be called after update")
	}
}

func TestDelete_NotFound(t *testing.T) {
	repo := &mockRepository{}
	err := newTestService(repo, &mockCache{}).Delete(context.Background(), uuid.New())
	if err == nil {
		t.Fatal("expected error")
	}
	if repo.deleteCalled {
		t.Fatal("repository Delete should not be called")
	}
}

func TestDelete_InvalidatesCache(t *testing.T) {
	repo := &mockRepository{books: []book.Book{{Title: "1984"}}}
	c := &mockCache{}
	err := newTestService(repo, c).Delete(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !c.invalidateCalled {
		t.Fatal("cache Invalidate should be called after delete")
	}
}
