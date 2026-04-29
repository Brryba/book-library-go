package author

import (
	"time"

	"book-library-go/internal/book"
	"github.com/google/uuid"
)

type Author struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthorWithBooks struct {
	Author
	Books []book.Book `json:"books"`
}

type CreateRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type UpdateRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}
