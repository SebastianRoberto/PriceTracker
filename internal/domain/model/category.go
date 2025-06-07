package model

import (
	"time"

	"gorm.io/gorm"
)

// Category representa una categoría de productos (ej: Portátiles, GPUs, etc.)
type Category struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"uniqueIndex;not null;size:100" json:"name"`
	Slug         string    `gorm:"uniqueIndex;not null;size:100" json:"slug"` // Para URLs amigables
	Products     []Product `gorm:"foreignKey:CategoryID"`
	ProductCount int       `json:"product_count"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
