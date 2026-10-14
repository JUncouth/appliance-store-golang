CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY,
    country VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    street VARCHAR(200) NOT NULL
);

CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY,
    client_name VARCHAR(100) NOT NULL,
    client_surname VARCHAR(100) NOT NULL,
    birthday DATE NOT NULL,
    gender VARCHAR(30) NOT NULL,
    registration_date TIMESTAMPTZ NOT NULL,
    address_id UUID NOT NULL UNIQUE REFERENCES addresses(id) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS suppliers (
    id UUID PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    address_id UUID NOT NULL UNIQUE REFERENCES addresses(id) ON DELETE RESTRICT,
    phone_number VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS images (
    id UUID PRIMARY KEY,
    image BYTEA NOT NULL
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(100) NOT NULL,
    price NUMERIC(12,2) NOT NULL CHECK (price >= 0),
    available_stock INTEGER NOT NULL CHECK (available_stock >= 0),
    last_update_date TIMESTAMPTZ NOT NULL,
    supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE RESTRICT,
    image_id UUID NULL REFERENCES images(id) ON DELETE SET NULL
);

CREATE USER readonly WITH PASSWORD 'readonly';

GRANT CONNECT ON DATABASE storeapi TO readonly;

GRANT USAGE ON SCHEMA public TO readonly;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly;