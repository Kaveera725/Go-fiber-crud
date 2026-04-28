package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-fiber-crud/internal/models"
)

type DB struct {
	ORM *gorm.DB
}

// Connect initializes a database pool from environment configuration.
func Connect(ctx context.Context) (*DB, error) {
	dsn, err := buildDSN()
	if err != nil {
		return nil, err
	}

	orm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := orm.WithContext(ctx).AutoMigrate(&models.Product{}); err != nil {
		return nil, err
	}

	return &DB{ORM: orm}, nil
}

func buildDSN() (string, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || dbname == "" {
		return "", fmt.Errorf("missing DB env vars: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME")
	}

	// Default to local development without TLS.
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname), nil
}

// Close shuts down the database pool if it exists.
func (db *DB) Close() {
	if db == nil || db.ORM == nil {
		return
	}

	sqlDB, err := db.ORM.DB()
	if err != nil {
		return
	}
	sqlDB.Close()
}

// ListProducts returns products ordered by newest first, optionally filtered by name.
func (db *DB) ListProducts(ctx context.Context, nameFilter string) ([]models.Product, error) {
	var products []models.Product
	query := db.ORM.WithContext(ctx).Order("id desc")
	if strings.TrimSpace(nameFilter) != "" {
		query = query.Where("name ILIKE ?", "%"+nameFilter+"%")
	}
	err := query.Find(&products).Error
	return products, err
}

// GetProduct fetches a single product by id.
func (db *DB) GetProduct(ctx context.Context, id int) (models.Product, error) {
	var p models.Product
	err := db.ORM.WithContext(ctx).First(&p, id).Error
	return p, err
}

// CreateProduct inserts a new product record.
func (db *DB) CreateProduct(ctx context.Context, name string, price float64, quantity int) error {
	product := models.Product{Name: name, Price: price, Quantity: quantity}
	return db.ORM.WithContext(ctx).Create(&product).Error
}

// UpdateProduct updates an existing product record.
func (db *DB) UpdateProduct(ctx context.Context, id int, name string, price float64, quantity int) error {
	updates := map[string]any{
		"name":     name,
		"price":    price,
		"quantity": quantity,
	}
	return db.ORM.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteProduct removes a product by id.
func (db *DB) DeleteProduct(ctx context.Context, id int) error {
	return db.ORM.WithContext(ctx).Delete(&models.Product{}, id).Error
}
