# 📝 Repositorios de Dominio (`/repositories`)

Este directorio es una pieza clave de la Arquitectura Limpia/Hexagonal del proyecto. No contiene implementaciones, sino **interfaces** de Go. Cada interfaz define un **contrato** de lo que la aplicación debe ser capaz de hacer con respecto a la persistencia de datos para una entidad de dominio específica.

El propósito de esta capa es definir las operaciones de datos desde el punto de vista de la lógica de negocio, abstrayendo completamente la tecnología de base de datos subyacente. Los casos de uso (`usecase`) dependen de estas interfaces, no de la implementación concreta.

---

## Contracts Definidos

A continuación se detallan las interfaces de repositorio y los métodos que exponen.

### `UserRepository`
Define las operaciones para la entidad [`User`](../model/readme.md).

| Método | Descripción |
| :--- | :--- |
| `Create` | Crea un nuevo usuario. |
| `FindByID` | Busca un usuario por su ID. |
| `FindByEmail` | Busca un usuario por su email. |
| `FindByUsername` | Busca un usuario por su nombre de usuario. |
| `FindByVerifyToken` | Busca un usuario por su token de verificación. |
| `Update` | Actualiza los datos de un usuario. |
| `Delete` | Elimina un usuario. |

### `ProductRepository`
Define las operaciones para la entidad [`Product`](../model/readme.md).

| Método | Descripción |
| :--- | :--- |
| `Create`, `Update`, `Delete` | Operaciones CRUD básicas. |
| `FindByID`, `FindBySlug` | Buscan un producto por ID o por su URL amigable (slug). |
| `FindByCategory`, `FindFilteredProductsByCategory` | Buscan productos dentro de una categoría, con y sin filtros avanzados. |
| `CountByCategory`, `CountFilteredProductsByCategory` | Cuentan productos en una categoría, con y sin filtros. |
| `FindBestDeals`, `FindSimilarProducts` | Lógica de negocio para encontrar ofertas y productos relacionados. |
| `ExistsBySlug` | Comprueba si un producto con un slug dado ya existe. |

### `CategoryRepository`
Define las operaciones para la entidad [`Category`](../model/readme.md).

| Método | Descripción |
| :--- | :--- |
| `Create`, `Update`, `Delete` | Operaciones CRUD básicas. |
| `FindByID`, `FindBySlug`, `GetAll` | Métodos de búsqueda para categorías. |
| `GetCategoryWithProductCount`, `GetAllCategoriesWithProductCount` | Obtienen categorías junto con el número de productos que contienen. |

### `PriceRepository`
Define las operaciones para la entidad [`Price`](../model/readme.md).

| Método | Descripción |
| :--- | :--- |
| `Create`, `Update`, `Delete` | Operaciones CRUD básicas. |
| `FindByID`, `FindByProductID` | Buscan precios por su ID o asociados a un producto. |
| `FindBestPriceByProductID`, `FindTopOffersByProductID` | Buscan la mejor oferta o una lista de las mejores ofertas para un producto. |
| `DeleteOldPrices` | Elimina registros de precios antiguos para mantenimiento. |

### `PriceAlertRepository` & `NotificationRepository`
Definen las operaciones para las entidades [`PriceAlert`](../model/readme.md) y [`Notification`](../model/readme.md).

| Repositorio | Método Destacado | Descripción |
| :--- | :--- | :--- |
| `PriceAlertRepository` | `FindActiveAlertsForPrice` | Encuentra todas las alertas que se cumplen para un producto y un nuevo precio. |
| `NotificationRepository`| `CountUnreadByUserID`| Cuenta las notificaciones no leídas de un usuario. |
| `NotificationRepository`| `MarkAllAsRead` | Marca todas las notificaciones de un usuario como leídas. |

### `WatchlistRepository` & `WatchlistItemRepository`
Definen las operaciones para las entidades [`Watchlist`](../model/readme.md) y [`WatchlistItem`](../model/readme.md).

| Repositorio | Método Destacado | Descripción |
| :--- | :--- | :--- |
| `WatchlistRepository` | `FindByUserID` | Busca (o crea si no existe) la lista de seguimiento de un usuario. |
| `WatchlistItemRepository` | `IsProductInWatchlist` | Comprueba si un usuario ya tiene un producto en su lista. |

La implementación concreta de estas interfaces se encuentra en [`/internal/infrastructure/persistance/`](../../infrastructure/persistance/readme.md). 