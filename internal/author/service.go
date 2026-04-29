package author

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"book-library-go/internal/book"
)

var ErrNotFound = errors.New("author not found")

type repository interface {
	FindAll(ctx context.Context) ([]Author, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Author, error)
	Create(ctx context.Context, req CreateRequest) (*Author, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Author, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo        repository
	bookService *book.Service
}

func NewService(repo repository, bookService *book.Service) *Service {
	return &Service{repo: repo, bookService: bookService}
}

func (s *Service) GetAll(ctx context.Context) ([]Author, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) GetAllWithBooks(ctx context.Context) ([]AuthorWithBooks, error) {
	authors, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(authors))
	for i, a := range authors {
		ids[i] = a.ID
	}

	books, err := s.bookService.GetByAuthorIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	booksByAuthor := make(map[uuid.UUID][]book.Book)
	for _, b := range books {
		booksByAuthor[b.AuthorID] = append(booksByAuthor[b.AuthorID], b)
	}

	result := make([]AuthorWithBooks, len(authors))
	for i, a := range authors {
		result[i] = AuthorWithBooks{
			Author: a,
			Books:  booksByAuthor[a.ID],
		}
	}
	return result, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Author, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

func (s *Service) Create(ctx context.Context, req CreateRequest) (*Author, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	return s.repo.Create(ctx, req)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Author, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s.repo.Update(ctx, id, req)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
