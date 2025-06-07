-- Script de inicialización de la base de datos para el Comparador de Precios
-- Este script crea la estructura de la base de datos para el proyecto, el web scrapping se 
--encargara de rellenar las tablas products y prices y las otras tablas son informacion relacionada al usuario
-- (Sus notificaciones, watchlists, etc)


-- Crear la base de datos si no existe
CREATE DATABASE IF NOT EXISTS comparador_precios;

-- Usar la base de datos
USE comparador_precios;

-- Crear tabla de categorías
CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    product_count INT DEFAULT 0, -- Contador de productos
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category_slug (slug) -- Índice para búsquedas por slug
);

-- Crear tabla de productos
CREATE TABLE IF NOT EXISTS products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    image_url VARCHAR(255),
    category_id INT,
    image_hash BIGINT UNSIGNED NULL,
    specifications JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL, -- Para soft delete
    FOREIGN KEY (category_id) REFERENCES categories(id),
    INDEX idx_product_slug (slug), -- Índice para búsquedas por slug
    INDEX idx_product_category (category_id), -- Índice para búsquedas por categoría
    INDEX idx_product_deleted (deleted_at) -- Índice para filtrar por deleted_at
);

-- Crear tabla de precios
CREATE TABLE IF NOT EXISTS prices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    store VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'EUR',
    url VARCHAR(512) NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id),
    INDEX idx_price_product (product_id), -- Índice para búsquedas por producto
    INDEX idx_price_store (store), -- Índice para búsquedas por tienda
    INDEX idx_price_product_store (product_id, store), -- Índice compuesto para producto y tienda
    INDEX idx_price_value (price) -- Índice para búsquedas por precio
);

-- Crear tabla de usuarios
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_email (email), -- Índice para búsquedas por email
    INDEX idx_user_username (username) -- Índice para búsquedas por username
);

-- Crear tabla de watchlists (lista de deseos del usuario)
CREATE TABLE IF NOT EXISTS watchlists (
    user_id INT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Crear tabla de alertas de precio
CREATE TABLE IF NOT EXISTS price_alerts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    target_price DECIMAL(10, 2) NOT NULL,
    notify_by_email BOOLEAN DEFAULT TRUE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    INDEX idx_alert_user (user_id), -- Índice para búsquedas por usuario
    INDEX idx_alert_product (product_id) -- Índice para búsquedas por producto
);

-- Crear tabla de watchlist_items (elementos de la lista de deseos)
CREATE TABLE IF NOT EXISTS watchlist_items (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    target_price DECIMAL(10, 2) NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES watchlists(user_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    UNIQUE KEY unique_watchlist_item (user_id, product_id),
    INDEX idx_watchlist_user (user_id), -- Índice para búsquedas por usuario
    INDEX idx_watchlist_product (product_id) -- Índice para búsquedas por producto
);

-- Crear tabla de notificaciones
CREATE TABLE IF NOT EXISTS notifications (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    product_id INT NOT NULL,
    alert_id INT DEFAULT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (alert_id) REFERENCES price_alerts(id) ON DELETE SET NULL,
    INDEX idx_notification_user (user_id),
    INDEX idx_notification_read (is_read)
);

-- Crear triggers para mantener actualizado el contador de productos en categorías

-- Trigger para incrementar el contador cuando se añade un producto
DELIMITER //
CREATE TRIGGER IF NOT EXISTS trg_product_insert AFTER INSERT ON products
FOR EACH ROW
BEGIN
    UPDATE categories SET product_count = product_count + 1 WHERE id = NEW.category_id;
END //
DELIMITER ;

-- Trigger para decrementar el contador cuando se elimina un producto
DELIMITER //
CREATE TRIGGER IF NOT EXISTS trg_product_delete AFTER DELETE ON products
FOR EACH ROW
BEGIN
    UPDATE categories SET product_count = product_count - 1 WHERE id = OLD.category_id;
END //
DELIMITER ;

-- Trigger para actualizar el contador cuando se cambia la categoría de un producto
DELIMITER //
CREATE TRIGGER IF NOT EXISTS trg_product_update AFTER UPDATE ON products
FOR EACH ROW
BEGIN
    IF OLD.category_id != NEW.category_id THEN
        UPDATE categories SET product_count = product_count - 1 WHERE id = OLD.category_id;
        UPDATE categories SET product_count = product_count + 1 WHERE id = NEW.category_id;
    END IF;
END //
DELIMITER ;

-- Insertar las categorías iniciales si no existen
INSERT INTO categories (name, slug) VALUES
('Portátiles', 'portatiles'),
('Tarjetas Gráficas', 'tarjetas-graficas'),
('Auriculares', 'auriculares'),
('Teclados', 'teclados'),
('Monitores', 'monitores'),
('Discos SSD', 'ssd')
ON DUPLICATE KEY UPDATE name = VALUES(name);

-- Actualizar los contadores de productos para las categorías existentes
UPDATE categories c
SET product_count = (SELECT COUNT(*) FROM products p WHERE p.category_id = c.id AND p.deleted_at IS NULL);

-- Los productos y precios serán añadidos por el scraper, Este script solo establece la estructura inicial de la base de datos si es que no existe
