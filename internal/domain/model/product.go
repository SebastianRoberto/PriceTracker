package model

import (
	"time"

	"gorm.io/gorm"
)

// Product representa un producto que será scrapeado de diferentes tiendas
type Product struct {
	ID             uint              `gorm:"primaryKey"`
	Name           string            `gorm:"not null;size:200"`
	Slug           string            `gorm:"uniqueIndex;size:100"`
	Description    string            `gorm:"type:text"`
	ImageURL       string            `gorm:"size:255"`
	CategoryID     uint              `gorm:"index"`
	Category       Category          `gorm:"foreignKey:CategoryID"`
	Specifications map[string]string `gorm:"-:all"` // Se ignora en GORM para simplificar
	Prices         []Price           `gorm:"foreignKey:ProductID"`
	ImageHash      *uint64           `gorm:"column:image_hash;type:BIGINT UNSIGNED NULL"` // Hash de percepción de la imagen para deduplicación
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// BeforeCreate hook para GORM para generar Slug antes de crear un producto
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	// No generar slug aquí, se genera en el Usecase para asegurar unicidad
	return
}

// ProductFilterOptions contiene opciones para filtrar y ordenar productos
type ProductFilterOptions struct {
	CategorySlug string  // Slug de la categoría
	Limit        int     // Número máximo de productos a devolver
	Offset       int     // Desplazamiento para paginación
	StoreFilter  string  // Filtrar por tienda
	SortOrder    string  // Orden de clasificación (asc/desc)
	MinPrice     float64 // Precio mínimo (opcional)
	MaxPrice     float64 // Precio máximo (opcional)
}
