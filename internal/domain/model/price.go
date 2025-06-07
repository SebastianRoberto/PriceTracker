package model

import (
	"time"

	"gorm.io/gorm"
)

// Price representa una oferta de precio para un producto específico
type Price struct {
	ID          uint      `gorm:"primaryKey"`
	ProductID   uint      `gorm:"index;not null"`
	Product     Product   `gorm:"foreignKey:ProductID"`
	Store       string    `gorm:"not null;size:50"` // Tienda: PCComponentes, MercadoLibre, eBay
	Price       float64   `gorm:"not null"`
	Currency    string    `gorm:"size:3;default:'EUR'"` // EUR, USD, MXN, etc.
	URL         string    `gorm:"not null;size:1024"`   // URL para comprar el producto
	IsAvailable bool      `gorm:"default:true"`
	RetrievedAt time.Time `gorm:"not null"` // Cuándo se obtuvo este precio
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
