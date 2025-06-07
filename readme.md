# 🚀 PriceTracker - Comparador de Precios de Tecnología

<p align="center">
  <img src="https://mioti.es/wp-content/uploads/2023/05/AdobeStock_474211244-2.jpeg" alt="PriceTracker Screenshot" width="800"/>
</p>

PriceTracker es una aplicación web desarrollada en Go que rastrea precios de productos tecnológicos en múltiples tiendas online mediante web scraping, permitiendo a los usuarios establecer alertas para recibir notificaciones cuando los productos alcanzan un precio objetivo. La siguiente documentación tiene como objetivo describir de forma general su funcionamiento; para información más detallada sobre apartados concretos, puedes consultar la documentación específica de los distintos apartados en este mismo repositorio.

---

## ✨ Características

-   **Comparación de precios en tiempo real**: Datos actualizados regularmente desde eBay, Coolmod y Aussar.
-   **Categorías especializadas**: Portátiles, GPUs, auriculares, teclados, monitores y SSDs.
-   **Alertas personalizadas**: Notificaciones en la plataforma y por correo electrónico cuando los productos alcanzan un precio objetivo.
-   **Sistema de usuarios completo**: Registro, verificación por email, login, perfil de usuario y recuperación de contraseña.
-   **Validación de productos por categoría**: Un sistema de reglas con palabras clave para asegurar que los productos extraídos vayan a sus categorías correspondientes o se excluyan del sistema en caso de no pertenecer a ninguna de las categorías para las que se da soporte.
-   **Seguridad**: Contraseñas hasheadas con `bcrypt`, tokens de seguridad para verificación de usuario y restablecimiento de contraseña.
-   **Interfaz de usuario interactiva**: Validaciones de formulario en tiempo real, notificaciones dinámicas y animaciones para una experiencia de usuario fluida.

---

## 🛠️ Tecnologías Utilizadas

| Categoría      | Tecnología                                           |
| :------------- | :--------------------------------------------------- |
| **Backend**    | Go, Gin Framework, GORM, Colly, Goquery, MySQL       |
| **Frontend**   | HTML, CSS, JavaScript, Bootstrap 5, jQuery           |
| **Herramientas** | `robfig/cron` (Tareas programadas), `gin-sessions` (Gestión de sesiones) |

---

## 🚀 Instalación y Configuración

Sigue estos pasos para poner en marcha el proyecto en tu entorno local.

### Requisitos
- Go 1.18 o superior
- MySQL 5.7 o superior

### Pasos

1.  **Clonar el Repositorio**:
    ```bash
    git clone https://github.com/SebastianRoberto/ProyectoFinalBase.git
    cd ProyectoFinalBase
    ```

2.  **Configurar las Variables de Entorno (MUY IMPORTANTE)**:
    Para que la aplicación funcione, necesitamos credenciales para la base de datos y para el servicio de envío de correos. Por seguridad, el archivo con las credenciales reales (`configs/config.yaml`) se ignora en el control de versiones. Deberás crearlo a partir de la plantilla proporcionada.

    **a. Renombrar el archivo de ejemplo:**
    Dentro de la carpeta `configs/`, cambia el nombre de `config.example.yaml` a `config.yaml`.
    
    **b. Editar `config.yaml`:**
    Abre el archivo `configs/config.yaml` y rellena los datos de tu base de datos en la sección `database` incluyendo tu correo y la contraseña de aplicacion(no es lo mismo que la contraseña que usas para iniciar sesion).

    **¿Por qué es necesario esto y qué es una "Contraseña de Aplicación"?**
    La aplicación envía correos automáticos (para verificar cuentas, enviar alertas, etc.). Para hacer esto de forma segura, no puede usar tu contraseña principal de Gmail. Si lo hiciera, tus credenciales estarían guardadas en un archivo y sería un riesgo de seguridad enorme.
    
    Una **Contraseña de Aplicación** es una contraseña de 16 dígitos que le da a una aplicación específica permiso para acceder a tu cuenta de Google, pero de forma muy limitada (solo para enviar correos). Si alguna vez quieres quitarle el acceso, simplemente borras esa contraseña de aplicación en tu cuenta de Google, sin necesidad de cambiar tu contraseña principal. Es la forma moderna y segura de permitir que programas automatizados usen tu cuenta.

    **a. Activar la Verificación en Dos Pasos (Requisito de Google):**
    - Ve a la [página de seguridad de tu cuenta de Google](https://myaccount.google.com/security).
    - Busca la sección "Cómo inicias sesión en Google" y asegúrate de que la **"Verificación en dos pasos" esté Activada**. Si no lo está, actívala. Google no te permitirá crear contraseñas de aplicación sin esto.

    **b. Generar la Contraseña de Aplicación:**
    - En la misma página de seguridad, haz clic en **"Contraseñas de aplicaciones"**.
    - En el desplegable "Seleccionar aplicación", elige **"Otra (nombre personalizado)"**.
    - Dale un nombre que recuerdes, como "PriceTracker Go App", y haz clic en **"Generar"**.
    - Google te mostrará una **contraseña de 16 caracteres** en un recuadro amarillo. **Copia esta contraseña; no la podrás volver a ver.**

    **c. Añadir credenciales de email al `config.yaml`:**
    Vuelve a tu archivo `configs/config.yaml`(al copiar el repo lo tendras como config.example.yaml) y rellena la sección `email` con tu correo y la contraseña de 16 caracteres que acabas de copiar:
    ```yaml
    email:
      smtp_host: "smtp.gmail.com"
      smtp_port: 587
      smtp_user: "tu_email@gmail.com"
      smtp_pass: "aqui_va_la_contraseña_de_16_letras" 
      smtp_from: "tu_email@gmail.com"
    ```

    **d. Añadir credenciales de base de datos al `config.yaml`:**
    De la misma manera, en el archivo `configs/config.yaml`, rellena la sección `database:` con los datos de tu servidor MySQL. El `username` y `password` deben ser los que usas para conectarte a tu gestor de base de datos.
    ```yaml
    database:
      driver: "mysql"
      host: "localhost"
      port: 3306
      username: "root"
      password: "TU_CONTRASENA_DE_BD" # <-- REEMPLAZAR
      name: "comparador_precios"
    ```


3.  **Instalar Dependencias**:
    Desde la raíz del proyecto, ejecuta:
    ```bash
    go mod tidy
    ```

4.  **Ejecutar la Aplicación**:
    ```bash
    go run cmd/main.go
    ```
    Al arrancar, verás logs en la consola. La aplicación creará la base de datos `comparador_precios` con la estructura de tablas si detecta que esta no existe.
    
5.  **Acceder en el Navegador**:
    La aplicación estará disponible en `http://localhost:8080`.

6   **Aviso**:
    Si el firewall te empieza a dar problemas y pedir permisos cada vez que intentes ejecutar el programa haz uso del setup_firewall.bat que esta ubicado en /scripts, ve al explorador de archivos y ejecutalo como administrador.

---

## 🏗️ Arquitectura del Sistema

El sistema se divide en una serie de apartados, cada uno con tareas muy específicas para facilitar la escalabilidad del sistema y una correcta cohesión entre sus distintos componentes, algunos modulos clave son:

-   `internal/domain`: El núcleo de la aplicación. Contiene los **modelos** y las **interfaces de los repositorios**.
-   `internal/usecase`: Contiene la lógica de negocio pura y los casos de uso.
-   `internal/infrastructure`: Implementaciones concretas de las interfaces (Base de Datos, Scrapers, Email).
-   `internal/interface`: Adaptadores que conectan el mundo exterior con los casos de uso (Handlers Web, Tareas Cron).
-   `cmd/main.go`: Punto de entrada que inicializa y conecta todos los componentes.

```
Directorio_Raiz/
├── cmd/                 # Punto de entrada de la aplicación
├── configs/             # Archivos de configuración
├── internal/
│   ├── domain/
│   │   ├── model/       # Entidades de dominio
│   │   └── repositories/ # Interfaces de repositorio
│   ├── infrastructure/
│   │   ├── email/
│   │   ├── persistance/
│   │   └── scraper/
│   ├── interface/
│   │   ├── cron/        # Tareas programadas
│   │   └── web/         # Handlers, middleware, router
│   └── usecase/         # Lógica de negocio
├── pkg/
│   └── utils/           # Paquetes de utilidad
├── scripts/             # Scripts de utilidad (setup.sql)
└── web/
    ├── static/          # CSS, JS, imágenes
    └── templates/       # Plantillas HTML
```

---

## 🗂️ Modelo de Datos

La base de datos se estructura en torno a los siguientes modelos principales, definidos en `internal/domain/model`:

-   **User**: Almacena los datos de los usuarios registrados, incluyendo credenciales y estado de verificación.
-   **Category**: Define las categorías de los productos (ej. "Portátiles", "Monitores") para la organización.
-   **Product**: Contiene la información general de un producto, como nombre, descripción e imagen.
-   **Price**: Guarda el historial de precios de un producto en una tienda y fecha específicas.
-   **PriceAlert**: Representa las alertas que un usuario configura para un producto a un precio objetivo.
-   **Notification**: Almacena las notificaciones generadas para los usuarios (ej. una alerta de precio alcanzada).
-   **Watchlist / WatchlistItem**: Modela la "cesta" o lista de seguimiento de un usuario, que contiene los productos que le interesan.

---

## 🌐 Endpoints del Sistema

A continuación se documentan los principales endpoints de la aplicación.

<details>
<summary><strong>👤 Gestión de Usuarios</strong></summary>

-   **Registro de Usuario**
    -   `GET /registro`: Muestra el formulario de registro.
    -   `POST /registro`: Procesa los datos del nuevo usuario.
-   **Inicio y Cierre de Sesión**
    -   `GET /login`: Muestra el formulario de inicio de sesión.
    -   `POST /login`: Autentica al usuario.
    -   `GET /logout`: Cierra la sesión del usuario.
-   **Verificación de Email**
    -   `GET /verificar`: Valida la cuenta del usuario a través de un token.
    -   `GET /verificar`: Valida la cuenta del usuario a través de un token.
-   **Gestión de Perfil (Requiere autenticación)**
    -   `GET /perfil`: Muestra la página del perfil del usuario.
    -   `POST /cambiar-password`: Permite al usuario cambiar su contraseña.
    -   `POST /borrar-cuenta`: Permite al usuario eliminar su cuenta.
-   **Recuperación de Contraseña**
    -   `GET /forgot-password`: Muestra el formulario para solicitar el restablecimiento.
    -   `POST /forgot-password`: Envía el email con el enlace de restablecimiento.
    -   `GET /restablecer-password`: Muestra el formulario para introducir la nueva contraseña (requiere token).
    -   `POST /restablecer-password`: Procesa el cambio de la nueva contraseña.

</details>

<details>
<summary><strong>🔍 Navegación de Productos</strong></summary>

-   `GET /`: Página principal con productos destacados.
-   `GET /categoria/{slug}`: Muestra todos los productos de una categoría específica.
-   `GET /producto/{id}`: Muestra la página de detalle de un producto, con su historial de precios.
-   `GET /api/categoria/{slug}`: Endpoint JSON para obtener los productos de una categoría (usado para filtrado y paginación dinámica).

</details>

<details>
<summary><strong>📊 Gestión de Alertas y "Mi Cesta"</strong></summary>

-   `GET /watchlist`: Muestra la "cesta" del usuario con todos los productos que está siguiendo.
-   `POST /price-alert/set`: (API) Añade un producto a la cesta o crea una alerta de precio.
-   `GET /price-alert/update`: Actualiza el precio objetivo de una alerta existente.
-   `GET /price-alert/delete`: Elimina un producto/alerta de la cesta.

</details>

<details>
<summary><strong>🔔 Sistema de Notificaciones</strong></summary>

-   `GET /notificaciones`: Muestra todas las notificaciones del usuario.
-   `POST /notificaciones/marcar-leida`: Marca una notificación específica como leída.
-   `POST /notificaciones/marcar-leidas`: Marca todas las notificaciones como leídas.
-   `POST /api/notifications/delete-read`: (API) Elimina todas las notificaciones que ya han sido leídas.

</details>

<details>
<summary><strong>⚠️ Errores Comunes</strong></summary>

-   **Autenticación**: Credenciales inválidas, sesión expirada, acceso no autorizado.
-   **Registro**: Usuario o email ya existente, contraseña no cumple los requisitos.
-   **Alertas**: Precio objetivo inválido, producto no encontrado, alerta duplicada.

</details>

---

## 🔐 Seguridad

-   **Hash de contraseñas**: Se utiliza `bcrypt` para almacenar las contraseñas de forma segura.
-   **Validación de formularios**: Se valida la entrada del usuario tanto en el frontend como en el backend.
-   **Tokens seguros**: Se usan tokens criptográficamente seguros para la verificación de email y el restablecimiento de contraseña.
-   **Protección de rutas**: Se utilizan middlewares para proteger las rutas que requieren autenticación.

## ⚙️ Tareas Programadas (Cron Jobs)

El sistema ejecuta las siguientes tareas en segundo plano de forma automática:

-   **Scraping completo (Cada 48 horas):** Descubre nuevos productos en todas las tiendas.
-   **Verificación de Alertas (Cada 6 horas):** Comprueba si se ha alcanzado algún precio objetivo y envía notificaciones.
-   **Limpieza de precios (Cada 72 horas):** Elimina registros de precios antiguos para mantener la base de datos optimizada.
