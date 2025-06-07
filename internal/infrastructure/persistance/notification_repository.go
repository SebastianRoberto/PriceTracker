package persistance

import (
	"context"
	"errors"
	"time"

	"app/internal/domain/model"
	"app/internal/domain/repositories"

	"gorm.io/gorm"
)

// notificationRepository implementa la interfaz NotificationRepository
type notificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository crea una nueva instancia del repositorio de notificaciones
func NewNotificationRepository(db *gorm.DB) repositories.NotificationRepository {
	return &notificationRepository{
		db: db,
	}
}

// Create crea una nueva notificación en la base de datos
func (r *notificationRepository) Create(ctx context.Context, notification *model.Notification) error {
	return r.db.WithContext(ctx).Create(notification).Error
}

// Update actualiza una notificación existente
func (r *notificationRepository) Update(ctx context.Context, notification *model.Notification) error {
	return r.db.WithContext(ctx).Save(notification).Error
}

// Delete elimina una notificación por su ID
func (r *notificationRepository) Delete(ctx context.Context, notificationID uint) error {
	return r.db.WithContext(ctx).Delete(&model.Notification{}, notificationID).Error
}

// MarkAsRead marca una notificación como leída
func (r *notificationRepository) MarkAsRead(ctx context.Context, notificationID uint) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("id = ?", notificationID).
		Update("is_read", true).Error
}

// MarkAllAsRead marca todas las notificaciones de un usuario como leídas
func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&model.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

// FindByID busca una notificación por su ID
func (r *notificationRepository) FindByID(ctx context.Context, notificationID uint) (*model.Notification, error) {
	var notification model.Notification
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		Preload("PriceAlert").
		First(&notification, notificationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notificación no encontrada")
		}
		return nil, err
	}
	return &notification, nil
}

// FindByUserID busca notificaciones por ID de usuario
func (r *notificationRepository) FindByUserID(ctx context.Context, userID uint, limit, offset int) ([]*model.Notification, error) {
	var notifications []*model.Notification
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Product").
		Preload("PriceAlert").
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, err
	}
	return notifications, nil
}

// CountUnreadByUserID cuenta las notificaciones no leídas para un usuario
func (r *notificationRepository) CountUnreadByUserID(ctx context.Context, userID uint) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

// DeleteOldNotifications elimina notificaciones antiguas
func (r *notificationRepository) DeleteOldNotifications(ctx context.Context, olderThan time.Time) error {
	return r.db.WithContext(ctx).
		Where("created_at < ?", olderThan).
		Delete(&model.Notification{}).Error
}

// FindUnreadByUserID busca todas las notificaciones no leídas de un usuario
func (r *notificationRepository) FindUnreadByUserID(ctx context.Context, userID uint) ([]*model.Notification, error) {
	var notifications []*model.Notification
	result := r.db.WithContext(ctx).Where("user_id = ? AND is_read = ?", userID, false).Find(&notifications)
	return notifications, result.Error
}

// FindByProductID busca las notificaciones para un producto específico
func (r *notificationRepository) FindByProductID(ctx context.Context, productID uint, limit, offset int) ([]*model.Notification, error) {
	var notifications []*model.Notification
	result := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications)
	return notifications, result.Error
}

// FindByAlertID busca las notificaciones relacionadas con una alerta específica
func (r *notificationRepository) FindByAlertID(ctx context.Context, alertID uint) ([]*model.Notification, error) {
	var notifications []*model.Notification
	result := r.db.WithContext(ctx).
		Where("alert_id = ?", alertID).
		Find(&notifications)
	return notifications, result.Error
}
