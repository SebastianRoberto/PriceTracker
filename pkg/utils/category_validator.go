package utils

import (
	"fmt"
	"regexp"
	"strings"

	"app/internal/domain/model"
)

// Palabras clave generales que excluyen un producto de todas las categorías
var globalExcludeKeywords = []string{
	"funda", "mochila", "bolsa", "estuche", "protector", "carcasa", "soporte",
	"adaptador", "cable", "cargador", "batería", "bateria", "maletín", "maletin",
	"pegatina", "sticker", "skin", "vinilo", "accesorio", "accesorio para", "para laptop",
	"para portátil", "para portatil", "para monitor", "para teclado", "para auriculares",
	"para tarjeta", "limpiador", "almohadilla", "pad", "reposamuñecas", "reposamanos",
	"reposa muñecas", "reposa manos", "elevador", "base para", "base de", "forro",
	"cubierta", "sleeve", "bag", "case", "cover", "protector de", "protección para",
	"proteccion para", "accesorios para", "kit de limpieza", "cleaning kit",
	"extension", "extensión", "extensor", "conversor", "convertidor", "hub usb",
	"usb hub", "splitter", "divisor", "dock", "docking", "estación", "estacion",
	"refrigerador para", "cooler para", "ventilador para", "cooling pad", "cooling stand",
}

// Definimos palabras clave excluyentes por categoría
// Si un producto contiene alguna de estas palabras, NO pertenece a la categoría
var categoryExcludeKeywords = map[uint][]string{
	1: { // Portátiles
		"teclado", "ratón", "mouse", "tarjeta gráfica", "tarjeta grafica", "gpu", "monitor",
		"auricular", "auriculares", "headset", "micrófono", "microfono", "disco duro externo",
		"procesador", "cpu", "tablet", "smartphone", "móvil", "movil", "cámara", "camara",
		"altavoz", "speaker", "impresora", "escáner", "scanner", "router", "switch",
		"keyboard", "headphone", "cascos", "webcam", "pantalla externa", "motherboard",
		"placa base", "fuente alimentación", "mando", "controller", "refrigeración", "ventilador",
		"cooling", "fan", "refrigeracion", "adaptador", "adapter", "hub",
		"docking", "station", "lector", "tarjeta", "card", "cartucho", "tinta", "soporte",
		"nvidia", "radeon", "rtx", "gtx", "rx", "geforce", "fuente", "psu", "kit", "liquid", "cooler",
		"funda portátil", "mochila portátil", "soporte portátil", "cargador portátil",
		"funda portatil", "mochila portatil", "soporte portatil", "cargador portatil",
		"laptop sleeve", "laptop bag", "laptop stand", "laptop cooler", "laptop cooling",
		"protector portátil", "protector portatil", "base refrigeradora", "base refrigerante",
		"alfombrilla", "mouse pad", "mousepad", "almohadilla", "reposamuñecas", "reposa muñecas",
		"coolpc gamer",
	},
	2: { // Tarjetas Gráficas
		"portátil", "portatil", "laptop", "notebook", "teclado", "ratón", "mouse",
		"monitor", "auricular", "headset", "micrófono", "microfono", "disco duro",
		"ssd", "ram", "procesador", "cpu", "tablet", "smartphone", "móvil", "movil",
		"cámara", "camara", "altavoz", "speaker", "impresora", "escáner", "scanner",
		"router", "switch", "cable", "adaptador", "carcasa", "funda",
		"keyboard", "headphone", "auriculares", "cascos", "webcam", "mando", "controller",
		"placa base", "motherboard", "fuente alimentación", "power supply",
		"soporte", "base", "refrigeración líquida", "refrigeracion liquida",
		"refrigerador", "dock", "controlador", "hub", "usb", "patín", "patin", "silla", "chair",
		"gaming", "gamer", "juego", "play", "touchpad", "raton",
		"soporte para tarjeta", "soporte gpu", "gpu support", "graphics card holder",
		"bracket", "adaptador gpu", "gpu adapter", "riser", "extensor pcie", "pcie riser",
		"cable extensión", "cable extension", "cable alargador", "cable extensor",
	},
	3: { // Auriculares
		"portátil", "portatil", "laptop", "notebook", "teclado", "ratón", "mouse",
		"tarjeta gráfica", "tarjeta grafica", "gpu", "monitor", "disco duro", "ssd",
		"ram", "procesador", "cpu", "tablet", "smartphone", "móvil", "movil",
		"cámara", "camara", "altavoz", "speaker", "impresora", "escáner", "scanner",
		"router", "switch", "keyboard", "webcam", "ventilador", "cooling",
		"refrigeración", "placa base", "motherboard", "fuente alimentación",
		"adapter", "cable", "hub", "docking", "station", "mando", "controller",
		"nvidia", "geforce", "rtx", "gtx", "radeon", "rx", "graphics card",
		"soporte auriculares", "headset stand", "headphone stand", "headphone hook",
		"colgador auriculares", "gancho auriculares", "almohadillas", "ear pads",
		"espuma", "foam", "repuesto", "replacement", "cable auriculares", "headphone cable",
		"cable para auriculares", "cable para headset", "adaptador jack", "jack adapter",
	},
	4: { // Teclados
		"portátil", "portatil", "laptop", "notebook",
		"tarjeta gráfica", "tarjeta grafica", "gpu", "monitor", "disco duro",
		"ssd", "ram", "procesador", "cpu", "tablet", "smartphone", "móvil", "movil",
		"cámara", "camara", "impresora", "escáner", "scanner", "router", "switch",
		"pantalla", "display", "webcam", "ventilador", "cooling", "refrigeración",
		"placa base", "motherboard", "fuente alimentación", "power supply",
		"graphic card", "graphic", "memoria", "memory", "card", "tarjeta", "nvidia", "amd",
		"geforce", "radeon", "rtx", "gtx", "rx", "fuente", "fan", "led strip", "tira led",
		"auricular", "auriculares", "headset", "headphone", "cascos", "earphone", "earbud",
		"reposamuñecas", "wrist rest", "reposa muñecas", "keycaps", "teclas", "switches",
		"funda teclado", "keyboard cover", "protector teclado", "keyboard protector",
		"almohadilla teclado", "keyboard pad", "soporte teclado", "keyboard stand",
		"extractor teclas", "keycap puller", "extractor keycaps", "keycap remover",
		"coolpc gamer",
	},
	5: { // Monitores
		"portátil", "portatil", "laptop", "notebook", "tarjeta gráfica", "tarjeta grafica",
		"gpu", "auricular", "headset", "micrófono", "microfono", "disco duro",
		"ssd", "ram", "procesador", "cpu", "tablet", "smartphone", "móvil", "movil",
		"cámara", "camara", "altavoz", "speaker", "impresora", "escáner", "scanner",
		"router", "switch", "headphone", "auriculares", "cascos", "teclado", "keyboard",
		"ventilador", "cooling", "refrigeración", "placa base", "motherboard",
		"fuente alimentación", "power supply", "graphic card", "graphic", "tarjeta",
		"psu", "cooler", "mechanical keyboard", "mechanical gaming keyboard",
		"soporte monitor", "monitor stand", "monitor arm", "brazo monitor", "monitor mount",
		"base monitor", "monitor riser", "elevador monitor", "vesa mount", "soporte vesa",
		"adaptador monitor", "monitor adapter", "protector pantalla", "screen protector",
		"filtro monitor", "monitor filter", "filtro luz azul", "blue light filter",
	},
	6: { // Discos SSD
		"portátil", "portatil", "laptop", "notebook", "teclado", "ratón", "mouse",
		"tarjeta gráfica", "tarjeta grafica", "gpu", "monitor", "auricular", "headset",
		"micrófono", "microfono", "procesador", "cpu", "tablet", "smartphone", "móvil",
		"movil", "cámara", "camara", "altavoz", "speaker", "impresora", "escáner",
		"scanner", "router", "switch", "headphone", "auriculares", "cascos", "webcam",
		"ventilador", "cooling", "refrigeración", "placa base", "motherboard",
		"fuente alimentación", "power supply", "keyboard", "mando", "controller",
		"psu", "headset",
		"carcasa ssd", "ssd enclosure", "adaptador ssd", "ssd adapter", "caddy ssd",
		"conversor ssd", "ssd converter", "soporte ssd", "ssd bracket", "ssd mount",
		"cable sata", "cable nvme", "cable m.2", "extension ssd", "extensión ssd",
		"coolpc gamer",
	},
	// Añade más categorías según sea necesario
}

// También podemos tener palabras clave obligatorias para algunas categorías específicas
var categoryRequiredKeywords = map[uint][]string{
	1: { // Portátiles - al menos una de estas palabras debe estar presente
		"portátil", "portatil", "laptop", "notebook", "gaming laptop", "ordenador portatil",
		"ordenador portátil", "portátil gaming", "portatil gaming", "gaming portatil", "gaming portátil",
		"coolpc laptop", "pc portatil", "ordenador", "computer", "dell", "hp", "lenovo", "asus", "acer", "msi",
		"macbook", "surface", "thinkpad", "ideapad", "pavilion", "inspiron", "latitude", "precision",
	},
	2: { // Tarjetas gráficas
		"tarjeta gráfica", "tarjeta grafica", "gpu", "geforce", "radeon", "rtx", "gtx", "rx",
		"nvidia", "amd", "graphics card", "gráfica", "grafica", "video card", "tarjeta de video",
		"vga", "pcie", "gddr", "ddr", "gddr5", "gddr6", "hbm", "cuda", "ti", "super",
	},
	3: { // Auriculares
		"auricular", "auriculares", "headset", "headphone", "cascos", "earphone", "earbud",
		"gaming headset", "micrófono", "microfono", "surround", "sonido", "sound", "7.1",
		"5.1", "estéreo", "estereo", "stereo", "wireless", "bluetooth", "inalámbrico", "inalambrico",
		"on-ear", "over-ear", "in-ear", "noise cancelling", "cancelación de ruido",
	},
	4: { // Teclados
		"teclado", "keyboard", "gaming keyboard", "mechanical keyboard", "mecánico", "mecanico", "mechanical",
		"switches", "gaming keyboard", "rgb keyboard", "retroiluminado", "backlit", "cherry mx", "membrane",
		"membrana", "qwerty", "macro", "keycaps", "teclas", "keyboard layout", "tkl keyboard",
		"razer keyboard", "corsair keyboard", "logitech keyboard", "hyperx keyboard",
		"60%", "75%", "87%", "104 keys", "108 keys",
	},
	5: { // Monitores
		"monitor", "pantalla", "display", "screen", "lcd", "led", "ips", "gaming monitor",
		"curved", "curvo", "panel", "freesync", "gsync", "g-sync", "hdmi", "displayport",
		"144hz", "165hz", "240hz", "120hz", "ultrawide", "ultraancho", "4k", "2k", "qhd",
		"pulgadas", "inch", "inches", "monitor coolpc", "monitor coolmod",
	},
	6: { // Discos SSD
		"ssd", "disco", "nvme", "m.2", "sata", "almacenamiento", "storage", "solid state",
		"estado sólido", "estado solido", "drive", "pcie", "tlc", "qlc", "mlc",
		"gen3", "gen4", "nand", "flash", "gb", "tb",
	},
	// Añade más categorías según sea necesario
}

// ValidateProductCategory verifica si un producto realmente pertenece a la categoría asignada
func ValidateProductCategory(product *model.Product) bool {
	// Convertir nombre y descripción a minúsculas para comparaciones insensibles a mayúsculas
	name := strings.ToLower(product.Name)
	description := strings.ToLower(product.Description)

	// Suma de puntos para determinar categoría (mayor puntuación = mayor probabilidad)
	matchScore := 0

	// 0. VERIFICACIÓN DE PALABRAS GLOBALES EXCLUYENTES
	// Estas palabras excluyen un producto de cualquier categoría
	for _, keyword := range globalExcludeKeywords {
		if strings.Contains(name, keyword) {
			// Si encontramos una palabra clave excluyente global en el nombre, rechazamos inmediatamente
			return false
		}
	}

	// Verificación especial para teclados mecánicos
	// Esta verificación tiene precedencia sobre las demás
	if strings.Contains(name, "keyboard") || strings.Contains(name, "teclado") {
		// Si contiene palabras clave específicas de teclados mecánicos
		if strings.Contains(name, "mechanical") || strings.Contains(name, "mecánico") ||
			strings.Contains(name, "mecanico") || strings.Contains(name, "gaming keyboard") ||
			product.CategoryID == 4 {
			// Este es claramente un teclado
			product.CategoryID = 4 // Asegurarnos de que esté en la categoría correcta
			return true
		}
	}

	// Caso especial para productos COOLPC/Coolmod
	if strings.Contains(name, "coolpc") || strings.Contains(name, "coolmod") {
		// Verificar primero si es un teclado
		if product.CategoryID == 4 ||
			strings.Contains(name, "keyboard") ||
			strings.Contains(name, "teclado") {
			product.CategoryID = 4
			return true
		}

		// Si menciona laptops o portátiles, va a portátiles
		if product.CategoryID == 1 ||
			strings.Contains(name, "laptop") ||
			strings.Contains(name, "portátil") ||
			strings.Contains(name, "portatil") ||
			strings.Contains(name, "notebook") {
			product.CategoryID = 1
			return true
		}

		// Si menciona monitores o pantallas, va a monitores
		if product.CategoryID == 5 ||
			strings.Contains(name, "monitor") ||
			strings.Contains(name, "pantalla") ||
			strings.Contains(name, "display") ||
			strings.Contains(name, "screen") ||
			strings.Contains(name, "pulgadas") ||
			strings.Contains(name, "inch") {
			product.CategoryID = 5
			return true
		}
	}

	// 1. VERIFICACIÓN DE PALABRAS EXCLUYENTES (menos estricta)
	if excludeKeywords, exists := categoryExcludeKeywords[product.CategoryID]; exists {
		exclusionCount := 0.0
		for _, keyword := range excludeKeywords {
			// Verificar en nombre con delimitadores de palabra
			regex := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(keyword))
			if matched, _ := regexp.MatchString(regex, name); matched {
				exclusionCount++
				// Si hay más de 2 palabras excluyentes, rechazamos
				if exclusionCount >= 2 {
					return false
				}
			}

			// Verificar en descripción si existe (menos peso)
			if description != "" {
				if matched, _ := regexp.MatchString(regex, description); matched {
					exclusionCount += 0.5
					if exclusionCount >= 3 {
						return false
					}
				}
			}
		}
	}

	// 2. VERIFICACIÓN DE PALABRAS REQUERIDAS
	if requiredKeywords, exists := categoryRequiredKeywords[product.CategoryID]; exists {
		// Contador de palabras requeridas encontradas
		requiredWordsFound := 0

		for _, keyword := range requiredKeywords {
			regex := fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(keyword))

			// Buscar en nombre (mayor prioridad)
			if matched, _ := regexp.MatchString(regex, name); matched {
				requiredWordsFound++
				matchScore += 3 // Dar más peso a coincidencias en el nombre
				continue
			}

			// Buscar en descripción si existe
			if description != "" && len(description) > 5 {
				if matched, _ := regexp.MatchString(regex, description); matched {
					requiredWordsFound++
					matchScore += 1 // Menor peso a coincidencias en descripción
				}
			}
		}

		// Si no encontramos ninguna palabra requerida, rechazar
		if requiredWordsFound == 0 {
			return false
		}

		// Si encontramos al menos una palabra requerida, es probable que sea correcta
		if requiredWordsFound >= 1 {
			return true
		}
	}

	// 3. VERIFICACIÓN POR CATEGORÍA ESPECÍFICA
	switch product.CategoryID {
	case 1: // Portátiles
		// Excluir accesorios para portátiles que pueden tener la palabra "portátil" pero no ser un portátil
		for _, exclude := range []string{"funda portátil", "mochila portátil", "soporte portátil", "cargador portátil",
			"funda portatil", "mochila portatil", "soporte portatil", "cargador portatil",
			"laptop sleeve", "laptop bag", "laptop stand", "laptop cooler"} {
			if strings.Contains(name, exclude) {
				return false
			}
		}

		// Patrones para identificar portátiles
		portátilPatterns := []string{
			`\blaptop\b`, `\bnotebook\b`, `\bportatil\b`, `\bportátil\b`,
			`\bordenador\s+port[aá]til\b`, `\bgaming\s+laptop\b`, `\bpc\b`,
			`\bdell\b`, `\bhp\b`, `\blenovo\b`, `\basus\b`, `\bacer\b`, `\bmsi\b`,
			`\bideapad\b`, `\bthinkpad\b`, `\binspirón\b`, `\blatitude\b`,
		}

		for _, pattern := range portátilPatterns {
			if matched, _ := regexp.MatchString(pattern, name); matched {
				matchScore += 2
			}
		}

		// Palabras clave adicionales que aumentan la probabilidad
		additionalKeywords := []string{"ram", "ssd", "intel", "amd", "ryzen", "core i"}
		for _, keyword := range additionalKeywords {
			if strings.Contains(name, keyword) {
				matchScore++
			}
		}

		return matchScore >= 2

	case 2: // Tarjetas Gráficas
		// Las tarjetas gráficas generalmente tienen números de modelo específicos
		gpuPatterns := []string{
			`\brtx\s*\d{4}\b`,
			`\bgtx\s*\d{3,4}\b`,
			`\brx\s*\d{4}\b`,
			`\bradeon\b`,
			`\bgeforce\b`,
			`\bnvidia\b`,
			`\bgráfica\b`,
			`\bgrafica\b`,
			`\bgpu\b`,
			`\btarjeta\b`,
		}

		for _, pattern := range gpuPatterns {
			if matched, _ := regexp.MatchString(pattern, name); matched {
				matchScore += 3
			}
		}

		// Si tiene alto puntaje, es muy probable que sea una tarjeta gráfica
		return matchScore >= 3

	case 3: // Auriculares
		// Patrones específicos para auriculares
		auricularePatterns := []string{
			`\bauricular\b`,
			`\bauriculares\b`,
			`\bcascos\b`,
			`\bheadset\b`,
			`\bheadphone\b`,
			`\bearphone\b`,
			`\bmic\b`,
			`\bmicrófono\b`,
			`\bmicrofono\b`,
			`\bstereo\b`,
			`\bestéreo\b`,
			`\bsonido\b`,
			`\bsound\b`,
			`\bwireless\b`,
			`\binhalámbrico\b`,
			`\binalambrico\b`,
		}

		for _, pattern := range auricularePatterns {
			if matched, _ := regexp.MatchString(pattern, name); matched {
				matchScore += 2
			}
		}

		return matchScore >= 2

	case 4: // Teclados
		// Patrones específicos para teclados
		tecladoPatterns := []string{
			`\bteclado\b`,
			`\bkeyboard\b`,
			`\bmecánico\b`,
			`\bmecanico\b`,
			`\bmechanical\b`,
			`\bswitches\b`,
			`\brgb\b`,
			`\bteclas\b`,
			`\bkeycaps\b`,
			`\b60%\b`,
			`\b75%\b`,
			`\b87%\b`,
			`\b104\s*keys\b`,
			`\b108\s*keys\b`,
			`\btkl\b`,
		}

		// Dar mucha prioridad a patrones que indican claramente un teclado
		for _, pattern := range tecladoPatterns {
			if matched, _ := regexp.MatchString(pattern, name); matched {
				matchScore += 3
			}
		}

		// Verificar marcas conocidas de teclados
		keyboardBrands := []string{"razer", "corsair", "logitech", "hyperx", "steelseries", "ducky"}
		for _, brand := range keyboardBrands {
			if strings.Contains(name, brand) &&
				(strings.Contains(name, "keyboard") || strings.Contains(name, "teclado")) {
				matchScore += 4
			}
		}

		return matchScore >= 3

	case 5: // Monitores
		// Patrones específicos para monitores
		monitorPatterns := []string{
			`\bmonitor\b`,
			`\bdisplay\b`,
			`\bpantalla\b`,
			`\bpanel\b`,
			`\bultrawide\b`,
			`\bcurved\b`,
			`\b\d{2,3}hz\b`, // Patrones de refresco (144hz, etc)
			`\bips\b`,
			`\b4k\b`,
			`\b2k\b`,
			`\bqhd\b`,
			`\bhd\b`,
			`\bfhd\b`,
			`\d+[\.,]?\d*\s*["'´]`, // Patrones para pulgadas (24", 27", etc)
			`\d+[\.,]?\d*\s*pulgadas\b`,
			`\binch\b`,
		}

		for _, pattern := range monitorPatterns {
			if matched, _ := regexp.MatchString(pattern, name); matched {
				matchScore += 2
			}
		}

		// Descartar si tiene palabras clave de teclado
		if strings.Contains(name, "keyboard") || strings.Contains(name, "teclado") {
			if strings.Contains(name, "mechanical") || strings.Contains(name, "mecánico") ||
				strings.Contains(name, "mecanico") || strings.Contains(name, "gaming keyboard") {
				return false
			}
		}

		return matchScore >= 2

	case 6: // Discos SSD
		// Patrones específicos para SSD
		ssdPatterns := []string{
			`\bssd\b`,
			`\bnvme\b`,
			`\bm\.2\b`,
			`\bestado\s+s[oó]lido\b`,
			`\bsolid\s+state\b`,
			`\b\d+\s*gb\b`,
			`\b\d+\s*tb\b`,
			`\bpcie\b`,
			`\bgen\d\b`,
			`\bsata\b`,
			`\bdisco\b`,
			`\bdrive\b`,
			`\balmacenamiento\b`,
			`\bstorage\b`,
		}

		for _, pattern := range ssdPatterns {
			if matched, _ := regexp.MatchString(pattern, name); matched {
				matchScore += 2
			}
		}

		return matchScore >= 2
	}

	// Si pasa todas las validaciones, el producto pertenece a la categoría
	// Pero requerimos un puntaje mínimo para mayor seguridad
	return matchScore >= 1
}
