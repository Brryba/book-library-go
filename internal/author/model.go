package author

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type UpdateRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}
