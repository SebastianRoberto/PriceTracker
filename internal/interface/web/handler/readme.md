# 🔌 Capa de Handlers

Los `handlers` son el punto de entrada para todas las peticiones HTTP que llegan a la aplicación, actuando como el conector principal entre el mundo exterior (el navegador del usuario) y la lógica de negocio interna. Se encuentran en la capa de `Interfaz`, y su rol es orquestar el flujo de datos para cada solicitud.

---

## 🌊 Flujo Típico de una Petición

Cada `handler` sigue un patrón de trabajo consistente que asegura la separación de responsabilidades y la claridad del código:

1.  **Recepción**: La petición llega desde el `router` (ver `internal/interface/web/router/readme.md`), que la dirige al método del `handler` apropiado.
2.  **Middleware**: Antes de que el `handler` se ejecute, los `middlewares` (ver `internal/interface/web/middleware/readme.md`) pueden procesar la petición para realizar tareas como la autenticación de usuarios o la carga de datos comunes.
3.  **Análisis y Validación**: El `handler` extrae y valida los datos de la petición, como parámetros de la URL, datos de formularios (`c.PostForm`) o cuerpos JSON.
4.  **Invocación del Caso de Uso**: El `handler` llama a los métodos correspondientes en la capa de `Usecase` (ej: `userUseCase.Register(...)`), pasándoles los datos validados. Nunca contiene lógica de negocio directamente.
5.  **Procesamiento de la Respuesta**:
    -   **Éxito**: Si el caso de uso devuelve datos, el `handler` los empaqueta en un `ViewModel` y utiliza el `TemplateRenderer` (ver `internal/interface/web/views/readme.md`) para renderizar una plantilla HTML. Para las rutas de API, devuelve una respuesta JSON.
    -   **Error**: Si el caso de uso devuelve un error, el `handler` renderiza una página de error o devuelve un JSON con el código de estado y mensaje apropiados.

---

## 📋 Handlers Implementados

Cada archivo en este directorio agrupa la lógica para una entidad o funcionalidad específica del sistema.

| Archivo                        | Responsabilidad Principal                                                                                                        |
| :----------------------------- | :------------------------------------------------------------------------------------------------------------------------------- |
| **`auth_handler.go`**          | Gestiona todo el ciclo de vida del usuario: registro, verificación por email, inicio de sesión, cierre de sesión y recuperación de contraseña. También maneja la lógica de la página de perfil para cambiar contraseña y eliminar la cuenta. |
| **`category_handler.go`**      | Muestra la página de una categoría de productos. Incluye una versión para renderizado en servidor (`GetCategory`) y una API (`GetCategoryAPI`) para el filtrado dinámico y paginación con JavaScript. |
| **`home_handler.go`**          | Controla la página de inicio de la aplicación, obteniendo y mostrando los productos destacados o las mejores ofertas.               |
| **`notification_handler.go`**  | Gestiona la visualización y las acciones sobre las notificaciones del usuario, como marcarlas como leídas o eliminarlas.              |
| **`price_alert_handler.go`**   | Maneja toda la lógica relacionada con "Mi Cesta" (Watchlist) y las alertas de precio. Permite a los usuarios añadir, actualizar y eliminar productos de su lista de seguimiento. |
| **`product_handler.go`**       | Muestra la página de detalle para un producto específico, incluyendo su información, historial de precios y productos similares.     |
| **`user_handler.go`**          | Contiene lógica adicional del perfil de usuario. Aunque gran parte de la gestión de perfil está en `auth_handler.go` por cohesión con la autenticación, este handler podría expandirse en el futuro. |

---

## 🔗 Relaciones con Otros Módulos

-   **`Router`**: Es el encargado de dirigir las peticiones a estos handlers. La configuración de rutas se encuentra en `internal/interface/web/router/router.go`.
-   **`Usecase`**: Los handlers dependen directamente de los casos de uso para ejecutar la lógica de negocio. Son sus principales consumidores.
-   **`Views`**: Para las respuestas HTML, los handlers utilizan el `TemplateRenderer` y los `ViewModels` definidos en la capa de vistas para construir y enviar la página final al usuario.
-   **`Middleware`**: La funcionalidad de los handlers es extendida y protegida por los middlewares, que se aplican a nivel de `router`. 