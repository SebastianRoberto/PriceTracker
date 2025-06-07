package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TemplateRenderer es el componente encargado de renderizar las plantillas HTML
type TemplateRenderer struct {
	builder *TemplateBuilder
}

// NewTemplateRenderer crea una nueva instancia del renderizador de plantillas
func NewTemplateRenderer() (*TemplateRenderer, error) {
	// Construir las plantillas utilizando el nuevo TemplateBuilder
	builder, err := BuildTemplates()
	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{
		builder: builder,
	}, nil
}

// Render renderiza una plantilla específica con los datos proporcionados
func (r *TemplateRenderer) Render(c *gin.Context, status int, name string, data interface{}) {
	r.builder.Render(c, status, name, data)
}

// RenderError renderiza la plantilla de error con un mensaje y código de estado
func (r *TemplateRenderer) RenderError(c *gin.Context, status int, message string) {
	data := gin.H{
		"Title":   "Error",
		"Message": message,
		"Status":  status,
	}
	r.Render(c, status, "error.html", data)
}

// RenderNotFound renderiza la página de error 404
func (r *TemplateRenderer) RenderNotFound(c *gin.Context) {
	r.RenderError(c, http.StatusNotFound, "La página que buscas no existe")
}

// RenderServerError renderiza la página de error 500
func (r *TemplateRenderer) RenderServerError(c *gin.Context, err error) {
	r.RenderError(c, http.StatusInternalServerError, "Error interno del servidor: "+err.Error())
}

// SetupTemplates configura las plantillas HTML para el motor Gin
func SetupTemplates(engine *gin.Engine) (*TemplateRenderer, error) {
	// Construir las plantillas utilizando el nuevo TemplateBuilder
	builder, err := BuildTemplates()
	if err != nil {
		return nil, err
	}

	// No necesitamos configurar el motor de plantillas de Gin
	// porque ahora estamos gestionando las plantillas nosotros mismos

	// Devolver el renderizador
	return &TemplateRenderer{
		builder: builder,
	}, nil
}
