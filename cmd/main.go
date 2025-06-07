package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/domain/model"
	"app/internal/infrastructure/email"
	"app/internal/infrastructure/persistance"
	"app/internal/interface/cron"
	"app/internal/interface/web/router"
	"app/internal/usecase"
	"app/pkg/config"

	"github.com/joho/godotenv"
)

// go run ./cmd/main.go
// go run ./cmd/main.go --build
// go run ./cmd/main.go --swagger
// go run ./cmd/main.go -test
// go run ./cmd/main.go -test -product-url="https://www.pccomponentes.com/producto"

// @title						V2
// @description				Documentación usando Swagger
// @contact.email				spalomino@netelcomunicaciones.es
// @BasePath					/
// @schemes http https
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						X-API-Key
func main() {
	// Parsear argumentos de línea de comandos
	testMode := flag.Bool("test", false, "Ejecutar en modo prueba sin iniciar el servidor web")
	productURL := flag.String("product-url", "", "URL de un producto específico para hacer scraping (solo con -test)")
	flag.Parse()

	// Cargar variables de entorno
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Println("Archivo .env no encontrado, usando variables de entorno del sistema")
	}

	// Cargar configuración
	config.InitConfig()

	// Verificar configuración de email cargada
	log.Printf("Configuración de email cargada: Host: %s, Puerto: %d, Usuario: %s",
		config.Config.Email.SMTPHost,
		config.Config.Email.SMTPPort,
		config.Config.Email.SMTPUser)

	// Conectar a la base de datos usando la configuración
	db, err := persistance.NewDatabase(
		config.Config.Database.Username,
		config.Config.Database.Password,
		config.Config.Database.Host,
		fmt.Sprintf("%d", config.Config.Database.Port),
		config.Config.Database.Name,
	)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Ejecutar migraciones
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Error al migrar la base de datos: %v", err)
	}

	// Inicializar servicio de email
	mailer := email.NewMailer()

	// --------------------------------------
	// Repositorios y casos de uso
	// --------------------------------------
	userRepo := persistance.NewUserRepository(db.DB)
	productRepo := persistance.NewProductRepository(db.DB)
	categoryRepo := persistance.NewCategoryRepository(db.DB)
	priceRepo := persistance.NewPriceRepository(db.DB)
	priceAlertRepo := persistance.NewPriceAlertRepository(db.DB)
	watchlistRepo := persistance.NewWatchlistRepository(db.DB)
	watchlistItemRepo := persistance.NewWatchlistItemRepository(db.DB)
	notificationRepo := persistance.NewNotificationRepository(db.DB)

	// Crear casos de uso
	productUseCase := usecase.NewProductUseCase(productRepo, categoryRepo, priceRepo)
	userUseCase := usecase.NewUserUseCase(userRepo, mailer)
	scraperUseCase := usecase.NewScraperUseCase(categoryRepo, productRepo, priceRepo)
	priceAlertUseCase := usecase.NewPriceAlertUseCase(
		priceAlertRepo,
		notificationRepo,
		productRepo,
		priceRepo,
		userRepo,
		mailer,
	)

	ctx := context.Background()

	// Modo de prueba para scraping
	if *testMode {
		if *productURL != "" {
			// Si se proporciona una URL de producto específica, hacer scraping solo de ese producto
			log.Printf("Iniciando scraping del producto individual: %s", *productURL)
			// Usar la primera categoría (ID=1, portátiles) por defecto
			product, err := scraperUseCase.ScrapeProductDetails(ctx, *productURL, 1)
			if err != nil {
				log.Printf("Error al hacer scraping del producto: %v", err)
			} else {
				log.Printf("Producto scrapeado exitosamente: %s", product.Name)
				log.Printf("Precio: %.2f EUR", product.Prices[0].Price)
				log.Printf("URL de imagen: %s", product.ImageURL)
				log.Printf("Descripción: %s", product.Description)
			}
		} else {
			// Ejecutar scraping normal
			if err := scraperUseCase.ScrapeAllCategories(ctx); err != nil {
				log.Printf("Error al ejecutar scraping: %v", err)
			}
		}
		log.Println("Scraping finalizado. Saliendo...")
		return
	}

	// Crear categorías predefinidas si no existen
	createInitialCategories(ctx, db)

	// --------------------------------------
	// Configurar router
	// --------------------------------------
	r := router.SetupRouter(productUseCase, userUseCase, priceAlertUseCase, watchlistRepo, watchlistItemRepo)

	// --------------------------------------
	// Scheduler de scraping
	// --------------------------------------
	scheduler := cron.NewScraperScheduler(productRepo, priceRepo, categoryRepo, priceAlertUseCase)
	scheduler.Start()
	defer scheduler.Stop()

	// Iniciar servidor HTTP
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = fmt.Sprintf("%d", config.Config.App.Port)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// Manejar graceful shutdown
	go func() {
		log.Printf("Servidor iniciado en http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Capturar señales para shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Apagando servidor...")

	// Dar tiempo para que se completen las solicitudes activas
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error en el apagado del servidor: %v", err)
	}

	log.Println("Servidor apagado correctamente")
}

// createInitialCategories crea las categorías iniciales si no existen
func createInitialCategories(ctx context.Context, db *persistance.Database) {
	// Definir las categorías según el archivo de configuración
	categories := []model.Category{
		{Name: "Portátiles", Slug: "portatiles"},
		{Name: "Tarjetas Gráficas", Slug: "tarjetas-graficas"},
		{Name: "Auriculares", Slug: "auriculares"},
		{Name: "Teclados", Slug: "teclados"},
		{Name: "Monitores", Slug: "monitores"},
		{Name: "Discos SSD", Slug: "ssd"},
	}

	for _, cat := range categories {
		var count int64
		db.Model(&model.Category{}).Where("slug = ?", cat.Slug).Count(&count)
		if count == 0 {
			db.Create(&cat)
			log.Printf("Categoría creada: %s", cat.Name)
		}
	}
}
