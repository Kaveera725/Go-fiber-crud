package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-fiber-crud/internal/models"
)

type DB struct {
	Pool *pgxpool.Pool
}

// Connect initializes a database pool from environment configuration.
func Connect(ctx context.Context) (*DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || dbname == "" {
		return nil, fmt.Errorf("missing DB env vars: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME")
	}

	// Default to local development without TLS.
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	cfg.MaxConns = 5
	cfg.MinConns = 1
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &DB{Pool: pool}, nil
}

// Close shuts down the database pool if it exists.
func (db *DB) Close() {
	if db != nil && db.Pool != nil {
		db.Pool.Close()
	}
}

// ListProducts returns all products ordered by newest first.
func (db *DB) ListProducts(ctx context.Context) ([]models.Product, error) {
	rows, err := db.Pool.Query(ctx, `SELECT id, name, price FROM products ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

// GetProduct fetches a single product by id.
func (db *DB) GetProduct(ctx context.Context, id int) (models.Product, error) {
	var p models.Product
	err := db.Pool.QueryRow(ctx, `SELECT id, name, price FROM products WHERE id=$1`, id).Scan(&p.ID, &p.Name, &p.Price)
	return p, err
}

// CreateProduct inserts a new product record.
func (db *DB) CreateProduct(ctx context.Context, name string, price float64) error {
	_, err := db.Pool.Exec(ctx, `INSERT INTO products (name, price) VALUES ($1, $2)`, name, price)
	return err
}

// UpdateProduct updates an existing product record.
func (db *DB) UpdateProduct(ctx context.Context, id int, name string, price float64) error {
	_, err := db.Pool.Exec(ctx, `UPDATE products SET name=$1, price=$2 WHERE id=$3`, name, price, id)
	return err
}

// DeleteProduct removes a product by id.
func (db *DB) DeleteProduct(ctx context.Context, id int) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM products WHERE id=$1`, id)
	return err
}
