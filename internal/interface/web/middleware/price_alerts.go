package middleware

import (
	"app/internal/domain/model"
	"app/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

// IncludePriceAlerts agrega las alertas de precio del usuario a todas las vistas protegidas
func IncludePriceAlerts(priceAlertUseCase *usecase.PriceAlertUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Por defecto, inicializar PriceAlerts como un slice vacío de punteros.
		// Esto asegura que la variable siempre esté disponible en las plantillas,
		// incluso si el usuario no está logueado o hay un error.
		c.Set("PriceAlerts", []*model.PriceAlert{})

		userInterface, exists := c.Get("user")
		if !exists {
			c.Next()
			return
		}

		user, ok := userInterface.(*model.User)
		if !ok {
			c.Next()
			return
		}

		alerts, err := priceAlertUseCase.GetUserAlerts(c.Request.Context(), user.ID)
		if err != nil {
			log.Printf("Error en middleware IncludePriceAlerts al obtener alertas para userID %d: %v", user.ID, err)
			// Ya hemos seteado un slice vacío, así que continuamos.
			c.Next()
			return
		}

		// Si alerts no es nil (incluso si es un slice vacío devuelto por el usecase), lo seteamos.
		// Si es nil (por alguna razón inesperada aunque sin error), el valor por defecto (slice vacío) se mantendrá.
		if alerts != nil {
			c.Set("PriceAlerts", alerts)
		}

		c.Next()
	}
}
