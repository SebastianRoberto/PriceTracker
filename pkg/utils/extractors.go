package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ExtractPrice extrae un precio de un texto
// Por ejemplo: "€29,99" -> 29.99
func ExtractPrice(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("texto vacío")
	}

	// Eliminar símbolos de moneda y espacios
	s = strings.ReplaceAll(s, "€", "")
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, " ", "")

	// Si contiene tanto puntos como comas, asumimos formato europeo: "1.349,95"
	if strings.Contains(s, ".") && strings.Contains(s, ",") {
		s = strings.ReplaceAll(s, ".", "")  // eliminar puntos de miles
		s = strings.ReplaceAll(s, ",", ".") // convertir coma decimal a punto
	} else if strings.Contains(s, ",") {
		// Solo coma: asume que es decimal
		s = strings.ReplaceAll(s, ",", ".")
	} // si solo contiene puntos, asumimos que son decimales válidos ya

	// Extraer número flotante
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	match := re.FindString(s)
	if match == "" {
		return 0, fmt.Errorf("no se encontró un precio válido en: %s", s)
	}

	price, err := strconv.ParseFloat(match, 64)
	if err != nil {
		return 0, fmt.Errorf("error al convertir '%s' a float: %v", match, err)
	}

	return price, nil
}

// WriteDebugFile escribe datos binarios a un archivo para depuración
func WriteDebugFile(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

// GetRandomUserAgent devuelve un User-Agent aleatorio
func GetRandomUserAgent() string {
	agents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
	}
	return agents[int(time.Now().UnixNano())%len(agents)]
}
