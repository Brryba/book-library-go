package author

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

func (r *Repository) FindAll(ctx context.Context) ([]Author, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, name, bio, created_at FROM authors ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var a Author
		if err := rows.Scan(&a.ID, &a.Name, &a.Bio, &a.CreatedAt); err != nil {
			return nil, err
		}
		authors = append(authors, a)
	}
	return authors, nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*Author, error) {
	var a Author
	err := r.db.QueryRow(ctx, `
		SELECT id, name, bio, created_at FROM authors WHERE id = $1
	`, id).Scan(&a.ID, &a.Name, &a.Bio, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) Create(ctx context.Context, req CreateRequest) (*Author, error) {
	var a Author
	err := r.db.QueryRow(ctx, `
		INSERT INTO authors (name, bio) VALUES ($1, $2)
		RETURNING id, name, bio, created_at
	`, req.Name, req.Bio).Scan(&a.ID, &a.Name, &a.Bio, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (*Author, error) {
	var a Author
	err := r.db.QueryRow(ctx, `
		UPDATE authors SET name = $1, bio = $2 WHERE id = $3
		RETURNING id, name, bio, created_at
	`, req.Name, req.Bio, id).Scan(&a.ID, &a.Name, &a.Bio, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM authors WHERE id = $1`, id)
	return err
}
