package utils

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	regExpNonAlphanumeric = regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	regExpMultipleSpaces  = regexp.MustCompile(`\s+`)
)

// GenerateSlug crea un slug a partir de un texto
// Convierte "Tarjeta Gráfica ASUS TUF" en "tarjeta-grafica-asus-tuf"
func GenerateSlug(text string) string {
	// Pasar a minúsculas
	slug := strings.ToLower(text)

	// Remover acentos
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	slug, _, _ = transform.String(t, slug)

	// Reemplazar caracteres no alfanuméricos con espacios
	slug = regExpNonAlphanumeric.ReplaceAllString(slug, " ")

	// Reemplazar múltiples espacios con uno solo
	slug = regExpMultipleSpaces.ReplaceAllString(slug, " ")

	// Reemplazar espacios con guiones
	slug = strings.ReplaceAll(slug, " ", "-")

	// Recortar espacios al inicio y final
	slug = strings.Trim(slug, "-")

	// Limitar longitud (la columna slug tiene size:100)
	if len(slug) > 95 { // Reducir a 95 para dejar espacio para sufijo
		// Conservar los primeros 75 caracteres y añadir un hash único basado en el texto completo
		// para minimizar colisiones al truncar
		hashValue := fmt.Sprintf("%x", sha256.Sum256([]byte(text)))[:6] // usar 6 caracteres del hash
		slug = slug[:95-len(hashValue)-1] + "-" + hashValue
	}

	return slug
}

// GenerateUniqueSlug genera un slug único asegurándose de que no exista en la base de datos
// Si el slug ya existe, se le añade un sufijo numérico
func GenerateUniqueSlug(text string, checkExists func(string) bool) string {
	baseSlug := GenerateSlug(text)
	slug := baseSlug
	counter := 1

	// Si el slug ya existe, añadir sufijo numérico
	for checkExists(slug) {
		// Asegurar que la longitud total con el sufijo no exceda 100 caracteres
		suffix := "-" + string(rune(counter+'0'))
		if len(baseSlug)+len(suffix) > 100 {
			baseSlug = baseSlug[:100-len(suffix)]
		}
		slug = baseSlug + suffix
		counter++
	}

	return slug
}
