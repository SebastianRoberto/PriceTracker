package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"app/internal/domain/model"
	"app/pkg/utils"

	"github.com/gocolly/colly/v2"
)

// EbayScraper implementa el scraping para eBay
type EbayScraper struct {
	BaseURL string
}

// NewEbayScraper crea una nueva instancia del scraper de eBay
func NewEbayScraper() *EbayScraper {
	return &EbayScraper{
		BaseURL: "https://www.ebay.com",
	}
}

// ScrapCategory obtiene productos de una categoría específica
func (s *EbayScraper) ScrapCategory(category *model.Category) ([]*model.Product, error) {
	var products []*model.Product

	// Definir URL según la categoría
	var searchTerm string
	switch strings.ToLower(category.Slug) {
	case "portatiles":
		searchTerm = "laptop+computers+notebooks"
	case "tarjetas-graficas":
		searchTerm = "graphics+card+gpu+nvidia+amd"
	case "auriculares":
		searchTerm = "gaming+headphones+headset"
	case "teclados":
		searchTerm = "gaming+keyboard+mechanical"
	case "monitores":
		searchTerm = "computer+monitor+gaming"
	case "ssd":
		searchTerm = "ssd+solid+state+drive"
	default:
		return nil, fmt.Errorf("categoría no soportada: %s", category.Slug)
	}

	// Construir URL de búsqueda
	searchURL := fmt.Sprintf("%s/sch/i.html?_nkw=%s&_sacat=0", s.BaseURL, searchTerm)

	// Configurar el collector de colly
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		colly.MaxDepth(1),
	)

	// Primero capturamos las imágenes precargadas (nuevo formato de eBay)
	prefetchImages := make(map[string]string)
	c.OnHTML("div.s-prefetch-image img", func(e *colly.HTMLElement) {
		imgSrc := e.Attr("src")
		if imgSrc != "" && !utils.IsPlaceholderImage(imgSrc) {
			// Extraer el ID de imagen de la URL (por ejemplo, g/kpQAAOSwnCdmMmCt)
			imgID := extractImageIDFromURL(imgSrc)
			if imgID != "" {
				prefetchImages[imgID] = imgSrc
			}
		}
	})

	// Registrar callbacks para los elementos HTML
	c.OnHTML("li.s-item", func(e *colly.HTMLElement) {
		// Ignorar elementos promocionales
		if e.ChildText(".s-item__title--tagblock") != "" {
			return
		}

		name := e.ChildText(".s-item__title")
		// Si el nombre está vacío o es un banner, saltar
		if name == "" || name == "Shop on eBay" {
			return
		}

		url := e.ChildAttr("a.s-item__link", "href")

		// ------------------------------
		// 1. Extraer thumbnail del listado
		// ------------------------------
		imageURL := strings.TrimSpace(e.ChildAttr("img.s-item__image-img", "src"))

		// Detectar placeholders
		if utils.IsPlaceholderImage(imageURL) {
			// Intentar data-src o data-lazyimg
			candidates := []string{"data-src", "data-lazyimg", "data-srcset"}
			for _, attr := range candidates {
				candidate := strings.TrimSpace(e.ChildAttr("img.s-item__image-img", attr))
				if !utils.IsPlaceholderImage(candidate) {
					imageURL = candidate
					break
				}
			}
		}

		// Intentar extraer la imagen de cualquier elemento imagen hijo
		if utils.IsPlaceholderImage(imageURL) {
			e.ForEach("img", func(_ int, img *colly.HTMLElement) {
				if imgSrc := img.Attr("src"); imgSrc != "" && !utils.IsPlaceholderImage(imgSrc) {
					imageURL = imgSrc
				}
			})
		}

		// ------------------------------
		// 2. Intentar usar imágenes del prefetch
		// ------------------------------
		if utils.IsPlaceholderImage(imageURL) {
			// Buscar el ID de imagen en la miniatura actual para intentar hacer match con prefetch
			imgIDFromPlaceholder := extractImageIDFromURL(imageURL)
			if imgIDFromPlaceholder != "" && prefetchImages[imgIDFromPlaceholder] != "" {
				imageURL = prefetchImages[imgIDFromPlaceholder]
			} else {
				// Si no podemos hacer match exacto, intentar con cualquier imagen del prefetch
				// que contenga parte del ID (menos preciso pero mejor que nada)
				for id, img := range prefetchImages {
					if strings.Contains(imageURL, id) || strings.Contains(id, imgIDFromPlaceholder) {
						imageURL = img
						break
					}
				}
			}
		}

		// ------------------------------
		// 3. Fallback: visitar página de detalle
		// ------------------------------
		if utils.IsPlaceholderImage(imageURL) {
			detailCollector := c.Clone()
			detailImageURL := ""

			// meta og:image suele tener la imagen principal de alta resolución
			detailCollector.OnHTML("meta[property='og:image']", func(he *colly.HTMLElement) {
				if v := he.Attr("content"); v != "" {
					detailImageURL = v
				}
			})

			// Backup: img#icImg tradicional
			detailCollector.OnHTML("img#icImg", func(he *colly.HTMLElement) {
				if v := he.Attr("src"); v != "" {
					detailImageURL = v
				}
			})

			// Visitar la URL de detalle (síncrono porque collector no es Async)
			_ = detailCollector.Visit(url)

			if !utils.IsPlaceholderImage(detailImageURL) {
				imageURL = detailImageURL
			}
		}

		// ------------------------------
		// 4. Normalizar imagen a alta resolución (s-l500)
		// ------------------------------
		if !utils.IsPlaceholderImage(imageURL) {
			replacements := []string{"s-l64", "s-l75", "s-l96", "s-l140", "s-l160", "s-l180", "s-l200", "s-l225", "s-l300", "s-l400"}
			for _, r := range replacements {
				imageURL = strings.Replace(imageURL, r, "s-l500", 1)
			}
		}

		// Extraer y limpiar el precio
		priceText := e.ChildText(".s-item__price")

		// Detectar y manejar rangos de precios
		if strings.Contains(priceText, " to ") {
			// Si hay un rango, tomamos solo el primer precio
			priceParts := strings.Split(priceText, " to ")
			priceText = priceParts[0]
		} else if strings.Contains(priceText, "$") && strings.Count(priceText, "$") > 1 {
			// Si hay múltiples símbolos de dólar, es probablemente un rango sin "to"
			// Ejemplo: "$150.00$210.00", tomamos solo hasta el segundo "$"
			parts := strings.SplitN(priceText, "$", 3) // "$150.00$210.00" -> ["", "150.00", "210.00"]
			if len(parts) >= 2 {
				priceText = "$" + parts[1] // Usamos solo el primer precio
			}
		}

		// Eliminar símbolo de moneda y espacios
		priceText = strings.ReplaceAll(priceText, "$", "")
		priceText = strings.ReplaceAll(priceText, "US ", "") // Eliminar prefijo "US "
		priceText = strings.ReplaceAll(priceText, ",", "")   // Eliminar comas
		priceText = strings.TrimSpace(priceText)

		// Verificar que el texto resultante sea un número válido
		priceRegex := regexp.MustCompile(`^[\d.]+$`)
		if !priceRegex.MatchString(priceText) {
			log.Printf("Texto de precio no válido después de limpieza: '%s'", priceText)
			return
		}

		var price float64
		if priceText != "" {
			var err error
			price, err = strconv.ParseFloat(priceText, 64)
			if err != nil {
				log.Printf("Error al convertir precio '%s': %v\n", priceText, err)
				return
			}
		}

		// Ignorar productos sin precio válido
		if price <= 0 {
			return
		}

		product := &model.Product{
			Name:        name,
			Slug:        utils.GenerateSlug(name), // Generar slug a partir del nombre
			ImageURL:    imageURL,
			CategoryID:  category.ID,
			Description: "", //
		}

		priceModel := &model.Price{
			Store:       "eBay",
			Price:       price,
			Currency:    "USD",
			URL:         url,
			IsAvailable: true,
			RetrievedAt: time.Now(),
		}

		product.Prices = []model.Price{*priceModel}
		products = append(products, product)
	})

	// Manejar errores
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error al scrapear %s: %v\n", r.Request.URL, err)
	})

	// Visitar la página de búsqueda
	if err := c.Visit(searchURL); err != nil {
		return nil, fmt.Errorf("error al visitar %s: %w", searchURL, err)
	}

	// Estadísticas de imágenes para depuración
	log.Printf("[EBAY] Imágenes precargadas encontradas: %d", len(prefetchImages))

	// Contar productos con y sin imágenes
	withImage := 0
	withoutImage := 0
	for _, p := range products {
		if !utils.IsPlaceholderImage(p.ImageURL) {
			withImage++
		} else {
			withoutImage++
			log.Printf("[EBAY] Producto sin imagen: %s", p.Name)
		}
	}

	log.Printf("[EBAY] Productos con imagen: %d, sin imagen: %d (total: %d) - %s",
		withImage, withoutImage, len(products), category.Name)

	return products, nil
}

// extractImageIDFromURL extrae el ID de imagen de una URL de eBay
// Ejemplo: https://i.ebayimg.com/images/g/kpQAAOSwnCdmMmCt/s-l500.webp -> g/kpQAAOSwnCdmMmCt
func extractImageIDFromURL(url string) string {
	re := regexp.MustCompile(`images/([^/]+/[^/]+)/`)
	matches := re.FindStringSubmatch(url)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}
