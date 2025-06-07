package views

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// TemplateBuilder es el encargado de construir las plantillas HTML
type TemplateBuilder struct {
	baseTemplate   string
	templates      map[string]*template.Template
	templatesFuncs template.FuncMap
}

// NewTemplateBuilder crea una nueva instancia del constructor de plantillas
func NewTemplateBuilder() *TemplateBuilder {
	return &TemplateBuilder{
		baseTemplate:   "layout.html",
		templates:      make(map[string]*template.Template),
		templatesFuncs: make(template.FuncMap),
	}
}

// AddFunc añade una función personalizada a los templates
func (tb *TemplateBuilder) AddFunc(name string, fn interface{}) {
	tb.templatesFuncs[name] = fn
}

// LoadTemplates carga todas las plantillas HTML
func (tb *TemplateBuilder) LoadTemplates() error {
	// Definir funciones personalizadas
	funcMap := template.FuncMap{
		"formatPrice": func(price float64) string {
			return fmt.Sprintf("%.2f €", price)
		},
		// Funciones matemáticas para la watchlist
		"subFloat": func(a, b float64) float64 {
			return a - b
		},
		"divideFloat": func(a, b, mul float64) float64 {
			if b == 0 {
				return 0
			}
			return (a / b) * mul
		},
		"mul": func(a, b float64) float64 {
			return a * b
		},
		// Funciones de utilidad para paginación
		"sub": func(a, b int) int {
			return a - b
		},
		"add": func(a, b int) int {
			return a + b
		},
		"sequence": func(start, end int) []int {
			length := end - start + 1
			if length <= 0 {
				return []int{}
			}
			seq := make([]int, length)
			for i := 0; i < length; i++ {
				seq[i] = start + i
			}
			return seq
		},
		// Función para truncar texto
		"truncate": func(s string, limit int) string {
			if len(s) <= limit {
				return s
			}
			// Truncar al límite y añadir puntos suspensivos
			runes := []rune(s)
			if len(runes) > limit {
				return string(runes[0:limit-3]) + "..."
			}
			return s
		},
		// Puedes añadir más funciones si lo necesitas
	}

	// Añade las funciones al mapa de funciones
	for name, fn := range funcMap {
		tb.AddFunc(name, fn)
	}

	// Ruta base de las plantillas
	templatesDir := "web/templates"

	// Cargar la plantilla base
	baseTemplatePath := filepath.Join(templatesDir, tb.baseTemplate)
	baseContent, err := tb.readTemplateFile(baseTemplatePath)
	if err != nil {
		return fmt.Errorf("error leyendo la plantilla base: %w", err)
	}

	// Listar todas las plantillas que queremos cargar
	templateFiles := []string{
		"home.html",
		"login.html",
		"register.html",
		"register_success.html",
		"verify_success.html",
		"error.html",
		"profile.html",
		"category.html",
		"product_detail.html",
		"watchlist.html",
		"notifications.html",
		"reset_password.html",
		"forgot_password.html",
	}

	// Crear y compilar cada plantilla
	for _, fileName := range templateFiles {
		// Crear nuevo template con funciones
		tmpl := template.New(fileName).Funcs(tb.templatesFuncs)

		// Parsear primero la plantilla base
		tmpl, err = tmpl.Parse(baseContent)
		if err != nil {
			return fmt.Errorf("error al parsear la plantilla base para %s: %w", fileName, err)
		}

		// Leer y parsear la plantilla específica
		templatePath := filepath.Join(templatesDir, fileName)
		content, err := tb.readTemplateFile(templatePath)
		if err != nil {
			return fmt.Errorf("error leyendo la plantilla %s: %w", fileName, err)
		}

		// Parsear el contenido de la plantilla específica
		tmpl, err = tmpl.Parse(content)
		if err != nil {
			return fmt.Errorf("error al parsear la plantilla %s: %w", fileName, err)
		}

		// Guardar la plantilla compilada
		tb.templates[fileName] = tmpl
	}

	return nil
}

// Render renderiza una plantilla específica
func (tb *TemplateBuilder) Render(c *gin.Context, status int, name string, data interface{}) {
	tmpl, exists := tb.templates[name]
	if !exists {
		c.String(http.StatusInternalServerError, "Error: plantilla %s no encontrada", name)
		return
	}

	// Crear el mapa de datos final para pasar a la plantilla
	finalData := make(gin.H)

	// 1. Copiar todos los datos del contexto de Gin (establecidos por middlewares)
	if c.Keys != nil {
		for key, value := range c.Keys {
			finalData[key] = value
		}
	}

	// 2. Fusionar (y sobrescribir si hay colisión) con los datos específicos pasados por el handler
	// Se asume que 'data' es usualmente gin.H o map[string]interface{} para los handlers.
	if handlerData, ok := data.(gin.H); ok {
		for key, value := range handlerData {
			finalData[key] = value
		}
	} else if handlerDataMap, ok := data.(map[string]interface{}); ok {
		for key, value := range handlerDataMap {
			finalData[key] = value
		}
	} else if data != nil {
		// Si 'data' no es un mapa pero no es nil (por ejemplo, un struct),
		// y tenemos claves de contexto, la forma de fusionar depende de cómo
		// la plantilla espera acceder a 'data'. Por ahora, si 'data' no es un mapa,
		// y hay finalData (de c.Keys), 'data' no se fusiona directamente como campos de nivel superior.
		// Si no hay c.Keys y 'data' no es un mapa, se pasa 'data' directamente.
		// Esta lógica asume que los handlers usan gin.H o similar para la compatibilidad de fusión.
		// Si data es un struct y se pasa directamente, no se fusiona con c.Keys a nivel raíz.
		// Dado que los handlers usan gin.H, el bloque anterior es el principal.
		// Si data no es gin.H y no hay c.Keys, pasamos data directamente.
		if len(finalData) == 0 { // Solo si no hay nada de c.Keys
			// Configurar la respuesta
			c.Status(status)
			c.Header("Content-Type", "text/html; charset=utf-8")
			if err := tmpl.Execute(c.Writer, data); err != nil { // Usar 'data' original
				c.String(http.StatusInternalServerError, "Error al renderizar la plantilla: %v", err)
			}
			return
		}
		// Si 'data' no es un mapa y hay 'finalData', 'data' no se fusiona como campos raíz.
		// Podríamos añadirlo bajo una clave específica si fuera necesario, ej: finalData["PageSpecificData"] = data
		// pero por ahora, si no es un mapa, no se fusiona directamente.
	}

	// Configurar la respuesta
	c.Status(status)
	c.Header("Content-Type", "text/html; charset=utf-8")

	// Renderizar la plantilla
	if err := tmpl.Execute(c.Writer, finalData); err != nil {
		c.String(http.StatusInternalServerError, "Error al renderizar la plantilla: %v", err)
	}
}

// readTemplateFile lee el contenido de un archivo de plantilla
func (tb *TemplateBuilder) readTemplateFile(path string) (string, error) {
	// En lugar de usar templates hardcoded, vamos a leer los archivos reales
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("error al leer el template %s: %w", path, err)
	}
	return string(content), nil
}

// BuildTemplates configura todas las plantillas y devuelve un renderizador
func BuildTemplates() (*TemplateBuilder, error) {
	builder := NewTemplateBuilder()
	if err := builder.LoadTemplates(); err != nil {
		return nil, err
	}
	return builder, nil
}
