package utils

import (
	"strings"
)

// IsEbayURL verifica si una URL pertenece al dominio de eBay
func IsEbayURL(url string) bool {
	return strings.Contains(url, "ebay.com")
}

// IsCoolmodURL verifica si una URL pertenece al dominio de Coolmod
func IsCoolmodURL(url string) bool {
	return strings.Contains(url, "coolmod.com")
}

// IsAussarURL verifica si una URL pertenece al dominio de Aussar
func IsAussarURL(url string) bool {
	return strings.Contains(url, "aussar.es")
}

// IsMercadoLibreURL verifica si una URL pertenece al dominio de MercadoLibre
// MercadoLibre no esta implementado en el scraper por cierto, era muy complicado extraer los datos de aqui
func IsMercadoLibreURL(url string) bool {
	return strings.Contains(url, "mercadolibre.com") ||
		strings.Contains(url, "mercadolivre.com") ||
		strings.Contains(url, "mercadolibre.com.mx")
}

// Carrefour tampoco esta implementado
func IsCarrefourURL(url string) bool {
	return strings.Contains(url, "carrefour.es")
}
