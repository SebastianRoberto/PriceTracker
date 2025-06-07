# 🛠️ Paquetes de Utilidades (pkg/utils)

Este directorio contiene un conjunto de paquetes y funciones de ayuda (`helpers`) que realizan tareas comunes y autocontenidas. Estas utilidades están diseñadas para ser reutilizadas en diferentes partes de la aplicación, promoviendo un código más limpio y evitando la duplicación.

---

## 📋 Utilidades Disponibles

A continuación se detalla el propósito y funcionamiento de cada archivo de utilidad.

### `category_validator.go`

Esta es una de las utilidades más importantes del sistema, responsable de garantizar la **calidad y relevancia de los datos** obtenidos mediante web scraping.

-   **Propósito**: Valida si un producto scrapeado realmente pertenece a la categoría a la que fue asignado inicialmente. Esto es crucial para filtrar productos irrelevantes, como "funda para portátil" en la categoría "Portátiles".
-   **Funcionamiento**: Utiliza un sistema sofisticado basado en palabras clave y puntuaciones:
    -   `globalExcludeKeywords`: Palabras que descartan un producto de **todas** las categorías (ej: "cable", "adaptador").
    -   `categoryExcludeKeywords`: Palabras que descartan un producto de una categoría **específica**.
    -   `categoryRequiredKeywords`: Palabras que **deben** estar presentes para que un producto sea considerado válido para una categoría.
-   **Función Principal**:
    ```go
    func ValidateProductCategory(product *model.Product) bool
    ```

### `extractors.go`

Contiene funciones para extraer datos estructurados a partir de texto plano.

-   **Propósito**: Simplificar la conversión de datos no estructurados a formatos utilizables por la aplicación.
-   **Funciones Principales**:
    -   `ExtractPrice(s string) (float64, error)`: Parsea strings de precios que pueden venir en múltiples formatos (ej: `"1.299,95€"`, `"$349.99"`) y los convierte a un `float64` estándar.
    -   `GetRandomUserAgent() string`: Devuelve una cabecera `User-Agent` de navegador aleatoria de una lista predefinida. Esencial para que los scrapers eviten ser bloqueados.

### `image.go`

Utilidades para el procesamiento y análisis de imágenes, enfocadas en el proceso de scraping.

-   **Propósito**: Gestionar las imágenes de los productos, desde la descarga hasta la detección de duplicados.
-   **Funciones Principales**:
    -   `IsPlaceholderImage(url string) bool`: Detecta si una URL de imagen corresponde a un *placeholder* (imagen de carga, pixel transparente, etc.), especialmente en plataformas como eBay.
    -   `DownloadImage(url string) (image.Image, error)`: Descarga y decodifica una imagen desde una URL.
    -   `CalculatePerceptionHash(...)` y `ComparePerceptionHashes(...)`: Calculan y comparan un hash de percepción (pHash) de las imágenes. Esto permite identificar productos duplicados que usan la misma imagen, incluso si el nombre del producto es ligeramente diferente.

### `slug.go`

Funciones para generar slugs amigables para las URLs a partir de texto.

-   **Propósito**: Crear identificadores únicos y legibles para las URLs de productos y categorías.
-   **Función Principal**:
    -   `GenerateSlug(text string) string`: Convierte un string como `"Tarjeta Gráfica ASUS TUF"` en un slug como `"tarjeta-grafica-asus-tuf"`. El proceso incluye:
        1.  Convertir a minúsculas.
        2.  Eliminar acentos y diacríticos.
        3.  Reemplazar cualquier caracter no alfanumérico por guiones.
        4.  Limitar la longitud y añadir un hash para evitar colisiones.

### `url.go`

Funciones de ayuda muy simples para identificar la tienda de origen a partir de una URL.

-   **Propósito**: Comprobar rápidamente a qué tienda (`Aussar`, `Coolmod`, `eBay`) pertenece una URL de producto.
-   **Funciones Principales**:
    -   `IsAussarURL(url string) bool`
    -   `IsCoolmodURL(url string) bool`
    -   `IsEbayURL(url string) bool` 