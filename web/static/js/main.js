/**
 * Comparador de Precios - Main JavaScript
 * Funcionalidades principales para la interfaz de usuario
 */
document.addEventListener('DOMContentLoaded', function() {
    // Inicialización de componentes de Bootstrap
    initBootstrapComponents();
    
    // Inicializar funcionalidades específicas de páginas
    initPageSpecificFunctionality();
    
    // Configurar eventos globales
    setupGlobalEvents();
    
    // Configurar manejo de visibilidad de página
    setupPageVisibilityHandling();
});

/**
 * Inicializa componentes de Bootstrap
 */
function initBootstrapComponents() {
    // Asegurar que Bootstrap está disponible
    if (typeof bootstrap === 'undefined') {
        console.error('Bootstrap no está cargado correctamente');
        // Intentar cargar Bootstrap manualmente como última opción
        const script = document.createElement('script');
        script.src = 'https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/js/bootstrap.bundle.min.js';
        script.integrity = 'sha384-kenU1KFdBIe4zVF0s0G1M5b4hcpxyD9F7jL+jjXkk+Q2h455rYXK/7HAuoJl+0I4';
        script.crossOrigin = 'anonymous';
        document.head.appendChild(script);
        
        script.onload = function() {
            console.log('Bootstrap cargado manualmente');
            initDropdowns();
        };
        return;
    }
    
    // Tooltips
    const tooltipElements = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltipElements.forEach(el => {
        try {
            new bootstrap.Tooltip(el);
        } catch (e) {
            console.warn('Error al inicializar tooltip:', e);
        }
    });
    
    // Popovers
    const popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'));
    popoverTriggerList.forEach(el => {
        try {
            new bootstrap.Popover(el);
        } catch (e) {
            console.warn('Error al inicializar popover:', e);
        }
    });
    
    // Inicializar toasts
    initializeToasts();
    
    // Inicializar dropdowns manualmente para asegurar que funcionen en todas las páginas
    initDropdowns();

    // Inicializar tooltips de Bootstrap
    const bootstrapTooltips = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
    bootstrapTooltips.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl)
    });

    // Inicializar dropdowns animados
    document.querySelectorAll('.dropdown-menu-animated').forEach(menu => {
        const parent = menu.parentElement;
        
        parent.addEventListener('show.bs.dropdown', function () {
            menu.classList.add('dropdown-menu-animated-enter');
        });
        
        parent.addEventListener('hide.bs.dropdown', function () {
            menu.classList.add('dropdown-menu-animated-leave');
            setTimeout(() => {
                menu.classList.remove('dropdown-menu-animated-enter', 'dropdown-menu-animated-leave');
            }, 300);
        });
    });
}

/**
 * Inicializa todos los toasts de Bootstrap en la página
 */
function initializeToasts() {
    if (typeof bootstrap === 'undefined' || !bootstrap.Toast) {
        console.warn('Bootstrap Toast no está disponible');
        return;
    }
    
    // Inicializar todos los toasts
    document.querySelectorAll('.toast').forEach(toastEl => {
        try {
            // Crear opciones con valores predeterminados
            const options = {
                animation: true,
                autohide: true,
                delay: parseInt(toastEl.getAttribute('data-bs-delay')) || 5000
            };
            
            // Crear instancia del toast
            new bootstrap.Toast(toastEl, options);
        } catch (e) {
            console.warn('Error al inicializar toast:', e);
        }
    });
}

/**
 * Inicializa los dropdowns manualmente
 */
function initDropdowns() {
    const dropdownElementList = document.querySelectorAll('.dropdown-toggle');
    dropdownElementList.forEach(el => {
        if (!el.hasAttribute('data-bs-toggle')) {
            el.setAttribute('data-bs-toggle', 'dropdown');
        }
        
        try {
            new bootstrap.Dropdown(el);
        } catch (e) {
            console.warn('Error al inicializar dropdown:', e);
        }
    });
}

/**
 * Inicializa funcionalidades específicas de cada página
 */
function initPageSpecificFunctionality() {
    // Página de Inicio - Filtro de productos
    initProductFilter();
    
    // Página de Registro - Validación de formulario
    initRegisterFormValidation();
    
    // Página de detalle de producto - Efectos al añadir a la cesta
    initProductDetailPage();
    
    // Página de watchlist (Mi Cesta)
    initWatchlistPage();
    
    // Página de perfil
    initProfilePage();
    
    // Página de notificaciones
    initNotificationsPage();
    
    // Página de cambio de contraseña
    initChangePasswordPage();
    
    // Página de edición de perfil
    initEditProfilePage();
}

/**
 * Configura eventos globales
 */
function setupGlobalEvents() {
    // Configurar desplegables para que no se cierren inmediatamente
    document.querySelectorAll('.dropdown-toggle').forEach(toggle => {
        toggle.addEventListener('click', function(e) {
            const dropdownMenu = this.nextElementSibling;
            if (dropdownMenu && dropdownMenu.classList.contains('dropdown-menu')) {
                const dropdownInstance = bootstrap.Dropdown.getInstance(this);
                if (dropdownInstance) {
                    e.stopPropagation();
                    dropdownInstance.toggle();
                    
                    // Asegurar que el menú no se cierre al hacer clic en él
                    dropdownMenu.addEventListener('click', function(e) {
                        e.stopPropagation();
                    });
                }
            }
        });
    });
    
    // Mejorar manejo de enlaces en dropdown
    document.querySelectorAll('.dropdown-menu a.dropdown-item').forEach(link => {
        link.addEventListener('click', function(e) {
            window.location.href = this.getAttribute('href');
        });
    });
    
    // Mejorar específicamente el dropdown de usuario
    const userDropdown = document.querySelector('.nav-item.dropdown:not(.categories-dropdown)');
    if (userDropdown) {
        const userToggle = userDropdown.querySelector('.dropdown-toggle');
        const userMenu = userDropdown.querySelector('.dropdown-menu');
        
        if (userToggle && userMenu) {
            // Evitar que el menú de usuario se cierre al hacer clic dentro
            userMenu.addEventListener('click', function(e) {
                e.stopPropagation();
            });
            
            // Añadir manejo manual del dropdown de usuario
            userToggle.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                
                // Usar la API del dropdown de Bootstrap si disponible
                if (typeof bootstrap !== 'undefined' && bootstrap.Dropdown) {
                    const instance = bootstrap.Dropdown.getInstance(userToggle) || 
                                    new bootstrap.Dropdown(userToggle);
                    instance.toggle();
                } else {
                    // Fallback manual si bootstrap no está disponible
                    userMenu.classList.toggle('show');
                }
            });
        }
    }
    
    // Asegurar que los enlaces a productos funcionen correctamente
    document.querySelectorAll('a[href^="/producto/"]').forEach(link => {
        link.addEventListener('click', function(e) {
            // Prevenir que el hash cause problemas
            if (this.getAttribute('href').endsWith('#')) {
                e.preventDefault();
                window.location.href = this.getAttribute('href').slice(0, -1);
            }
        });
    });
    
    // Animación al hacer hover en tarjetas de producto
    document.querySelectorAll('.product-card').forEach(card => {
        card.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-10px)';
        });
        
        card.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0)';
        });
    });

    // Añadir efecto de ripple a botones primarios
    document.querySelectorAll('.btn-primary').forEach(btn => {
        btn.addEventListener('click', createRippleEffect);
    });
    
    // Añadir feedback visual cuando se intenta añadir un producto ya en la cesta
    initDuplicateProductAlert();
    
    // Asegurarse que los documentos sean clicables en dispositivos móviles
    document.addEventListener('touchend', function(e) {
        // Evitar que los clicks en áreas no relacionadas con menús cierren los dropdowns
        if (!e.target.closest('.dropdown-menu') && !e.target.closest('.dropdown-toggle')) {
            // Cerrar cualquier dropdown abierto solo si se hace clic fuera de ellos
            document.querySelectorAll('.dropdown-menu.show').forEach(menu => {
                if (!menu.contains(e.target) && !menu.previousElementSibling.contains(e.target)) {
                    const toggle = menu.previousElementSibling;
                    if (toggle && bootstrap.Dropdown.getInstance(toggle)) {
                        bootstrap.Dropdown.getInstance(toggle).hide();
                    } else {
                        menu.classList.remove('show');
                    }
                }
            });
        }
    }, { passive: true });
}

/**
 * Inicializa el filtro de productos en la página de inicio
 */
function initProductFilter() {
    const productFilter = document.getElementById('productFilter');
    if (!productFilter) return;
    
    productFilter.addEventListener('keyup', function(e) {
        const searchText = e.target.value.toLowerCase().trim();
        const productCards = document.querySelectorAll('.product-card');
        
        let foundCount = 0;
        
        productCards.forEach(function(card) {
            const title = card.querySelector('.card-title')?.textContent.toLowerCase() || '';
            const description = card.querySelector('.card-text')?.textContent.toLowerCase() || '';
            const category = card.querySelector('.badge')?.textContent.toLowerCase() || '';
            
            if (title.includes(searchText) || description.includes(searchText) || category.includes(searchText)) {
                card.closest('.col').style.display = 'block';
                foundCount++;
            } else {
                card.closest('.col').style.display = 'none';
            }
        });
        
        // Mostrar mensaje si no hay resultados
        let noResultsMsg = document.getElementById('noResultsMsg');
        if (searchText && foundCount === 0) {
            if (!noResultsMsg) {
                noResultsMsg = document.createElement('div');
                noResultsMsg.id = 'noResultsMsg';
                noResultsMsg.className = 'alert alert-info col-12 mt-3';
                noResultsMsg.textContent = 'No se encontraron productos para tu búsqueda.';
                productCards[0].closest('.row').appendChild(noResultsMsg);
            }
        } else if (noResultsMsg) {
            noResultsMsg.remove();
        }
    });
}

/**
 * Inicializa validación del formulario de registro
 */
function initRegisterFormValidation() {
    const registerForm = document.getElementById('registerForm');
    if (!registerForm) return;
    
    // Validación de campos individuales mientras el usuario escribe
    const username = document.getElementById('username');
    const email = document.getElementById('email');
    const password = document.getElementById('password');
    const confirmPassword = document.getElementById('confirm_password');
    
    // Validar nombre de usuario (sin caracteres extraños)
    if (username) {
        username.addEventListener('input', function() {
            validateField(this, /^[a-zA-Z0-9_.-]+$/, 'El nombre de usuario solo puede contener letras, números, guiones, puntos y guiones bajos');
        });
    }
    
    // Validar email (debe contener @)
    if (email) {
        email.addEventListener('input', function() {
            validateField(this, /^[^\s@]+@[^\s@]+\.[^\s@]+$/, 'Introduce un correo electrónico válido');
        });
    }
    
    // Validar contraseña (sin caracteres muy extraños)
    if (password) {
        password.addEventListener('input', function() {
            validateField(this, /^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]+$/, 'La contraseña contiene caracteres no permitidos');
        });
    }
    
    // Función para validar un campo según una expresión regular
    function validateField(field, regex, errorMessage) {
        const isValid = regex.test(field.value) || field.value === '';
        
        // Eliminar mensaje de error previo si existe
        const existingError = field.parentNode.querySelector('.invalid-feedback');
        if (existingError) {
            existingError.remove();
        }
        
        // Actualizar clases de validación
        field.classList.toggle('is-invalid', !isValid);
        field.classList.toggle('is-valid', isValid && field.value !== '');
        
        // Si no es válido, mostrar mensaje de error
        if (!isValid) {
            const errorDiv = document.createElement('div');
            errorDiv.className = 'invalid-feedback';
            errorDiv.textContent = errorMessage;
            field.parentNode.appendChild(errorDiv);
        }
    }
    
    // Validación al enviar el formulario
        registerForm.addEventListener('submit', function(event) {
        let isValid = true;
        
        // Validar que el email contiene @
        if (email && !email.value.includes('@')) {
            event.preventDefault();
            isValid = false;
            
            // Mostrar error específico de email
            if (!email.classList.contains('is-invalid')) {
                email.classList.add('is-invalid');
                const errorDiv = document.createElement('div');
                errorDiv.className = 'invalid-feedback';
                errorDiv.textContent = 'El correo electrónico debe contener @';
                email.parentNode.appendChild(errorDiv);
            }
        }
        
        // Validar formato de nombre de usuario
        if (username && !/^[a-zA-Z0-9_.-]+$/.test(username.value)) {
            event.preventDefault();
            isValid = false;
            
            // Mostrar error específico de username si no existe ya
            if (!username.classList.contains('is-invalid')) {
                username.classList.add('is-invalid');
                const errorDiv = document.createElement('div');
                errorDiv.className = 'invalid-feedback';
                errorDiv.textContent = 'El nombre de usuario solo puede contener letras, números, guiones, puntos y guiones bajos';
                username.parentNode.appendChild(errorDiv);
            }
        }
        
        // Validar que las contraseñas coinciden
        if (password && confirmPassword && password.value !== confirmPassword.value) {
            event.preventDefault();
            isValid = false;
            
            // Mostrar error
            confirmPassword.classList.add('is-invalid');
            
            if (!confirmPassword.parentNode.querySelector('.invalid-feedback')) {
                const errorDiv = document.createElement('div');
                errorDiv.className = 'invalid-feedback';
                errorDiv.textContent = 'Las contraseñas no coinciden';
                confirmPassword.parentNode.appendChild(errorDiv);
            }
        }
        
        if (!isValid) {
            event.preventDefault();
            showToast('Por favor, corrige los errores en el formulario', 'error');
            
            // Hacer scroll suave al primer error
            const firstError = registerForm.querySelector('.is-invalid');
            if (firstError) {
                firstError.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }
    });
}

/**
 * Inicializa funcionalidades en la página de detalle de producto
 */
function initProductDetailPage() {
    // Botón para añadir alerta de precio (añadir a la cesta)
    const alertForm = document.querySelector('#priceAlertForm');
    if (!alertForm) return;
    
    // Hacer que la etiqueta de categoría sea clicable
    const categoryBadge = document.querySelector('.badge.bg-primary');
    if (categoryBadge) {
        const categoryName = categoryBadge.textContent.trim();
        
        // Mapeo de nombres de categorías a slugs correctos
        const categoryMapping = {
            'Portátiles': 'portatiles',
            'Tarjetas Gráficas': 'tarjetas-graficas',
            'Auriculares': 'auriculares',
            'Teclados': 'teclados',
            'Monitores': 'monitores',
            'Discos SSD': 'ssd'
        };
        
        // Usar el slug predefinido o generar uno basado en el nombre
        const categorySlug = categoryMapping[categoryName] || 
                            categoryName.toLowerCase()
                                .normalize('NFD').replace(/[\u0300-\u036f]/g, '') // Eliminar acentos
                                .replace(/\s+/g, '-') // Convertir espacios a guiones
                                .replace(/[^\w\-]+/g, ''); // Eliminar caracteres especiales
        
        categoryBadge.style.cursor = 'pointer';
        categoryBadge.title = 'Ver todos los productos en esta categoría';
        categoryBadge.addEventListener('click', function() {
            window.location.href = `/categoria/${categorySlug}`;
        });
    }
    
    // Elementos para el feedback visual
    const alertFeedback = document.getElementById('alertFeedback');
    const alertFeedbackText = document.getElementById('alertFeedbackText');
    const productToast = document.getElementById('productToast');
    const toastMessage = document.getElementById('toastMessage');
    
    // Inicializar el toast de Bootstrap si existe
    let bsToast;
    if (productToast && typeof bootstrap !== 'undefined') {
        bsToast = new bootstrap.Toast(productToast, {
            autohide: true,
            delay: 5000
        });
    }
    
    alertForm.addEventListener('submit', function(event) {
        event.preventDefault(); // Prevenir el envío normal del formulario
        console.log('[MAIN.JS] Formulario de alerta de precio enviado.');
        
        // Determinar si es añadir o actualizar
        const isUpdate = this.querySelector('button[type="submit"]').textContent.includes('Actualizar');
        const formData = new FormData(this);
        console.log('[MAIN.JS] FormData:', formData);
        for (let [key, value] of formData.entries()) {
            console.log(`[MAIN.JS] FormData ${key}:`, value);
        }
        
        // Mostrar indicador de carga inmediatamente
        const submitButton = this.querySelector('button[type="submit"]');
        const originalButtonText = submitButton.textContent;
        submitButton.innerHTML = '<span class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>Procesando...';
        submitButton.disabled = true;
        
        // Enviar el formulario con fetch
        console.log('[MAIN.JS] Iniciando fetch a:', this.action);
        fetch(this.action, {
            method: 'POST',
            body: formData,
            headers: {
                'X-Requested-With': 'XMLHttpRequest'
            }
        })
        .then(response => {
            console.log('[MAIN.JS] Respuesta recibida del fetch:', response);
            console.log('[MAIN.JS] Status de la respuesta:', response.status);
            console.log('[MAIN.JS] StatusText de la respuesta:', response.statusText);
            console.log('[MAIN.JS] Headers de la respuesta - Content-Type:', response.headers.get('Content-Type'));
            
            // Restaurar el botón inmediatamente después de recibir la respuesta, antes de procesarla.
            if (submitButton) {
            submitButton.innerHTML = originalButtonText;
            submitButton.disabled = false;
            }
            
            if (!response.ok) {
                console.error('[MAIN.JS] La respuesta del servidor no fue OK (status >= 400). Status:', response.status);
                // Intentar leer el cuerpo del error como texto si no es OK
                return response.text().then(text => {
                    console.error('[MAIN.JS] Cuerpo del error (texto):', text);
                    throw new Error(`Error del servidor: ${response.status} ${response.statusText}. Cuerpo: ${text}`);
                });
            }
            console.log('[MAIN.JS] Intentando parsear respuesta como JSON...');
            return response.json();
        })
        .then(data => {
            console.log('[MAIN.JS] Datos JSON parseados exitosamente:', data);
            // Determinar si la operación fue exitosa basado en la respuesta
            let success = data.success === true;
            console.log('[MAIN.JS] ¿Operación exitosa según data.success?:', success);
            
            // Definir mensaje según si era actualización o creación
            let message = '';
            if (data.is_update) {
                message = data.message || 'Este producto ya estaba en tu cesta. Precio objetivo actualizado correctamente';
            } else {
                message = data.message || 'Producto añadido a tu cesta correctamente';
            }
            
            // Mostrar mensaje de éxito o error en el formulario
            if (alertFeedback && alertFeedbackText) {
                alertFeedbackText.textContent = message;
                
                // Cambiar la clase del alertFeedback según el resultado
                alertFeedback.classList.remove('alert-success', 'alert-danger');
                alertFeedback.classList.add(success ? 'alert-success' : 'alert-danger');
                alertFeedback.classList.remove('d-none');
                
                // Animar el mensaje
                alertFeedback.style.opacity = '0';
                alertFeedback.style.transform = 'translateY(10px)';
                
                setTimeout(() => {
                    alertFeedback.style.transition = 'all 0.3s ease';
                    alertFeedback.style.opacity = '1';
                    alertFeedback.style.transform = 'translateY(0)';
                }, 10);
                
                // Ocultar después de 5 segundos
                setTimeout(() => {
                    alertFeedback.style.opacity = '0';
                    alertFeedback.style.transform = 'translateY(-10px)';
                    
                    setTimeout(() => {
                        alertFeedback.classList.add('d-none');
                        alertFeedback.removeAttribute('style');
                    }, 300);
                }, 5000);
            }
            
            // Mostrar toast de notificación
            if (bsToast && toastMessage) {
                toastMessage.textContent = message;
                
                // Cambiar el estilo del toast según el resultado
                const toastHeader = productToast.querySelector('.toast-header');
                if (toastHeader) {
                    toastHeader.classList.remove('bg-success', 'bg-danger');
                    toastHeader.classList.add(success ? 'bg-success' : 'bg-danger');
                    
                    // Cambiar el icono según el resultado
                    const toastIcon = toastHeader.querySelector('i');
                    if (toastIcon) {
                        toastIcon.className = success ? 'bi bi-cart-check-fill me-2' : 'bi bi-exclamation-triangle-fill me-2';
                    }
                }
                
                bsToast.show();
            }
            
            // Si fue exitoso, animar el icono del carrito y actualizar el botón
            if (success) {
                // Animar el icono del carrito si es una adición nueva
                if (!isUpdate) {
                    animateCartIcon();
                }
                
                // Actualizar el botón para reflejar que ahora sería una actualización
                if (!isUpdate) {
                    if (submitButton) {
                        submitButton.textContent = 'Actualizar precio';
                    }
                }
            }
        })
        .catch(error => {
            console.error('[MAIN.JS] Error en fetch o en el procesamiento .then():', error);
            console.error('[MAIN.JS] Detalle del error:', error.message, error.stack);
            
            // Restaurar el botón si no se hizo antes
            submitButton.innerHTML = originalButtonText;
            submitButton.disabled = false;
            
            // Mostrar mensaje de error
            const errorMessage = 'Error al procesar la solicitud. Inténtalo de nuevo.';
            
            if (alertFeedback && alertFeedbackText) {
                alertFeedback.classList.remove('alert-success');
                alertFeedback.classList.add('alert-danger');
                alertFeedbackText.textContent = errorMessage;
                alertFeedback.classList.remove('d-none');
                
                // Animar el mensaje
                alertFeedback.style.opacity = '0';
                alertFeedback.style.transform = 'translateY(10px)';
                
                setTimeout(() => {
                    alertFeedback.style.transition = 'all 0.3s ease';
                    alertFeedback.style.opacity = '1';
                    alertFeedback.style.transform = 'translateY(0)';
                }, 10);
                
                // Ocultar después de 5 segundos
                setTimeout(() => {
                    alertFeedback.style.opacity = '0';
                    alertFeedback.style.transform = 'translateY(-10px)';
                    
                    setTimeout(() => {
                        alertFeedback.classList.add('d-none');
                        alertFeedback.removeAttribute('style');
                        alertFeedback.classList.remove('alert-danger');
                        alertFeedback.classList.add('alert-success');
                    }, 300);
                }, 5000);
            }
            
            // Mostrar toast de error
            if (bsToast && toastMessage) {
                toastMessage.textContent = errorMessage;
                const toastHeader = productToast.querySelector('.toast-header');
                if (toastHeader) {
                    toastHeader.classList.remove('bg-success');
                    toastHeader.classList.add('bg-danger');
                    
                    // Cambiar el icono
                    const toastIcon = toastHeader.querySelector('i');
                    if (toastIcon) {
                        toastIcon.className = 'bi bi-exclamation-triangle-fill me-2';
                    }
                }
                bsToast.show();
            }
        });
    });
}

/**
 * Inicializa funcionalidades de la página de watchlist (Mi Cesta)
 */
function initWatchlistPage() {
    // Aplicar animación de entrada escalonada a los elementos de la cesta
    const watchlistItems = document.querySelectorAll('.watchlist-item-container');
    if (watchlistItems.length > 0) {
        watchlistItems.forEach((item, index) => {
            // Asignar delay escalonado para cada elemento
            item.style.animationDelay = `${index * 0.1}s`;
        });
    }

    // Mejorar la experiencia de usuario en la cesta vacía
    const emptyWatchlist = document.querySelector('.watchlist-empty');
    if (emptyWatchlist) {
        emptyWatchlist.style.opacity = '0';
        emptyWatchlist.style.transform = 'translateY(20px)';
        
        setTimeout(() => {
            emptyWatchlist.style.transition = 'all 0.5s ease';
            emptyWatchlist.style.opacity = '1';
            emptyWatchlist.style.transform = 'translateY(0)';
        }, 200);
    }

    // Animación suave para el contador de productos
    const watchlistCount = document.querySelector('.watchlist-count');
    if (watchlistCount) {
        watchlistCount.style.opacity = '0';
        watchlistCount.style.transform = 'scale(0.8)';
        
        setTimeout(() => {
            watchlistCount.style.transition = 'all 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275)';
            watchlistCount.style.opacity = '1';
            watchlistCount.style.transform = 'scale(1)';
        }, 300);
    }
    
    // Popovers mejorados para editar precio objetivo
    const editPriceBtns = document.querySelectorAll('.edit-price-btn');
    
    // Cerrar todos los popovers primero
    function closeAllPopovers() {
        document.querySelectorAll('.edit-price-popover').forEach(popover => {
            if (popover.classList.contains('show')) {
                // Añadir animación de salida
                popover.style.opacity = '0';
                popover.style.transform = 'scale(0.8)';
                
                // Eliminar la clase show después de la animación
                setTimeout(() => {
                    popover.classList.remove('show');
                    // Resetear estilos para la próxima apertura
                    popover.style = '';
                }, 200);
            }
        });
    }
    
    // Cerrar popovers al hacer clic fuera
    document.addEventListener('click', function(e) {
        if (!e.target.closest('.edit-price-container')) {
            closeAllPopovers();
        }
    });
    
    // Configurar cada botón de edición
    editPriceBtns.forEach(btn => {
        const alertId = btn.getAttribute('data-alert-id');
        const targetPrice = btn.getAttribute('data-price');
        const popover = document.getElementById(`edit-popover-${alertId}`);
        const closeBtn = popover.querySelector('.btn-close');
        
        // Abrir popover al hacer clic en el botón de editar
        btn.addEventListener('click', function(e) {
            e.preventDefault();
            e.stopPropagation();
            
            // Si ya está abierto, cerrarlo
            if (popover.classList.contains('show')) {
                closeAllPopovers();
                return;
            }
            
            // Cerrar otros popovers abiertos
            closeAllPopovers();
            
            // Mostrar este popover con una animación suave
            popover.classList.add('show');
            popover.style.display = 'block';
            popover.style.opacity = '0';
            popover.style.transform = 'scale(0.8)';
            
            // Aplicar animación de entrada
            setTimeout(() => {
                popover.style.transition = 'all 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275)';
                popover.style.opacity = '1';
                popover.style.transform = 'scale(1)';
                
                // Enfocar automáticamente el campo de precio
                popover.querySelector('input[name="price"]').focus();
                popover.querySelector('input[name="price"]').select();
            }, 50);
            
            // Añadir efecto de pulsación al botón
            btn.classList.add('pulse');
            setTimeout(() => btn.classList.remove('pulse'), 1000);
        });
        
        // Cerrar al hacer clic en el botón de cerrar
        if (closeBtn) {
            closeBtn.addEventListener('click', function(e) {
                e.preventDefault();
                closeAllPopovers();
            });
        }
        
        // Validación de entrada de precio en tiempo real
        const priceInput = popover.querySelector('input[name="price"]');
        if (priceInput) {
            priceInput.addEventListener('input', function() {
                // Validar que el precio sea un número válido
                const price = parseFloat(this.value);
                const submitBtn = this.closest('form').querySelector('button[type="submit"]');
                
                if (isNaN(price) || price <= 0) {
                    submitBtn.disabled = true;
                    this.classList.add('is-invalid');
                } else {
                    submitBtn.disabled = false;
                    this.classList.remove('is-invalid');
                    
                    // Cambiar el color del botón si el precio es diferente del actual
                    if (price !== parseFloat(targetPrice)) {
                        submitBtn.classList.remove('btn-success');
                        submitBtn.classList.add('btn-primary');
                    } else {
                        submitBtn.classList.remove('btn-primary');
                        submitBtn.classList.add('btn-success');
                    }
                }
            });
        }
        
        // Manejar envío del formulario dentro del popover
        const form = popover.querySelector('form');
        form.addEventListener('submit', function(e) {
            const priceInput = this.querySelector('input[name="price"]');
            const price = parseFloat(priceInput.value);
            
            if (isNaN(price) || price <= 0) {
                e.preventDefault();
                priceInput.classList.add('is-invalid');
                return;
            }
            
            // Añadir efecto de carga al botón de submit
            const submitBtn = this.querySelector('button[type="submit"]');
            submitBtn.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>';
            submitBtn.disabled = true;
            
            // Cerrar el popover con una animación suave
            closeAllPopovers();
        });
    });
    
    // Animación mejorada en barras de progreso
    const progressBars = document.querySelectorAll('.progress-bar');
    if (progressBars.length === 0) return;
    
    progressBars.forEach((bar, index) => {
        const finalWidth = bar.style.width || bar.classList.contains('w-100') ? '100%' : 
                          bar.classList.contains('w-75') ? '75%' : 
                          bar.classList.contains('w-50') ? '50%' : 
                          bar.classList.contains('w-25') ? '25%' : '10%';
        
        bar.style.width = '0';
        
        // Animar la barra progresivamente con delay escalonado
        setTimeout(() => {
            bar.style.transition = 'width 1s cubic-bezier(0.165, 0.84, 0.44, 1)';
            bar.style.width = finalWidth;
        }, 300 + (index * 150)); // Delay escalonado por barra
    });
    
    // Mejorar los efectos de hover en los botones de acción
    const actionButtons = document.querySelectorAll('.edit-price-btn, .delete-btn');
    actionButtons.forEach(btn => {
        btn.addEventListener('mouseenter', function() {
            this.style.transition = 'all 0.3s ease';
        });
    });
    
    // Mejorar la animación de eliminación
    const deleteButtons = document.querySelectorAll('.delete-btn');
    deleteButtons.forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.preventDefault();
            const href = this.getAttribute('href');
            
            Swal.fire({
                title: '¿Estás seguro?',
                text: "¿Quieres eliminar este producto de tu cesta?",
                icon: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#dc3545',
                cancelButtonColor: '#6c757d',
                confirmButtonText: 'Sí, eliminar',
                cancelButtonText: 'Cancelar'
            }).then((result) => {
                if (result.isConfirmed) {
                    const card = this.closest('.watchlist-card');
                    if (card) {
                        card.style.transition = 'all 0.5s ease';
                        card.style.transform = 'translateX(100px) rotate(5deg)';
                        card.style.opacity = '0';
                        
                        setTimeout(() => {
                            window.location.href = href;
                        }, 500);
                    } else {
                        window.location.href = href;
                    }
                }
            });
        });
    });
}

/**
 * Inicializa la página de perfil
 */
function initProfilePage() {
    const editProfileButton = document.querySelector('a[href="/editar-perfil"]');
    const changePasswordButton = document.querySelector('a[href="/cambiar-contrasena"]');
    
    // Añadir iconos a los botones si existen
    if (editProfileButton) {
        editProfileButton.innerHTML = '<i class="bi bi-person-gear"></i> ' + editProfileButton.innerHTML;
    }
    
    if (changePasswordButton) {
        changePasswordButton.innerHTML = '<i class="bi bi-key"></i> ' + changePasswordButton.innerHTML;
    }
    
    // Mejorar la interactividad del formulario de cambio de contraseña
    const changePasswordForm = document.getElementById('changePasswordForm');
    if (changePasswordForm) {
        // Mejorar validación de contraseñas
        const newPasswordInput = changePasswordForm.querySelector('#new_password');
        const confirmPasswordInput = changePasswordForm.querySelector('#confirm_password');
        
        if (newPasswordInput && confirmPasswordInput) {
            // Validar coincidencia de contraseñas en tiempo real
            confirmPasswordInput.addEventListener('input', function() {
                const isMatch = this.value === newPasswordInput.value;
                if (this.value.length > 0) {
                    if (isMatch) {
                        this.classList.remove('is-invalid');
                        this.classList.add('is-valid');
                    } else {
                        this.classList.remove('is-valid');
                        this.classList.add('is-invalid');
                    }
                } else {
                    this.classList.remove('is-valid', 'is-invalid');
                }
            });
            
            // También validar cuando cambias la contraseña nueva
            newPasswordInput.addEventListener('input', function() {
                if (confirmPasswordInput.value.length > 0) {
                    const isMatch = this.value === confirmPasswordInput.value;
                    if (isMatch) {
                        confirmPasswordInput.classList.remove('is-invalid');
                        confirmPasswordInput.classList.add('is-valid');
                    } else {
                        confirmPasswordInput.classList.remove('is-valid');
                        confirmPasswordInput.classList.add('is-invalid');
                    }
                }
                
                // Validar la fuerza de la contraseña
                if (this.value.length > 0) {
                    const strength = checkPasswordStrength(this.value);
                    // Actualizar estilo según la fuerza
                    if (strength === 'strong') {
                        this.classList.remove('is-invalid', 'is-warning');
                        this.classList.add('is-valid');
                    } else if (strength === 'medium') {
                        this.classList.remove('is-invalid', 'is-valid');
                        this.classList.add('is-warning'); // Clase personalizada
                    } else {
                        this.classList.remove('is-valid', 'is-warning');
                        this.classList.add('is-invalid');
                    }
                } else {
                    this.classList.remove('is-valid', 'is-invalid', 'is-warning');
                }
            });
        }
        
        // Validación en envío
        changePasswordForm.addEventListener('submit', function(e) {
            if (newPasswordInput && confirmPasswordInput) {
                if (newPasswordInput.value !== confirmPasswordInput.value) {
                    e.preventDefault();
                    showToast('Las contraseñas no coinciden', 'danger');
                    confirmPasswordInput.classList.add('is-invalid');
                }
            }
            
            // Efecto visual al enviar
            const submitBtn = this.querySelector('button[type="submit"]');
            if (submitBtn) {
                submitBtn.classList.add('btn-pulsing');
            }
        });
    }
    
    // Efecto de revelación para las secciones del perfil
    document.querySelectorAll('.profile-section').forEach((section, index) => {
        section.style.opacity = '0';
        section.style.transform = 'translateY(20px)';
        section.style.transition = 'opacity 0.5s ease, transform 0.5s ease';
        
        // Aplicar efecto escalonado
        setTimeout(() => {
            section.style.opacity = '1';
            section.style.transform = 'translateY(0)';
        }, 100 + (index * 150)); // Cada sección aparece con un pequeño retraso
    });
    
    // Efecto hover para los botones
    document.querySelectorAll('.profile-card button').forEach(button => {
        button.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-2px)';
        });
        
        button.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0)';
        });
    });
    
    // Efecto visual para campos de contraseña
    document.querySelectorAll('.password-input').forEach(input => {
        // Crear botón para mostrar/ocultar contraseña
        const toggleBtn = document.createElement('button');
        toggleBtn.type = 'button';
        toggleBtn.className = 'password-toggle';
        toggleBtn.innerHTML = '<i class="bi bi-eye"></i>';
        toggleBtn.style.position = 'absolute';
        toggleBtn.style.right = '10px';
        toggleBtn.style.top = '50%';
        toggleBtn.style.transform = 'translateY(-50%)';
        toggleBtn.style.border = 'none';
        toggleBtn.style.background = 'transparent';
        toggleBtn.style.cursor = 'pointer';
        toggleBtn.style.zIndex = '5';
        
        // Posicionar relativamente al contenedor
        const inputGroup = input.closest('.input-group');
        if (inputGroup) {
            inputGroup.style.position = 'relative';
            inputGroup.appendChild(toggleBtn);
        }
        
        // Manejar el evento click para mostrar/ocultar contraseña
        toggleBtn.addEventListener('click', function() {
            const type = input.getAttribute('type') === 'password' ? 'text' : 'password';
            input.setAttribute('type', type);
            this.innerHTML = type === 'password' ? '<i class="bi bi-eye"></i>' : '<i class="bi bi-eye-slash"></i>';
        });
    });
    
    // Configurar el formulario de configuración para usar AJAX
    const configForm = document.getElementById('configForm');
    if (configForm) {
        configForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            // Obtener el estado del checkbox
            const emailNotifications = this.querySelector('input[name="email_notifications"]').checked;
            
            // Crear datos del formulario
            const formData = new FormData();
            formData.append('email_notifications', emailNotifications ? 'on' : 'off');
            
            // Mostrar indicador de carga
            const submitButton = this.querySelector('button[type="submit"]');
            const originalText = submitButton.innerHTML;
            submitButton.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Guardando...';
            submitButton.disabled = true;
            
            // Enviar petición AJAX
            fetch('/actualizar-configuracion', {
                method: 'POST',
                body: formData,
                headers: {
                    'X-Requested-With': 'XMLHttpRequest'
                }
            })
            .then(response => response.json())
            .then(data => {
                // Restaurar botón
                submitButton.innerHTML = originalText;
                submitButton.disabled = false;
                
                if (data.success) {
                    // Mostrar mensaje de éxito
                    showToast('Configuración guardada correctamente', 'success');
                    
                    // También añadir un mensaje visible en el formulario
                    let successMsg = configForm.closest('.card-body').querySelector('.alert-success');
                    if (!successMsg) {
                        successMsg = document.createElement('div');
                        successMsg.className = 'alert alert-success mt-3';
                        configForm.insertAdjacentElement('afterend', successMsg);
                    }
                    successMsg.textContent = 'Configuración actualizada correctamente';
                    
                    // Ocultar después de unos segundos
                    setTimeout(() => {
                        successMsg.style.opacity = '0';
                        successMsg.style.transition = 'opacity 0.5s';
                        setTimeout(() => {
                            successMsg.remove();
                        }, 500);
                    }, 3000);
                } else {
                    // Mostrar error
                    showToast(data.error || 'Error al guardar la configuración', 'danger');
                }
            })
            .catch(error => {
                // Restaurar botón
                submitButton.innerHTML = originalText;
                submitButton.disabled = false;
                
                // Mostrar error
                showToast('Error de conexión al guardar la configuración', 'danger');
                console.error('Error:', error);
            });
        });
    }
}

/**
 * Verifica la fuerza de una contraseña
 * @param {string} password - La contraseña a verificar
 * @returns {string} - 'weak', 'medium', o 'strong'
 */
function checkPasswordStrength(password) {
    const hasLetter = /[a-zA-Z]/.test(password);
    const hasNumber = /[0-9]/.test(password);
    const hasSpecial = /[^a-zA-Z0-9]/.test(password);
    
    if (password.length < 6) {
        return 'weak';
    } else if (password.length >= 8 && hasLetter && hasNumber && hasSpecial) {
        return 'strong';
    } else if (password.length >= 6 && ((hasLetter && hasNumber) || (hasLetter && hasSpecial) || (hasNumber && hasSpecial))) {
        return 'medium';
    } else {
        return 'weak';
    }
}

/**
 * Inicializa la página de notificaciones
 */
function initNotificationsPage() {
    const markAllAsReadButton = document.querySelector('form[action="/notificaciones/marcar-leidas"] button');
    const notificationItems = document.querySelectorAll('.list-group-item');
    
    if (markAllAsReadButton) {
        // Solo añadimos el icono al botón de marcar todas como leídas
        markAllAsReadButton.innerHTML = '<i class="bi bi-check-all me-1"></i>' + markAllAsReadButton.textContent;
    }
    
    // Mejorar visualización de notificaciones leídas/no leídas
    notificationItems.forEach(item => {
        const markAsReadBtn = item.querySelector('.mark-read-btn');
        if (markAsReadBtn) {
            markAsReadBtn.innerHTML = '<i class="bi bi-check2 me-1"></i>' + markAsReadBtn.textContent;
        }
    });
}

/**
 * Inicializa la página de cambio de contraseña
 */
function initChangePasswordPage() {
    const changePasswordForm = document.querySelector('form[action="/cambiar-contrasena"]');
    if (!changePasswordForm) return;
    
    changePasswordForm.addEventListener('submit', function(e) {
        const newPassword = document.getElementById('new_password');
        const confirmPassword = document.getElementById('confirm_password');
        
        if (newPassword && confirmPassword && newPassword.value !== confirmPassword.value) {
            e.preventDefault();
            
            const errorMsg = document.createElement('div');
            errorMsg.className = 'alert alert-danger';
            errorMsg.textContent = 'Las contraseñas nuevas no coinciden';
            
            const prevError = changePasswordForm.querySelector('.alert');
            if (prevError) prevError.remove();
            
            changePasswordForm.prepend(errorMsg);
        }
    });
}

/**
 * Inicializa la página de edición de perfil
 */
function initEditProfilePage() {
    const editProfileForm = document.querySelector('form[action="/editar-perfil"]');
    if (!editProfileForm) return;
    
    // Código para la funcionalidad del formulario de edición de perfil
}

/**
 * Crea efecto ripple al hacer click en un botón
 */
function createRippleEffect(event) {
    const button = event.currentTarget;
    const ripple = document.createElement('span');
    
    const diameter = Math.max(button.clientWidth, button.clientHeight);
    const radius = diameter / 2;
    
    ripple.style.width = ripple.style.height = `${diameter}px`;
    ripple.style.left = `${event.clientX - (button.getBoundingClientRect().left + radius)}px`;
    ripple.style.top = `${event.clientY - (button.getBoundingClientRect().top + radius)}px`;
    ripple.className = 'ripple';
    
    // Eliminar el ripple anterior si existe
    const prevRipple = button.querySelector('.ripple');
    if (prevRipple) {
        prevRipple.remove();
    }
    
    button.appendChild(ripple);
    
    // Eliminar el ripple después de la animación
    setTimeout(() => {
        ripple.remove();
    }, 800);
}

/**
 * Anima el icono del carrito al añadir un producto
 */
function animateCartIcon() {
    const cartLink = document.getElementById('cart-link');
    if (!cartLink) return;
    
    cartLink.classList.add('cart-bounce');
    
    // Mostrar notificación bonita usando showToast que ya está definida en este archivo
    showToast('Producto añadido a tu cesta');
    
    // Eliminar clases de animación
    setTimeout(() => {
        cartLink.classList.remove('cart-bounce');
        const animatedCount = cartLink.querySelector('.cart-count-animate');
        if (animatedCount) {
            animatedCount.classList.remove('cart-count-animate');
        }
    }, 1000);
}

/**
 * Inicializa la alerta para productos duplicados en la cesta
 */
function initDuplicateProductAlert() {
    const alertForm = document.querySelector('form[action="/price-alert/set"]');
    if (!alertForm) return;
    
    alertForm.addEventListener('submit', function(event) {
        const productId = this.querySelector('input[name="product_id"]').value;
        const targetPrice = this.querySelector('input[name="target_price"]').value;
        
        // Comprobar si el producto ya está en la cesta por el texto del botón
        const isUpdate = this.querySelector('button[type="submit"]').textContent.includes('Actualizar');
        
        if (isUpdate) {
            // Si es una actualización, mostrar feedback pero permitir el envío
            showToast('Este producto ya estaba en tu cesta. Precio objetivo actualizado correctamente', 'success');
        } else {
            // En vez de comprobar por contador, lo hacemos correctamente con la interfaz
            // También dejamos que se envíe el formulario siempre, solo damos diferente feedback
            showToast('Producto añadido a tu cesta', 'success');
            animateCartIcon();
        }
    });
}

/**
 * Muestra un toast (notificación pequeña en la esquina)
 */
function showToast(message, type = 'success') {
    // Intentar encontrar el toast y el mensaje del DOM que usa initProductDetailPage
    const productToastEl = document.getElementById('productToast');
    const toastMessageEl = document.getElementById('toastMessage');

    if (productToastEl && toastMessageEl) {
        try {
            // Obtener o crear instancia de Bootstrap Toast
            let toastInstance = bootstrap.Toast.getInstance(productToastEl);
            if (!toastInstance) {
                toastInstance = new bootstrap.Toast(productToastEl, {
                    autohide: true,
                    delay: 5000 
                });
            }

            // Configurar mensaje
            toastMessageEl.textContent = message;

            // Configurar tipo (estilo)
            const toastHeader = productToastEl.querySelector('.toast-header');
            if (toastHeader) {
                toastHeader.classList.remove('bg-success', 'bg-danger', 'bg-warning', 'bg-info'); // Limpiar clases previas
                const toastIcon = toastHeader.querySelector('i');

                switch (type) {
                    case 'success':
                        toastHeader.classList.add('bg-success');
                        if (toastIcon) toastIcon.className = 'bi bi-check-circle-fill me-2';
                        break;
                    case 'error': // Alias para danger
                    case 'danger':
                        toastHeader.classList.add('bg-danger');
                        if (toastIcon) toastIcon.className = 'bi bi-exclamation-triangle-fill me-2';
                        break;
                    case 'warning':
                        toastHeader.classList.add('bg-warning');
                        if (toastIcon) toastIcon.className = 'bi bi-exclamation-circle-fill me-2';
                        break;
                    case 'info':
                        toastHeader.classList.add('bg-info');
                        if (toastIcon) toastIcon.className = 'bi bi-info-circle-fill me-2';
                        break;
                    default:
                        toastHeader.classList.add('bg-secondary'); // Un color por defecto
                        if (toastIcon) toastIcon.className = 'bi bi-bell-fill me-2';
                }
            }
            
            toastInstance.show();

        } catch (e) {
            console.error("Error al mostrar el toast:", e);
            // Fallback si el toast de Bootstrap falla, para al menos ver el mensaje en consola
            console.log(`Fallback Toast (${type}): ${message}`);
        }
    } else {
        console.warn("Elementos del Toast (productToast o toastMessage) no encontrados. Mensaje:", message);
        // Fallback si no se encuentran los elementos
        alert(`Notificación (${type}): ${message}`); 
    }
}

// Asegurar que los dropdowns se inicialicen en cada carga de página
window.addEventListener('load', function() {
    setTimeout(initDropdowns, 500); // Intentar inicializar después de que todo esté cargado
});

// Mejorar el comportamiento de los desplegables
document.addEventListener('DOMContentLoaded', function() {
    // Inicializar todos los dropdowns de Bootstrap
    initializeAllDropdowns();
    
    // Implementar una mejor gestión para todos los dropdowns
    function initializeAllDropdowns() {
        // 1. Asegurarse de que todos los dropdowns tengan los atributos correctos
        document.querySelectorAll('.dropdown-toggle').forEach(toggle => {
            if (!toggle.hasAttribute('data-bs-toggle')) {
                toggle.setAttribute('data-bs-toggle', 'dropdown');
            }
        });

        // 2. Inicializar todos los dropdowns manualmente
        if (typeof bootstrap !== 'undefined' && bootstrap.Dropdown) {
            document.querySelectorAll('.dropdown-toggle').forEach(toggle => {
                try {
                    new bootstrap.Dropdown(toggle);
                } catch (e) {
                    console.warn('Error al inicializar dropdown:', e);
                }
            });
        }

        // 3. Inicializar específicamente el dropdown de categorías para todas las páginas
        initializeCategoriesDropdown();
        
        // 4. Mejorar comportamiento: mantener abiertos al hacer hover
        const dropdownElements = document.querySelectorAll('.nav-item.dropdown, .btn-group.dropdown, .btn-group.dropup');
        
        dropdownElements.forEach(dropdown => {
            const toggle = dropdown.querySelector('.dropdown-toggle');
            const menu = dropdown.querySelector('.dropdown-menu');
            
            if (!toggle || !menu) return;
            
            // Cuando el ratón entra en el elemento dropdown
            dropdown.addEventListener('mouseenter', function() {
                if (window.innerWidth < 992) return; // Solo en dispositivos no móviles
                
                // Cerrar otros dropdowns
                document.querySelectorAll('.dropdown-menu.show').forEach(openMenu => {
                    if (openMenu !== menu) {
                        openMenu.classList.remove('show');
                    }
                });
                
                // Mostrar este dropdown
                menu.classList.add('show');
            });
            
            // Mantener abierto y prevenir cierre automático
            dropdown.addEventListener('click', function(e) {
                if (e.target.closest('.dropdown-menu')) {
                    e.stopPropagation(); // Evitar cierre al hacer clic dentro del menú
                }
            });
            
            // Cuando el ratón sale del elemento dropdown
            dropdown.addEventListener('mouseleave', function() {
                if (window.innerWidth < 992) return; // Solo en dispositivos no móviles
                
                // Añadir un pequeño retraso antes de cerrar
                setTimeout(() => {
                    if (!dropdown.matches(':hover')) {
                        menu.classList.remove('show');
                    }
                }, 300);
            });
            
            // Asegurarse de que el clic también funciona (para móviles)
            toggle.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                
                // Usar la API de Bootstrap si está disponible
                if (typeof bootstrap !== 'undefined' && bootstrap.Dropdown) {
                    const instance = bootstrap.Dropdown.getInstance(toggle);
                    if (instance) {
                        if (menu.classList.contains('show')) {
                            instance.hide();
                        } else {
                            instance.show();
                        }
                    } else {
                        menu.classList.toggle('show');
                    }
                } else {
                    // Fallback si Bootstrap no está disponible
                    menu.classList.toggle('show');
                }
            });
            
            // Cerrar al hacer click fuera
            document.addEventListener('click', function(e) {
                // No cerrar si el clic es dentro del dropdown
                if (dropdown.contains(e.target)) {
                    return;
                }
                
                // Cerrar este dropdown
                menu.classList.remove('show');
            });
        });
    }
    
    // Inicializar específicamente el dropdown de categorías para que funcione en todas las páginas
    function initializeCategoriesDropdown() {
        const categoriesDropdown = document.querySelector('.nav-item.dropdown');
        if (!categoriesDropdown) return;
        
        const categoriesToggle = categoriesDropdown.querySelector('.dropdown-toggle');
        const categoriesMenu = categoriesDropdown.querySelector('.dropdown-menu');
        
        if (!categoriesToggle || !categoriesMenu) return;
        
        // Añadir clase para identificar este dropdown específicamente
        categoriesDropdown.classList.add('categories-dropdown');
        
        // Hacer que el toggle muestre el menú al hacer clic
        categoriesToggle.addEventListener('click', function(e) {
            e.preventDefault();
            e.stopPropagation();
            
            if (typeof bootstrap !== 'undefined' && bootstrap.Dropdown) {
                const instance = bootstrap.Dropdown.getInstance(categoriesToggle);
                if (instance) {
                    if (categoriesMenu.classList.contains('show')) {
                        instance.hide();
                    } else {
                        instance.show();
                    }
                } else {
                    categoriesMenu.classList.toggle('show');
                }
            } else {
                categoriesMenu.classList.toggle('show');
            }
        });
        
        // Mantener abierto al hacer hover
        categoriesDropdown.addEventListener('mouseenter', function() {
            if (window.innerWidth >= 992) { // Solo en escritorio
                categoriesMenu.classList.add('show');
            }
        });
        
        categoriesDropdown.addEventListener('mouseleave', function() {
            if (window.innerWidth >= 992) { // Solo en escritorio
                setTimeout(() => {
                    if (!categoriesDropdown.matches(':hover')) {
                        categoriesMenu.classList.remove('show');
                    }
                }, 300);
            }
        });
    }
});

// Añadir tooltips informativos
document.addEventListener('DOMContentLoaded', function() {
    // Tooltips para los botones de acción
    tippy('.price-alert-btn', {
        content: 'Añadir alerta de precio',
        placement: 'top'
    });
    
    tippy('.edit-price-btn', {
        content: 'Editar precio objetivo',
        placement: 'top'
    });
    
    tippy('.delete-btn', {
        content: 'Eliminar de mi cesta',
        placement: 'top'
    });
    
    // Tooltips para los badges de tienda
    tippy('.store-badge', {
        content: 'Ver todos los productos de esta tienda',
        placement: 'top'
    });
    
    // Animaciones al hacer scroll
    document.querySelectorAll('.product-card').forEach((card, index) => {
        card.setAttribute('data-aos', 'fade-up');
        card.setAttribute('data-aos-delay', (index % 3) * 100);
    });
});

/**
 * Configura el manejo de visibilidad de la página para evitar la pantalla en blanco
 * cuando el usuario regresa a la pestaña después de un tiempo
 */
function setupPageVisibilityHandling() {
    // Detectar cuando la página pierde o recupera visibilidad
    document.addEventListener('visibilitychange', function() {
        if (document.visibilityState === 'visible') {
            // La página vuelve a ser visible
            console.log('Página visible de nuevo - verificando estado');
            
            // Verificar si el contenido principal está visible
            const mainContent = document.querySelector('main.content-area');
            if (mainContent && mainContent.offsetHeight === 0) {
                console.log('Contenido principal no visible - recargando página');
                window.location.reload();
            }
            
            // Verificar si hay productos visibles (si estamos en una página de productos)
            const productsContainer = document.querySelector('.products-container');
            if (productsContainer && productsContainer.children.length > 0) {
                let visibleProducts = false;
                productsContainer.querySelectorAll('.product-item').forEach(item => {
                    if (item.offsetHeight > 0 && !item.classList.contains('d-none')) {
                        visibleProducts = true;
                    }
                });
                
                if (!visibleProducts) {
                    console.log('Productos no visibles - recargando página');
                    window.location.reload();
                }
            }
            
            // Verificar si hay elementos críticos visibles
            const criticalElements = [
                document.querySelector('.navbar'),
                document.querySelector('.site-footer'),
                document.querySelector('.products-container'),
                document.querySelector('.product-detail'),
                document.querySelector('.watchlist-container'),
                document.querySelector('.profile-container')
            ];
            
            let visibleCriticalElements = false;
            criticalElements.forEach(element => {
                if (element && element.offsetHeight > 0) {
                    visibleCriticalElements = true;
                }
            });
            
            if (!visibleCriticalElements) {
                console.log('Elementos críticos no visibles - recargando página');
                window.location.reload();
            }
        }
    });
    
    // También verificar periódicamente el estado de la página
    // Esto es un respaldo adicional al evento visibilitychange
    setInterval(function() {
        // Solo verificar si la página está visible
        if (document.visibilityState === 'visible') {
            const mainContent = document.querySelector('main.content-area');
            const navbar = document.querySelector('.navbar');
            
            // Si el contenido principal o la barra de navegación no son visibles, recargar
            if ((mainContent && mainContent.offsetHeight === 0) || 
                (navbar && navbar.offsetHeight === 0)) {
                console.log('Verificación periódica: contenido no visible - recargando página');
                window.location.reload();
            }
        }
    }, 30000); // Verificar cada 30 segundos
} 