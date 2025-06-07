package repositories

import (
	"context"
	"time"

	"app/internal/domain/model"
)

// PriceAlertRepository define las operaciones para el repositorio de alertas de precio
type PriceAlertRepository interface {
	// Crear una nueva alerta de precio
	Create(ctx context.Context, alert *model.PriceAlert) error

	// Actualizar una alerta existente
	Update(ctx context.Context, alert *model.PriceAlert) error

	// Eliminar una alerta por su ID
	Delete(ctx context.Context, alertID uint) error

	// Buscar una alerta por su ID
	FindByID(ctx context.Context, alertID uint) (*model.PriceAlert, error)

	// Buscar alertas por ID de usuario
	FindByUserID(ctx context.Context, userID uint) ([]*model.PriceAlert, error)

	// Buscar alertas por ID de producto
	FindByProductID(ctx context.Context, productID uint) ([]*model.PriceAlert, error)

	// Buscar alertas activas para un precio específico
	// Devuelve alertas donde targetPrice >= nuevoPrice
	FindActiveAlertsForPrice(ctx context.Context, productID uint, newPrice float64) ([]*model.PriceAlert, error)
}

// NotificationRepository define las operaciones para gestionar notificaciones
type NotificationRepository interface {
	// Create crea una nueva notificación
	Create(ctx context.Context, notification *model.Notification) error

	// Update actualiza una notificación existente
	Update(ctx context.Context, notification *model.Notification) error

	// Delete elimina una notificación
	Delete(ctx context.Context, id uint) error

	// FindByID busca una notificación por su ID
	FindByID(ctx context.Context, id uint) (*model.Notification, error)

	// FindByUserID busca las notificaciones de un usuario
	FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]*model.Notification, error)

	// FindByProductID busca las notificaciones para un producto
	FindByProductID(ctx context.Context, productID uint, limit, offset int) ([]*model.Notification, error)

	// FindByAlertID busca las notificaciones relacionadas con una alerta específica
	FindByAlertID(ctx context.Context, alertID uint) ([]*model.Notification, error)

	// FindUnreadByUserID busca las notificaciones no leídas de un usuario
	FindUnreadByUserID(ctx context.Context, userID uint) ([]*model.Notification, error)

	// CountUnreadByUserID cuenta el número de notificaciones no leídas de un usuario
	CountUnreadByUserID(ctx context.Context, userID uint) (int, error)

	// Marcar notificación como leída
	MarkAsRead(ctx context.Context, notificationID uint) error

	// Marcar todas las notificaciones de un usuario como leídas
	MarkAllAsRead(ctx context.Context, userID uint) error

	// Eliminar notificaciones antiguas (más de cierto tiempo)
	DeleteOldNotifications(ctx context.Context, olderThan time.Time) error
}
