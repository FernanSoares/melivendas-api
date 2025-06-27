-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS melivendas;

-- Use the database
USE melivendas;

-- Create items table
CREATE TABLE IF NOT EXISTS items (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price BIGINT NOT NULL,
    stock BIGINT NOT NULL,
    status ENUM('ACTIVE', 'INACTIVE') NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
