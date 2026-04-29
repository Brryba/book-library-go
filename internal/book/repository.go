package book

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(ctx context.Context) ([]Book, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, author_id, title, year, created_at FROM books ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.AuthorID, &b.Title, &b.Year, &b.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*Book, error) {
	var b Book
	err := r.db.QueryRow(ctx, `
		SELECT id, author_id, title, year, created_at FROM books WHERE id = $1
	`, id).Scan(&b.ID, &b.AuthorID, &b.Title, &b.Year, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) FindByAuthorIDs(ctx context.Context, authorIDs []uuid.UUID) ([]Book, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, author_id, title, year, created_at FROM books WHERE author_id = ANY($1)
	`, authorIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.AuthorID, &b.Title, &b.Year, &b.CreatedAt); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *Repository) Create(ctx context.Context, authorID uuid.UUID, req CreateRequest) (*Book, error) {
	var b Book
	err := r.db.QueryRow(ctx, `
		INSERT INTO books (author_id, title, year) VALUES ($1, $2, $3)
		RETURNING id, author_id, title, year, created_at
	`, authorID, req.Title, req.Year).Scan(&b.ID, &b.AuthorID, &b.Title, &b.Year, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Book, error) {
	var b Book
	err := r.db.QueryRow(ctx, `
		UPDATE books SET title = $1, year = $2 WHERE id = $3
		RETURNING id, author_id, title, year, created_at
	`, req.Title, req.Year, id).Scan(&b.ID, &b.AuthorID, &b.Title, &b.Year, &b.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM books WHERE id = $1`, id)
	return err
}
