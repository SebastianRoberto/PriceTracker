package repositories

import (
	"context"
	"database/sql"

	"app/internal/domain/model"
)

// CategoryRepository define las operaciones de persistencia para las categorías
type CategoryRepository interface {
	// Create crea una nueva categoría en la base de datos
	Create(ctx context.Context, category *model.Category) error

	// FindByID busca una categoría por su ID
	FindByID(ctx context.Context, id uint) (*model.Category, error)

	// FindBySlug busca una categoría por su slug
	FindBySlug(ctx context.Context, slug string) (*model.Category, error)

	// GetAll obtiene todas las categorías
	GetAll(ctx context.Context) ([]*model.Category, error)

	// Update actualiza una categoría existente
	Update(ctx context.Context, category *model.Category) error

	// Delete elimina una categoría de la base de datos
	Delete(ctx context.Context, id uint) error

	// GetCategoryWithProductCount obtiene una categoría con la cantidad de productos
	GetCategoryWithProductCount(categoryID int) (*model.Category, error)

	// GetAllCategoriesWithProductCount obtiene todas las categorías con la cantidad de productos
	GetAllCategoriesWithProductCount() ([]model.Category, error)
}

type CategoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{db: db}
}

func (r *CategoryRepositoryImpl) GetCategoryWithProductCount(categoryID int) (*model.Category, error) {
	query := `
		SELECT c.id, c.name, c.slug, COUNT(p.id) as product_count
		FROM categories c
		LEFT JOIN products p ON c.id = p.category_id
		WHERE c.id = $1 AND p.active = true
		GROUP BY c.id, c.name, c.slug`

	var category model.Category
	err := r.db.QueryRow(query, categoryID).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.ProductCount,
	)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepositoryImpl) GetAllCategoriesWithProductCount() ([]model.Category, error) {
	query := `
		SELECT c.id, c.name, c.slug, COUNT(p.id) as product_count
		FROM categories c
		LEFT JOIN products p ON c.id = p.category_id AND p.active = true
		GROUP BY c.id, c.name, c.slug
		ORDER BY c.name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.ProductCount,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
