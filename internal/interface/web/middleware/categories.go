package middleware

import (
	"app/internal/usecase"

	"github.com/gin-gonic/gin"
)

// IncludeCategories agrega las categorías disponibles a todas las vistas
func IncludeCategories(productUseCase *usecase.ProductUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener las categorías
		categories, err := productUseCase.GetAllCategories(c.Request.Context())
		if err != nil {
			// Si hay error, continuamos sin categorías
			c.Next()
			return
		}

		// Agregar las categorías al contexto
		c.Set("allCategories", categories)

		// Continuar con la solicitud
		c.Next()
	}
}
