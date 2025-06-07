package model

import (
	"time"

	"gorm.io/gorm"
)

// User representa el modelo de usuario para autenticación y gestión de sesiones
type User struct {
	ID                 uint   `gorm:"primaryKey"`
	Username           string `gorm:"uniqueIndex;not null;size:100"`
	Email              string `gorm:"uniqueIndex;not null;size:100"`
	PasswordHash       string `gorm:"column:password_hash;not null;type:varchar(255)"`
	Verified           bool   `gorm:"default:false"`
	VerifyToken        string `gorm:"size:100"`
	EmailNotifications bool   `gorm:"default:true"`
	IsAdmin            bool   `gorm:"default:false"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}
