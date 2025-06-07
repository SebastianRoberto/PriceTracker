package repositories

import (
	"context"

	"app/internal/domain/model"
)

// WatchlistRepository define las operaciones para el repositorio de listas de seguimiento
type WatchlistRepository interface {
	// Crear una nueva lista de seguimiento
	Create(ctx context.Context, watchlist *model.Watchlist) error

	// Actualizar una lista de seguimiento existente
	Update(ctx context.Context, watchlist *model.Watchlist) error

	// Eliminar una lista de seguimiento por su ID
	Delete(ctx context.Context, watchlistID uint) error

	// Buscar una lista de seguimiento por su ID
	FindByID(ctx context.Context, watchlistID uint) (*model.Watchlist, error)

	// Buscar lista de seguimiento por ID de usuario
	FindByUserID(ctx context.Context, userID uint) (*model.Watchlist, error)
}

// WatchlistItemRepository define las operaciones para el repositorio de elementos en la lista de seguimiento
type WatchlistItemRepository interface {
	// Añadir un producto a la lista de seguimiento
	Create(ctx context.Context, item *model.WatchlistItem) error

	// Actualizar un elemento existente en la lista de seguimiento
	Update(ctx context.Context, item *model.WatchlistItem) error

	// Eliminar un elemento de la lista de seguimiento por su ID
	Delete(ctx context.Context, itemID uint) error

	// Verificar si un producto está en la lista de seguimiento de un usuario
	IsProductInWatchlist(ctx context.Context, userID, productID uint) (bool, error)

	// Buscar un elemento por su ID
	FindByID(ctx context.Context, itemID uint) (*model.WatchlistItem, error)

	// Buscar todos los elementos de la lista de seguimiento de un usuario
	FindByUserID(ctx context.Context, userID uint) ([]*model.WatchlistItem, error)

	// Buscar elementos por ID de producto
	FindByProductID(ctx context.Context, productID uint) ([]*model.WatchlistItem, error)
}
