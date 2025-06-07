package cron

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"app/internal/domain/model"
	"app/internal/domain/repositories"
	"app/internal/infrastructure/scraper"
	"app/internal/usecase"
	"app/pkg/utils"

	"github.com/robfig/cron/v3"
)

// Colores para los logs
const (
	// Colores ANSI para texto
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorGray   = "\033[90m"

	// Estilos
	colorBold = "\033[1m"
	colorDim  = "\033[2m"
)

// Funciones para imprimir logs con colores
func logInfo(format string, v ...interface{}) {
	log.Printf("%s%s%s", colorBlue, fmt.Sprintf(format, v...), colorReset)
}

func logSuccess(format string, v ...interface{}) {
	log.Printf("%s%s%s", colorGreen, fmt.Sprintf(format, v...), colorReset)
}

func logWarning(format string, v ...interface{}) {
	log.Printf("%s%s%s", colorYellow, fmt.Sprintf(format, v...), colorReset)
}

func logError(format string, v ...interface{}) {
	log.Printf("%s%s%s", colorRed, fmt.Sprintf(format, v...), colorReset)
}

func logDebug(format string, v ...interface{}) {
	log.Printf("%s%s%s", colorGray, fmt.Sprintf(format, v...), colorReset)
}

// ScraperScheduler gestiona la ejecución periódica de los scrapers
type ScraperScheduler struct {
	cron              *cron.Cron
	productRepo       repositories.ProductRepository
	priceRepo         repositories.PriceRepository
	categoryRepo      repositories.CategoryRepository
	priceAlertUseCase *usecase.PriceAlertUseCase
	ebayScraper       *scraper.EbayScraper
	coolmodScraper    *scraper.CoolmodScraper
	aussarScraper     *scraper.AussarScraper
}

// NewScraperScheduler crea una nueva instancia del planificador de tareas
func NewScraperScheduler(
	productRepo repositories.ProductRepository,
	priceRepo repositories.PriceRepository,
	categoryRepo repositories.CategoryRepository,
	priceAlertUseCase *usecase.PriceAlertUseCase,
) *ScraperScheduler {
	return &ScraperScheduler{
		cron:              cron.New(),
		productRepo:       productRepo,
		priceRepo:         priceRepo,
		categoryRepo:      categoryRepo,
		priceAlertUseCase: priceAlertUseCase,
		ebayScraper:       scraper.NewEbayScraper(),
		coolmodScraper:    scraper.NewCoolmodScraper(),
		aussarScraper:     scraper.NewAussarScraper(),
	}
}

// Start inicia el planificador
func (s *ScraperScheduler) Start() {
	// Ejecutar scraping cada 48 horas
	s.cron.AddFunc("@every 48h", func() {
		s.RunAllScrapers()
	})

	// Ejecutar limpieza de precios viejos cada 72 horas (3 días)
	s.cron.AddFunc("@every 72h", func() {
		s.CleanupOldPrices()
	})

	// Ejecutar verificación de alertas de precio cada 6 horas
	s.cron.AddFunc("@every 6h", func() {
		s.CheckPriceAlerts()
	})

	// También ejecutamos una vez al iniciar
	go s.RunAllScrapers()

	// Y verificamos alertas al iniciar
	go s.CheckPriceAlerts()

	s.cron.Start()
	logSuccess("[SISTEMA] Sistema de scraping iniciado correctamente")
	logInfo("[SISTEMA] Próxima ejecución completa en 48 horas")
}

// Stop detiene el planificador
func (s *ScraperScheduler) Stop() {
	s.cron.Stop()
	logWarning("[SISTEMA] Sistema de scraping detenido")
}

// CheckPriceAlerts ejecuta la verificación de alertas de precio
func (s *ScraperScheduler) CheckPriceAlerts() {
	logInfo("[ALERTAS] Iniciando verificación de alertas de precio...")
	ctx := context.Background()

	if err := s.priceAlertUseCase.CheckPriceAlerts(ctx); err != nil {
		logError("[ALERTAS] Error en la verificación de alertas: %v", err)
	} else {
		logSuccess("[ALERTAS] Verificación completada correctamente")
	}
}

// RunAllScrapers ejecuta todos los scrapers para todas las categorías
func (s *ScraperScheduler) RunAllScrapers() {
	logInfo("[SCRAPING] 🔎 Iniciando proceso de scraping...")

	// Crear un contexto
	ctx := context.Background()

	// Obtener todas las categorías
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		logError("[SCRAPING] No se pudieron obtener las categorías: %v", err)
		return
	}

	// Si no hay categorías, no hay nada que hacer
	if len(categories) == 0 {
		logWarning("[SCRAPING] No hay categorías definidas para realizar scraping")
		return
	}

	logInfo("[SCRAPING] Se procesarán %d categorías", len(categories))

	// Para cada categoría, ejecutar cada scraper
	for _, category := range categories {
		s.scrapCategory(ctx, category)
	}

	logSuccess("[SCRAPING] ✅ Solicitudes de scraping enviadas para todas las categorías")

	// Después de actualizar los precios, verificamos si se han activado alertas
	s.CheckPriceAlerts()
}

// scrapCategory ejecuta todos los scrapers para una categoría específica
func (s *ScraperScheduler) scrapCategory(ctx context.Context, category *model.Category) {
	logInfo("[SCRAPING] Procesando categoría: %s", category.Name)

	// Ejecutar eBay, Coolmod y Aussar
	go s.scrapWithEbay(ctx, category)
	go s.scrapWithCoolmod(ctx, category)
	go s.scrapWithAussar(ctx, category)
}

// scrapWithEbay ejecuta el scraper de eBay
func (s *ScraperScheduler) scrapWithEbay(ctx context.Context, category *model.Category) {
	products, err := s.ebayScraper.ScrapCategory(category)
	if err != nil {
		logError("[EBAY] Error en categoría %s: %v", category.Name, err)
		return
	}

	logInfo("[EBAY] Obtenidos %d productos para %s", len(products), category.Name)
	s.saveProducts(ctx, products, "eBay", category.Name)
}

// scrapWithCoolmod ejecuta el scraper de Coolmod
func (s *ScraperScheduler) scrapWithCoolmod(ctx context.Context, category *model.Category) {
	products, err := s.coolmodScraper.ScrapCategory(category)
	if err != nil {
		logError("[COOLMOD] Error en categoría %s: %v", category.Name, err)
		return
	}

	// Revisar si hay productos antes de guardar
	if len(products) > 0 {
		logInfo("[COOLMOD] Obtenidos %d productos para %s", len(products), category.Name)
		s.saveProducts(ctx, products, "Coolmod", category.Name)
	} else {
		logWarning("[COOLMOD] No se encontraron productos para %s", category.Name)
	}
}

// scrapWithAussar ejecuta el scraper de Aussar
func (s *ScraperScheduler) scrapWithAussar(ctx context.Context, category *model.Category) {
	products, err := s.aussarScraper.ScrapCategory(category)
	if err != nil {
		logError("[AUSSAR] Error en categoría %s: %v", category.Name, err)
		return
	}

	// Revisar si hay productos antes de guardar
	if len(products) > 0 {
		logInfo("[AUSSAR] Obtenidos %d productos para %s", len(products), category.Name)
		s.saveProducts(ctx, products, "Aussar", category.Name)
	} else {
		logWarning("[AUSSAR] No se encontraron productos para %s", category.Name)
	}
}

// checkSlugExists verifica si un slug ya existe en la base de datos
func (s *ScraperScheduler) checkSlugExists(ctx context.Context, slug string) bool {
	exists, err := s.productRepo.ExistsBySlug(ctx, slug)
	if err != nil {
		logDebug("Error al verificar existencia de slug: %v\n", err)
		return false
	}
	return exists
}

// saveProducts guarda los productos y precios en la base de datos
func (s *ScraperScheduler) saveProducts(ctx context.Context, products []*model.Product, source string, categoryName string) {
	logDebug("[%s] Procesando %d productos de %s", source, len(products), categoryName)

	// Activar logs detallados solo si está en modo debug
	debugLogsEnabled := false

	// Contar estadísticas
	validProductsCount := 0
	reclassifiedCount := 0
	discardedCount := 0

	// Crear un mapa para almacenar todos los IDs de productos que vamos a procesar
	processedProductIDs := make(map[uint]bool)

	for _, product := range products {
		// VALIDACIÓN DE CATEGORÍA: Verificar que el producto pertenezca a la categoría asignada
		originalCategoryID := product.CategoryID

		if !utils.ValidateProductCategory(product) {
			if debugLogsEnabled {
				logDebug("[CATEGORÍA] 🚫 '%s' NO pertenece a la categoría %d",
					truncateString(product.Name, 30), product.CategoryID)
			}

			// Obtener todas las categorías
			categories, err := s.categoryRepo.GetAll(ctx)
			if err != nil {
				logError("[ERROR] Error al obtener categorías: %v", err)
				discardedCount++
				continue // Saltamos este producto
			}

			// Probar cada categoría hasta encontrar una que sea válida para el producto
			matchFound := false
			for _, category := range categories {
				// Crear una copia temporal del producto para probar con otra categoría
				tempProduct := *product
				tempProduct.CategoryID = category.ID

				// Verificar si el producto es válido para esta categoría
				if utils.ValidateProductCategory(&tempProduct) {
					// Actualizar la categoría del producto original
					product.CategoryID = category.ID
					if debugLogsEnabled {
						logDebug("[CATEGORÍA] ✅ '%s' reclasificado de categoría %d a %s (ID: %d)",
							truncateString(product.Name, 30), originalCategoryID, category.Name, category.ID)
					}
					matchFound = true
					reclassifiedCount++
					break
				}
			}

			// Si no se encontró una categoría adecuada, descartamos el producto
			if !matchFound {
				if debugLogsEnabled {
					logDebug("[CATEGORÍA] ❌ No se encontró categoría para '%s'",
						truncateString(product.Name, 30))
				}
				discardedCount++
				continue
			}
		} else {
			if debugLogsEnabled {
				logDebug("[CATEGORÍA] ✅ '%s' validado para categoría %d",
					truncateString(product.Name, 30), product.CategoryID)
			}
			validProductsCount++
		}

		// CONTINUAR CON EL PROCESO DE GUARDADO NORMAL
		// Asegurarnos de que el producto tenga un slug válido
		if product.Slug == "" {
			product.Slug = utils.GenerateSlug(product.Name)
		}

		// Asegurarnos de que el slug sea único
		product.Slug = utils.GenerateUniqueSlug(product.Name, func(slug string) bool {
			return s.checkSlugExists(ctx, slug)
		})

		// Verificar si ya existe un producto similar por nombre
		// Pasamos un string vacío para storeFilter porque queremos todos los productos de la categoría para comparar,
		// la lógica de la tienda se maneja después al iterar sobre los precios.
		allProducts, err := s.productRepo.FindByCategory(ctx, product.CategoryID, 1000, 0, "")
		if err != nil {
			logDebug("Error al buscar productos existentes: %v\n", err)
			continue
		}

		var existingProduct *model.Product
		for _, p := range allProducts {
			// Comparamos nombres normalizados para evitar duplicados por variaciones menores
			normalizedName1 := strings.ToLower(strings.TrimSpace(p.Name))
			normalizedName2 := strings.ToLower(strings.TrimSpace(product.Name))

			// Si los nombres son muy similares, consideramos que es el mismo producto
			if normalizedName1 == normalizedName2 ||
				strings.Contains(normalizedName1, normalizedName2) ||
				strings.Contains(normalizedName2, normalizedName1) {
				existingProduct = p
				break
			}
		}

		if existingProduct != nil {
			// Registrar que hemos procesado este producto
			processedProductIDs[existingProduct.ID] = true

			// El producto ya existe, guardamos solo el nuevo precio
			price := product.Prices[0] // Asumimos que hay al menos un precio

			// Muy importante: asegurar que el ID sea cero para que GORM asigne uno nuevo
			price.ID = 0

			price.ProductID = existingProduct.ID

			// Verificamos si ya existe un precio para esta tienda
			prices, err := s.priceRepo.FindByProductID(ctx, existingProduct.ID)
			if err != nil {
				logDebug("Error al buscar precios existentes: %v\n", err)
				continue
			}

			var existingPrice *model.Price
			for _, p := range prices {
				if p.Store == price.Store {
					existingPrice = p
					break
				}
			}

			if existingPrice != nil {
				// Actualizamos el precio existente
				existingPrice.Price = price.Price
				existingPrice.URL = price.URL
				existingPrice.RetrievedAt = price.RetrievedAt

				if err := s.priceRepo.Update(ctx, existingPrice); err != nil {
					logDebug("Error al actualizar precio para %s: %v\n", existingProduct.Name, err)
				}
			} else {
				// Creamos un nuevo precio
				if err := s.priceRepo.Create(ctx, &price); err != nil {
					logDebug("Error al crear precio para %s: %v\n", existingProduct.Name, err)
				}
			}
		} else {
			// El producto no existe, lo creamos junto con su precio
			if err := s.productRepo.Create(ctx, product); err != nil {
				logDebug("Error al guardar producto %s: %v\n", product.Name, err)
				continue
			}

			// Registrar que hemos procesado este producto
			processedProductIDs[product.ID] = true

			// El ID del producto se ha generado automáticamente
			price := product.Prices[0]

			// Muy importante: asegurar que el ID sea cero para que GORM asigne uno nuevo
			price.ID = 0

			price.ProductID = product.ID

			if err := s.priceRepo.Create(ctx, &price); err != nil {
				logDebug("Error al crear precio para nuevo producto %s: %v\n", product.Name, err)
			}
		}
	}

	logSuccess("[%s-%s] Total: %d | Válidos: %d | Reclasificados: %d | Descartados: %d",
		source, categoryName, len(products), validProductsCount, reclassifiedCount, discardedCount)

	// Eliminar precios antiguos que no fueron actualizados en esta ejecución
	// Solo para productos procesados en esta ejecución
	deletedPricesCount := 0
	for productID := range processedProductIDs {
		// Obtenemos todos los precios de este producto
		prices, err := s.priceRepo.FindByProductID(ctx, productID)
		if err != nil {
			logDebug("Error al obtener precios para limpiar antiguos: %v\n", err)
			continue
		}

		// Identificamos precios que tienen más de 3 días
		now := time.Now()
		minFreshTime := now.AddDate(0, 0, -3) // Precios de hace 3 días o más son antiguos

		for _, price := range prices {
			if price.RetrievedAt.Before(minFreshTime) {
				// Eliminar precio antiguo
				if err := s.priceRepo.Delete(ctx, price.ID); err != nil {
					logDebug("Error al eliminar precio antiguo ID %d: %v\n", price.ID, err)
				} else {
					deletedPricesCount++
					if debugLogsEnabled {
						logDebug("Precio antiguo eliminado: ID %d, Producto ID %d, Tienda %s\n",
							price.ID, price.ProductID, price.Store)
					}
				}
			}
		}
	}

	if deletedPricesCount > 0 {
		logInfo("[LIMPIEZA-%s] Se eliminaron %d precios antiguos", source, deletedPricesCount)
	}
}

// CleanupOldPrices elimina precios antiguos de todos los productos
func (s *ScraperScheduler) CleanupOldPrices() {
	logInfo("[LIMPIEZA] Iniciando eliminación de precios antiguos...")
	ctx := context.Background()

	// Precios más antiguos de 7 días
	oldPriceDate := time.Now().AddDate(0, 0, -7)

	count, err := s.priceRepo.DeleteOldPrices(ctx, oldPriceDate)
	if err != nil {
		logError("[LIMPIEZA] Error al eliminar precios antiguos: %v", err)
		return
	}

	logSuccess("[LIMPIEZA] ✅ Eliminados %d precios antiguos con éxito", count)
}

// Función auxiliar para truncar strings largos en los logs
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}
