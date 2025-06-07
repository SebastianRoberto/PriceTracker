package persistance

import (
	"context"
	"errors"

	"app/internal/domain/model"
	"app/internal/domain/repositories"

	"gorm.io/gorm"
)

// priceAlertRepository implementa la interfaz PriceAlertRepository
type priceAlertRepository struct {
	db *gorm.DB
}

// NewPriceAlertRepository crea una nueva instancia del repositorio de alertas de precio
func NewPriceAlertRepository(db *gorm.DB) repositories.PriceAlertRepository {
	return &priceAlertRepository{
		db: db,
	}
}

// Create crea una nueva alerta de precio en la base de datos
func (r *priceAlertRepository) Create(ctx context.Context, alert *model.PriceAlert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

// Update actualiza una alerta de precio existente
func (r *priceAlertRepository) Update(ctx context.Context, alert *model.PriceAlert) error {
	return r.db.WithContext(ctx).Save(alert).Error
}

// Delete elimina una alerta de precio por su ID
func (r *priceAlertRepository) Delete(ctx context.Context, alertID uint) error {
	return r.db.WithContext(ctx).Delete(&model.PriceAlert{}, alertID).Error
}

// FindByID busca una alerta de precio por su ID
func (r *priceAlertRepository) FindByID(ctx context.Context, alertID uint) (*model.PriceAlert, error) {
	var alert model.PriceAlert
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		First(&alert, alertID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alerta de precio no encontrada")
		}
		return nil, err
	}
	return &alert, nil
}

// FindByUserID busca alertas de precio por ID de usuario
func (r *priceAlertRepository) FindByUserID(ctx context.Context, userID uint) ([]*model.PriceAlert, error) {
	var alerts []*model.PriceAlert
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Preload("Product").
		Find(&alerts).Error; err != nil {
		return nil, err
	}
	return alerts, nil
}

// FindByProductID busca alertas de precio por ID de producto
func (r *priceAlertRepository) FindByProductID(ctx context.Context, productID uint) ([]*model.PriceAlert, error) {
	var alerts []*model.PriceAlert
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_active = ?", productID, true).
		Preload("User").
		Find(&alerts).Error; err != nil {
		return nil, err
	}
	return alerts, nil
}

// FindActiveAlertsForPrice busca alertas activas para un precio específico
// Devuelve alertas donde targetPrice >= nuevoPrice
func (r *priceAlertRepository) FindActiveAlertsForPrice(ctx context.Context, productID uint, newPrice float64) ([]*model.PriceAlert, error) {
	var alerts []*model.PriceAlert
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_active = ? AND target_price >= ?", productID, true, newPrice).
		Preload("User").
		Preload("Product").
		Find(&alerts).Error; err != nil {
		return nil, err
	}
	return alerts, nil
}
