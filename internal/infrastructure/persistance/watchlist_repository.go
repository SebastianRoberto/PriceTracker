package persistance

import (
	"context"
	"errors"

	"app/internal/domain/model"
	"app/internal/domain/repositories"

	"gorm.io/gorm"
)

// watchlistRepository implementa la interfaz WatchlistRepository
type watchlistRepository struct {
	db *gorm.DB
}

// NewWatchlistRepository crea una nueva instancia del repositorio de listas de seguimiento
func NewWatchlistRepository(db *gorm.DB) repositories.WatchlistRepository {
	return &watchlistRepository{
		db: db,
	}
}

// Create crea una nueva lista de seguimiento en la base de datos
func (r *watchlistRepository) Create(ctx context.Context, watchlist *model.Watchlist) error {
	return r.db.WithContext(ctx).Create(watchlist).Error
}

// Update actualiza una lista de seguimiento existente
func (r *watchlistRepository) Update(ctx context.Context, watchlist *model.Watchlist) error {
	return r.db.WithContext(ctx).Save(watchlist).Error
}

// Delete elimina una lista de seguimiento por su ID
func (r *watchlistRepository) Delete(ctx context.Context, watchlistID uint) error {
	return r.db.WithContext(ctx).Delete(&model.Watchlist{}, watchlistID).Error
}

// FindByID busca una lista de seguimiento por su ID
func (r *watchlistRepository) FindByID(ctx context.Context, watchlistID uint) (*model.Watchlist, error) {
	var watchlist model.Watchlist
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		First(&watchlist, watchlistID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("lista de seguimiento no encontrada")
		}
		return nil, err
	}
	return &watchlist, nil
}

// FindByUserID busca lista de seguimiento por ID de usuario
func (r *watchlistRepository) FindByUserID(ctx context.Context, userID uint) (*model.Watchlist, error) {
	var watchlist model.Watchlist
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Items").
		Preload("Items.Product").
		First(&watchlist).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Si no existe, creamos una nueva lista de seguimiento
			watchlist = model.Watchlist{
				UserID: userID,
				Name:   "Mi lista de seguimiento",
				Items:  []model.WatchlistItem{},
			}
			if createErr := r.Create(ctx, &watchlist); createErr != nil {
				return nil, createErr
			}
			return &watchlist, nil
		}
		return nil, err
	}
	return &watchlist, nil
}

// watchlistItemRepository implementa la interfaz WatchlistItemRepository
type watchlistItemRepository struct {
	db *gorm.DB
}

// NewWatchlistItemRepository crea una nueva instancia del repositorio de elementos en la lista de seguimiento
func NewWatchlistItemRepository(db *gorm.DB) repositories.WatchlistItemRepository {
	return &watchlistItemRepository{
		db: db,
	}
}

// Create añade un producto a la lista de seguimiento
func (r *watchlistItemRepository) Create(ctx context.Context, item *model.WatchlistItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// Update actualiza un elemento existente en la lista de seguimiento
func (r *watchlistItemRepository) Update(ctx context.Context, item *model.WatchlistItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// Delete elimina un elemento de la lista de seguimiento por su ID
func (r *watchlistItemRepository) Delete(ctx context.Context, itemID uint) error {
	return r.db.WithContext(ctx).Delete(&model.WatchlistItem{}, itemID).Error
}

// IsProductInWatchlist verifica si un producto está en la lista de seguimiento de un usuario
func (r *watchlistItemRepository) IsProductInWatchlist(ctx context.Context, userID, productID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.WatchlistItem{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByID busca un elemento por su ID
func (r *watchlistItemRepository) FindByID(ctx context.Context, itemID uint) (*model.WatchlistItem, error) {
	var item model.WatchlistItem
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Product").
		First(&item, itemID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("elemento no encontrado")
		}
		return nil, err
	}
	return &item, nil
}

// FindByUserID busca todos los elementos de la lista de seguimiento de un usuario
func (r *watchlistItemRepository) FindByUserID(ctx context.Context, userID uint) ([]*model.WatchlistItem, error) {
	var items []*model.WatchlistItem
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Product").
		Preload("Product.Category").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindByProductID busca elementos por ID de producto
func (r *watchlistItemRepository) FindByProductID(ctx context.Context, productID uint) ([]*model.WatchlistItem, error) {
	var items []*model.WatchlistItem
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Preload("User").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
