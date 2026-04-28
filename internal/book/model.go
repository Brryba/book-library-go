package book

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID        uuid.UUID `json:"id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRequest struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
}

type UpdateRequest struct {
	Title string `json:"title"`
	Year  int    `json:"year"`
}
