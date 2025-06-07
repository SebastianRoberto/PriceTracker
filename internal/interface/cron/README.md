# Módulo de Tareas Programadas (Cron)

Este directorio contiene la lógica para la ejecución de tareas programadas en segundo plano, como el web scraping periódico y la verificación de alertas de precios. Utiliza la librería `github.com/robfig/cron/v3` para gestionar los trabajos (jobs).

## Componentes Principales

### `scheduler.go`

Es el orquestador central de todas las tareas programadas. Se encarga de:

1.  **Inicializar el planificador**: Crea una nueva instancia de `cron.Cron`.
2.  **Registrar los trabajos**: Añade las funciones que deben ejecutarse y define su periodicidad.
3.  **Iniciar y detener el ciclo de vida** del planificador.

## Tareas Programadas

El planificador ejecuta las siguientes tareas a intervalos definidos:

1.  **Scraping Completo de Productos (`@every 48h`)**
    -   **Disparador**: Se ejecuta cada 48 horas.
    -   **Acción**: Llama a `RunAllScrapers()`, que obtiene todas las categorías de la base de datos y lanza una goroutine por cada scraper (Ebay, Coolmod, Aussar) y por cada categoría.
    -   **Post-Acción**: Una vez finalizado el scraping, invoca `CheckPriceAlerts()` para notificar inmediatamente sobre cualquier oferta que se haya activado con los nuevos precios.
    -   **Nota**: También se ejecuta una vez al iniciar la aplicación para asegurar que hay datos desde el principio.

2.  **Limpieza de Precios Antiguos (`@every 72h`)**
    -   **Disparador**: Se ejecuta cada 3 días.
    -   **Acción**: Llama a `CleanupOldPrices()`, que elimina de la base de datos los registros de precios que no se han actualizado en los últimos 7 días. Esto mantiene la base de datos relevante y con un tamaño controlado.

3.  **Verificación de Alertas de Precio (`@every 6h`)**
    -   **Disparador**: Se ejecuta cada 6 horas.
    -   **Acción**: Invoca `CheckPriceAlerts()`, que utiliza el `PriceAlertUseCase` para comprobar si el precio actual de algún producto ha caído por debajo del precio objetivo fijado por un usuario en su "cesta". Si es así, se crea una notificación y se envía un correo electrónico.
    -   **Nota**: También se ejecuta una vez al iniciar la aplicación.

## Flujo de Trabajo

1.  Al arrancar la aplicación, se crea una instancia del `ScraperScheduler`.
2.  El método `Start()` registra todas las tareas y pone en marcha el planificador.
3.  Los logs de la consola mostrarán la actividad de cada tarea, usando un sistema de colores para diferenciar entre información, éxitos, advertencias y errores.
4.  El planificador se detiene de forma segura cuando la aplicación termina su ejecución mediante el método `Stop()`.

## Dependencias

-   **`app/internal/usecase`**: Interactúa con los casos de uso (ej. `PriceAlertUseCase`) para ejecutar la lógica de negocio.
-   **`app/internal/domain/repositories`**: Depende de las interfaces de los repositorios para acceder a la base de datos (productos, categorías, precios).
-   **`app/internal/infrastructure/scraper`**: Utiliza los scrapers específicos de cada tienda.

---
*Para más información sobre la lógica de negocio, consulta la documentación en  `internal/usecase`.* 