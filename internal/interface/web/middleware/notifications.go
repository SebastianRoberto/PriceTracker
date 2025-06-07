package middleware

import (
	"app/internal/domain/model"
	"app/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

// IncludeUnreadNotificationsCount agrega el contador de notificaciones no leídas al contexto.
func IncludeUnreadNotificationsCount(uc *usecase.PriceAlertUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Por defecto, inicializar UnreadNotifications a 0.
		// Esto asegura que la variable siempre esté disponible en las plantillas.
		c.Set("UnreadNotifications", 0)

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

		count, err := uc.GetUnreadNotificationsCount(c.Request.Context(), user.ID)
		if err != nil {
			log.Printf("Error en middleware IncludeUnreadNotificationsCount para userID %d: %v", user.ID, err)
			// Ya hemos seteado 0, así que continuamos.
			c.Next()
			return
		}

		c.Set("UnreadNotifications", count)
		c.Next()
	}
}
