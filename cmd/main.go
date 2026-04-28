package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	htmltmpl "github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"

	"go-fiber-crud/internal/database"
	"go-fiber-crud/internal/handlers"
)

// main loads config, initializes dependencies, and starts the HTTP server.
func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer db.Close()

	engine := htmltmpl.New("./views", ".html")

	app := fiber.New(fiber.Config{Views: engine})

	h := handlers.NewProductHandler(db)

	app.Get("/", h.List)
	app.Get("/products/new", h.NewForm)
	app.Post("/products", h.Create)
	app.Get("/products/:id/edit", h.EditForm)
	app.Post("/products/:id/update", h.Update)
	app.Post("/products/:id/delete", h.Delete)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("listening on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
