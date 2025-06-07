package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"app/internal/domain/model"
	"app/pkg/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// CoolmodScraper implementa el scraping para Coolmod
type CoolmodScraper struct {
	BaseURL string
}

// NewCoolmodScraper crea una nueva instancia del scraper de Coolmod
func NewCoolmodScraper() *CoolmodScraper {
	return &CoolmodScraper{
		BaseURL: "https://www.coolmod.com",
	}
}

// mapCategoryToURL mapea las categorías de nuestro sistema a URLs de Coolmod
func (s *CoolmodScraper) mapCategoryToURL(slug string) (string, error) {
	switch strings.ToLower(slug) {
	case "portatiles":
		return fmt.Sprintf("%s/#01cc/fullscreen/m=n&q=portatiles", s.BaseURL), nil
	case "ssd":
		return fmt.Sprintf("%s/#01cc/fullscreen/m=and&q=ssd", s.BaseURL), nil
	case "auriculares":
		return fmt.Sprintf("%s/#01cc/fullscreen/m=and&q=auriculares", s.BaseURL), nil
	case "teclados":
		return fmt.Sprintf("%s/#01cc/fullscreen/m=and&q=teclados", s.BaseURL), nil
	case "monitores":
		return fmt.Sprintf("%s/#01cc/fullscreen/m=and&q=monitores", s.BaseURL), nil
	case "tarjetas-graficas":
		return fmt.Sprintf("%s/#01cc/fullscreen/m=and&q=tarjetas+graficas", s.BaseURL), nil
	default:
		return "", fmt.Errorf("categoría no soportada para Coolmod: %s", slug)
	}
}

// ScrapCategory realiza el scraping de productos para una categoría específica
func (s *CoolmodScraper) ScrapCategory(category *model.Category) ([]*model.Product, error) {
	// Mapear categoría a URL de Coolmod
	categoryURL, err := s.mapCategoryToURL(category.Slug)
	if err != nil {
		return nil, err
	}

	log.Printf("Scraping Coolmod - Categoría: %s, URL: %s", category.Name, categoryURL)

	// Configurar el collector de colly
	c := colly.NewCollector(
		colly.UserAgent(utils.GetRandomUserAgent()),
		colly.MaxDepth(1),
	)

	// Permitir cookies y cambiar la política de redirección
	c.AllowURLRevisit = false
	c.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
		return nil // permitir redirecciones
	})

	var products []*model.Product

	// Procesar cada producto encontrado - actualizado con los selectores correctos
	c.OnHTML("article.product-card", func(e *colly.HTMLElement) {
		// Extraer nombre del producto
		name := strings.TrimSpace(e.ChildText("p.card-title a"))
		if name == "" {
			// Intentar con selector alternativo
			name = strings.TrimSpace(e.ChildText(".card-title a"))
			if name == "" {
				return // Si no hay nombre, ignorar
			}
		}

		// Extraer URL del producto
		productURL := e.ChildAttr("p.card-title a", "href")
		if productURL == "" {
			productURL = e.ChildAttr(".card-title a", "href")
		}

		if !strings.HasPrefix(productURL, "http") {
			productURL = s.BaseURL + productURL
		}

		// Extraer URL de la imagen
		imageURL := e.ChildAttr("figure a img", "src")

		// Si la URL de la imagen no comienza con http, añadir el BaseURL
		if imageURL != "" && !strings.HasPrefix(imageURL, "http") {
			imageURL = s.BaseURL + imageURL
		}

		// Extraer precio del producto - parte entera y decimal
		priceInt := strings.TrimSpace(e.ChildText("span.product_price.int_price"))
		priceDec := strings.TrimSpace(e.ChildText("span.dec_price"))

		// Combinar parte entera y decimal
		priceText := priceInt
		if priceDec != "" {
			priceText = priceInt + "," + priceDec // Usar coma como separador decimal (formato europeo)
		}

		// Usar la función ExtractPrice para manejar correctamente el formato
		price, err := utils.ExtractPrice(priceText)
		if err != nil {
			log.Printf("Error al convertir precio para %s: %v", name, err)
			price = 0
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
				Store:       "Coolmod",
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

	log.Printf("Se scrapearon %d productos de Coolmod - %s", len(products), category.Name)
	return products, nil
}

// ScrapProductDetails obtiene los detalles completos de un producto específico
func (s *CoolmodScraper) ScrapProductDetails(productURL string) (*model.Product, error) {
	// Verificar que la URL sea de Coolmod
	if !utils.IsCoolmodURL(productURL) {
		return nil, fmt.Errorf("la URL no pertenece a Coolmod: %s", productURL)
	}

	// Configurar el collector de colly
	c := colly.NewCollector(
		colly.UserAgent(utils.GetRandomUserAgent()),
	)

	var product model.Product
	var price model.Price

	// Extraer nombre del producto
	c.OnHTML("h1.card-title, .product-name", func(e *colly.HTMLElement) {
		product.Name = strings.TrimSpace(e.Text)
		product.Slug = utils.GenerateSlug(product.Name)
	})

	// Extraer descripción del producto
	c.OnHTML(".product-description, .desc-det", func(e *colly.HTMLElement) {
		product.Description = strings.TrimSpace(e.Text)
	})

	// Extraer imagen del producto
	c.OnHTML(".swiper-slide img, figure a img", func(e *colly.HTMLElement) {
		imageURL := e.Attr("src")

		// Si la URL de la imagen no comienza con http, añadir el BaseURL
		if imageURL != "" && !strings.HasPrefix(imageURL, "http") {
			imageURL = s.BaseURL + imageURL
		}

		product.ImageURL = imageURL
	})

	// Extraer precio del producto
	c.OnHTML("span.product_price.int_price", func(e *colly.HTMLElement) {
		priceInt := strings.TrimSpace(e.Text)

		// Buscar la parte decimal
		var priceDec string
		e.DOM.Parent().Parent().Find("span.dec_price").Each(func(i int, s *goquery.Selection) {
			priceDec = strings.TrimSpace(s.Text())
		})

		// Combinar parte entera y decimal
		priceText := priceInt
		if priceDec != "" {
			priceText = priceInt + "," + priceDec // Usar coma como separador decimal (formato europeo)
		}

		// Usar la función ExtractPrice para manejar correctamente el formato
		extractedPrice, err := utils.ExtractPrice(priceText)
		if err != nil {
			log.Printf("Error al convertir precio: %v", err)
		} else {
			price.Price = extractedPrice
		}
	})

	// Extraer disponibilidad
	c.OnHTML(".card-text.text-xs.text-cool-green, .text-delivered", func(e *colly.HTMLElement) {
		stockText := strings.ToLower(strings.TrimSpace(e.Text))
		price.IsAvailable = true // Por defecto asumimos disponible si hay mensaje de entrega

		// Si hay mensajes específicos de agotado, cambiamos a no disponible
		if strings.Contains(stockText, "agotado") || strings.Contains(stockText, "no disponible") {
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
	price.Store = "Coolmod"
	price.Currency = "EUR"
	price.URL = productURL
	price.RetrievedAt = time.Now()

	// Asignar el precio al producto
	product.Prices = []model.Price{price}

	return &product, nil
}
