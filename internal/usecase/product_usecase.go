package usecase

import (
	"context"
	"fmt"
	"sort"

	"app/internal/domain/model"
	"app/internal/domain/repositories"
)

// ProductUseCase implementa la lógica de negocio para los productos
type ProductUseCase struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
	priceRepo    repositories.PriceRepository
}

// NewProductUseCase crea una nueva instancia del caso de uso para productos
func NewProductUseCase(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
	priceRepo repositories.PriceRepository,
) *ProductUseCase {
	return &ProductUseCase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		priceRepo:    priceRepo,
	}
}

// GetBestDeals obtiene los productos con las mejores ofertas
func (uc *ProductUseCase) GetBestDeals(ctx context.Context, limit int) ([]*model.Product, error) {
	products, err := uc.productRepo.FindBestDeals(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("error al obtener las mejores ofertas: %w", err)
	}

	// Enriquecemos cada producto con su mejor precio
	for _, product := range products {
		bestPrice, err := uc.priceRepo.FindBestPriceByProductID(ctx, product.ID)
		if err != nil {
			return nil, fmt.Errorf("error al obtener el mejor precio del producto %d: %w", product.ID, err)
		}

		if bestPrice != nil {
			product.Prices = []model.Price{*bestPrice}
		}
	}

	return products, nil
}

// GetProductsByCategory obtiene productos por categoría
func (uc *ProductUseCase) GetProductsByCategory(ctx context.Context, categorySlug string, limit, offset int, storeFilter string) ([]*model.Product, error) {
	category, err := uc.categoryRepo.FindBySlug(ctx, categorySlug)
	if err != nil {
		return nil, fmt.Errorf("error al buscar categoría %s: %w", categorySlug, err)
	}

	products, err := uc.productRepo.FindByCategory(ctx, category.ID, limit, offset, storeFilter)
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos de categoría %s: %w", categorySlug, err)
	}

	// Enriquecemos cada producto con su mejor precio
	for _, product := range products {
		bestPrice, err := uc.priceRepo.FindBestPriceByProductID(ctx, product.ID)
		if err != nil {
			return nil, fmt.Errorf("error al obtener el mejor precio del producto %d: %w", product.ID, err)
		}

		if bestPrice != nil {
			product.Prices = []model.Price{*bestPrice}
		}
	}

	return products, nil
}

// GetProductDetail obtiene el detalle de un producto con sus mejores ofertas
func (uc *ProductUseCase) GetProductDetail(ctx context.Context, productID uint) (*model.Product, error) {
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar producto %d: %w", productID, err)
	}

	// Obtenemos las mejores ofertas para este producto
	bestPrices, err := uc.priceRepo.FindTopOffersByProductID(ctx, productID, 3)
	if err != nil {
		return nil, fmt.Errorf("error al obtener mejores ofertas para el producto %d: %w", productID, err)
	}

	// Convertimos el slice de punteros a un slice de valores
	prices := make([]model.Price, len(bestPrices))
	for i, p := range bestPrices {
		prices[i] = *p
	}

	product.Prices = prices

	return product, nil
}

// GetSimilarProducts obtiene productos similares a un producto dado
func (uc *ProductUseCase) GetSimilarProducts(ctx context.Context, productID uint, limit int) ([]*model.Product, error) {
	// Obtener productos similares
	similarProducts, err := uc.productRepo.FindSimilarProducts(ctx, productID, limit)
	if err != nil {
		return nil, fmt.Errorf("error al buscar productos similares: %v", err)
	}

	// Para cada producto similar, obtenemos su mejor precio
	for _, similarProduct := range similarProducts {
		bestPrice, err := uc.priceRepo.FindBestPriceByProductID(ctx, similarProduct.ID)
		if err == nil && bestPrice != nil {
			similarProduct.Prices = []model.Price{*bestPrice}
		}
	}

	return similarProducts, nil
}

// GetAllCategories obtiene todas las categorías
func (uc *ProductUseCase) GetAllCategories(ctx context.Context) ([]*model.Category, error) {
	return uc.categoryRepo.GetAll(ctx)
}

// CountProductsByCategory cuenta el número de productos en una categoría
// Esta función es más general y no aplica filtro de tienda por defecto.
// Si se necesita contar con filtro de tienda, usar GetTotalProductsInCategory.
func (uc *ProductUseCase) CountProductsByCategory(ctx context.Context, categoryID uint) (int, error) {
	// Pasamos un string vacío para storeFilter, ya que esta función es genérica.
	products, err := uc.productRepo.FindByCategory(ctx, categoryID, 0, 0, "")
	if err != nil {
		return 0, fmt.Errorf("error al contar productos en categoría %d: %w", categoryID, err)
	}
	return len(products), nil
}

// GetFeaturedProducts obtiene productos destacados
func (uc *ProductUseCase) GetFeaturedProducts(ctx context.Context, limit int) ([]model.Product, error) {
	// Por ahora, simplemente reutilizamos el método GetBestDeals
	// En el futuro, podríamos tener un criterio diferente para productos destacados
	productsPtr, err := uc.GetBestDeals(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos destacados: %w", err)
	}

	// Convertir []*model.Product a []model.Product
	products := make([]model.Product, len(productsPtr))
	for i, p := range productsPtr {
		products[i] = *p
	}

	return products, nil
}

// GetTotalProductsInCategory obtiene el número total de productos en una categoría
func (uc *ProductUseCase) GetTotalProductsInCategory(ctx context.Context, categorySlug string, storeFilter string) (int, error) {
	category, err := uc.categoryRepo.FindBySlug(ctx, categorySlug)
	if err != nil {
		return 0, fmt.Errorf("error al buscar categoría %s: %w", categorySlug, err)
	}

	// Primero intentamos obtener el total usando una función específica del repositorio si existe
	count, err := uc.productRepo.CountByCategory(ctx, category.ID, storeFilter)
	if err == nil {
		return count, nil
	}

	// Si no existe o falla, obtenemos todos los productos y contamos (menos eficiente)
	// Nota: Esta ruta alternativa también debería idealmente aplicar el storeFilter si es posible,
	// pero la función FindByCategory que se usa aquí también fue actualizada para recibir storeFilter.
	products, err := uc.productRepo.FindByCategory(ctx, category.ID, 0, 0, storeFilter)
	if err != nil {
		return 0, fmt.Errorf("error al contar productos en categoría %s: %w", categorySlug, err)
	}
	return len(products), nil
}

// GetFilteredProductsByCategory obtiene productos filtrados y ordenados por categoría
func (uc *ProductUseCase) GetFilteredProductsByCategory(ctx context.Context, options model.ProductFilterOptions) ([]*model.Product, error) {
	// Validar el orden de clasificación
	if options.SortOrder != "asc" && options.SortOrder != "desc" {
		options.SortOrder = "asc" // Valor predeterminado
	}

	// Utilizar el nuevo método del repositorio que aplica filtros directamente en SQL
	products, err := uc.productRepo.FindFilteredProductsByCategory(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos filtrados de categoría %s: %w", options.CategorySlug, err)
	}

	// Asegurarnos de que cada producto tenga su precio
	for _, product := range products {
		// Si el producto ya tiene precios, no los buscamos de nuevo
		if len(product.Prices) > 0 {
			continue
		}

		// Obtener el mejor precio para este producto
		bestPrice, err := uc.priceRepo.FindBestPriceByProductID(ctx, product.ID)
		if err == nil && bestPrice != nil {
			product.Prices = []model.Price{*bestPrice}
		}
	}

	return products, nil
}

// GetTotalFilteredProductsInCategory obtiene el número total de productos filtrados en una categoría
func (uc *ProductUseCase) GetTotalFilteredProductsInCategory(ctx context.Context, options model.ProductFilterOptions) (int, error) {
	// Utilizar el nuevo método del repositorio que cuenta productos filtrados directamente en SQL
	count, err := uc.productRepo.CountFilteredProductsByCategory(ctx, options)
	if err != nil {
		return 0, fmt.Errorf("error al contar productos filtrados de categoría %s: %w", options.CategorySlug, err)
	}

	return count, nil
}

// GetProductWithPrices obtiene un producto con sus precios asociados
func (uc *ProductUseCase) GetProductWithPrices(ctx context.Context, productID uint) (*model.Product, error) {
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener producto: %v", err)
	}

	// Cargar los precios del producto
	pricesPtr, err := uc.priceRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("error al obtener precios: %v", err)
	}

	// Convertir []*model.Price a []model.Price
	prices := make([]model.Price, len(pricesPtr))
	for i, p := range pricesPtr {
		prices[i] = *p
	}

	// Ordenar precios por valor ascendente (el más bajo primero)
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Price < prices[j].Price
	})

	product.Prices = prices
	return product, nil
}
