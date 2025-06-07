# Plantillas HTML (Templates)

Este directorio alberga todos los archivos de plantilla HTML que componen la interfaz de usuario de la aplicación. Las plantillas están escritas utilizando el motor `html/template` nativo de Go, lo que permite renderizar HTML dinámico inyectando datos desde el backend.

## Estructura y Funcionamiento

El sistema de plantillas se organiza en torno a un diseño base (`layout.html`) que es extendido por el resto de vistas.

### `layout.html`
Es la plantilla maestra o esqueleto de todas las páginas. Define la estructura HTML principal, incluyendo:
-   `<!DOCTYPE>`, `<html>`, `<head>`, `<body>`.
-   Inclusión de todos los assets comunes:
    -   **CSS**: Bootstrap 5, Bootstrap Icons, Google Fonts y hojas de estilo personalizadas (`styles.css`, `toast.css`).
    -   **JavaScript**: Bootstrap Bundle, SweetAlert2, AOS (Animate On Scroll), Tippy.js y el script principal `main.js`.
-   La cabecera (header) con la barra de navegación y el pie de página (footer).
-   Define un bloque de contenido principal `{{ block "content" . }}{{ end }}` que las plantillas específicas llenarán con su contenido único.

### Plantillas de Contenido
Cada archivo `.html` (excepto `layout.html`) define una página o un componente específico. Utilizan la directiva `{{ define "content" }}` para inyectar su HTML dentro del `layout.html`.

-   **`home.html`**: Página de inicio que muestra los productos destacados.
-   **`category.html`**: Muestra la lista de productos de una categoría específica. Incluye una lógica de JavaScript compleja para el filtrado del lado del cliente (por precio, tienda) y la carga perezosa (`load more`).
-   **`product_detail.html`**: Vista detallada de un solo producto. Muestra el mejor precio, una lista de precios, productos relacionados y el formulario para añadir a la "cesta" (crear alerta de precio).
-   **`login.html`**, **`register.html`**: Formularios de inicio de sesión y registro de usuarios.
-   **`register_success.html`**: Página que se muestra tras un registro exitoso, instruyendo al usuario a verificar su email.
-   **`verify_success.html`**: Confirma que la cuenta ha sido verificada correctamente después de que el usuario haga clic en el enlace del email.
-   **`forgot_password.html`**, **`reset_password.html`**: Vistas para el flujo de restablecimiento de contraseña.
-   **`profile.html`**: Página de perfil de usuario donde puede cambiar su contraseña o eliminar su cuenta.
-   **`edit_profile.html`**: Formulario para editar detalles del perfil del usuario, como el nombre de usuario.
-   **`change_password.html`**: Vista con el formulario dedicado exclusivamente a cambiar la contraseña.
-   **`watchlist.html`**: La "cesta" del usuario, que lista todos los productos para los que ha creado una alerta de precio.
-   **`notifications.html`**: Muestra las notificaciones generadas por el sistema (ej. alertas de precio activadas).
-   **`error.html`**: Página genérica para mostrar mensajes de error.

## Inyección de Datos

Las plantillas son renderizadas por los `Handlers` (ver `internal/interface/web/handler`), que les pasan una estructura de datos. Se puede acceder a estos datos dentro de las plantillas usando la sintaxis de punto (`.`), por ejemplo:
-   `{{ .User.Username }}`: Accede al nombre del usuario logueado.
-   `{{ range .Products }}`: Itera sobre una lista de productos.
-   `{{ .Category.Name }}`: Muestra el nombre de la categoría actual.

## Frontend

La interactividad y el estilo se basan en:
-   **Bootstrap 5**: Para el layout, componentes y diseño responsive.
-   **Vanilla JavaScript / Alpine.js**: Para la manipulación del DOM, eventos, y peticiones `fetch` a la API interna (ej. en `category.html` para la carga dinámica de productos).
-   **Librerías externas**: Como `SweetAlert2` para notificaciones y `AOS` para animaciones de scroll.

---
*La lógica que prepara los datos y renderiza estas plantillas se encuentra en `internal/interface/web/handler` y `internal/interface/web/middleware`.* 