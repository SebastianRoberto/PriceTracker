package handler

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"app/internal/domain/model"
	"app/internal/interface/web/views"
	"app/internal/usecase"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthHandler maneja las peticiones de autenticación
type AuthHandler struct {
	userUseCase      *usecase.UserUseCase
	templateRenderer *views.TemplateRenderer
}

// NewAuthHandler crea una nueva instancia del AuthHandler
func NewAuthHandler(userUseCase *usecase.UserUseCase, templateRenderer *views.TemplateRenderer) *AuthHandler {
	return &AuthHandler{
		userUseCase:      userUseCase,
		templateRenderer: templateRenderer,
	}
}

// ShowLoginForm muestra el formulario de inicio de sesión
func (h *AuthHandler) ShowLoginForm(c *gin.Context) {
	// Obtener categorías del contexto
	categories, _ := c.Get("allCategories")

	h.templateRenderer.Render(c, http.StatusOK, "login.html", gin.H{
		"Title":      "Iniciar Sesión - Comparador de Precios",
		"Categories": categories,
	})
}

// Login procesa el formulario de inicio de sesión
func (h *AuthHandler) Login(c *gin.Context) {
	// Obtener categorías del contexto
	categories, _ := c.Get("allCategories")

	// Obtener los datos del formulario
	var form struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "login.html", gin.H{
			"Title":      "Iniciar Sesión - Comparador de Precios",
			"Error":      "Por favor, completa todos los campos correctamente",
			"Email":      form.Email,
			"Categories": categories,
		})
		return
	}

	// Verificar las credenciales
	ctx := c.Request.Context()
	user, err := h.userUseCase.AuthenticateUser(ctx, form.Email, form.Password)
	if err != nil {
		h.templateRenderer.Render(c, http.StatusUnauthorized, "login.html", gin.H{
			"Title": "Iniciar Sesión - Comparador de Precios",
			"Error": "Email o contraseña incorrectos",
			"Email": form.Email,
		})
		return
	}

	// Verificar si el usuario está verificado
	if !user.Verified {
		h.templateRenderer.Render(c, http.StatusUnauthorized, "login.html", gin.H{
			"Title": "Iniciar Sesión - Comparador de Precios",
			"Error": "Tu cuenta no ha sido verificada. Por favor, revisa tu correo electrónico.",
			"Email": form.Email,
		})
		return
	}

	// Iniciar sesión
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	// Redireccionar a la página principal
	c.Redirect(http.StatusFound, "/")
}

// ShowRegisterForm muestra el formulario de registro
func (h *AuthHandler) ShowRegisterForm(c *gin.Context) {
	// Obtener categorías del contexto
	categories, _ := c.Get("allCategories")

	h.templateRenderer.Render(c, http.StatusOK, "register.html", gin.H{
		"Title":      "Registro - Comparador de Precios",
		"Categories": categories,
	})
}

// RegisterHandler maneja la solicitud POST de registro
func (h *AuthHandler) RegisterHandler(c *gin.Context) {
	var input struct {
		Username        string `form:"username" binding:"required"`
		Email           string `form:"email" binding:"required,email"`
		Password        string `form:"password" binding:"required,min=6"`
		ConfirmPassword string `form:"confirm_password" binding:"required,eqfield=Password"`
		NotifyByEmail   bool   `form:"notify_by_email"`
	}

	if err := c.ShouldBind(&input); err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "register.html", gin.H{
			"Error":       "Error en el formulario de registro",
			"ErrorDetail": err.Error(),
		})
		return
	}

	// Mostrar mensaje de carga mientras se procesa
	c.Writer.Header().Set("HX-Trigger", `{"showLoading": true}`)

	// Crear el usuario usando el caso de uso de autenticación
	ctx := c.Request.Context()
	user, err := h.userUseCase.Register(ctx, input.Username, input.Email, input.Password, input.NotifyByEmail)
	if err != nil {
		// Personalizar el mensaje de error
		errorMsg := "Error al registrar usuario"
		if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "email") {
			errorMsg = "Este correo electrónico ya está registrado"
		} else if strings.Contains(err.Error(), "Duplicate entry") && strings.Contains(err.Error(), "username") {
			errorMsg = "Este nombre de usuario ya está en uso"
		}

		// Enviar el error al cliente
		h.templateRenderer.Render(c, http.StatusBadRequest, "register.html", gin.H{
			"Error":       errorMsg,
			"ErrorDetail": err.Error(),
			"Username":    input.Username,
			"Email":       input.Email,
		})
		return
	}

	// Enviar correo de verificación en una goroutine separada para no bloquear la respuesta
	go func() {
		err = h.userUseCase.SendVerificationEmail(ctx, user)
		if err != nil {
			log.Printf("Error al enviar correo de verificación: %v", err)
		}
	}()

	// Redirigir a la página de registro exitoso
	c.Redirect(http.StatusSeeOther, "/registro-exitoso?email="+url.QueryEscape(user.Email))
}

// VerifyEmail verifica la cuenta de un usuario
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	// Obtener el token de verificación
	token := c.Query("token")
	if token == "" {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Title":   "Error de Verificación - Comparador de Precios",
			"Message": "Token de verificación no proporcionado",
		})
		return
	}

	// Verificar el token
	ctx := c.Request.Context()
	user, err := h.userUseCase.VerifyUser(ctx, token)
	if err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Title":   "Error de Verificación - Comparador de Precios",
			"Message": "Token de verificación inválido o expirado",
			"Error":   err.Error(),
		})
		return
	}

	// Mostrar página de éxito
	h.templateRenderer.Render(c, http.StatusOK, "verify_success.html", gin.H{
		"Title":   "Verificación Exitosa - Comparador de Precios",
		"Message": "Tu cuenta ha sido verificada correctamente. Ahora puedes iniciar sesión.",
		"User":    user,
	})
}

// Logout cierra la sesión del usuario
func (h *AuthHandler) Logout(c *gin.Context) {
	// Eliminar la sesión
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	// Redireccionar a la página principal
	c.Redirect(http.StatusFound, "/")
}

// ShowProfile muestra el perfil del usuario
func (h *AuthHandler) ShowProfile(c *gin.Context) {
	// Obtener el usuario de la sesión
	user, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	// Obtener categorías para el menú
	categories, _ := c.Get("allCategories")

	// Obtener alertas de precio para el contador de Mi Cesta, si el caso de uso está disponible a través de middleware
	priceAlerts, _ := c.Get("priceAlerts")

	// Renderizar la página de perfil
	h.templateRenderer.Render(c, http.StatusOK, "profile.html", gin.H{
		"Title":       "Mi Perfil - Comparador de Precios",
		"User":        user,
		"Categories":  categories,
		"PriceAlerts": priceAlerts,
		"Success":     c.Query("success"),
	})
}

// RequestPasswordReset envía correo con enlace de restablecimiento y redirige al perfil
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	userInterface, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	user := userInterface.(*model.User)

	// Iniciar proceso
	if err := h.userUseCase.InitiatePasswordReset(c.Request.Context(), user.ID); err != nil {
		h.templateRenderer.Render(c, http.StatusInternalServerError, "error.html", gin.H{
			"Message": "No se pudo enviar el correo de restablecimiento",
			"Error":   err.Error(),
		})
		return
	}

	// Redirigir al perfil con mensaje
	c.Redirect(http.StatusFound, "/perfil?success=reset")
}

// ShowPasswordResetForm muestra el formulario para restablecer la contraseña
func (h *AuthHandler) ShowPasswordResetForm(c *gin.Context) {
	// Obtener token de la URL
	token := c.Query("token")
	if token == "" {
		h.templateRenderer.Render(c, http.StatusBadRequest, "error.html", gin.H{
			"Title":   "Error - Comparador de Precios",
			"Message": "Token de restablecimiento no proporcionado",
		})
		return
	}

	// Obtener categorías para el menú
	categories, _ := c.Get("allCategories")

	h.templateRenderer.Render(c, http.StatusOK, "reset_password.html", gin.H{
		"Title":      "Restablecer Contraseña - Comparador de Precios",
		"Token":      token,
		"Categories": categories,
	})
}

// ProcessPasswordReset procesa el formulario de restablecimiento de contraseña
func (h *AuthHandler) ProcessPasswordReset(c *gin.Context) {
	// Obtener categorías del contexto
	categories, _ := c.Get("allCategories")

	// Obtener los datos del formulario
	var form struct {
		Token           string `form:"token" binding:"required"`
		Password        string `form:"password" binding:"required,min=6"`
		ConfirmPassword string `form:"confirm_password" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "reset_password.html", gin.H{
			"Title":      "Restablecer Contraseña - Comparador de Precios",
			"Error":      "Por favor, completa todos los campos correctamente",
			"Token":      form.Token,
			"Categories": categories,
		})
		return
	}

	// Verificar que las contraseñas coinciden
	if form.Password != form.ConfirmPassword {
		h.templateRenderer.Render(c, http.StatusBadRequest, "reset_password.html", gin.H{
			"Title":      "Restablecer Contraseña - Comparador de Precios",
			"Error":      "Las contraseñas no coinciden",
			"Token":      form.Token,
			"Categories": categories,
		})
		return
	}

	// Restablecer la contraseña
	ctx := c.Request.Context()
	user, err := h.userUseCase.ResetPassword(ctx, form.Token, form.Password)
	if err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "reset_password.html", gin.H{
			"Title":      "Restablecer Contraseña - Comparador de Precios",
			"Error":      "No se pudo restablecer la contraseña: " + err.Error(),
			"Token":      form.Token,
			"Categories": categories,
		})
		return
	}

	// Iniciar sesión automáticamente
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Save()

	// Redireccionar al perfil con mensaje de éxito
	c.Redirect(http.StatusFound, "/perfil?success=password_changed")
}

// ChangePassword procesa el formulario de cambio de contraseña
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	log.Printf("[INFO] AuthHandler.ChangePassword - Solicitud recibida")
	// Obtener el usuario de la sesión
	userModel, exists := c.Get("user")
	if !exists {
		log.Printf("[WARN] AuthHandler.ChangePassword - Usuario no encontrado en sesión. Redirigiendo a /login")
		c.Redirect(http.StatusFound, "/login")
		return
	}
	user := userModel.(*model.User)
	log.Printf("[DEBUG] AuthHandler.ChangePassword - Usuario en sesión: ID %d, Username %s", user.ID, user.Username)

	// Obtener categorías para el menú
	categories, _ := c.Get("allCategories")

	// Obtener los datos del formulario
	var form struct {
		CurrentPassword string `form:"old_password" binding:"required"` // Mapea 'old_password' del form a CurrentPassword
		NewPassword     string `form:"new_password" binding:"required,min=6"`
		ConfirmPassword string `form:"confirm_password" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		log.Printf("[ERROR] AuthHandler.ChangePassword - Error al bidear formulario para usuario ID %d: %v. Datos del formulario: old_password presente: %t, new_password presente: %t, confirm_password presente: %t",
			user.ID, err, c.PostForm("old_password") != "", c.PostForm("new_password") != "", c.PostForm("confirm_password") != "")
		h.templateRenderer.Render(c, http.StatusBadRequest, "profile.html", gin.H{
			"Title":      "Mi Perfil - Comparador de Precios",
			"User":       user,
			"Categories": categories,
			"Error":      "Por favor, completa todos los campos correctamente",
		})
		return
	}
	log.Printf("[DEBUG] AuthHandler.ChangePassword - Formulario bindeado para usuario ID %d. CurrentPassword(len): %d, NewPassword(len): %d, ConfirmPassword(len): %d",
		user.ID, len(form.CurrentPassword), len(form.NewPassword), len(form.ConfirmPassword))

	// Verificar que las contraseñas nuevas coinciden
	if form.NewPassword != form.ConfirmPassword {
		log.Printf("[WARN] AuthHandler.ChangePassword - Las nuevas contraseñas no coinciden para usuario ID %d", user.ID)
		h.templateRenderer.Render(c, http.StatusBadRequest, "profile.html", gin.H{
			"Title":      "Mi Perfil - Comparador de Precios",
			"User":       user,
			"Categories": categories,
			"Error":      "Las contraseñas nuevas no coinciden",
		})
		return
	}
	log.Printf("[DEBUG] AuthHandler.ChangePassword - Nuevas contraseñas coinciden para usuario ID %d.", user.ID)

	// Cambiar la contraseña
	ctx := c.Request.Context()
	log.Printf("[INFO] AuthHandler.ChangePassword - Llamando a userUseCase.ChangePassword para usuario ID %d", user.ID)
	err := h.userUseCase.ChangePassword(ctx, user.ID, form.CurrentPassword, form.NewPassword)
	if err != nil {
		log.Printf("[ERROR] AuthHandler.ChangePassword - Error desde userUseCase.ChangePassword para usuario ID %d: %v", user.ID, err)
		h.templateRenderer.Render(c, http.StatusBadRequest, "profile.html", gin.H{
			"Title":      "Mi Perfil - Comparador de Precios",
			"User":       user,
			"Categories": categories,
			"Error":      "No se pudo cambiar la contraseña: " + err.Error(),
		})
		return
	}

	log.Printf("[INFO] AuthHandler.ChangePassword - Contraseña cambiada exitosamente para usuario ID %d. Redirigiendo a /perfil", user.ID)
	// Redireccionar al perfil con mensaje de éxito
	c.Redirect(http.StatusFound, "/perfil?success=password_changed")
}

// DeleteAccount procesa la solicitud de eliminación de cuenta
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	// Obtener el usuario de la sesión
	userModel, exists := c.Get("user")
	if !exists {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	user := userModel.(*model.User)

	// Obtener categorías para el menú
	categories, _ := c.Get("allCategories")

	// Obtener la contraseña del formulario
	var form struct {
		Password string `form:"password_confirm" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "profile.html", gin.H{
			"Title":      "Mi Perfil - Comparador de Precios",
			"User":       user,
			"Categories": categories,
			"Error":      "Por favor, introduce tu contraseña para confirmar",
		})
		return
	}

	// Eliminar la cuenta
	ctx := c.Request.Context()
	err := h.userUseCase.DeleteAccount(ctx, user.ID, form.Password)
	if err != nil {
		h.templateRenderer.Render(c, http.StatusBadRequest, "profile.html", gin.H{
			"Title":      "Mi Perfil - Comparador de Precios",
			"User":       user,
			"Categories": categories,
			"Error":      "No se pudo eliminar la cuenta: " + err.Error(),
		})
		return
	}

	// Cerrar la sesión
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	// Redireccionar a la página de inicio con mensaje
	c.Redirect(http.StatusFound, "/login?message=account_deleted")
}

// ShowForgotPasswordForm muestra el formulario para solicitar el restablecimiento de contraseña
func (h *AuthHandler) ShowForgotPasswordForm(c *gin.Context) {
	categories, _ := c.Get("allCategories")
	h.templateRenderer.Render(c, http.StatusOK, "forgot_password.html", gin.H{
		"Title":      "Restablecer Contraseña - Comparador de Precios",
		"Categories": categories,
	})
}

// ProcessForgotPasswordForm procesa la solicitud de restablecimiento de contraseña
func (h *AuthHandler) ProcessForgotPasswordForm(c *gin.Context) {
	categories, _ := c.Get("allCategories")
	var form struct {
		Email string `form:"email" binding:"required,email"`
	}

	if err := c.ShouldBind(&form); err != nil {
		log.Printf("[WARN] Error de validación en formulario de recuperación de contraseña: %v", err)
		h.templateRenderer.Render(c, http.StatusBadRequest, "forgot_password.html", gin.H{
			"Title":      "Restablecer Contraseña - Comparador de Precios",
			"Error":      "Por favor, introduce una dirección de correo electrónico válida.",
			"Email":      form.Email,
			"Categories": categories,
		})
		return
	}

	log.Printf("[INFO] Procesando solicitud de recuperación de contraseña para email: %s", form.Email)
	ctx := c.Request.Context()
	user, err := h.userUseCase.GetUserByEmail(ctx, form.Email)

	// Por seguridad, no revelamos si el email existe o no, o si está verificado.
	// Simplemente enviamos el correo si las condiciones se cumplen.
	if err == nil && user != nil && user.Verified {
		log.Printf("[INFO] Usuario encontrado y verificado con ID=%d, email=%s", user.ID, user.Email)
		if err := h.userUseCase.InitiatePasswordReset(ctx, user.ID); err != nil {
			// Loguear el error internamente
			log.Printf("[ERROR] Error al iniciar el reseteo de contraseña para ID=%d, email=%s: %v", user.ID, form.Email, err)
		}
	} else {
		// Log para depuración sin exponer información sensible al usuario
		if err != nil {
			log.Printf("[WARN] Usuario no encontrado para email=%s: %v", form.Email, err)
		} else if user == nil {
			log.Printf("[WARN] Usuario no encontrado para email=%s: usuario es nil", form.Email)
		} else if !user.Verified {
			log.Printf("[WARN] Usuario encontrado pero no verificado para email=%s, ID=%d", form.Email, user.ID)
		}
	}

	// Siempre mostrar un mensaje genérico para evitar enumeración de usuarios.
	h.templateRenderer.Render(c, http.StatusOK, "forgot_password.html", gin.H{
		"Title":          "Restablecer Contraseña - Comparador de Precios",
		"SuccessMessage": "Si tu dirección de correo electrónico está registrada y verificada, recibirás un enlace para restablecer tu contraseña.",
		"Categories":     categories,
	})
}

// ShowRegisterSuccessPage muestra la página de éxito después de un registro
func (h *AuthHandler) ShowRegisterSuccessPage(c *gin.Context) {
	// Obtener categorías del contexto
	categories, _ := c.Get("allCategories")

	// Obtener el email del query string
	email := c.Query("email")

	h.templateRenderer.Render(c, http.StatusOK, "register_success.html", gin.H{
		"Title":      "Registro Exitoso - Comparador de Precios",
		"Email":      email,
		"Message":    "Te hemos enviado un correo de verificación. Por favor, revisa tu bandeja de entrada.",
		"Categories": categories,
	})
}
