package persistance

import (
	"context"
	"errors"
	"time"

	"app/internal/domain/model"
	"app/internal/domain/repositories"

	"gorm.io/gorm"
)

// priceRepository implementa la interfaz PriceRepository
type priceRepository struct {
	db *gorm.DB
}

// NewPriceRepository crea una nueva instancia del repositorio de precios
func NewPriceRepository(db *gorm.DB) repositories.PriceRepository {
	return &priceRepository{
		db: db,
	}
}

// Create crea un nuevo precio en la base de datos
func (r *priceRepository) Create(ctx context.Context, price *model.Price) error {
	return r.db.WithContext(ctx).Create(price).Error
}

// FindByID busca un precio por su ID
func (r *priceRepository) FindByID(ctx context.Context, id uint) (*model.Price, error) {
	var price model.Price
	if err := r.db.WithContext(ctx).First(&price, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("precio no encontrado")
		}
		return nil, err
	}
	return &price, nil
}

// FindByProductID busca precios por ID de producto
func (r *priceRepository) FindByProductID(ctx context.Context, productID uint) ([]*model.Price, error) {
	var prices []*model.Price
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&prices).Error; err != nil {
		return nil, err
	}
	return prices, nil
}

// FindBestPriceByProductID busca el mejor precio para un producto
func (r *priceRepository) FindBestPriceByProductID(ctx context.Context, productID uint) (*model.Price, error) {
	var price model.Price

	// Verificar si el contexto ya está cancelado antes de iniciar la consulta
	if ctx.Err() != nil {
		// Si el contexto ya está cancelado, devolvemos nil sin error
		// para que no se propague el error de contexto cancelado
		return nil, nil
	}

	err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_available = ?", productID, true).
		Order("price asc").
		Limit(1).
		First(&price).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Devolver nil si no hay precio disponible
		}
		if errors.Is(err, context.Canceled) {
			// Si el contexto fue cancelado durante la consulta, también devolvemos
			// nil sin error para no propagar este error
			return nil, nil
		}
		return nil, err
	}

	return &price, nil
}

// FindTopOffersByProductID busca las mejores ofertas para un producto
func (r *priceRepository) FindTopOffersByProductID(ctx context.Context, productID uint, limit int) ([]*model.Price, error) {
	var prices []*model.Price

	// Verificar si el contexto ya está cancelado antes de iniciar la consulta
	if ctx.Err() != nil {
		return nil, nil
	}

	err := r.db.WithContext(ctx).
		Where("product_id = ? AND is_available = ?", productID, true).
		Order("price asc").
		Limit(limit).
		Find(&prices).Error

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, nil
		}
		return nil, err
	}

	return prices, nil
}

// Update actualiza un precio existente
func (r *priceRepository) Update(ctx context.Context, price *model.Price) error {
	return r.db.WithContext(ctx).Save(price).Error
}

// Delete elimina un precio de la base de datos
func (r *priceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Price{}, id).Error
}

// DeleteOldPrices elimina precios más antiguos que una fecha dada
// Devuelve el número de precios eliminados y un error si hubo problemas
func (r *priceRepository) DeleteOldPrices(ctx context.Context, olderThan time.Time) (int, error) {
	result := r.db.WithContext(ctx).Where("retrieved_at < ?", olderThan).Delete(&model.Price{})
	return int(result.RowsAffected), result.Error
}
