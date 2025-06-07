package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"app/internal/domain/model"
	"app/internal/domain/repositories"
	"app/internal/infrastructure/email"
)

// PriceAlertUseCase gestiona las alertas de precio
type PriceAlertUseCase struct {
	priceAlertRepo   repositories.PriceAlertRepository
	notificationRepo repositories.NotificationRepository
	productRepo      repositories.ProductRepository
	priceRepo        repositories.PriceRepository
	userRepo         repositories.UserRepository
	mailer           *email.Mailer
}

// NewPriceAlertUseCase crea una nueva instancia del caso de uso de alertas de precio
func NewPriceAlertUseCase(
	priceAlertRepo repositories.PriceAlertRepository,
	notificationRepo repositories.NotificationRepository,
	productRepo repositories.ProductRepository,
	priceRepo repositories.PriceRepository,
	userRepo repositories.UserRepository,
	mailer *email.Mailer,
) *PriceAlertUseCase {
	return &PriceAlertUseCase{
		priceAlertRepo:   priceAlertRepo,
		notificationRepo: notificationRepo,
		productRepo:      productRepo,
		priceRepo:        priceRepo,
		userRepo:         userRepo,
		mailer:           mailer,
	}
}

// CreateAlert crea una nueva alerta de precio
func (uc *PriceAlertUseCase) CreateAlert(ctx context.Context, userID, productID uint, targetPrice float64, notifyByEmail bool) (*model.PriceAlert, error) {
	// Verificar que el usuario existe
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("usuario no encontrado: %w", err)
	}

	// Verificar que el producto existe
	product, err := uc.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("producto no encontrado: %w", err)
	}

	// Crear la alerta
	alert := &model.PriceAlert{
		UserID:        userID,
		ProductID:     productID,
		TargetPrice:   targetPrice,
		NotifyByEmail: notifyByEmail,
		IsActive:      true,
	}

	// Guardar la alerta en la base de datos
	if err := uc.priceAlertRepo.Create(ctx, alert); err != nil {
		return nil, fmt.Errorf("error al crear alerta: %w", err)
	}

	// Verificar inmediatamente si el precio actual ya cumple con la alerta
	bestPrice, err := uc.priceRepo.FindBestPriceByProductID(ctx, productID)
	if err != nil {
		log.Printf("Error al obtener mejor precio para verificación inicial: %v", err)
		return alert, nil
	}

	// Si ya existe un precio menor o igual al objetivo, notificar inmediatamente
	if bestPrice != nil && bestPrice.Price <= targetPrice {
		uc.createNotification(ctx, alert, product, user, bestPrice)
	}

	return alert, nil
}

// UpdateAlert actualiza una alerta existente
func (uc *PriceAlertUseCase) UpdateAlert(ctx context.Context, alertID, userID uint, targetPrice float64, notifyByEmail bool, isActive bool) (*model.PriceAlert, error) {
	// Buscar la alerta
	alert, err := uc.priceAlertRepo.FindByID(ctx, alertID)
	if err != nil {
		return nil, fmt.Errorf("alerta no encontrada: %w", err)
	}

	// Verificar que el usuario es el propietario de la alerta
	if alert.UserID != userID {
		return nil, fmt.Errorf("no tienes permiso para modificar esta alerta")
	}

	// Actualizar los campos
	alert.TargetPrice = targetPrice
	alert.NotifyByEmail = notifyByEmail
	alert.IsActive = isActive

	// Guardar los cambios
	if err := uc.priceAlertRepo.Update(ctx, alert); err != nil {
		return nil, fmt.Errorf("error al actualizar alerta: %w", err)
	}

	// Verificar inmediatamente si el precio actual ya cumple con la alerta
	// Obtener el producto para la notificación
	product, err := uc.productRepo.FindByID(ctx, alert.ProductID)
	if err != nil {
		log.Printf("Error al obtener producto para verificación inmediata: %v", err)
		return alert, nil
	}

	// Obtener el usuario para la notificación
	user, err := uc.userRepo.FindByID(ctx, alert.UserID)
	if err != nil {
		log.Printf("Error al obtener usuario para verificación inmediata: %v", err)
		return alert, nil
	}

	// Verificar el precio actual
	bestPrice, err := uc.priceRepo.FindBestPriceByProductID(ctx, alert.ProductID)
	if err != nil {
		log.Printf("Error al obtener mejor precio para verificación inmediata: %v", err)
		return alert, nil
	}

	// Si ya existe un precio menor o igual al objetivo, notificar inmediatamente
	if bestPrice != nil && bestPrice.Price <= targetPrice {
		log.Printf("El precio actual (%.2f) ya está por debajo del objetivo (%.2f), notificando", bestPrice.Price, targetPrice)
		uc.createNotification(ctx, alert, product, user, bestPrice)
	}

	return alert, nil
}

// DeleteAlert elimina una alerta
func (uc *PriceAlertUseCase) DeleteAlert(ctx context.Context, alertID, userID uint) error {
	// Buscar la alerta
	alert, err := uc.priceAlertRepo.FindByID(ctx, alertID)
	if err != nil {
		return fmt.Errorf("alerta no encontrada: %w", err)
	}

	// Verificar que el usuario es el propietario de la alerta
	if alert.UserID != userID {
		return fmt.Errorf("no tienes permiso para eliminar esta alerta")
	}

	// Eliminar la alerta
	if err := uc.priceAlertRepo.Delete(ctx, alertID); err != nil {
		return fmt.Errorf("error al eliminar alerta: %w", err)
	}

	return nil
}

// GetUserAlerts obtiene todas las alertas de un usuario
func (uc *PriceAlertUseCase) GetUserAlerts(ctx context.Context, userID uint) ([]*model.PriceAlert, error) {
	return uc.priceAlertRepo.FindByUserID(ctx, userID)
}

// CheckPriceAlerts verifica todas las alertas activas contra los precios actuales
// Este método es el que debería ser llamado por un scheduler periódicamente
func (uc *PriceAlertUseCase) CheckPriceAlerts(ctx context.Context) error {
	// Este proceso puede ser pesado, así que usamos una hora límite
	startTime := time.Now()
	log.Printf("Iniciando verificación de alertas de precio a las %s", startTime.Format("15:04:05"))

	// Obtener todos los productos que tienen al menos un precio
	// En un sistema real, podríamos querer hacer esto más eficientemente filtrando
	// solo a productos con cambios de precio recientes

	// Para cada producto, buscar el mejor precio actual
	// Para cada precio que ha cambiado, buscar alertas que se activen con este precio
	// Esto se puede optimizar de muchas maneras según el volumen de datos

	// Ejecutamos la búsqueda de alertas activadas para cada producto con su mejor precio
	// En un sistema real, esta lógica podría implementarse con una consulta SQL más eficiente

	// Vamos a buscar los precios más bajos primero
	// Esto es una simplificación, en producción podríamos necesitar paginar y procesar por lotes
	bestDeals, err := uc.productRepo.FindBestDeals(ctx, 100) // Limitamos a 100 productos
	if err != nil {
		return fmt.Errorf("error al obtener productos con mejores precios: %w", err)
	}

	for _, product := range bestDeals {
		// Obtener el mejor precio actual para el producto
		bestPrice, err := uc.priceRepo.FindBestPriceByProductID(ctx, product.ID)
		if err != nil || bestPrice == nil {
			log.Printf("No se pudo obtener el mejor precio para producto %d: %v", product.ID, err)
			continue
		}

		// Buscar alertas que se activen con este precio
		alertsToNotify, err := uc.priceAlertRepo.FindActiveAlertsForPrice(ctx, product.ID, bestPrice.Price)
		if err != nil {
			log.Printf("Error al buscar alertas para producto %d: %v", product.ID, err)
			continue
		}

		// Si hay alertas, crear notificaciones
		for _, alert := range alertsToNotify {
			// Obtener el usuario para la notificación
			user, err := uc.userRepo.FindByID(ctx, alert.UserID)
			if err != nil {
				log.Printf("Error al obtener usuario %d para notificación: %v", alert.UserID, err)
				continue
			}

			// Crear notificación
			uc.createNotification(ctx, alert, product, user, bestPrice)
		}
	}

	log.Printf("Verificación de alertas de precio completada en %v por favor espere un momento en breves se estara haciendo el scraping", time.Since(startTime))
	return nil
}

// GetUserNotifications obtiene todas las notificaciones de un usuario
func (uc *PriceAlertUseCase) GetUserNotifications(ctx context.Context, userID uint) ([]*model.Notification, error) {
	// Obtener las últimas 50 notificaciones sin paginación
	const limit = 50
	const offset = 0
	return uc.notificationRepo.FindByUserID(ctx, userID, limit, offset)
}

// MarkNotificationAsRead marca una notificación como leída
func (uc *PriceAlertUseCase) MarkNotificationAsRead(ctx context.Context, notificationID, userID uint) error {
	// Buscar la notificación
	notification, err := uc.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("notificación no encontrada: %w", err)
	}

	// Verificar que el usuario es el propietario de la notificación
	if notification.UserID != userID {
		return fmt.Errorf("no tienes permiso para marcar esta notificación")
	}

	// Marcar como leída
	notification.IsRead = true

	// Guardar los cambios
	if err := uc.notificationRepo.Update(ctx, notification); err != nil {
		return fmt.Errorf("error al actualizar notificación: %w", err)
	}

	return nil
}

// MarkAllNotificationsAsRead marca todas las notificaciones de un usuario como leídas
func (uc *PriceAlertUseCase) MarkAllNotificationsAsRead(ctx context.Context, userID uint) error {
	// Obtener todas las notificaciones no leídas del usuario
	notifications, err := uc.notificationRepo.FindUnreadByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("error al obtener notificaciones: %w", err)
	}

	// Marcar cada una como leída
	for _, notification := range notifications {
		notification.IsRead = true
		if err := uc.notificationRepo.Update(ctx, notification); err != nil {
			log.Printf("Error al marcar notificación %d como leída: %v", notification.ID, err)
		}
	}

	return nil
}

// createNotification crea una notificación y envía correo si está configurado
func (uc *PriceAlertUseCase) createNotification(ctx context.Context, alert *model.PriceAlert, product *model.Product, user *model.User, price *model.Price) {
	// Crear texto para la notificación
	title := fmt.Sprintf("¡Alerta de precio para %s!", product.Name)
	message := fmt.Sprintf("El precio actual es %.2f€ en %s, por debajo de tu objetivo de %.2f€.",
		price.Price, price.Store, alert.TargetPrice)

	// Crear notificación en la base de datos
	notification := &model.Notification{
		UserID:    user.ID,
		ProductID: product.ID,
		AlertID:   &alert.ID,
		Title:     title,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := uc.notificationRepo.Create(ctx, notification); err != nil {
		log.Printf("Error al crear notificación: %v", err)
		// Continuamos para intentar enviar el correo de todas formas
	}

	// Si está configurado el envío de correo, enviar
	if alert.NotifyByEmail && user.Email != "" {
		if err := uc.mailer.SendPriceAlertEmail(
			user.Email,
			user.Username,
			product.Name,
			product.ID,
			alert.TargetPrice,
			price.Price,
			price.Store,
			price.URL,
		); err != nil {
			log.Printf("Error al enviar correo de alerta de precio: %v", err)
		} else {
			log.Printf("Correo de alerta enviado con éxito a %s para producto %s",
				user.Email, product.Name)
		}
	}

	// Desactivar la alerta después de notificar (opcional)
	// Si queremos que la alerta sea de un solo uso y se desactive
	// después de notificar, descomentar las siguientes líneas:
	/*
		alert.IsActive = false
		if err := uc.priceAlertRepo.Update(ctx, alert); err != nil {
			log.Printf("Error al desactivar alerta: %v", err)
		}
	*/
}

// GetUnreadNotificationsCount obtiene el número de notificaciones no leídas para un usuario.
func (uc *PriceAlertUseCase) GetUnreadNotificationsCount(ctx context.Context, userID uint) (int, error) {
	count, err := uc.notificationRepo.CountUnreadByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error al obtener contador de notificaciones no leídas para userID %d: %v", userID, err)
		return 0, fmt.Errorf("error al contar notificaciones no leídas: %w", err)
	}
	return count, nil
}

// PrepareDeleteAlert elimina las notificaciones relacionadas con una alerta antes de eliminarla
// para evitar problemas de clave foránea
func (uc *PriceAlertUseCase) PrepareDeleteAlert(ctx context.Context, alertID uint) error {
	// Buscar todas las notificaciones relacionadas con esta alerta específica
	notifications, err := uc.notificationRepo.FindByAlertID(ctx, alertID)
	if err != nil {
		return fmt.Errorf("error al buscar notificaciones relacionadas: %w", err)
	}

	// Eliminar todas las notificaciones relacionadas con esta alerta
	for _, notification := range notifications {
		if err := uc.notificationRepo.Delete(ctx, notification.ID); err != nil {
			log.Printf("Error al eliminar notificación relacionada ID=%d: %v", notification.ID, err)
		} else {
			log.Printf("Eliminada notificación ID=%d relacionada con alerta ID=%d", notification.ID, alertID)
		}
	}

	return nil
}

// DeleteNotification elimina una notificación específica
func (uc *PriceAlertUseCase) DeleteNotification(ctx context.Context, notificationID, userID uint) error {
	// Verificar que la notificación pertenece al usuario
	notification, err := uc.notificationRepo.FindByID(ctx, notificationID)
	if err != nil {
		return fmt.Errorf("notificación no encontrada: %w", err)
	}

	// Seguridad: verificar que el usuario es el propietario
	if notification.UserID != userID {
		return fmt.Errorf("no tienes permiso para eliminar esta notificación")
	}

	// Eliminar la notificación
	if err := uc.notificationRepo.Delete(ctx, notificationID); err != nil {
		return fmt.Errorf("error al eliminar notificación: %w", err)
	}

	return nil
}
