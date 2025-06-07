package model

import (
	"time"
)

// WatchlistItem representa un producto en la lista de seguimiento de un usuario
type WatchlistItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;uniqueIndex:idx_watchlist_item_user_product" json:"user_id"`
	ProductID   uint      `gorm:"not null;uniqueIndex:idx_watchlist_item_user_product" json:"product_id"`
	TargetPrice float64   `gorm:"type:decimal(10,2)" json:"target_price"`
	Notes       string    `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relaciones
	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// Watchlist representa la lista de seguimiento completa de un usuario
type Watchlist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;uniqueIndex" json:"user_id"`
	Name      string    `gorm:"size:100;not null;default:'Mi Lista de Seguimiento'" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relaciones
	User  User            `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Items []WatchlistItem `gorm:"foreignKey:UserID;references:UserID" json:"items"`
}
