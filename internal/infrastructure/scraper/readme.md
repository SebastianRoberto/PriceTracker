# ⚙️ Sistema de Web Scraping

Este directorio contiene las implementaciones concretas de los `scrapers` para las diferentes tiendas online. Su responsabilidad es navegar por los sitios web externos, extraer información de los productos y devolverla en un formato estructurado que el resto de la aplicación pueda procesar.

Esta capa es una implementación de la lógica definida en el caso de uso de scraping y es invocada periódicamente por el planificador de tareas.

---

## 🏗️ Estructura y Scrapers Implementados

El sistema está diseñado para ser modular, donde cada scraper es un componente independiente enfocado en una única tienda.

| Archivo            | Tienda   | Descripción                                                                                                                              |
| :----------------- | :------- | :--------------------------------------------------------------------------------------------------------------------------------------- |
| **`aussar.go`**    | Aussar   | Implementa el scraping para Aussar.es. Mapea categorías internas a URLs de la tienda y extrae la información básica de los listados.          |
| **`coolmod.go`**   | Coolmod  | Implementa el scraping para Coolmod.com. Maneja la estructura específica de su catálogo y la forma en que presentan los precios.            |
| **`ebay.go`**      | eBay     | Implementa el scraping para eBay.com. Incluye lógica avanzada para manejar la variabilidad de los listados y extraer imágenes de alta calidad, evitando los *placeholders* comunes de la plataforma. |

<br/>

> 📚 **Nota sobre la implementación:** Todos los scrapers utilizan la biblioteca [**Colly**](https://github.com/gocolly/colly), un framework de scraping rápido y elegante para Go.

---

## ⚙️ Funcionamiento del Sistema

El proceso de scraping es una tarea fundamental que se ejecuta en segundo plano para mantener la base de datos de productos actualizada.

1.  **Invocación Programada**: El planificador de tareas, definido en `internal/interface/cron/readme.md`, invoca al `ScraperUseCase` cada 48 horas.

2.  **Ejecución por Categoría**: El `ScraperUseCase` itera sobre todas las categorías activas del sistema (Portátiles, Tarjetas Gráficas, etc.) y ejecuta cada uno de los scrapers implementados para esa categoría.

3.  **Extracción de Datos**: Cada scraper visita la URL correspondiente de la tienda y extrae una lista de productos con su información esencial:
    -   Nombre del producto
    -   Precio
    -   URL de la página de detalle
    -   URL de la imagen

4.  **Validación de Relevancia**: Una vez que un scraper devuelve una lista de productos, estos se pasan por un validador (`pkg/utils/category_validator.go`). Este es un paso **crítico** que utiliza un sistema de palabras clave para asegurar que un producto extraído (ej: "funda para portátil") no sea incorrectamente asignado a una categoría principal (ej: "Portátiles"). Esto garantiza una alta calidad y relevancia de los datos. Para más detalles, consulta la documentación en `pkg/utils/readme.md`.

5.  **Persistencia de Datos**: Los productos validados son procesados por el `ScraperUseCase` para ser guardados en la base de datos. El sistema comprueba si el producto ya existe para actualizar su precio, o lo crea si es nuevo.

---

