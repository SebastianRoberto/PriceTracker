# Casos de Uso (Lógica de Negocio)

Este directorio representa la **capa de servicio** o **lógica de negocio** de la aplicación, siguiendo los principios de la Arquitectura Limpia (Clean Architecture). Contiene la orquestación de las operaciones y flujos de trabajo del dominio, actuando como intermediario entre la capa de interfaz (handlers, cron jobs) y la capa de dominio (repositorios).

## Principios de Diseño

-   **Separación de Responsabilidades**: Cada archivo `*_usecase.go` agrupa la lógica de un dominio de negocio específico (usuarios, productos, alertas, etc.), manteniendo el código organizado y cohesivo.
-   **Inversión de Dependencias**: Los casos de uso dependen de **interfaces** de repositorio (definidas en `internal/domain/repositories`), no de sus implementaciones concretas. Esto desacopla la lógica de negocio de la infraestructura de persistencia (la base de datos), facilitando las pruebas unitarias (mediante mocks) y la flexibilidad para cambiar la tecnología de la base de datos si fuera necesario.

## Casos de Uso Implementados

### `user_usecase.go`

-   **Responsabilidad**: Gestiona todo el ciclo de vida y las operaciones relacionadas con los usuarios.
-   **Funciones Clave**:
    -   `AuthenticateUser`: Verifica las credenciales de un usuario.
    -   `CreateUser`: Registra un nuevo usuario, hashea su contraseña y genera un token de verificación.
    -   `SendVerificationEmail`: Envía un correo para activar la cuenta.
    -   `VerifyUser`: Valida un token y marca la cuenta como verificada.
    -   `ChangePassword`, `InitiatePasswordReset`, `ResetPassword`: Gestionan todos los flujos de cambio de contraseña.
    -   `DeleteAccount`: Elimina una cuenta de usuario de forma segura.

### `product_usecase.go`

-   **Responsabilidad**: Proporciona métodos para consultar información de productos de una manera que sea útil para la UI.
-   **Funciones Clave**:
    -   `GetBestDeals`, `GetFeaturedProducts`: Obtiene listas de productos para la página de inicio.
    -   `GetProductsByCategory`: Devuelve productos filtrados y paginados para las vistas de categoría.
    -   `GetProductDetail`, `GetSimilarProducts`: Recupera toda la información para la página de detalle de un producto, incluyendo sus precios y productos relacionados.
    -   `GetFilteredProductsByCategory`: Orquesta la búsqueda avanzada de productos aplicando filtros de precio, tienda y ordenación.

### `price_alert_usecase.go`

-   **Responsabilidad**: Contiene toda la lógica de la "cesta" de seguimiento y el sistema de notificaciones.
-   **Funciones Clave**:
    -   `CreateAlert`, `UpdateAlert`, `DeleteAlert`: Permite a los usuarios añadir, modificar o eliminar productos de su cesta.
    -   `CheckPriceAlerts`: Lógica central que es llamada por el `cron`. Compara los precios actuales con los objetivos de los usuarios y, si se cumple una condición, dispara la creación de notificaciones.
    -   `GetUserNotifications`, `MarkNotificationAsRead`: Gestiona la visualización y el estado de las notificaciones para el usuario.
    -   `createNotification`: Proceso interno que guarda una notificación en la base de datos y (si el usuario lo desea) envía un correo electrónico a través del `Mailer`.

### `scraper_usecase.go`

-   **Responsabilidad**: Orquesta el proceso completo de web scraping. Es uno de los componentes más complejos.
-   **Funciones Clave**:
    -   `ScrapeAllCategories`, `ScrapeCategory`: Inicia el proceso de scraping para todas o una categoría específica, invocando a los scrapers de la capa de `infrastructure`.
    -   `saveProducts`, `saveProduct`: Contiene la lógica crucial para procesar los productos scrapeados antes de guardarlos:
        1.  **Validación y Reclasificación**: Utiliza `utils.ValidateProductCategory` para asegurar que un producto pertenece a la categoría correcta. Si no, intenta reclasificarlo.
        2.  **Deduplicación**: Esto no esta completamente implementado pero el sistema esta pensado para utilizar un sistema para evitar duplicados a futuro utilizando un:
            -   **Hash de Imagen (pHash)**: Calcula un hash perceptual de la imagen del producto y lo compara con los existentes para encontrar duplicados visuales.
            -   **Slug**: Si no hay coincidencia por imagen, recurre a la comparación por `slug`.
        3.  **Persistencia**: Decide si crear un nuevo producto o actualizar uno existente con un nuevo precio.

## Flujo de Datos Típico

Un `Handler` HTTP recibe una petición -> llama a un método del `UseCase` apropiado -> el `UseCase` ejecuta la lógica, posiblemente llamando a varios `Repositories` para leer o escribir datos -> el `UseCase` devuelve el resultado al `Handler` -> el `Handler` renderiza una `Template` o devuelve una respuesta JSON.

---
*Este directorio es el núcleo de la aplicación. Para entender cómo se almacenan y recuperan los datos, consulta la documentación de `internal/domain/repositories`.* 