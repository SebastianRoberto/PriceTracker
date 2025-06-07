package model

import (
	"time"
)

// PriceAlert representa una alerta configurada por un usuario para recibir notificaciones
// cuando un producto alcance un precio igual o menor al establecido
type PriceAlert struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index:idx_alert_user" json:"user_id"`
	ProductID     uint      `gorm:"not null;index:idx_alert_product" json:"product_id"`
	TargetPrice   float64   `gorm:"not null" json:"target_price"`
	NotifyByEmail bool      `gorm:"default:true" json:"notify_by_email"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Relaciones
	User    User    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// Notification representa una notificación enviada a un usuario sobre un cambio de precio
type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_notification_user" json:"user_id"`
	ProductID uint      `gorm:"not null;index:idx_notification_product" json:"product_id"`
	AlertID   *uint     `gorm:"index:idx_notification_alert" json:"alert_id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Message   string    `gorm:"size:1000;not null" json:"message"`
	IsRead    bool      `gorm:"default:false;index:idx_notification_read" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`

	// Relaciones
	User       User        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Product    Product     `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	PriceAlert *PriceAlert `gorm:"foreignKey:AlertID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}
