package handler

import (
	"net/http"
	"strconv"
	"time"

	"app/internal/domain/model"
	"app/internal/interface/web/views"
	"app/internal/usecase"

	"github.com/gin-gonic/gin"
)

// ProductHandler maneja las peticiones para las páginas de producto
type ProductHandler struct {
	productUseCase   *usecase.ProductUseCase
	templateRenderer *views.TemplateRenderer
}

// NewProductHandler crea una nueva instancia del ProductHandler
func NewProductHandler(productUseCase *usecase.ProductUseCase, templateRenderer *views.TemplateRenderer) *ProductHandler {
	return &ProductHandler{
		productUseCase:   productUseCase,
		templateRenderer: templateRenderer,
	}
}

// GetProduct maneja la petición GET para una página de detalle de producto
func (h *ProductHandler) GetProduct(c *gin.Context) {
	// Obtener el usuario de la sesión (puede ser nil si no está autenticado)
	user, _ := c.Get("user")

	// Obtener el ID del producto de la URL
	idStr := c.Param("id")
	if idStr == "" {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Message": "ID de producto no especificado",
		})
		return
	}

	// Convertir el ID a uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Message": "ID de producto inválido",
			"Error":   err.Error(),
		})
		return
	}

	// Obtener detalle del producto
	ctx := c.Request.Context()
	product, err := h.productUseCase.GetProductDetail(ctx, uint(id))
	if err != nil {
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al obtener el producto",
			"Error":   err.Error(),
		})
		return
	}

	// Verificar si el producto existe
	if product == nil {
		h.templateRenderer.Render(c, http.StatusNotFound, "error.html", gin.H{
			"Message": "Producto no encontrado",
		})
		return
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

	// Preparar el mejor precio
	var bestPrice *model.Price
	if len(product.Prices) > 0 {
		bestPrice = &product.Prices[0]
	}

	// Obtener productos similares
	similarProducts, err := h.productUseCase.GetSimilarProducts(ctx, uint(id), 3)
	if err != nil {
		// Si hay error, continuamos sin productos similares
		similarProducts = []*model.Product{}
	}

	// Preparar viewmodels para las categorías
	categoryVMs := make([]views.CategoryViewModel, 0, len(categories))
	for _, cat := range categories {
		categoryVMs = append(categoryVMs, views.ToCategoryViewModel(*cat, 0))
	}

	// Preparar viewmodel para el producto
	productVM := views.ToProductViewModel(product, bestPrice)

	// Preparar viewmodel para el mejor precio
	var bestPriceVM views.PriceViewModel
	if bestPrice != nil {
		bestPriceVM = views.ToPriceViewModel(*bestPrice)
	}

	// Preparar viewmodels para productos similares
	similarProductsVM := make([]views.SimilarProductViewModel, 0, len(similarProducts))
	for _, sp := range similarProducts {
		// Obtener el mejor precio para cada producto similar
		var similarBestPrice *model.Price
		if len(sp.Prices) > 0 {
			similarBestPrice = &sp.Prices[0]
		}

		spVM := views.SimilarProductViewModel{
			ID:       sp.ID,
			Name:     sp.Name,
			ImageURL: sp.ImageURL,
		}

		if similarBestPrice != nil {
			spVM.Price = similarBestPrice.Price
			spVM.Store = similarBestPrice.Store
			spVM.URL = similarBestPrice.URL
		}

		similarProductsVM = append(similarProductsVM, spVM)
	}

	h.templateRenderer.Render(c, http.StatusOK, "product_detail.html", gin.H{
		"Title":               product.Name + " - Comparador de Precios",
		"Product":             productVM,
		"BestPrice":           bestPriceVM,
		"SimilarProducts":     similarProductsVM,
		"Categories":          categoryVMs,
		"LastUpdated":         time.Now().Format("02/01/2006 15:04"),
		"User":                user,
		"UnreadNotifications": 0,
		"IsFollowing":         false,
		"PriceAlert":          nil,
		"RelatedProducts":     []views.ProductViewModel{},
	})
}
