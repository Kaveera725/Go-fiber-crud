-- Run with: psql -U postgres -f setup.sql

SELECT 'CREATE DATABASE productsdb'
WHERE NOT EXISTS (
  SELECT 1 FROM pg_database WHERE datname = 'productsdb'
)\gexec

\connect productsdb

CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  price NUMERIC(10,2) NOT NULL DEFAULT 0,
  quantity INTEGER NOT NULL DEFAULT 0
);
