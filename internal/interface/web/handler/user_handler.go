package handler

import (
	"net/http"

	"app/internal/interface/web/views"
	"app/internal/usecase"

	"github.com/gin-gonic/gin"
)

// UserHandler maneja las peticiones relacionadas con los usuarios
type UserHandler struct {
	userUseCase       *usecase.UserUseCase
	priceAlertUseCase *usecase.PriceAlertUseCase
	renderer          *views.TemplateRenderer
}

// NewUserHandler crea una nueva instancia del manejador de usuarios
func NewUserHandler(userUseCase *usecase.UserUseCase, priceAlertUseCase *usecase.PriceAlertUseCase, renderer *views.TemplateRenderer) *UserHandler {
	return &UserHandler{
		userUseCase:       userUseCase,
		priceAlertUseCase: priceAlertUseCase,
		renderer:          renderer,
	}
}

// GetProfile muestra el perfil del usuario
func (h *UserHandler) GetProfile(c *gin.Context) {
	// Obtener el ID del usuario desde la sesión
	userID, exists := c.Get("userID")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	ctx := c.Request.Context()

	// Obtener los datos del usuario
	user, err := h.userUseCase.GetUserByID(ctx, userID.(uint))
	if err != nil {
		h.renderer.RenderError(c, http.StatusInternalServerError, "Error al obtener datos del usuario")
		return
	}

	// Obtener categorías para el menú desplegable
	allCategories, _ := c.Get("allCategories")

	// Renderizar la plantilla de perfil (simplificada)
	h.renderer.Render(c, http.StatusOK, "profile.html", gin.H{
		"Title":      "Mi Perfil",
		"User":       user,
		"Categories": allCategories,
		"Success":    c.Query("success"),
	})
}
