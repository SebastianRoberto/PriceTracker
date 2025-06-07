# 💾 Capa de Persistencia (`/persistance`)

Este directorio contiene la **implementación concreta** de las interfaces de repositorio definidas en [`/internal/domain/repositories`](../../domain/repositories/readme.md). Es el "adaptador" de la base de datos para nuestra aplicación.

Mientras que la capa de `domain` define el **QUÉ** (el contrato), esta capa de `persistance` define el **CÓMO** (la lógica específica para interactuar con la base de datos).

La tecnología utilizada para esta implementación es [**GORM**](https://gorm.io/), un ORM (Object-Relational Mapper) para Go que simplifica enormemente las operaciones con la base de datos.

---

## 🧩 Componentes Clave

### `db.go`
Este archivo es el punto de entrada y configuración de la base de datos.
- **`NewDatabase(...)`**: Se encarga de establecer la conexión con la base de datos MySQL. Crucialmente, primero se conecta sin especificar una base de datos para ejecutar `CREATE DATABASE IF NOT EXISTS`, asegurando que la base de datos exista antes de intentar usarla.
- **`AutoMigrate()`**: Utiliza la función de GORM para crear o actualizar automáticamente el esquema de la base de datos (tablas, columnas, índices) basándose en las definiciones de los [modelos de dominio](../../domain/model/readme.md). Esto asegura que la estructura de la base de datos siempre esté sincronizada con el código.

### Implementaciones de Repositorios (`*_repository.go`)

Cada archivo `*_repository.go` en este directorio corresponde a una entidad de dominio y cumple un contrato específico.

| Archivo | Interfaz Implementada | Descripción de la Implementación |
| :--- | :--- | :--- |
| `user_repository.go` | [`UserRepository`](../../domain/repositories/readme.md#userrepository) | Implementa las funciones para gestionar usuarios (`Create`, `FindByID`, etc.) utilizando métodos de GORM como `db.Create()` y `db.First()`. |
| `product_repository.go`| [`ProductRepository`](../../domain/repositories/readme.md#productrepository) | Contiene la lógica para interactuar con productos. Incluye consultas complejas con `JOINs` y subconsultas para filtros avanzados y búsqueda de ofertas. |
| `category_repository.go`|[`CategoryRepository`](../../domain/repositories/readme.md#categoryrepository)| Implementa las operaciones para categorías, incluyendo consultas SQL `Raw` para obtener el conteo de productos de manera eficiente. |
| `price_repository.go`| [`PriceRepository`](../../domain/repositories/readme.md#pricerepository) | Gestiona los precios de los productos, con funciones clave como `FindBestPriceByProductID` que utiliza `ORDER BY price asc` para encontrar la mejor oferta. |
| `price_alert_repository.go`|[`PriceAlertRepository`](../../domain/repositories/readme.md#pricealertrepository--notificationrepository)| Implementa las operaciones para las alertas de precio. |
| `notification_repository.go`|[`NotificationRepository`](../../domain/repositories/readme.md#pricealertrepository--notificationrepository)| Gestiona la creación, búsqueda y actualización de notificaciones para los usuarios. |
| `watchlist_repository.go`|[`Watchlist...`](../../domain/repositories/readme.md#watchlistrepository--watchlistitemrepository)| Implementa la lógica para la "Cesta". Destaca la función `FindByUserID` que crea una lista de seguimiento para un usuario si no tiene una, asegurando que cada usuario siempre tenga una lista disponible. |

Gracias a esta estructura, si en el futuro se decidiera cambiar de MySQL a otra base de datos como PostgreSQL, solo habría que modificar el código dentro de esta carpeta (`persistance`) y, potencialmente, el conector en `db.go`, sin afectar a ninguna otra parte del sistema. 