package handler

import (
	"log"
	"net/http"

	"app/internal/domain/model"
	"app/internal/interface/web/views"
	"app/internal/usecase"

	"github.com/gin-gonic/gin"
)

// HomeHandler maneja las peticiones para la página principal
type HomeHandler struct {
	productUseCase   *usecase.ProductUseCase
	templateRenderer *views.TemplateRenderer
}

// NewHomeHandler crea una nueva instancia del HomeHandler
func NewHomeHandler(productUseCase *usecase.ProductUseCase, templateRenderer *views.TemplateRenderer) *HomeHandler {
	return &HomeHandler{
		productUseCase:   productUseCase,
		templateRenderer: templateRenderer,
	}
}

// GetHome maneja la petición GET para la página principal
func (h *HomeHandler) GetHome(c *gin.Context) {
	// Obtener el usuario de la sesión (puede ser nil si no está autenticado)
	user, _ := c.Get("user")

	// Obtener los mejores productos
	ctx := c.Request.Context()
	products, err := h.productUseCase.GetBestDeals(ctx, 12) // Mostrar 12 productos en la página principal
	if err != nil {
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al obtener productos",
			"Error":   err.Error(),
		})
		return
	}

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
		// Contar productos en esta categoría
		productCount, err := h.productUseCase.CountProductsByCategory(ctx, cat.ID)
		if err != nil {
			log.Printf("Error al contar productos para categoría %s: %v", cat.Name, err)
			productCount = 0
		}
		categoryVMs = append(categoryVMs, views.ToCategoryViewModel(*cat, productCount))
	}

	// Renderizar la página principal
	h.templateRenderer.Render(c, http.StatusOK, "home.html", gin.H{
		"Title":      "Comparador de Precios - Mejores Ofertas",
		"Products":   productVMs,
		"Categories": categoryVMs,
		"User":       user,
	})
}
