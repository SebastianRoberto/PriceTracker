package persistance

import (
	"context"
	"errors"

	"app/internal/domain/model"
	"app/internal/domain/repositories"

	"gorm.io/gorm"
)

// categoryRepository implementa la interfaz CategoryRepository
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository crea una nueva instancia del repositorio de categorías
func NewCategoryRepository(db *gorm.DB) repositories.CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

// Create crea una nueva categoría en la base de datos
func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// FindByID busca una categoría por su ID
func (r *categoryRepository) FindByID(ctx context.Context, id uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("categoría no encontrada")
		}
		return nil, err
	}
	return &category, nil
}

// FindBySlug busca una categoría por su slug
func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*model.Category, error) {
	var category model.Category
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("categoría no encontrada")
		}
		return nil, err
	}
	return &category, nil
}

// GetAll obtiene todas las categorías
func (r *categoryRepository) GetAll(ctx context.Context) ([]*model.Category, error) {
	var categories []*model.Category
	if err := r.db.WithContext(ctx).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// Update actualiza una categoría existente
func (r *categoryRepository) Update(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete elimina una categoría de la base de datos
func (r *categoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Category{}, id).Error
}

// GetAllCategoriesWithProductCount obtiene todas las categorías con su conteo de productos
func (r *categoryRepository) GetAllCategoriesWithProductCount() ([]model.Category, error) {
	var categories []model.Category

	// Consulta SQL que cuenta productos activos por categoría
	err := r.db.Raw(`
		SELECT c.*, COUNT(p.id) as product_count
		FROM categories c
		LEFT JOIN products p ON c.id = p.category_id AND p.deleted_at IS NULL
		GROUP BY c.id
		ORDER BY c.name
	`).Scan(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryWithProductCount obtiene una categoría específica con su conteo de productos
func (r *categoryRepository) GetCategoryWithProductCount(categoryID int) (*model.Category, error) {
	var category model.Category

	// Consulta SQL que cuenta productos activos para una categoría específica
	err := r.db.Raw(`
		SELECT c.*, COUNT(p.id) as product_count
		FROM categories c
		LEFT JOIN products p ON c.id = p.category_id AND p.deleted_at IS NULL
		WHERE c.id = ?
		GROUP BY c.id
	`, categoryID).Scan(&category).Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}
