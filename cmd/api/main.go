package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"

	"book-library-go/config"
	"book-library-go/internal/author"
	"book-library-go/internal/book"
)

func main() {
	cfg := config.Load()

	db, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	bookRepo := book.NewRepository(db)
	bookService := book.NewService(bookRepo)
	bookHandler := book.NewHandler(bookService)

	authorRepo := author.NewRepository(db)
	authorService := author.NewService(authorRepo, bookService)
	authorHandler := author.NewHandler(authorService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/authors", authorHandler.Routes)
	r.Route("/authors/{authorID}/books", bookHandler.Routes)

	log.Printf("server started on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
