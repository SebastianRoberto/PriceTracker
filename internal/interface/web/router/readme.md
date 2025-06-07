# 🌐 Documentación de Endpoints del Sistema

A continuación se documentan los principales endpoints de la aplicación PriceTracker, organizados por funcionalidad.

---

### 👤 Gestión de Usuarios

#### Registro de Usuario
- **`GET /registro`**
  > Muestra el formulario de registro.

- **`POST /registro`**
  > Procesa los datos del formulario de registro.
  >
  > **Cuerpo del Formulario:**
  >
  > | Parámetro          | Descripción                  |
  > |:-------------------|:-----------------------------|
  > | `username`         | Nombre de usuario.           |
  > | `email`            | Correo electrónico.          |
  > | `password`         | Contraseña.                  |
  > | `confirm_password` | Confirmación de la contraseña. |
  >
  > ✅ **Respuesta Exitosa**: Redirección a `/registro-exitoso`.

#### Página de Registro Exitoso
- **`GET /registro-exitoso`**
  > Muestra una página informando al usuario que revise su email para la verificación.

#### Verificación de Email
- **`GET /verificar`**
  > Activa la cuenta de un usuario a través de un token recibido por email.
  >
  > **Parámetros de la URL (Query):**
  >
  > | Parámetro | Descripción                   |
  > |:----------|:------------------------------|
  > | `token`   | Token de verificación único. |
  >
  > ✅ **Respuesta Exitosa**: Redirección a `/login`.

#### Inicio de Sesión
- **`GET /login`**
  > Muestra el formulario de inicio de sesión.

- **`POST /login`**
  > Procesa las credenciales de inicio de sesión del usuario.
  >
  > **Cuerpo del Formulario:**
  >
  > | Parámetro  | Descripción          |
  > |:-----------|:---------------------|
  > | `email`    | Correo electrónico.  |
  > | `password` | Contraseña.          |
  >
  > ✅ **Respuesta Exitosa**: Redirección a la página principal.

#### Cierre de Sesión
- **`GET /logout`**
  > Cierra la sesión del usuario actual. (Requiere autenticación).

#### Perfil de Usuario
- **`GET /perfil`**
  > Muestra la página con los datos del usuario autenticado. (Requiere autenticación).

#### Cambiar Contraseña
- **`POST /cambiar-password`**
  > Actualiza la contraseña del usuario autenticado. (Requiere autenticación).
  >
  > **Cuerpo del Formulario:**
  >
  > | Parámetro              | Descripción                       |
  > |:-----------------------|:----------------------------------|
  > | `current_password`     | Contraseña actual del usuario.    |
  > | `new_password`         | Nueva contraseña deseada.         |
  > | `confirm_new_password` | Confirmación de la nueva contraseña.|
  >
  > ✅ **Respuesta Exitosa**: Redirección al perfil con mensaje de éxito.

#### Borrar Cuenta
- **`POST /borrar-cuenta`**
  > Elimina permanentemente la cuenta del usuario autenticado. (Requiere autenticación).
  >
  > **Cuerpo del Formulario:**
  >
  > | Parámetro  | Descripción                         |
  > |:-----------|:------------------------------------|
  > | `password` | Contraseña actual para confirmar. |
  >
  > ✅ **Respuesta Exitosa**: Cierre de sesión y redirección a la página principal.

---

### 🔑 Flujo "He Olvidado Mi Contraseña" (Público)

#### Formulario para Solicitar Restablecimiento
- **`GET /forgot-password`**
  > Muestra un formulario para que el usuario introduzca su email y solicitar un enlace de restablecimiento.

#### Procesar Solicitud de Restablecimiento
- **`POST /forgot-password`**
  > Procesa la solicitud. Si el email existe, envía un correo con el enlace.
  >
  > **Cuerpo del Formulario:**
  >
  > | Parámetro | Descripción                     |
  > |:----------|:--------------------------------|
  > | `email`   | Correo del usuario que olvidó su contraseña. |
  >
  > ✅ **Respuesta Exitosa**: Envía el correo (no revela si el email existe por seguridad).

#### Formulario para Nueva Contraseña (con token)
- **`GET /restablecer-password`**
  > Muestra el formulario para introducir la nueva contraseña, validando el token desde la URL.
  >
  > **Parámetros de la URL (Query):**
  >
  > | Parámetro | Descripción                     |
  > |:----------|:--------------------------------|
  > | `token`   | Token de restablecimiento.      |
  >
  > ✅ **Respuesta Exitosa**: Muestra el formulario.

#### Procesar Nueva Contraseña
- **`POST /restablecer-password`**
  > Procesa el formulario con la nueva contraseña y el token.
  >
  > **Cuerpo del Formulario:**
  >
  > | Parámetro          | Descripción                     |
  > |:-------------------|:--------------------------------|
  > | `token`            | Token de restablecimiento.      |
  > | `password`         | Nueva contraseña.               |
  > | `confirm_password` | Confirmación de la nueva contraseña. |
  >
  > ✅ **Respuesta Exitosa**: Redirección a `/login` con mensaje de éxito.

#### Solicitud de Reset de Contraseña (Usuario Logueado)
- **`GET /solicitar-reset`**
  > Envía un email con un enlace para restablecer la contraseña al usuario autenticado. (Requiere autenticación).

---

### 🔍 Navegación de Productos

#### Página Principal
- **`GET /`**
  > Muestra la página principal con las categorías destacadas.

#### Listado por Categoría
- **`GET /categoria/{slug}`**
  > Muestra el listado de productos para una categoría específica.

#### Detalle de Producto
- **`GET /producto/{id}`**
  > Muestra la página de detalle de un producto, incluyendo su historial de precios.

#### API de Categoría (JSON)
- **`GET /api/categoria/{slug}`**
  > Devuelve los datos de los productos de una categoría en formato JSON.

---

### 📊 Gestión de Alertas y Seguimiento
(Corresponde a la sección "Mi Cesta" en la interfaz de usuario)

#### Ver "Mi Cesta" (Alertas de Precio)
- **`GET /watchlist`**
  > Muestra la página con productos en seguimiento y alertas configuradas. (Requiere autenticación).

#### Añadir Producto a "Mi Cesta" (Crear Alerta de Precio)
- **`POST /price-alert/set`**
  > API llamada desde JavaScript para crear/actualizar una alerta de precio. (Requiere autenticación).
  >
  > **Cuerpo de la Petición (JSON):**
  >
  > | Campo          | Descripción                      |
  > |:---------------|:---------------------------------|
  > | `product_id`   | ID del producto a seguir.        |
  > | `target_price` | Precio objetivo para la alerta. |
  >
  > ✅ **Respuesta Exitosa (JSON)**: `{ "success": true, "message": "¡Producto añadido a tu cesta!" }`
  >
  > ❌ **Respuesta de Error (JSON)**: `{ "success": false, "message": "Mensaje de error." }`

#### Eliminar Alerta de Precio
- **`GET /price-alert/delete`**
  > Elimina una alerta de precio de la lista del usuario. (Requiere autenticación).
  >
  > **Parámetros de la URL (Query):**
  >
  > | Parámetro | Descripción        |
  > |:----------|:-------------------|
  > | `id`      | ID de la alerta a eliminar. |
  >
  > ✅ **Respuesta Exitosa**: Redirección a `/watchlist`.

#### Actualizar Alerta de Precio
- **`GET /price-alert/update`**
  > Actualiza el precio objetivo de una alerta existente. (Requiere autenticación).
  >
  > **Parámetros de la URL (Query):**
  >
  > | Parámetro      | Descripción             |
  > |:---------------|:------------------------|
  > | `id`           | ID de la alerta.        |
  > | `target_price` | Nuevo precio objetivo.  |
  >
  > ✅ **Respuesta Exitosa**: Redirección a `/watchlist`.

#### Alias de la Lista de Seguimiento
- **`GET /price-alerts`**
  > Redirección a `/watchlist` por compatibilidad. (Requiere autenticación).

---

### 🔔 Sistema de Notificaciones

#### Listar Notificaciones
- **`GET /notificaciones`**
  > Muestra la página con el listado de notificaciones del usuario. (Requiere autenticación).

#### Marcar Notificación como Leída
- **`POST /notificaciones/marcar-leida`**
  > Marca una notificación específica como leída. (Requiere autenticación).
  >
  > **Cuerpo del Formulario:**
  >
  > | Parámetro         | Descripción                  |
  > |:------------------|:-----------------------------|
  > | `notification_id` | ID de la notificación a marcar. |
  >
  > ✅ **Respuesta Exitosa**: Actualización del estado de la notificación.

#### Marcar Todas como Leídas
- **`POST /notificaciones/marcar-leidas`**
  > Marca todas las notificaciones del usuario como leídas. (Requiere autenticación).

#### Eliminar Notificaciones Leídas (API)
- **`POST /api/notifications/delete-read`**
  > Elimina todas las notificaciones que ya han sido leídas. (Requiere autenticación).
  >
  > ✅ **Respuesta Exitosa (JSON)**: `{ "success": true, "message": "Notificaciones leídas eliminadas." }`

---

