package repositories

import (
	"context"
	"time"

	"app/internal/domain/model"
)

// PriceRepository define las operaciones de persistencia para los precios
type PriceRepository interface {
	// Create crea un nuevo precio en la base de datos
	Create(ctx context.Context, price *model.Price) error

	// FindByID busca un precio por su ID
	FindByID(ctx context.Context, id uint) (*model.Price, error)

	// FindByProductID busca precios por ID de producto
	FindByProductID(ctx context.Context, productID uint) ([]*model.Price, error)

	// FindBestPriceByProductID busca el mejor precio para un producto
	FindBestPriceByProductID(ctx context.Context, productID uint) (*model.Price, error)

	// FindTopOffersByProductID busca las mejores ofertas para un producto
	FindTopOffersByProductID(ctx context.Context, productID uint, limit int) ([]*model.Price, error)

	// Update actualiza un precio existente
	Update(ctx context.Context, price *model.Price) error

	// Delete elimina un precio de la base de datos
	Delete(ctx context.Context, id uint) error

	// DeleteOldPrices elimina precios más antiguos que una fecha dada
	// Devuelve el número de precios eliminados y un error si hubo problemas
	DeleteOldPrices(ctx context.Context, olderThan time.Time) (int, error)
}
