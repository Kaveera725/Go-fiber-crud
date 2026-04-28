# ProductVault — Go CRUD App

A full-stack CRUD application built with **Go + Fiber + PostgreSQL + Tailwind CSS**.

## Stack

| Layer    | Technology              |
|----------|-------------------------|
| Backend  | Go 1.21 + Fiber v2      |
| Database | PostgreSQL + pgx/v5     |
| Frontend | HTML Go Templates + Tailwind CSS CDN |

## Project Structure

```
products-crud/
├── cmd/
│   └── main.go              # Entry point, routes
├── internal/
│   ├── database/
│   │   └── db.go            # Connection + CRUD queries
│   ├── handlers/
│   │   └── product.go       # HTTP handlers
│   └── models/
│       └── product.go       # Product struct
├── views/
│   ├── layouts/
│   │   └── base.html        # Base layout
│   ├── index.html           # Product list
│   ├── new.html             # Create form
│   └── edit.html            # Edit form
├── .env                     # Environment config
├── go.mod
└── setup.sql                # DB setup script
```

## Setup

### 1. PostgreSQL

```bash
psql -U postgres -f setup.sql
```

### 2. Configure .env

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=productsdb
APP_PORT=3000
```

### 3. Install dependencies & run

```bash
cd products-crud
go mod tidy
go run cmd/main.go
```

Open http://localhost:3000

## Routes

| Method | Path                    | Action         |
|--------|-------------------------|----------------|
| GET    | /                       | List products  |
| GET    | /products/new           | Create form    |
| POST   | /products               | Create product |
| GET    | /products/:id/edit      | Edit form      |
| POST   | /products/:id/update    | Update product |
| POST   | /products/:id/delete    | Delete product |
