package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"

	"github.com/corona10/goimagehash"
)

// IsPlaceholderImage determina si la URL de la imagen hace referencia a un placeholder
// devuelto por eBay (gif de loading, dominio ir.ebaystatic.com, data URI, etc.).
// Esto nos permite descartar imágenes no válidas y forzar la búsqueda de la imagen
// principal real en la página de detalle del producto.
func IsPlaceholderImage(url string) bool {
	if url == "" {
		return true
	}

	l := strings.ToLower(url)

	// eBay usa gifs de loading y el dominio ir.ebaystatic.com para placeholders
	if strings.HasSuffix(l, ".gif") || strings.Contains(l, "ir.ebaystatic.com") {
		return true
	}

	// data URIs (base64) también son placeholders
	if strings.HasPrefix(l, "data:") {
		return true
	}

	// Placeholders y pixeles transparentes comunes
	placeholderPatterns := []string{
		"placeholder", "transparent", "blank", "no-image",
		"noimage", "pixel.gif", "1x1", "spacer",
		"s-l1-", "s-l5-", "s-l10-", // Imágenes demasiado pequeñas
	}

	for _, pattern := range placeholderPatterns {
		if strings.Contains(l, pattern) {
			return true
		}
	}

	// Para imágenes de eBay, detectar las imágenes muy pequeñas
	// Las miniaturas reales suelen ser al menos s-l64
	smallSizePatterns := []string{
		"s-l16", "s-l24", "s-l32",
	}

	if strings.Contains(l, "ebayimg.com") {
		for _, pattern := range smallSizePatterns {
			if strings.Contains(l, pattern) {
				return true
			}
		}
	}

	return false
}

// DownloadImage descarga una imagen desde una URL
func DownloadImage(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al descargar imagen: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error al descargar imagen: estado %d", resp.StatusCode)
	}

	// Leer el cuerpo de la respuesta
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer cuerpo de la imagen: %w", err)
	}

	// Detectar el formato de la imagen y decodificar
	// Intentar decodificar como JPEG, luego como PNG
	img, err := jpeg.Decode(bytes.NewReader(bodyBytes))
	if err == nil {
		return img, nil
	}

	img, err = png.Decode(bytes.NewReader(bodyBytes))
	if err == nil {
		return img, nil
	}

	return nil, fmt.Errorf("formato de imagen no soportado")
}

// CalculatePerceptionHash calcula el hash de percepción (pHash) de una imagen
func CalculatePerceptionHash(img image.Image) (*goimagehash.ImageHash, error) {
	// goimagehash.PerceptionHash espera una imagen en escala de grises
	// Sin embargo, la librería parece manejar la conversión interna si se le pasa una imagen de color.
	// Si esto causa problemas, se podría añadir una conversión explícita.

	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		return nil, fmt.Errorf("error al calcular pHash: %w", err)
	}

	return hash, nil
}

// ComparePerceptionHashes compara dos hashes de percepción y devuelve true si son similares
// Basado en un umbral de distancia Hamming. Un umbral bajo (ej. 0-5) indica alta similitud.
func ComparePerceptionHashes(hash1, hash2 *goimagehash.ImageHash, threshold int) (bool, error) {
	if hash1 == nil || hash2 == nil {
		return false, fmt.Errorf("uno o ambos hashes son nulos")
	}

	distance, err := hash1.Distance(hash2)
	if err != nil {
		return false, fmt.Errorf("error al calcular distancia entre hashes: %w", err)
	}

	return distance <= threshold, nil
}
