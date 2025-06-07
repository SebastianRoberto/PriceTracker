package persistance

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"app/internal/domain/model"
	"app/internal/domain/repositories"

	"gorm.io/gorm"
)

// productRepository implementa la interfaz ProductRepository
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository crea una nueva instancia del repositorio de productos
func NewProductRepository(db *gorm.DB) repositories.ProductRepository {
	return &productRepository{
		db: db,
	}
}

// Create crea un nuevo producto en la base de datos
func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

// FindByID busca un producto por su ID
func (r *productRepository) FindByID(ctx context.Context, id uint) (*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).Preload("Category").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("producto no encontrado")
		}
		return nil, err
	}
	return &product, nil
}

// FindBySlug busca un producto por su slug
func (r *productRepository) FindBySlug(ctx context.Context, slug string) (*model.Product, error) {
	var product model.Product
	// Asegurarnos de cargar todos los campos incluyendo ImageHash
	if err := r.db.WithContext(ctx).Preload("Category").Where("slug = ?", slug).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("producto no encontrado")
		}
		return nil, err
	}
	return &product, nil
}

// FindByCategory busca productos por categoría, opcionalmente filtrados por tienda
func (r *productRepository) FindByCategory(ctx context.Context, categoryID uint, limit, offset int, storeFilter string) ([]*model.Product, error) {
	var products []*model.Product
	query := r.db.WithContext(ctx).Model(&model.Product{}).Preload("Category") // Siempre cargar la categoría y ahora también image_hash

	if storeFilter != "" {
		// Seleccionar explícitamente columnas de 'products' para evitar ambigüedad y asegurar Preload.
		query = query.Joins("JOIN prices ON prices.product_id = products.id").
			Where("products.category_id = ? AND prices.store = ?", categoryID, storeFilter).
			Distinct("products.*") // Asegura que cada producto se liste solo una vez.
	} else {
		query = query.Where("products.category_id = ?", categoryID)
	}

	// Cargar también el ImageHash
	query = query.Select("products.*, products.image_hash")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// FindBestDeals obtiene los productos con mejores ofertas
func (r *productRepository) FindBestDeals(ctx context.Context, limit int) ([]*model.Product, error) {
	// Esta consulta es más compleja, necesitamos encontrar productos con los precios más bajos
	// Usamos un subquery para encontrar el precio mínimo por producto
	var products []*model.Product

	// Primero obtenemos los IDs de los productos con los precios más bajos
	type Result struct {
		ProductID uint
	}

	var results []Result

	// Subconsulta para obtener los productos con los precios más bajos
	subQuery := r.db.WithContext(ctx).
		Table("prices").
		Select("product_id, MIN(price) as min_price").
		Where("is_available = ?", true).
		Group("product_id").
		Order("min_price asc").
		Limit(limit)

	// Consulta principal para obtener los IDs de los productos
	if err := r.db.WithContext(ctx).
		Table("(?) as p", subQuery).
		Select("product_id").
		Find(&results).Error; err != nil {
		return nil, err
	}

	// Convertir resultados en slice de IDs
	var productIDs []uint
	for _, result := range results {
		productIDs = append(productIDs, result.ProductID)
	}

	// Si no hay productos, devolver un slice vacío
	if len(productIDs) == 0 {
		return products, nil
	}

	// Obtener los productos completos por sus IDs
	if err := r.db.WithContext(ctx).
		Preload("Category").
		Where("id IN ?", productIDs).
		Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// Update actualiza un producto existente
func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	// Asegurarnos de que se actualizan todos los campos, incluyendo ImageHash
	return r.db.WithContext(ctx).Save(product).Error
}

// Delete elimina un producto de la base de datos
func (r *productRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Product{}, id).Error
}

// ExistsBySlug verifica si un slug ya existe (true) o no (false)
func (r *productRepository) ExistsBySlug(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Product{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindSimilarProducts encuentra productos similares a un producto dado
func (r *productRepository) FindSimilarProducts(ctx context.Context, productID uint, limit int) ([]*model.Product, error) {
	var product model.Product
	if err := r.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		return nil, err
	}

	var similarProducts []*model.Product
	err := r.db.WithContext(ctx).
		Where("category_id = ? AND id != ?", product.CategoryID, productID).
		Limit(limit).
		Preload("Category").
		Find(&similarProducts).Error

	if err != nil {
		return nil, err
	}

	return similarProducts, nil
}

// CountByCategory cuenta el número total de productos en una categoría
func (r *productRepository) CountByCategory(ctx context.Context, categoryID uint, storeFilter string) (int, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Product{})

	if storeFilter != "" {
		// Si hay un filtro de tienda, necesitamos unir con la tabla de precios
		query = query.Joins("JOIN prices ON prices.product_id = products.id").
			Where("products.category_id = ? AND prices.store = ?", categoryID, storeFilter).
			Distinct("products.id") // Contar productos únicos
	} else {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

// FindFilteredProductsByCategory busca productos por categoría con filtros de precio y ordenación
func (r *productRepository) FindFilteredProductsByCategory(ctx context.Context, options model.ProductFilterOptions) ([]*model.Product, error) {
	startTime := time.Now()
	log.Printf("[SQL_DEBUG] Iniciando consulta de productos filtrados para categoría: %s", options.CategorySlug)
	log.Printf("[SQL_DEBUG] Filtros: minPrice=%.2f, maxPrice=%.2f, orden=%s, store=%s, limit=%d, offset=%d",
		options.MinPrice, options.MaxPrice, options.SortOrder, options.StoreFilter, options.Limit, options.Offset)

	var products []*model.Product

	// Construir una subconsulta para obtener el mejor precio por producto
	// Modificamos para no seleccionar 'store' directamente en la subconsulta, ya que causa problemas con GROUP BY
	subQuery := r.db.WithContext(ctx).
		Table("prices").
		Select("product_id, MIN(price) as min_price").
		Group("product_id")

	if options.StoreFilter != "" {
		subQuery = subQuery.Where("store = ?", options.StoreFilter)
	}

	subQuery = subQuery.Order("min_price " + options.SortOrder)

	// Consulta principal uniendo productos con la subconsulta de mejores precios
	query := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Joins("JOIN (?) AS best_prices ON products.id = best_prices.product_id", subQuery).
		Joins("JOIN categories ON products.category_id = categories.id").
		Where("categories.slug = ?", options.CategorySlug)

	// Aplicar filtros de precio
	if options.MinPrice > 0 {
		query = query.Where("best_prices.min_price >= ?", options.MinPrice)
	}

	if options.MaxPrice > 0 {
		query = query.Where("best_prices.min_price <= ?", options.MaxPrice)
	}

	// Aplicar ordenamiento
	query = query.Order("best_prices.min_price " + options.SortOrder)

	// Aplicar paginación
	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	// Cargar la categoría y asegurar que seleccionamos todos los campos del producto
	query = query.Preload("Category").
		Select("products.*, categories.name as category_name, categories.slug as category_slug")

	// Generar descripción de la consulta para depuración
	queryDescription := fmt.Sprintf("SELECT productos con categoría='%s'", options.CategorySlug)
	if options.MinPrice > 0 {
		queryDescription += fmt.Sprintf(", precio >= %.2f", options.MinPrice)
	}
	if options.MaxPrice > 0 {
		queryDescription += fmt.Sprintf(", precio <= %.2f", options.MaxPrice)
	}
	if options.StoreFilter != "" {
		queryDescription += fmt.Sprintf(", tienda='%s'", options.StoreFilter)
	}
	queryDescription += fmt.Sprintf(", ordenado por precio %s", options.SortOrder)
	queryDescription += fmt.Sprintf(", límite=%d, offset=%d", options.Limit, options.Offset)

	log.Printf("[SQL_DEBUG] Consulta SQL: %s", queryDescription)

	if err := query.Find(&products).Error; err != nil {
		log.Printf("[SQL_ERROR] Error en consulta: %v", err)
		return nil, err
	}

	log.Printf("[SQL_DEBUG] Consulta exitosa. Encontrados %d productos en %.2f ms",
		len(products), float64(time.Since(startTime).Microseconds())/1000.0)

	return products, nil
}

// CountFilteredProductsByCategory cuenta el número total de productos en una categoría con filtros
func (r *productRepository) CountFilteredProductsByCategory(ctx context.Context, options model.ProductFilterOptions) (int, error) {
	startTime := time.Now()
	log.Printf("[SQL_DEBUG] Contando productos filtrados para categoría: %s", options.CategorySlug)
	log.Printf("[SQL_DEBUG] Filtros: minPrice=%.2f, maxPrice=%.2f, store=%s",
		options.MinPrice, options.MaxPrice, options.StoreFilter)

	var count int64

	// Construir una subconsulta para obtener el mejor precio por producto
	// Modificamos para no seleccionar 'store' directamente en la subconsulta, ya que causa problemas con GROUP BY
	subQuery := r.db.WithContext(ctx).
		Table("prices").
		Select("product_id, MIN(price) as min_price").
		Group("product_id")

	if options.StoreFilter != "" {
		subQuery = subQuery.Where("store = ?", options.StoreFilter)
	}

	// Consulta principal uniendo productos con la subconsulta de mejores precios
	query := r.db.WithContext(ctx).
		Model(&model.Product{}).
		Joins("JOIN (?) AS best_prices ON products.id = best_prices.product_id", subQuery).
		Joins("JOIN categories ON products.category_id = categories.id").
		Where("categories.slug = ?", options.CategorySlug)

	// Aplicar filtros de precio
	if options.MinPrice > 0 {
		query = query.Where("best_prices.min_price >= ?", options.MinPrice)
	}

	if options.MaxPrice > 0 {
		query = query.Where("best_prices.min_price <= ?", options.MaxPrice)
	}

	// Generar descripción de la consulta para depuración
	queryDescription := fmt.Sprintf("COUNT productos con categoría='%s'", options.CategorySlug)
	if options.MinPrice > 0 {
		queryDescription += fmt.Sprintf(", precio >= %.2f", options.MinPrice)
	}
	if options.MaxPrice > 0 {
		queryDescription += fmt.Sprintf(", precio <= %.2f", options.MaxPrice)
	}
	if options.StoreFilter != "" {
		queryDescription += fmt.Sprintf(", tienda='%s'", options.StoreFilter)
	}

	log.Printf("[SQL_DEBUG] Consulta SQL de conteo: %s", queryDescription)

	if err := query.Count(&count).Error; err != nil {
		log.Printf("[SQL_ERROR] Error en consulta de conteo: %v", err)
		return 0, err
	}

	log.Printf("[SQL_DEBUG] Conteo exitoso. Total: %d productos en %.2f ms",
		count, float64(time.Since(startTime).Microseconds())/1000.0)

	return int(count), nil
}
