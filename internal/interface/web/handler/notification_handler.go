package handler

import (
	"context"
	"net/http"
	"strconv"

	"app/internal/interface/web/views"
	"app/internal/usecase"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// NotificationHandler maneja las peticiones relacionadas con notificaciones
type NotificationHandler struct {
	notificationUseCase *usecase.PriceAlertUseCase
	templateRenderer    *views.TemplateRenderer
}

// NewNotificationHandler crea una nueva instancia del NotificationHandler
func NewNotificationHandler(notificationUseCase *usecase.PriceAlertUseCase, templateRenderer *views.TemplateRenderer) *NotificationHandler {
	return &NotificationHandler{
		notificationUseCase: notificationUseCase,
		templateRenderer:    templateRenderer,
	}
}

// ShowNotifications muestra la página de notificaciones del usuario
func (h *NotificationHandler) ShowNotifications(c *gin.Context) {
	// Obtener el usuario de la sesión
	userID := sessions.Default(c).Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener notificaciones del usuario
	ctx := c.Request.Context()
	notifications, err := h.notificationUseCase.GetUserNotifications(ctx, userID.(uint))
	if err != nil {
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al obtener notificaciones",
			"Error":   err.Error(),
		})
		return
	}

	// Obtener categorías del contexto (añadidas por el middleware)
	categories, _ := c.Get("allCategories")

	// Obtener el usuario del contexto
	user, _ := c.Get("user")

	// Obtener alertas de precio para el contador de Mi Cesta
	priceAlerts, _ := c.Get("priceAlerts")

	// Renderizar página de notificaciones
	h.templateRenderer.Render(c, http.StatusOK, "notifications.html", gin.H{
		"Title":               "Notificaciones - Comparador de Precios",
		"Notifications":       notifications,
		"User":                user,
		"Categories":          categories,
		"PriceAlerts":         priceAlerts,
		"UnreadNotifications": h.getUnreadNotificationsCount(userID.(uint)),
	})
}

// MarkAsRead marca una notificación como leída
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	// Obtener el usuario de la sesión
	userID := sessions.Default(c).Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Acceso no autorizado"})
		return
	}

	// Obtener ID de notificación
	notificationIDStr := c.PostForm("notification_id")
	if notificationIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de notificación no proporcionado"})
		return
	}

	notificationID, err := strconv.ParseUint(notificationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de notificación inválido"})
		return
	}

	// Marcar notificación como leída
	ctx := c.Request.Context()
	if err := h.notificationUseCase.MarkNotificationAsRead(ctx, uint(notificationID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al marcar notificación: " + err.Error()})
		return
	}

	// Si se solicitó como AJAX, devolver éxito
	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	// Si no, redirigir a la página de notificaciones
	c.Redirect(http.StatusFound, "/notificaciones")
}

// MarkAllAsRead marca todas las notificaciones del usuario como leídas
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	// Obtener el usuario de la sesión
	userID := sessions.Default(c).Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Acceso no autorizado"})
		return
	}

	// Marcar todas las notificaciones como leídas
	ctx := c.Request.Context()
	if err := h.notificationUseCase.MarkAllNotificationsAsRead(ctx, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al marcar notificaciones: " + err.Error()})
		return
	}

	// Si se solicitó como AJAX, devolver éxito
	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	// Si no, redirigir a la página de notificaciones
	c.Redirect(http.StatusFound, "/notificaciones")
}

// DeleteReadNotifications elimina todas las notificaciones leídas del usuario
// Esta función es llamada vía AJAX para no recargar la página
func (h *NotificationHandler) DeleteReadNotifications(c *gin.Context) {
	// Obtener el usuario de la sesión
	userID := sessions.Default(c).Get("user_id")
	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Acceso no autorizado"})
		return
	}

	// Obtener notificaciones del usuario
	ctx := c.Request.Context()
	notifications, err := h.notificationUseCase.GetUserNotifications(ctx, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener notificaciones"})
		return
	}

	// Eliminar las notificaciones leídas una a una
	deletedCount := 0
	for _, notification := range notifications {
		if notification.IsRead {
			if err := h.notificationUseCase.DeleteNotification(ctx, notification.ID, userID.(uint)); err == nil {
				deletedCount++
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"deleted": deletedCount,
	})
}

// getUnreadNotificationsCount obtiene el número de notificaciones no leídas para un usuario
func (h *NotificationHandler) getUnreadNotificationsCount(userID uint) int {
	notifications, err := h.notificationUseCase.GetUserNotifications(context.Background(), userID)
	if err != nil {
		return 0
	}

	unreadCount := 0
	for _, notification := range notifications {
		if !notification.IsRead {
			unreadCount++
		}
	}
	return unreadCount
}
