package book

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

const cacheKey = "books:all"

var ErrNotFound = errors.New("book not found")

type repository interface {
	FindAll(ctx context.Context) ([]Book, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Book, error)
	FindByAuthorIDs(ctx context.Context, authorIDs []uuid.UUID) ([]Book, error)
	Create(ctx context.Context, authorID uuid.UUID, req CreateRequest) (*Book, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Book, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type Service struct {
	repo  repository
	cache *Cache
}

func NewService(repo repository, cache *Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

func (s *Service) GetAll(ctx context.Context) ([]Book, error) {
	if books, ok := s.cache.get(cacheKey); ok {
		return books, nil
	}

	books, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	s.cache.add(cacheKey, books)
	return books, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Book, error) {
	b, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return b, nil
}

func (s *Service) GetByAuthorIDs(ctx context.Context, authorIDs []uuid.UUID) ([]Book, error) {
	return s.repo.FindByAuthorIDs(ctx, authorIDs)
}

func (s *Service) Create(ctx context.Context, authorID uuid.UUID, req CreateRequest) (*Book, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Year < 1 {
		return nil, errors.New("year must be positive")
	}
	b, err := s.repo.Create(ctx, authorID, req)
	if err != nil {
		return nil, err
	}
	s.cache.invalidate(cacheKey)
	return b, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Book, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if req.Year < 1 {
		return nil, errors.New("year must be positive")
	}
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	b, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	s.cache.invalidate(cacheKey)
	return b, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.cache.invalidate(cacheKey)
	return nil
}
