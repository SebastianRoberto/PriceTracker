package repositories

import (
	"context"

	"app/internal/domain/model"
)

// ProductRepository define las operaciones de persistencia para los productos
type ProductRepository interface {
	// Create crea un nuevo producto en la base de datos
	Create(ctx context.Context, product *model.Product) error

	// FindByID busca un producto por su ID
	FindByID(ctx context.Context, id uint) (*model.Product, error)

	// FindBySlug busca un producto por su slug
	FindBySlug(ctx context.Context, slug string) (*model.Product, error)

	// FindByCategory busca productos por categoría
	FindByCategory(ctx context.Context, categoryID uint, limit, offset int, storeFilter string) ([]*model.Product, error)

	// FindFilteredProductsByCategory busca productos por categoría con filtros avanzados
	FindFilteredProductsByCategory(ctx context.Context, options model.ProductFilterOptions) ([]*model.Product, error)

	// FindBestDeals obtiene los productos con mejores ofertas (precio más bajo)
	FindBestDeals(ctx context.Context, limit int) ([]*model.Product, error)

	// Update actualiza un producto existente
	Update(ctx context.Context, product *model.Product) error

	// Delete elimina un producto de la base de datos
	Delete(ctx context.Context, id uint) error

	// ExistsBySlug verifica si un slug existente ya está en la base de datos
	ExistsBySlug(ctx context.Context, slug string) (bool, error)

	// FindSimilarProducts encuentra productos similares a un producto dado
	FindSimilarProducts(ctx context.Context, productID uint, limit int) ([]*model.Product, error)

	// CountByCategory cuenta el número total de productos en una categoría
	CountByCategory(ctx context.Context, categoryID uint, storeFilter string) (int, error)

	// CountFilteredProductsByCategory cuenta el número total de productos en una categoría con filtros
	CountFilteredProductsByCategory(ctx context.Context, options model.ProductFilterOptions) (int, error)
}
