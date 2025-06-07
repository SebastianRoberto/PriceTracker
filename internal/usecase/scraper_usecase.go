package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"app/internal/domain/model"
	"app/internal/domain/repositories"
	"app/internal/infrastructure/scraper"
	"app/pkg/utils"

	"github.com/corona10/goimagehash"
)

// ScraperUseCase implementa la lógica para el scraper de productos
type ScraperUseCase struct {
	categoryRepo   repositories.CategoryRepository
	productRepo    repositories.ProductRepository
	priceRepo      repositories.PriceRepository
	ebayScraper    *scraper.EbayScraper
	coolmodScraper *scraper.CoolmodScraper
	aussarScraper  *scraper.AussarScraper
}

// NewScraperUseCase crea una nueva instancia del caso de uso para scraping
func NewScraperUseCase(
	categoryRepo repositories.CategoryRepository,
	productRepo repositories.ProductRepository,
	priceRepo repositories.PriceRepository,
) *ScraperUseCase {
	return &ScraperUseCase{
		categoryRepo:   categoryRepo,
		productRepo:    productRepo,
		priceRepo:      priceRepo,
		ebayScraper:    scraper.NewEbayScraper(),
		coolmodScraper: scraper.NewCoolmodScraper(),
		aussarScraper:  scraper.NewAussarScraper(),
	}
}

// ScrapeAllCategories ejecuta el scraping de todas las categorías
func (uc *ScraperUseCase) ScrapeAllCategories(ctx context.Context) error {
	categories, err := uc.categoryRepo.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("error al obtener categorías: %w", err)
	}

	for _, category := range categories {
		if err := uc.ScrapeCategory(ctx, category.ID); err != nil {
			log.Printf("Error al scrapear categoría %s: %v\n", category.Name, err)
			continue
		}
	}

	return nil
}

// ScrapeCategory ejecuta el scraping de una categoría específica
func (uc *ScraperUseCase) ScrapeCategory(ctx context.Context, categoryID uint) error {
	category, err := uc.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("error al buscar categoría %d: %w", categoryID, err)
	}

	// Scraping de eBay
	log.Printf("Iniciando scraping de Ebay para categoría: %s (ID: %d)", category.Name, category.ID)
	ebayProducts, err := uc.ebayScraper.ScrapCategory(category)
	if err != nil {
		log.Printf("Error al scrapear Ebay para %s: %v\n", category.Name, err)
	} else {
		log.Printf("Scraping de Ebay para %s completado. %d productos encontrados.", category.Name, len(ebayProducts))
		if len(ebayProducts) > 0 {
			log.Printf("Guardando productos de Ebay para %s...", category.Name)
			if err := uc.saveProducts(ctx, ebayProducts); err != nil {
				log.Printf("Error al guardar productos de Ebay para %s: %v\n", category.Name, err)
			} else {
				log.Printf("Productos de Ebay para %s guardados/actualizados.", category.Name)
			}
		}
	}

	// Scraping de Coolmod
	log.Printf("Iniciando scraping de Coolmod para categoría: %s (ID: %d)", category.Name, category.ID)
	coolmodProducts, err := uc.coolmodScraper.ScrapCategory(category)
	if err != nil {
		log.Printf("Error al scrapear Coolmod para %s: %v\n", category.Name, err)
	} else {
		log.Printf("Scraping de Coolmod para %s completado. %d productos encontrados.", category.Name, len(coolmodProducts))
		if len(coolmodProducts) > 0 {
			log.Printf("Guardando productos de Coolmod para %s...", category.Name)
			if err := uc.saveProducts(ctx, coolmodProducts); err != nil {
				log.Printf("Error al guardar productos de Coolmod para %s: %v\n", category.Name, err)
			} else {
				log.Printf("Productos de Coolmod para %s guardados/actualizados.", category.Name)
			}
		}
	}

	// Scraping de Aussar
	log.Printf("Iniciando scraping de Aussar para categoría: %s (ID: %d)", category.Name, category.ID)
	aussarProducts, err := uc.aussarScraper.ScrapCategory(category)
	if err != nil {
		log.Printf("Error al scrapear Aussar para %s: %v\n", category.Name, err)
	} else {
		log.Printf("Scraping de Aussar para %s completado. %d productos encontrados.", category.Name, len(aussarProducts))
		if len(aussarProducts) > 0 {
			log.Printf("Guardando productos de Aussar para %s...", category.Name)
			if err := uc.saveProducts(ctx, aussarProducts); err != nil {
				log.Printf("Error al guardar productos de Aussar para %s: %v\n", category.Name, err)
			} else {
				log.Printf("Productos de Aussar para %s guardados/actualizados.", category.Name)
			}
		}
	}

	return nil
}

// ScrapeProductDetails obtiene los detalles completos de un producto específico
func (uc *ScraperUseCase) ScrapeProductDetails(ctx context.Context, productURL string, categoryID uint) (*model.Product, error) {
	// Obtener la categoría
	category, err := uc.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error al buscar categoría %d: %w", categoryID, err)
	}

	var product *model.Product

	// Determinar qué scraper usar según la URL
	if utils.IsCoolmodURL(productURL) {
		// Utilizar el scraper de Coolmod
		log.Printf("Iniciando scraping de detalles de producto desde Coolmod: %s", productURL)
		product, err = uc.coolmodScraper.ScrapProductDetails(productURL)
	} else if utils.IsAussarURL(productURL) {
		// Utilizar el scraper de Aussar
		log.Printf("Iniciando scraping de detalles de producto desde Aussar: %s", productURL)
		product, err = uc.aussarScraper.ScrapProductDetails(productURL)
	} else {
		return nil, fmt.Errorf("la URL no pertenece a una tienda soportada: %s", productURL)
	}

	if err != nil {
		return nil, fmt.Errorf("error al obtener detalles del producto: %w", err)
	}

	// Asignar la categoría al producto
	product.CategoryID = category.ID

	// Guardar el producto en la base de datos
	log.Printf("Guardando producto '%s' en la base de datos...", product.Name)
	if err := uc.saveProduct(ctx, product); err != nil {
		return nil, fmt.Errorf("error al guardar producto: %w", err)
	}

	return product, nil
}

// saveProducts guarda o actualiza productos scrapeados en la base de datos
func (uc *ScraperUseCase) saveProducts(ctx context.Context, products []*model.Product) error {
	log.Printf("[SCRAPER] Iniciando proceso de guardado para %d productos", len(products))
	validProductsCount := 0
	reclassifiedCount := 0
	discardedCount := 0

	for _, product := range products {
		// Paso 1: Validar que el producto pertenezca a la categoría asignada
		// Si no es válido para la categoría asignada, intentamos encontrar una categoría más adecuada
		originalCategoryID := product.CategoryID

		if !utils.ValidateProductCategory(product) {
			log.Printf("[CATEGORIZADOR] 🚫 Producto '%s' NO pertenece a la categoría %d, buscando categoría adecuada...",
				product.Name, product.CategoryID)

			// Obtener todas las categorías
			categories, err := uc.categoryRepo.GetAll(ctx)
			if err != nil {
				log.Printf("[ERROR] Error al obtener categorías para reclasificación: %v", err)
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
					log.Printf("[CATEGORIZADOR] ✅ Producto '%s' reclasificado de categoría %d a categoría: %s (ID: %d)",
						product.Name, originalCategoryID, category.Name, category.ID)
					matchFound = true
					reclassifiedCount++
					break
				}
			}

			// Si no se encontró una categoría adecuada, descartamos el producto
			if !matchFound {
				log.Printf("[CATEGORIZADOR] ❌ No se encontró categoría adecuada para '%s', producto descartado", product.Name)
				discardedCount++
				continue
			}
		} else {
			log.Printf("[CATEGORIZADOR] ✅ Producto '%s' validado correctamente para categoría %d",
				product.Name, product.CategoryID)
			validProductsCount++
		}

		// Paso 2: Guardar el producto con la categoría correcta (aquí se aplicará la lógica de deduplicación por imagen)
		if err := uc.saveProduct(ctx, product); err != nil {
			log.Printf("[ERROR] Error al guardar producto '%s': %v", product.Name, err)
			return err
		}
	}

	log.Printf("[RESUMEN CATEGORIZACIÓN] Total: %d | Válidos inicialmente: %d | Reclasificados: %d | Descartados: %d",
		len(products), validProductsCount, reclassifiedCount, discardedCount)

	return nil
}

// saveProduct guarda o actualiza un único producto en la base de datos, aplicando lógica de deduplicación por imagen y slug
func (uc *ScraperUseCase) saveProduct(ctx context.Context, product *model.Product) error {
	// Verificación final de que el producto realmente pertenece a la categoría asignada
	if !utils.ValidateProductCategory(product) {
		log.Printf("Descartando producto '%s' por no pertenecer a la categoría %d (verificación final)",
			product.Name, product.CategoryID)
		return nil // No es un error, simplemente no guardamos el producto
	}

	// --- Lógica de deduplicación --- Primero por imagen, luego por slug

	var currentPHashValue uint64 = 0
	var currentImageHash *goimagehash.ImageHash // Mantener el objeto hash para comparación

	// 1. Descargar imagen y calcular pHash
	if product.ImageURL != "" && !utils.IsPlaceholderImage(product.ImageURL) {
		img, err := utils.DownloadImage(product.ImageURL)
		if err != nil {
			log.Printf("Error al descargar o decodificar imagen para '%s': %v", product.Name, err)
			// Continuamos sin hash de imagen si hay error
		} else {
			hash, err := utils.CalculatePerceptionHash(img)
			if err != nil {
				log.Printf("Error al calcular pHash para '%s': %v", product.Name, err)
				// Continuamos sin hash de imagen si hay error
			} else {
				currentImageHash = hash
				currentPHashValue = hash.GetHash()     // Obtener el valor uint64 para almacenar usando GetHash()
				product.ImageHash = &currentPHashValue // Asignar el hash (puntero a uint64) al modelo del producto
			}
		}
	}

	// 2. Buscar productos existentes con pHash similar en la misma categoría (si se calculó el hash del producto actual)
	var existingProduct *model.Product
	if currentImageHash != nil {
		// Recuperamos productos de la misma categoría para comparar hashes en memoria.
		// Establecemos un límite razonable (ej. 200). Nos aseguramos de cargar el ImageHash.
		existingProducts, err := uc.productRepo.FindByCategory(ctx, product.CategoryID, 200, 0, "") // Ajusta el límite si es necesario
		if err != nil {
			log.Printf("Error al buscar productos existentes por categoría para deduplicación por pHash: %v", err)
			// Continuamos con la búsqueda por slug si falla la búsqueda por categoría
		} else {
			// Comparar pHash con productos existentes
			const PHASH_THRESHOLD = 5 // Umbral de similitud: 5 (ajustar si es necesario)
			for _, p := range existingProducts {
				// Solo comparar si el producto existente tiene un hash válido
				if p.ImageHash != nil {
					// Convertir uint64 almacenado a ImageHash para comparar
					existingImageHash := goimagehash.NewImageHash(*p.ImageHash, goimagehash.PHash) // Cambiado de AHash a PHash

					isSimilar, err := utils.ComparePerceptionHashes(existingImageHash, currentImageHash, PHASH_THRESHOLD)
					if err != nil {
						log.Printf("Error al comparar hashes para productos '%s' vs '%s': %v", product.Name, p.Name, err)
						continue
					}

					if isSimilar {
						// ¡Producto similar encontrado por hash de imagen!
						log.Printf("✅ Producto similar encontrado por pHash: '%s' es similar a '%s'",
							product.Name, p.Name)
						existingProduct = p // Encontramos el producto existente
						break               // Salir del bucle una vez que encontramos una coincidencia
					}
				}
			}
		}
	}

	// 3. Si no se encontró producto similar por hash, buscar por slug (lógica existente)
	if existingProduct == nil {
		// Buscar si ya existe un producto con el mismo slug
		existingProductBySlug, err := uc.productRepo.FindBySlug(ctx, product.Slug)
		// Nota: FindBySlug ya fue añadido a la interfaz y repositorio
		if err != nil && err.Error() != "producto no encontrado" {
			return fmt.Errorf("error al buscar producto existente por slug: %w", err)
		}

		if existingProductBySlug != nil {
			existingProduct = existingProductBySlug // Usar el producto encontrado por slug
			log.Printf("✅ Producto existente encontrado por slug: '%s'", existingProductBySlug.Name)
		}
	}

	// --- Manejar el producto encontrado o crear uno nuevo ---
	var savedProductID uint // Para guardar el ID del producto final (nuevo o existente)

	if existingProduct != nil {
		// Producto existente encontrado (por hash o por slug) - actualizar información y añadir precio si no existe
		log.Printf("Actualizando producto existente '%s' (ID: %d)", existingProduct.Name, existingProduct.ID)

		// Actualizar información del producto existente si es relevante
		// Priorizar la imagen si la nueva no es placeholder y la existente sí lo es o está vacía
		if (existingProduct.ImageURL == "" || utils.IsPlaceholderImage(existingProduct.ImageURL)) &&
			product.ImageURL != "" && !utils.IsPlaceholderImage(product.ImageURL) {
			existingProduct.ImageURL = product.ImageURL
			log.Printf("Actualizada URL de imagen para producto existente '%s' a '%s'", existingProduct.Name, existingProduct.ImageURL)
		}
		// Podríamos añadir lógica para actualizar descripción, etc.

		// Asignar el pHash si se calculó y el producto existente no lo tiene
		// Esto cubre casos donde el producto ya existía por slug pero antes de implementar pHash
		if existingProduct.ImageHash == nil && product.ImageHash != nil {
			existingProduct.ImageHash = product.ImageHash
			log.Printf("Asignado pHash %d al producto existente '%s'", *existingProduct.ImageHash, existingProduct.Name)
		}

		// Guardar las actualizaciones del producto existente
		if err := uc.productRepo.Update(ctx, existingProduct); err != nil {
			log.Printf("Error al actualizar producto existente %d: %v", existingProduct.ID, err)
			// Continuamos guardando el precio aunque la actualización del producto falle
		}

		savedProductID = existingProduct.ID // Usar el ID del producto existente

	} else {
		// 4. Si no se encontró ni por hash similar ni por slug, crear nuevo producto
		log.Printf("✨ Creando nuevo producto: '%s'", product.Name)

		// Generar un slug único si el producto es nuevo
		product.Slug = utils.GenerateUniqueSlug(product.Name, func(slug string) bool {
			exists, err := uc.productRepo.ExistsBySlug(ctx, slug)
			if err != nil {
				log.Printf("Error verificando existencia de slug '%s': %v", slug, err)
				// En caso de error al verificar, asumimos que existe para evitar duplicados en un escenario de error
				return true
			}
			return exists
		})

		// Guardar el nuevo producto (el ImageHash ya debería estar asignado si se calculó)
		if err := uc.productRepo.Create(ctx, product); err != nil {
			return fmt.Errorf("error al crear producto '%s': %w", product.Name, err)
		}

		savedProductID = product.ID // Usar el ID del nuevo producto
		log.Printf("Nuevo producto '%s' creado con ID: %d", product.Name, savedProductID)
	}

	// --- Guardar el precio asociado al producto (existente o nuevo) ---
	// Asumimos que `product` scrapeado siempre tiene al menos un precio
	price := product.Prices[0]
	price.ProductID = savedProductID // Asociar el precio al ID del producto guardado

	// Buscar si ya existe un precio para esta tienda en el producto (esto es posible si se actualiza un producto existente)
	existingPricesForProduct, err := uc.priceRepo.FindByProductID(ctx, savedProductID)
	if err != nil {
		log.Printf("Error al buscar precios existentes para producto %d antes de guardar nuevo precio: %v", savedProductID, err)
		// Continuamos creando un nuevo precio si no podemos buscar existentes
	}

	var priceToUpdate *model.Price
	for _, ep := range existingPricesForProduct {
		if ep.Store == price.Store {
			priceToUpdate = ep
			break
		}
	}

	if priceToUpdate != nil {
		// Actualizar precio existente para esta tienda y producto
		priceToUpdate.Price = price.Price
		priceToUpdate.URL = price.URL
		priceToUpdate.IsAvailable = price.IsAvailable
		priceToUpdate.RetrievedAt = time.Now()

		if err := uc.priceRepo.Update(ctx, priceToUpdate); err != nil {
			return fmt.Errorf("error al actualizar precio para producto %d tienda %s: %w", savedProductID, price.Store, err)
		}
		log.Printf("Actualizado precio para producto ID %d tienda '%s': %.2f %s",
			savedProductID, price.Store, price.Price, price.Currency)
	} else {
		// Crear nuevo precio para esta tienda y producto
		if err := uc.priceRepo.Create(ctx, &price); err != nil {
			return fmt.Errorf("error al crear nuevo precio para producto %d tienda %s: %w", savedProductID, price.Store, err)
		}
		log.Printf("Creado nuevo precio para producto ID %d tienda '%s': %.2f %s",
			savedProductID, price.Store, price.Price, price.Currency)
	}

	return nil
}
