package router

import (
	"log"
	"net/http"
	"os"

	"app/internal/domain/repositories"
	"app/internal/interface/web/handler"
	"app/internal/interface/web/middleware"
	"app/internal/interface/web/views"
	"app/internal/usecase"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// SetupRouter configura las rutas y handlers de la aplicación
func SetupRouter(productUseCase *usecase.ProductUseCase, userUseCase *usecase.UserUseCase, priceAlertUseCase *usecase.PriceAlertUseCase, watchlistRepo repositories.WatchlistRepository, watchlistItemRepo repositories.WatchlistItemRepository) *gin.Engine {
	// Inicializar Gin
	r := gin.Default()

	// Configurar middleware de sesiones
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "default_secret_key"
	}
	store := cookie.NewStore([]byte(secret))
	r.Use(sessions.Sessions("pricehunter", store))

	// Middleware global para cargar usuario en cada solicitud si está autenticado
	r.Use(middleware.LoadUser(userUseCase))

	// Middleware global para incluir categorías en todas las vistas
	r.Use(middleware.IncludeCategories(productUseCase))

	// Middlewares para datos de usuario logueado
	r.Use(middleware.IncludePriceAlerts(priceAlertUseCase))
	r.Use(middleware.IncludeUnreadNotificationsCount(priceAlertUseCase))

	// Cargar archivos estáticos
	r.Static("/static", "./web/static")

	// Crear renderer para las plantillas
	var templateRenderer *views.TemplateRenderer
	var err error
	if templateRenderer, err = views.SetupTemplates(r); err != nil {
		log.Fatalf("Error al configurar las plantillas HTML: %v", err)
	}

	// Inicializar handlers
	homeHandler := handler.NewHomeHandler(productUseCase, templateRenderer)
	productHandler := handler.NewProductHandler(productUseCase, templateRenderer)
	categoryHandler := handler.NewCategoryHandler(productUseCase, templateRenderer)
	authHandler := handler.NewAuthHandler(userUseCase, templateRenderer)
	notificationHandler := handler.NewNotificationHandler(priceAlertUseCase, templateRenderer)
	priceAlertHandler := handler.NewPriceAlertHandler(priceAlertUseCase, productUseCase, watchlistRepo, watchlistItemRepo, templateRenderer)

	// Rutas públicas
	r.GET("/", homeHandler.GetHome)
	r.GET("/login", authHandler.ShowLoginForm)
	r.POST("/login", authHandler.Login)
	r.GET("/registro", authHandler.ShowRegisterForm)
	r.POST("/registro", authHandler.RegisterHandler)
	r.GET("/registro-exitoso", authHandler.ShowRegisterSuccessPage)
	r.GET("/verificar", authHandler.VerifyEmail)
	r.GET("/producto/:id", productHandler.GetProduct)
	r.GET("/categoria/:slug", categoryHandler.GetCategory)

	// Nuevas rutas públicas para el flujo de "He olvidado mi contraseña"
	r.GET("/forgot-password", authHandler.ShowForgotPasswordForm)
	r.POST("/forgot-password", authHandler.ProcessForgotPasswordForm)

	// Rutas públicas para restablecer contraseña cuando ya se tiene el token
	r.GET("/restablecer-password", authHandler.ShowPasswordResetForm)
	r.POST("/restablecer-password", authHandler.ProcessPasswordReset)

	// API endpoints
	api := r.Group("/api")
	{
		api.GET("/categoria/:slug", categoryHandler.GetCategoryAPI)
		api.POST("/notifications/delete-read", notificationHandler.DeleteReadNotifications)
	}

	// Rutas protegidas (requieren autenticación)
	authorized := r.Group("/")
	authorized.Use(middleware.AuthRequired())
	{
		authorized.GET("/perfil", authHandler.ShowProfile)
		authorized.GET("/logout", authHandler.Logout)

		// Gestión de contraseña y cuenta
		authorized.POST("/cambiar-password", authHandler.ChangePassword)
		authorized.POST("/borrar-cuenta", authHandler.DeleteAccount)

		// Solicitar restablecimiento de contraseña (para usuario LOGUEADO, si se quiere mantener)
		// Esta ruta es diferente al flujo de /forgot-password
		authorized.GET("/solicitar-reset", authHandler.RequestPasswordReset)

		// Notificaciones
		authorized.GET("/notificaciones", notificationHandler.ShowNotifications)
		authorized.POST("/notificaciones/marcar-leida", notificationHandler.MarkAsRead)
		authorized.POST("/notificaciones/marcar-leidas", notificationHandler.MarkAllAsRead)

		// Lista de deseos y alertas (unificado)
		authorized.GET("/watchlist", priceAlertHandler.ShowWatchlist)
		authorized.POST("/price-alert/set", priceAlertHandler.SetPriceAlert)
		authorized.GET("/price-alert/delete", priceAlertHandler.DeletePriceAlert)
		authorized.GET("/price-alert/update", priceAlertHandler.UpdatePriceAlert)

		// Mantener esta ruta por compatibilidad pero redirigir a /watchlist
		authorized.GET("/price-alerts", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/watchlist")
		})
	}

	// Ruta para páginas no encontradas
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Página no encontrada",
		})
	})

	return r
}
