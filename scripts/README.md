# Scripts de Configuración

Este directorio contiene scripts de utilidad diseñados para facilitar la configuración inicial y el mantenimiento del entorno de desarrollo del proyecto.

## Scripts Disponibles

### 1. `setup.sql`

-   **Propósito**: Este es un script de `MySQL` que inicializa la base de datos `comparador_precios`. Es fundamental para establecer la estructura de datos sobre la que operará toda la aplicación.
-   **Contenido**: El script realiza las siguientes acciones:
    -   Crea la base de datos `comparador_precios` si no existe.
    -   Define el esquema de todas las tablas necesarias:
        -   `categories`: Para las categorías de productos.
        -   `products`: Para la información de los productos.
        -   `prices`: Para los precios de los productos en diferentes tiendas.
        -   `users`: Para los datos de los usuarios.
        -   `price_alerts`: Para las alertas de precios que configuran los usuarios (la "cesta").
        -   `notifications`: Para las notificaciones generadas por el sistema.
    -   Establece relaciones entre tablas mediante `FOREIGN KEY` (ej. `products.category_id` -> `categories.id`).
    -   Crea `INDEX` en columnas clave para optimizar las consultas.
    -   Define `TRIGGERS` para mantener la integridad de los datos, como actualizar el contador `product_count` en la tabla `categories` automáticamente.
    -   Inserta las categorías iniciales (`Portátiles`, `Tarjetas Gráficas`, etc.) en la base de datos.
-   **Uso**: Para ejecutar este script, se puede usar un cliente de MySQL:
    ```bash
    mysql -u [tu_usuario] -p [tu_base_de_datos] < setup.sql
    ```
    O importarlo directamente desde una herramienta de gestión de bases de datos como DBeaver, HeidiSQL o MySQL Workbench.

### 2. `setup_firewall.bat`

-   **Propósito**: Este es un script de batch para sistemas operativos **Windows**. Su función es configurar el Firewall de Windows para permitir las conexiones entrantes a la aplicación web, no es extrictamente necesario ejecutarla pero en caso el sistema nos este dando alertas y pidiendo permisos constantes se recomienda su uso.
-   **Contenido**: El script utiliza la utilidad de línea de comandos `netsh advfirewall` para:
    -   Eliminar cualquier regla previa con el mismo nombre para evitar duplicados.
    -   Añadir una nueva regla que permite el tráfico `TCP` entrante en el puerto `8080`, que es el puerto por defecto en el que escucha el servidor web de Go.
    -   Añadir una regla que autoriza al ejecutable de Go (`go.exe`) a través del firewall.
-   **Uso**:
    -   Hacer clic derecho sobre el archivo `setup_firewall.bat`.
    -   Seleccionar **"Ejecutar como administrador"**.
    -   El script se ejecutará y configurará las reglas necesarias.

---