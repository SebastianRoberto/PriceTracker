# 🎨 Recursos Estáticos (`/web/static`)

Este directorio contiene todos los activos del lado del cliente (frontend) que se sirven directamente al navegador sin procesamiento del lado del servidor. Son la base de la apariencia visual y la interactividad de la aplicación.

---

## 📁 Estructura del Directorio

```
/web/static
├── 🎨 /css
│   ├── styles.css
│   └── toast.css
├── 🖼️ /img
│   └── webscraping.png
└── ⚡ /js
    └── main.js
```

---

### 🎨 `css/`

Contiene las hojas de estilo que definen la apariencia visual completa de la aplicación.

- **`styles.css`**: Es la hoja de estilo principal y más importante. Define:
    - La paleta de colores (tema oscuro moderno).
    - La tipografía y estilos base.
    - La estructura del layout (header, footer, etc.).
    - El diseño de componentes complejos como las tarjetas de producto, perfiles de usuario y la lista de seguimiento (`watchlist`).
    - Animaciones y efectos `hover` para una experiencia de usuario fluida.
    - Estilos responsivos para asegurar la correcta visualización en dispositivos móviles.

- **`toast.css`**: Contiene estilos específicos y dedicados para las notificaciones (toasts) que aparecen para dar feedback al usuario (e.g., "Producto añadido a la cesta").

---

### 🖼️ `img/`

Almacena todas las imágenes estáticas utilizadas en la interfaz.

- **`fempalogo.png`**: Logo de la aplicación.
- **`webscraping.png`**: Imagen ilustrativa utilizada en el `README` principal.

---

### ⚡ `js/`

Este directorio alberga la lógica del lado del cliente que dota de vida e interactividad a la aplicación.

- **`main.js`**: Es el corazón de la funcionalidad del frontend. Sus responsabilidades incluyen:
    - **Inicialización de Componentes**: Activa y configura los componentes de Bootstrap como `Tooltips`, `Dropdowns` y `Toasts`.
    - **Validación de Formularios**: Proporciona feedback en tiempo real al usuario en formularios como el de registro o cambio de contraseña.
    - **Interactividad y AJAX**: Maneja eventos de usuario (clics, etc.) para realizar acciones sin recargar la página, como:
        - Añadir productos a la cesta (`/price-alert/set`).
        - Marcar notificaciones como leídas.
        - Actualizar configuraciones de usuario.
    - **Animaciones y Efectos**: Controla animaciones de CSS y JavaScript para mejorar la experiencia de usuario (e.g., la animación del icono del carrito).
    - **Lógica Específica de Página**: Ejecuta código concreto dependiendo de la página en la que se encuentre el usuario (página de perfil, detalle de producto, etc.). 