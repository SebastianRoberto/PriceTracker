package handler

import (
	"log"
	"net/http"
	"sort"
	"strconv"

	"app/internal/domain/model"
	"app/internal/interface/web/views"
	"app/internal/usecase"

	"github.com/gin-gonic/gin"
)

// CategoryHandler maneja las peticiones para las páginas de categoría
type CategoryHandler struct {
	productUseCase   *usecase.ProductUseCase
	templateRenderer *views.TemplateRenderer
}

// NewCategoryHandler crea una nueva instancia del CategoryHandler
func NewCategoryHandler(productUseCase *usecase.ProductUseCase, templateRenderer *views.TemplateRenderer) *CategoryHandler {
	return &CategoryHandler{
		productUseCase:   productUseCase,
		templateRenderer: templateRenderer,
	}
}

// GetCategory maneja la petición GET para una página de categoría
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	// Obtener el usuario de la sesión (puede ser nil si no está autenticado)
	user, _ := c.Get("user")

	// Obtener el slug de la categoría de la URL
	slug := c.Param("slug")
	if slug == "" {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Message": "Categoría no especificada",
		})
		return
	}

	// Obtener productos por categoría
	ctx := c.Request.Context()

	// Parsear parámetros de paginación
	limit := 48 // Aumentado a 48 productos por página (anteriormente 24)
	page := 0   // Por defecto primera página

	// Si se proporciona el parámetro de página, convertirlo
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p - 1 // Convertir de página (1-indexed) a offset (0-indexed)
		}
	}

	// Calcular offset
	offset := page * limit

	// Leer el filtro de tienda de la query
	storeFilter := c.Query("store")

	// Obtener productos
	products, err := h.productUseCase.GetProductsByCategory(ctx, slug, limit, offset, storeFilter)
	if err != nil {
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al obtener productos",
			"Error":   err.Error(),
		})
		return
	}

	// Ordenar productos por precio ascendente (mejor oferta primero)
	sort.SliceStable(products, func(i, j int) bool {
		var priceI, priceJ float64
		if len(products[i].Prices) > 0 {
			priceI = products[i].Prices[0].Price
		}
		if len(products[j].Prices) > 0 {
			priceJ = products[j].Prices[0].Price
		}
		// Si alguno no tiene precio, lo mandamos al final
		if priceI == 0 {
			return false
		}
		if priceJ == 0 {
			return true
		}
		return priceI < priceJ
	})

	// Convertir productos a ViewModels
	productVMs := make([]views.ProductViewModel, 0, len(products))
	for _, p := range products {
		var bestPrice *model.Price
		if len(p.Prices) > 0 {
			bestPrice = &p.Prices[0]
		}
		productVMs = append(productVMs, views.ToProductViewModel(p, bestPrice))
	}

	// Obtener las categorías para el menú de navegación
	categories, err := h.productUseCase.GetAllCategories(ctx)
	if err != nil {
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al obtener categorías",
			"Error":   err.Error(),
		})
		return
	}

	// Convertir categorías a ViewModels
	categoryVMs := make([]views.CategoryViewModel, 0, len(categories))
	for _, cat := range categories {
		categoryVMs = append(categoryVMs, views.ToCategoryViewModel(*cat, 0))
	}

	// Buscar la categoría actual en la lista de categorías
	var currentCategoryVM *views.CategoryViewModel
	for _, cat := range categoryVMs {
		if cat.Slug == slug {
			currentCategoryVM = &cat
			break
		}
	}

	if currentCategoryVM == nil {
		h.templateRenderer.Render(c, http.StatusNotFound, "error.html", gin.H{
			"Message": "Categoría no encontrada",
		})
		return
	}

	// Calcular el número total de productos y páginas
	totalProducts, err := h.productUseCase.GetTotalProductsInCategory(ctx, slug, storeFilter)
	if err != nil {
		// Si hay error, no detenemos la página, simplemente no mostramos paginación
		totalProducts = 0
	}
	totalPages := 0
	if totalProducts > 0 {
		totalPages = (totalProducts + limit - 1) / limit
	}

	// Renderizar la página de categoría
	h.templateRenderer.Render(c, http.StatusOK, "category.html", gin.H{
		"Title":         currentCategoryVM.Name + " - Comparador de Precios",
		"Category":      currentCategoryVM,
		"Products":      productVMs,
		"Categories":    categoryVMs,
		"User":          user,
		"CurrentPage":   page + 1, // Convertir de nuevo a 1-indexed para la vista
		"TotalPages":    totalPages,
		"TotalProducts": totalProducts,
	})
}

// GetCategoryAPI maneja la petición GET a la API para obtener productos de una categoría
func (h *CategoryHandler) GetCategoryAPI(c *gin.Context) {
	// Obtener el slug de la categoría de la URL
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Categoría no especificada",
		})
		return
	}

	log.Printf("[API_DEBUG] Recibida petición a /api/categoria/%s con query params: %v", slug, c.Request.URL.Query())

	// Obtener productos por categoría
	ctx := c.Request.Context()

	// Parsear parámetros de paginación
	limit := 48 // Aumentado a 48 productos por página (anteriormente 24)
	page := 0   // Por defecto primera página

	// Si se proporciona el parámetro de página, convertirlo
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p - 1 // Convertir de página (1-indexed) a offset (0-indexed)
		}
	}

	// Calcular offset
	offset := page * limit

	// Leer el filtro de tienda de la query
	storeFilter := c.Query("store")

	// Leer el orden de la query (asc o desc)
	sortOrder := c.Query("sort")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc" // Valor por defecto
	}

	// Leer filtros de precio
	var minPrice, maxPrice float64
	var err error

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		minPrice, err = strconv.ParseFloat(minPriceStr, 64)
		if err != nil {
			minPrice = 0
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		maxPrice, err = strconv.ParseFloat(maxPriceStr, 64)
		if err != nil {
			maxPrice = 0
		}
	}

	log.Printf("[API_DEBUG] Parámetros procesados: página=%d, offset=%d, limite=%d, filtroTienda=%s, ordenación=%s, precioMin=%.2f, precioMax=%.2f",
		page+1, offset, limit, storeFilter, sortOrder, minPrice, maxPrice)

	// IMPORTANTE: Modificamos el enfoque para obtener los productos ya filtrados y ordenados desde la base de datos
	// Creamos un objeto de opciones de filtrado para pasar al usecase
	options := model.ProductFilterOptions{
		CategorySlug: slug,
		Limit:        limit,
		Offset:       offset,
		StoreFilter:  storeFilter,
		SortOrder:    sortOrder,
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
	}

	// Obtener productos filtrados y ordenados directamente desde la base de datos
	products, err := h.productUseCase.GetFilteredProductsByCategory(ctx, options)
	if err != nil {
		log.Printf("[API_ERROR] Error al obtener productos filtrados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al obtener productos: " + err.Error(),
		})
		return
	}

	log.Printf("[API_DEBUG] Obtenidos %d productos filtrados para la página %d", len(products), page+1)

	// Asegurarnos de que cada producto tenga su mejor precio cargado
	// Si los productos no tienen precios cargados, necesitamos cargarlos
	for _, product := range products {
		if len(product.Prices) == 0 {
			// Si no tiene precios cargados, intentamos cargarlos
			productWithPrices, err := h.productUseCase.GetProductWithPrices(ctx, product.ID)
			if err == nil && productWithPrices != nil && len(productWithPrices.Prices) > 0 {
				product.Prices = productWithPrices.Prices
			}
		}
	}

	// Preparar la respuesta
	type simpleProduct struct {
		ID         uint    `json:"id"`
		Name       string  `json:"name"`
		ImageURL   string  `json:"image_url"`
		BestPrice  float64 `json:"best_price,omitempty"`
		BestStore  string  `json:"best_store,omitempty"`
		CategoryID uint    `json:"category_id"`
		Category   struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"category"`
	}

	// Convertir productos al formato JSON simplificado
	jsonProducts := make([]simpleProduct, 0, len(products))
	for _, p := range products {
		prod := simpleProduct{
			ID:         p.ID,
			Name:       p.Name,
			ImageURL:   p.ImageURL,
			CategoryID: p.CategoryID,
			Category: struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
			}{
				ID:   p.Category.ID,
				Name: p.Category.Name,
				Slug: p.Category.Slug,
			},
		}

		// Añadir el mejor precio si existe
		if len(p.Prices) > 0 {
			// Encontrar el mejor precio (el más bajo)
			bestPrice := p.Prices[0] // Empezamos con el primer precio
			for i := 1; i < len(p.Prices); i++ {
				if p.Prices[i].Price < bestPrice.Price {
					bestPrice = p.Prices[i]
				}
			}

			prod.BestPrice = bestPrice.Price
			prod.BestStore = bestPrice.Store
		}

		jsonProducts = append(jsonProducts, prod)
	}

	// Calcular el número total de productos y páginas (considerando los filtros)
	totalProducts, err := h.productUseCase.GetTotalFilteredProductsInCategory(ctx, options)
	if err != nil {
		log.Printf("[API_ERROR] Error al contar productos filtrados: %v", err)
		totalProducts = 0
	}
	totalPages := 0
	if totalProducts > 0 {
		totalPages = (totalProducts + limit - 1) / limit
	}

	log.Printf("[API_DEBUG] Total de productos: %d, total de páginas: %d", totalProducts, totalPages)

	// Devolver respuesta
	c.JSON(http.StatusOK, gin.H{
		"products":       jsonProducts,
		"current_page":   page + 1,
		"total_pages":    totalPages,
		"total_products": totalProducts,
	})
}
