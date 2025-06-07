package persistance

import (
	"database/sql"
	"fmt"
	"log"

	"app/internal/domain/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database representa el acceso a la base de datos
type Database struct {
	*gorm.DB
}

// NewDatabase crea una nueva conexión a la base de datos
func NewDatabase(username, password, host, port, name string) (*Database, error) {
	// Primero, intentamos crear la base de datos si no existe
	rootDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
	)

	// Conectar a MySQL sin una base de datos específica
	sqlDB, err := sql.Open("mysql", rootDsn)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a MySQL: %w", err)
	}
	defer sqlDB.Close()

	// Crear base de datos si no existe
	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name))
	if err != nil {
		return nil, fmt.Errorf("error al crear la base de datos: %w", err)
	}
	log.Printf("Base de datos '%s' verificada/creada correctamente", name)

	// Ahora conectamos a la base de datos específica
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	log.Println("Conexión a la base de datos establecida")
	return &Database{DB: db}, nil
}

// AutoMigrate ejecuta las migraciones automáticas de GORM para todas las entidades
func (d *Database) AutoMigrate() error {
	if err := d.DB.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.Product{},
		&model.Price{},
		&model.PriceAlert{},
		&model.Notification{},
		&model.Watchlist{},
		&model.WatchlistItem{},
	); err != nil {
		return fmt.Errorf("error al migrar la base de datos: %w", err)
	}

	log.Println("Migraciones de la base de datos completadas")
	return nil
}

// Close cierra la conexión a la base de datos
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
