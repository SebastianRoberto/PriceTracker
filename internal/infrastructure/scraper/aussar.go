package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"app/internal/domain/model"
	"app/pkg/utils"

	"github.com/gocolly/colly/v2"
)

// AussarScraper implementa el scraping para Aussar
type AussarScraper struct {
	BaseURL string
}

// NewAussarScraper crea una nueva instancia del scraper de Aussar
func NewAussarScraper() *AussarScraper {
	return &AussarScraper{
		BaseURL: "https://www.aussar.es",
	}
}

// mapCategoryToURL mapea las categorías de nuestro sistema a URLs de Aussar
func (s *AussarScraper) mapCategoryToURL(slug string) (string, error) {
	switch strings.ToLower(slug) {
	case "portatiles":
		return fmt.Sprintf("%s/equipos/portatiles", s.BaseURL), nil
	case "tarjetas-graficas":
		return fmt.Sprintf("%s/tarjetas-graficas", s.BaseURL), nil
	case "auriculares":
		return fmt.Sprintf("%s/perifericos/auriculares", s.BaseURL), nil
	case "teclados":
		return fmt.Sprintf("%s/perifericos/teclados", s.BaseURL), nil
	case "monitores":
		return fmt.Sprintf("%s/monitores", s.BaseURL), nil
	case "ssd":
		return fmt.Sprintf("%s/almacenamiento/discos-ssd", s.BaseURL), nil
	default:
		return "", fmt.Errorf("categoría no soportada para Aussar: %s", slug)
	}
}

// ScrapCategory realiza el scraping de productos para una categoría específica
func (s *AussarScraper) ScrapCategory(category *model.Category) ([]*model.Product, error) {
	// Mapear categoría a URL de Aussar
	categoryURL, err := s.mapCategoryToURL(category.Slug)
	if err != nil {
		return nil, err
	}

	log.Printf("Scraping Aussar - Categoría: %s, URL: %s", category.Name, categoryURL)

	// Configurar el collector de colly
	c := colly.NewCollector(
		colly.UserAgent(utils.GetRandomUserAgent()),
		colly.MaxDepth(1),
	)

	var products []*model.Product

	// Procesar cada producto encontrado
	c.OnHTML(".product-miniature", func(e *colly.HTMLElement) {
		// Extraer nombre del producto
		name := strings.TrimSpace(e.ChildText(".product-title a"))
		if name == "" {
			return // Si no hay nombre, ignorar
		}

		// Extraer URL del producto
		productURL := e.ChildAttr(".product-title a", "href")
		if !strings.HasPrefix(productURL, "http") {
			productURL = s.BaseURL + productURL
		}

		// Extraer URL de la imagen
		imageURL := e.ChildAttr(".product-thumbnail img", "src")
		if imageURL == "" {
			imageURL = e.ChildAttr(".product-thumbnail img", "data-src")
		}

		// Extraer precio
		priceText := e.ChildText(".product-price-and-shipping .price")
		// Limpiar el precio
		priceText = strings.Replace(priceText, "€", "", -1)
		priceText = strings.Replace(priceText, ".", "", -1)  // Eliminar puntos de miles
		priceText = strings.Replace(priceText, ",", ".", -1) // Convertir coma a punto decimal
		priceText = strings.TrimSpace(priceText)

		// Intentar extraer el precio con una expresión regular para mayor seguridad
		re := regexp.MustCompile(`(\d+[\.,]?\d*)`)
		matches := re.FindStringSubmatch(priceText)

		var price float64
		if len(matches) > 0 {
			price, err = utils.ExtractPrice(matches[1])
			if err != nil {
				log.Printf("Error al convertir precio para %s: %v", name, err)
				price = 0
			}
		}

		// Crear el producto solo si tiene precio válido
		if price > 0 {
			// Generar slug desde el nombre del producto
			slug := utils.GenerateSlug(name)

			// Crear el producto
			product := &model.Product{
				Name:        name,
				Slug:        slug,
				ImageURL:    imageURL,
				CategoryID:  category.ID,
				Description: "", // Se completará en ScrapProductDetails si es necesario
			}

			// Crear el precio
			priceObj := model.Price{
				Store:       "Aussar",
				Price:       price,
				Currency:    "EUR",
				URL:         productURL,
				IsAvailable: true,
				RetrievedAt: time.Now(),
			}

			// Asignar el precio al producto
			product.Prices = []model.Price{priceObj}

			products = append(products, product)
		}
	})

	// Manejar errores
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error al scrapear %s: %v", r.Request.URL, err)
	})

	// Visitar la URL de la categoría
	if err := c.Visit(categoryURL); err != nil {
		return nil, fmt.Errorf("error al visitar %s: %w", categoryURL, err)
	}

	log.Printf("Se scrapearon %d productos de Aussar - %s", len(products), category.Name)
	return products, nil
}

// ScrapProductDetails obtiene los detalles completos de un producto específico
func (s *AussarScraper) ScrapProductDetails(productURL string) (*model.Product, error) {
	// Verificar que la URL sea de Aussar
	if !utils.IsAussarURL(productURL) {
		return nil, fmt.Errorf("la URL no pertenece a Aussar: %s", productURL)
	}

	// Configurar el collector de colly
	c := colly.NewCollector(
		colly.UserAgent(utils.GetRandomUserAgent()),
	)

	var product model.Product
	var price model.Price

	// Extraer nombre del producto
	c.OnHTML("h1.h1", func(e *colly.HTMLElement) {
		product.Name = strings.TrimSpace(e.Text)
		product.Slug = utils.GenerateSlug(product.Name)
	})

	// Extraer descripción del producto
	c.OnHTML(".product-description", func(e *colly.HTMLElement) {
		product.Description = strings.TrimSpace(e.Text)
	})

	// Extraer imagen del producto
	c.OnHTML(".product-cover img", func(e *colly.HTMLElement) {
		imageURL := e.Attr("src")
		if imageURL == "" {
			imageURL = e.Attr("data-src")
		}
		product.ImageURL = imageURL
	})

	// Extraer precio del producto
	c.OnHTML(".current-price .price", func(e *colly.HTMLElement) {
		priceText := strings.TrimSpace(e.Text)
		priceText = strings.Replace(priceText, "€", "", -1)
		priceText = strings.Replace(priceText, ".", "", -1)  // Eliminar puntos de miles
		priceText = strings.Replace(priceText, ",", ".", -1) // Convertir coma a punto decimal

		extractedPrice, err := utils.ExtractPrice(priceText)
		if err != nil {
			log.Printf("Error al convertir precio: %v", err)
		} else {
			price.Price = extractedPrice
		}
	})

	// Extraer disponibilidad
	c.OnHTML(".product-availability", func(e *colly.HTMLElement) {
		availText := strings.ToLower(strings.TrimSpace(e.Text))
		price.IsAvailable = true // Por defecto asumimos disponible

		// Si hay mensajes específicos de agotado, cambiamos a no disponible
		if strings.Contains(availText, "agotado") || strings.Contains(availText, "no disponible") {
			price.IsAvailable = false
		}
	})

	// Manejar errores
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error al scrapear detalles del producto %s: %v", r.Request.URL, err)
	})

	// Visitar la URL del producto
	if err := c.Visit(productURL); err != nil {
		return nil, fmt.Errorf("error al visitar %s: %w", productURL, err)
	}

	// Si no se pudo extraer el nombre, devolver error
	if product.Name == "" {
		return nil, fmt.Errorf("no se pudo extraer el nombre del producto")
	}

	// Completar los datos del precio
	price.Store = "Aussar"
	price.Currency = "EUR"
	price.URL = productURL
	price.RetrievedAt = time.Now()

	// Asignar el precio al producto
	product.Prices = []model.Price{price}

	return &product, nil
}
