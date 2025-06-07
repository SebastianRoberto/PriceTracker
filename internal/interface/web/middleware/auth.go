package middleware

import (
	"net/http"

	"app/internal/domain/model"
	"app/internal/usecase"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthRequired es un middleware que verifica si el usuario está autenticado
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener la sesión
		session := sessions.Default(c)
		userID := session.Get("user_id")

		// Verificar si existe un usuario_id en la sesión
		if userID == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Continuar con la petición
		c.Next()
	}
}

// LoadUser es un middleware que carga el usuario actual (si existe) y lo añade al contexto
func LoadUser(userUseCase *usecase.UserUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener la sesión
		session := sessions.Default(c)
		userID := session.Get("user_id")

		// Si hay un ID de usuario en la sesión, intentar cargar el usuario
		if userID != nil {
			// Convertir el ID a uint
			id, ok := userID.(uint)
			if !ok {
				// Si no se puede convertir, limpiar la sesión
				session.Clear()
				session.Save()
				c.Next()
				return
			}

			// Intentar obtener el usuario de la base de datos
			user, err := userUseCase.GetUserByID(c.Request.Context(), id)
			if err == nil && user != nil {
				// Si se encuentra el usuario, añadirlo al contexto
				c.Set("user", user)
			} else {
				// Si no se encuentra el usuario, limpiar la sesión
				session.Clear()
				session.Save()
			}
		}

		// Continuar con la petición
		c.Next()
	}
}

// AdminRequired verifica si el usuario es administrador
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el usuario del contexto
		userInterface, exists := c.Get("user")
		if !exists {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Convertir la interfaz a *model.User
		user, ok := userInterface.(*model.User)
		if !ok || !user.IsAdmin {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"Message": "No tienes permisos para acceder a esta página",
			})
			c.Abort()
			return
		}

		// Si el usuario es administrador, continuar
		c.Next()
	}
} 