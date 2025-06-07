package handler

import (
	"net/http"
	"strconv"

	"app/internal/domain/model"
	"app/internal/domain/repositories"
	"app/internal/interface/web/views"
	"app/internal/usecase"

	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// PriceAlertHandler maneja las peticiones relacionadas con las alertas de precio
type PriceAlertHandler struct {
	priceAlertUseCase *usecase.PriceAlertUseCase
	productUseCase    *usecase.ProductUseCase
	watchlistRepo     repositories.WatchlistRepository
	watchlistItemRepo repositories.WatchlistItemRepository
	templateRenderer  *views.TemplateRenderer
}

// NewPriceAlertHandler crea una nueva instancia del PriceAlertHandler
func NewPriceAlertHandler(
	priceAlertUseCase *usecase.PriceAlertUseCase,
	productUseCase *usecase.ProductUseCase,
	watchlistRepo repositories.WatchlistRepository,
	watchlistItemRepo repositories.WatchlistItemRepository,
	templateRenderer *views.TemplateRenderer,
) *PriceAlertHandler {
	return &PriceAlertHandler{
		priceAlertUseCase: priceAlertUseCase,
		productUseCase:    productUseCase,
		watchlistRepo:     watchlistRepo,
		watchlistItemRepo: watchlistItemRepo,
		templateRenderer:  templateRenderer,
	}
}

// SetPriceAlert maneja la creación o actualización de una alerta de precio
func (h *PriceAlertHandler) SetPriceAlert(c *gin.Context) {
	// Determinar si es una solicitud AJAX
	isAjax := c.GetHeader("X-Requested-With") == "XMLHttpRequest"

	// Obtener el usuario de la sesión
	userID := sessions.Default(c).Get("user_id")
	if userID == nil {
		if isAjax {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Debe iniciar sesión para añadir productos a la cesta",
			})
			return
		}
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener datos del formulario
	productIDStr := c.PostForm("product_id")
	targetPriceStr := c.PostForm("target_price")

	// Siempre enviar notificación por email
	notifyByEmail := true

	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		if isAjax {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "ID de producto inválido: " + err.Error(),
			})
			return
		}
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Message": "ID de producto inválido",
			"Error":   err.Error(),
		})
		return
	}

	targetPrice, err := strconv.ParseFloat(targetPriceStr, 64)
	if err != nil {
		if isAjax {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Precio objetivo inválido: " + err.Error(),
			})
			return
		}
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Message": "Precio objetivo inválido",
			"Error":   err.Error(),
		})
		return
	}

	// Verificar si ya existe una alerta para este producto y usuario
	ctx := c.Request.Context()
	alerts, err := h.priceAlertUseCase.GetUserAlerts(ctx, userID.(uint))
	if err != nil {
		if isAjax {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Error al verificar alertas existentes: " + err.Error(),
			})
			return
		}
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al verificar alertas existentes",
			"Error":   err.Error(),
		})
		return
	}

	var existingAlert *model.PriceAlert
	for _, alert := range alerts {
		if alert.ProductID == uint(productID) {
			existingAlert = alert
			break
		}
	}

	// Variable para almacenar si es una actualización o creación
	isUpdate := existingAlert != nil

	// Crear o actualizar la alerta
	var savedAlert *model.PriceAlert
	if existingAlert != nil {
		// Actualizar alerta existente
		savedAlert, err = h.priceAlertUseCase.UpdateAlert(
			ctx,
			existingAlert.ID,
			userID.(uint),
			targetPrice,
			notifyByEmail,
			true, // alerta activa
		)
	} else {
		// Crear nueva alerta
		savedAlert, err = h.priceAlertUseCase.CreateAlert(
			ctx,
			userID.(uint),
			uint(productID),
			targetPrice,
			notifyByEmail,
		)
	}

	if err != nil {
		if isAjax {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Error al guardar la alerta de precio: " + err.Error(),
			})
			return
		}
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al guardar la alerta de precio",
			"Error":   err.Error(),
		})
		return
	}

	// Asegurarnos de que el producto esté también en la tabla watchlist_items
	// (puede fallar silenciosamente sin afectar al flujo principal)
	if err == nil {
		// Primero asegurarnos de que el usuario tiene una watchlist
		_, errWatchlist := h.watchlistRepo.FindByUserID(ctx, userID.(uint))
		if errWatchlist != nil {
			// Crear la watchlist para este usuario si no existe
			newWatchlist := &model.Watchlist{
				UserID: userID.(uint),
				Name:   "Mi lista de seguimiento",
			}
			if errCreate := h.watchlistRepo.Create(ctx, newWatchlist); errCreate != nil {
				log.Printf("[Watchlist] Error al crear watchlist para usuario=%d: %v", userID.(uint), errCreate)
			} else {
				log.Printf("[Watchlist] Creada watchlist para usuario=%d", userID.(uint))
			}
		}

		// Comprobar si ya existe en la watchlist
		if exists, _ := h.watchlistItemRepo.IsProductInWatchlist(ctx, userID.(uint), uint(productID)); !exists {
			if errCreate := h.watchlistItemRepo.Create(ctx, &model.WatchlistItem{
				UserID:      userID.(uint),
				ProductID:   uint(productID),
				TargetPrice: targetPrice,
			}); errCreate != nil {
				log.Printf("[Watchlist] Error al insertar item usuario=%d producto=%d: %v", userID.(uint), productID, errCreate)
			} else {
				log.Printf("[Watchlist] Item añadido usuario=%d producto=%d precio=%v", userID.(uint), productID, targetPrice)
			}
		} else {
			// Si ya existe, actualizar el target_price
			items, errItems := h.watchlistItemRepo.FindByUserID(ctx, userID.(uint))
			if errItems == nil {
				for _, item := range items {
					if item.ProductID == uint(productID) {
						item.TargetPrice = targetPrice
						if errUpdate := h.watchlistItemRepo.Update(ctx, item); errUpdate != nil {
							log.Printf("[Watchlist] Error al actualizar precio objetivo de item usuario=%d producto=%d: %v",
								userID.(uint), productID, errUpdate)
						} else {
							log.Printf("[Watchlist] Precio objetivo actualizado para usuario=%d producto=%d precio=%v",
								userID.(uint), productID, targetPrice)
						}
						break
					}
				}
			}
		}
	}

	// Si llegamos aquí, la operación fue exitosa (creación o actualización)
	// y la sincronización con watchlist_items se intentó (los errores se loguearon pero no detuvieron el flujo principal de la alerta).

	// Determinar mensaje de éxito
	successMessage := "Producto añadido a tu cesta correctamente."
	if isUpdate {
		successMessage = "Este producto ya estaba en tu cesta. Precio objetivo actualizado correctamente."
	}

	if isAjax {
		log.Printf("[SetPriceAlert] Operación AJAX exitosa (isUpdate: %t) para producto %s. Devolviendo JSON.", isUpdate, productIDStr)
		responseData := gin.H{
			"success":   true,
			"message":   successMessage,
			"is_update": isUpdate,
		}
		if savedAlert != nil {
			responseData["alert_id"] = savedAlert.ID
		}
		c.JSON(http.StatusOK, responseData)
		return // Asegurar que la función termina aquí para la respuesta AJAX
	}

	// Flujo no-AJAX (si alguna vez se usa): Redirigir con un mensaje flash
	log.Printf("[SetPriceAlert] Operación NO-AJAX exitosa (isUpdate: %t) para producto %s. Redirigiendo.", isUpdate, productIDStr)
	session := sessions.Default(c)
	session.AddFlash(successMessage, "success_message") // Usar una clave consistente para mensajes flash
	if err := session.Save(); err != nil {
		log.Printf("[SetPriceAlert] Error al guardar sesión para flash message: %v", err)
		// No es crítico, continuar con la redirección
	}
	// Redirigir de vuelta a la página del producto. El productoIDStr ya fue validado.
	c.Redirect(http.StatusFound, "/producto/"+productIDStr)
}

// DeletePriceAlert maneja la eliminación de una alerta de precio
func (h *PriceAlertHandler) DeletePriceAlert(c *gin.Context) {
	// Obtener el usuario de la sesión
	userID := sessions.Default(c).Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener ID de la alerta
	alertIDStr := c.Query("id")
	if alertIDStr == "" {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Message": "ID de alerta no proporcionado",
		})
		return
	}

	alertID, err := strconv.ParseUint(alertIDStr, 10, 32)
	if err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Message": "ID de alerta inválido",
			"Error":   err.Error(),
		})
		return
	}

	ctx := c.Request.Context()

	// 1. Buscar la alerta para obtener información antes de eliminarla
	var alert *model.PriceAlert
	if userAlerts, errAlerts := h.priceAlertUseCase.GetUserAlerts(ctx, userID.(uint)); errAlerts == nil {
		for _, a := range userAlerts {
			if a.ID == uint(alertID) {
				alert = a
				break
			}
		}
	}

	if alert == nil {
		h.templateRenderer.Render(c, http.StatusNotFound, "error.html", gin.H{
			"Message": "Alerta no encontrada",
		})
		return
	}

	// 2. Eliminar las notificaciones relacionadas con esta alerta primero
	// Esto lo haría normalmente el notificationRepository, pero como no tenemos acceso
	// directo aquí, lo manejamos a través del price_alert_usecase

	if err := h.priceAlertUseCase.PrepareDeleteAlert(ctx, uint(alertID)); err != nil {
		log.Printf("[ERROR] Error preparando eliminación de alerta ID=%d: %v", alertID, err)
	}

	// 3. Ahora intentar eliminar la alerta
	if err := h.priceAlertUseCase.DeleteAlert(ctx, uint(alertID), userID.(uint)); err != nil {
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al eliminar la alerta de precio",
			"Error":   err.Error(),
		})
		return
	}

	// 4. También eliminar el elemento de la watchlist si existe
	if alert != nil {
		if items, errItems := h.watchlistItemRepo.FindByUserID(ctx, userID.(uint)); errItems == nil {
			for _, itm := range items {
				if itm.ProductID == alert.ProductID {
					if err := h.watchlistItemRepo.Delete(ctx, itm.ID); err != nil {
						log.Printf("[ERROR] Error eliminando watchlist item ID=%d: %v", itm.ID, err)
					} else {
						log.Printf("[INFO] Eliminado watchlist item ID=%d", itm.ID)
					}
					break
				}
			}
		}
	}

	// Redirigir al perfil del usuario
	c.Redirect(http.StatusFound, "/watchlist")
}

// ShowPriceAlerts muestra todas las alertas de precio del usuario
func (h *PriceAlertHandler) ShowPriceAlerts(c *gin.Context) {
	// Obtener el usuario de la sesión
	userID := sessions.Default(c).Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener las alertas del usuario
	ctx := c.Request.Context()
	alerts, err := h.priceAlertUseCase.GetUserAlerts(ctx, userID.(uint))
	if err != nil {
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al obtener alertas de precio",
			"Error":   err.Error(),
		})
		return
	}

	// Renderizar página de alertas
	h.templateRenderer.Render(c, http.StatusOK, "price_alerts.html", gin.H{
		"Title":       "Mis Alertas de Precio - Comparador de Precios",
		"PriceAlerts": alerts,
		"User":        c.MustGet("user"),
	})
}

// ShowWatchlist muestra la lista de deseos del usuario (productos seguidos con sus alertas de precio)
func (h *PriceAlertHandler) ShowWatchlist(c *gin.Context) {
	// Obtener el usuario de la sesión
	userID := sessions.Default(c).Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener las alertas del usuario (que también son los productos seguidos)
	ctx := c.Request.Context()
	alerts, err := h.priceAlertUseCase.GetUserAlerts(ctx, userID.(uint))
	if err != nil {
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al obtener tu lista de deseos",
			"Error":   err.Error(),
		})
		return
	}

	// Para cada alerta, obtener también el precio actual del producto
	type WatchlistItem struct {
		Alert        *model.PriceAlert
		Product      *model.Product
		CurrentPrice *model.Price
		PriceDiff    float64
	}

	var watchlistItems []WatchlistItem

	for _, alert := range alerts {
		// Obtener detalles del producto
		product, err := h.productUseCase.GetProductDetail(ctx, alert.ProductID)
		if err != nil {
			continue // Saltamos este producto si hay error
		}

		// Obtener precio actual (el primero en la lista si existe)
		var currentPrice *model.Price
		if len(product.Prices) > 0 {
			price := product.Prices[0]
			currentPrice = &price
		}

		if currentPrice == nil {
			// Añadir el item sin precio actual
			watchlistItems = append(watchlistItems, WatchlistItem{
				Alert:   alert,
				Product: product,
			})
			continue
		}

		// Calcular diferencia de precio (positivo = falta para alcanzar el objetivo)
		priceDiff := alert.TargetPrice - currentPrice.Price

		// Añadir a la lista
		watchlistItems = append(watchlistItems, WatchlistItem{
			Alert:        alert,
			Product:      product,
			CurrentPrice: currentPrice,
			PriceDiff:    priceDiff,
		})
	}

	// Obtener categorías para el menú desplegable
	allCategories, _ := c.Get("allCategories")

	// Obtener notificaciones no leídas para el contador
	unreadNotifications, _ := c.Get("unreadNotifications")

	// Renderizar página de lista de deseos
	h.templateRenderer.Render(c, http.StatusOK, "watchlist.html", gin.H{
		"Title":               "Mi Cesta - Comparador de Precios",
		"WatchlistItems":      watchlistItems,
		"User":                c.MustGet("user"),
		"Categories":          allCategories,
		"UnreadNotifications": unreadNotifications,
		"PriceAlerts":         alerts,
	})
}

// UpdatePriceAlert actualiza el precio objetivo de una alerta existente
func (h *PriceAlertHandler) UpdatePriceAlert(c *gin.Context) {
	// Obtener el usuario de la sesión
	userID := sessions.Default(c).Get("user_id")
	if userID == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener ID de la alerta y el nuevo precio objetivo
	alertIDStr := c.Query("id")
	targetPriceStr := c.Query("price")

	if alertIDStr == "" {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Message": "ID de alerta no proporcionado",
		})
		return
	}

	alertID, err := strconv.ParseUint(alertIDStr, 10, 32)
	if err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Message": "ID de alerta inválido",
			"Error":   err.Error(),
		})
		return
	}

	targetPrice, err := strconv.ParseFloat(targetPriceStr, 64)
	if err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Message": "Precio objetivo inválido",
			"Error":   err.Error(),
		})
		return
	}

	// Actualizar la alerta
	ctx := c.Request.Context()
	_, err = h.priceAlertUseCase.UpdateAlert(
		ctx,
		uint(alertID),
		userID.(uint),
		targetPrice,
		true, // notificar por email
		true, // alerta activa
	)

	if err != nil {
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "Error al actualizar la alerta de precio",
			"Error":   err.Error(),
		})
		return
	}

	// Redirigir a la watchlist
	c.Redirect(http.StatusFound, "/watchlist")
}
