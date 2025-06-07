# 🎨 Capa de Vistas (`/views`)

Este directorio es el responsable final de construir la interfaz de usuario HTML que se envía al navegador. Su función principal es tomar los datos proporcionados por los `handlers` y los `middlewares`, y fusionarlos con las plantillas HTML para generar la página final.

Esta capa desacopla por completo la lógica de negocio de la presentación, permitiendo que el diseño visual y la estructura del frontend puedan modificarse sin afectar al backend.

---

## ⚙️ Flujo de Renderizado

1.  Un `handler` recibe una petición y decide qué página mostrar.
2.  Llama a un método del `TemplateRenderer` (ej: `renderer.Render(c, "home.html", data)`).
3.  El `TemplateRenderer` utiliza el `TemplateBuilder` para realizar el trabajo pesado.
4.  El `TemplateBuilder` toma la plantilla solicitada (ej: `home.html`), que ya ha sido previamente parseada junto con la plantilla base (`layout.html`).
5.  Fusiona los datos globales (del contexto de Gin) con los datos específicos de la página (del handler).
6.  Ejecuta la plantilla final con los datos combinados, generando el HTML que se envía al usuario.

---

## 🔩 Componentes Clave

### `template_builder.go`
Es el motor principal del sistema de plantillas.
- **Carga y Parsing**: Descubre y carga todas las plantillas `.html` del directorio `web/templates` durante el arranque de la aplicación.
- **Plantilla Base**: Utiliza un sistema de herencia donde una plantilla base (`layout.html`) define la estructura común (header, footer, menús), y las plantillas específicas (`home.html`, `profile.html`, etc.) "rellenan" el contenido principal.
- **Fusión de Datos**: Su lógica de `Render` es crucial. Combinamos los datos que los `middlewares` inyectan en todas las peticiones (como la información del usuario o las categorías) con los datos que el `handler` pasa para esa página en concreto. Esto asegura que datos globales estén siempre disponibles sin tener que pasarlos manualmente en cada `handler`.
- **Funciones Personalizadas**: Inyecta una serie de funciones de ayuda (`FuncMap`) que pueden ser utilizadas directamente dentro de las plantillas para formatear datos.

| Función | Descripción | Ejemplo de Uso en Plantilla |
| :--- | :--- | :--- |
| `formatPrice` | Formatea un `float64` como una cadena de texto con el símbolo del euro. | `{{ .BestPrice | formatPrice }}` |
| `truncate` | Acorta una cadena de texto a una longitud máxima, añadiendo "..." al final. | `{{ .Product.Name | truncate 50 }}` |
| `sub`, `add` | Realizan operaciones aritméticas básicas (resta y suma). | `{{ sub .Page 1 }}` |
| `sequence` | Genera una secuencia de números, útil para bucles de paginación. | `{{ range sequence 1 .PageCount }}` |
| `subFloat`, `divideFloat`| Realizan operaciones con decimales, usadas en la `watchlist`.| `{{ divideFloat .CurrentPrice .TargetPrice 100 }}`|

### `view_models.go`
Define las estructuras (`structs`) de datos que se pasan a las plantillas. Son **DTOs (Data Transfer Objects)** específicos para la vista. Su propósito es adaptar y combinar la información de los modelos de dominio (`/domain/model`) a lo que la vista necesita exactamente.

- **Desacoplamiento**: Evitan que las plantillas dependan directamente de los modelos de la base de datos.
- **Agregación**: Un `ViewModel` puede combinar datos de múltiples fuentes. Por ejemplo, `HomePageViewModel` contiene al usuario, una lista de productos destacados y una lista de categorías, todo lo necesario para renderizar la página de inicio.

### `mapper.go`
Contiene las funciones "traductoras" que convierten los modelos de dominio en los `ViewModels` definidos arriba.

- **`ToUserViewModel(*model.User)`**: Convierte un usuario de dominio a su versión para la vista.
- **`ToProductViewModel(*model.Product, ...)`**: Convierte un producto de dominio, pero además recibe y añade información extra como su mejor precio actual.
- **`BuildHomePageViewModel(...)`**: Es un constructor de alto nivel que orquesta la creación del `ViewModel` completo para la página principal.

### `renderer.go`
Actúa como una fachada o un "wrapper" simplificado para el `TemplateBuilder`. Los `handlers` interactúan con este componente en lugar de hacerlo directamente con el `builder`, lo que simplifica su código.

- **API Limpia**: Ofrece métodos sencillos como `Render`, `RenderError`, `RenderNotFound` y `RenderServerError`.
- **Estandarización**: Asegura que todas las páginas de error se rendericen de la misma manera, pasando los datos correctos a la plantilla `error.html`. 