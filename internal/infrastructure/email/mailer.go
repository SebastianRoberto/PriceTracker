package email

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"path/filepath"
	"runtime"
	"strconv"

	"app/pkg/config"
)

// Mailer es un servicio para enviar correos electrónicos
type Mailer struct {
	smtpHost string
	smtpPort string
	smtpUser string
	smtpPass string
	from     string
	cssStyle string
}

// NewMailer crea una nueva instancia del servicio de correo
func NewMailer() *Mailer {
	// Usar la configuración del archivo config.yaml
	emailConfig := config.Config.Email

	// Cargar el archivo CSS
	_, currentFile, _, _ := runtime.Caller(0)
	cssPath := filepath.Join(filepath.Dir(currentFile), "email.css")
	cssContent, err := ioutil.ReadFile(cssPath)
	if err != nil {
		log.Printf("Error al cargar el archivo CSS para emails: %v", err)
		cssContent = []byte{} // Si hay error, usar un string vacío
	}

	return &Mailer{
		smtpHost: emailConfig.SMTPHost,
		smtpPort: strconv.Itoa(emailConfig.SMTPPort),
		smtpUser: emailConfig.SMTPUser,
		smtpPass: emailConfig.SMTPPass,
		from:     emailConfig.SMTPFrom,
		cssStyle: string(cssContent),
	}
}

// SendVerificationEmail envía un correo de verificación
func (m *Mailer) SendVerificationEmail(to string, token string, username string) error {
	// Construir el asunto y el cuerpo del correo
	subject := "Verificación de cuenta - Comparador de Precios"

	// URL de verificación
	verificationURL := fmt.Sprintf("%s/verificar?token=%s", config.Config.App.URL, token)

	// Construir el cuerpo del correo en formato HTML
	body := fmt.Sprintf(`
		<html>
		<head>
			<style>
				%s
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header header-primary">
					<h2>Verificación de Cuenta</h2>
				</div>
				<div class="content">
					<div class="icon">✉️</div>
					<h3>¡Hola <span class="highlight highlight-primary">%s</span>!</h3>
					<p>Gracias por registrarte en nuestro <b>Comparador de Precios</b>.</p>
					<p>Para activar tu cuenta y comenzar a ahorrar, haz clic en el siguiente botón:</p>
					
					<a href="%s" class="button button-primary">Verificar mi cuenta</a>
					
					<p><small>¿El botón no funciona? Copia y pega este enlace en tu navegador:</small></p>
					<p class="link">%s</p>
					<p><small>Este enlace expirará en 24 horas por seguridad.</small></p>
					<p><small>Si no te has registrado en nuestra plataforma, puedes ignorar este correo.</small></p>
				</div>
				<div class="footer">
					<p>© Comparador de Precios - Ahorra en tus compras online</p>
					<p>Este correo es automático, por favor no lo respondas.</p>
				</div>
			</div>
		</body>
		</html>
	`, m.cssStyle, username, verificationURL, verificationURL)

	// Enviar el correo
	return m.sendMail(to, subject, body)
}

// SendPasswordResetEmail envía un correo con un enlace para restablecer la contraseña
func (m *Mailer) SendPasswordResetEmail(to, token, username string) error {
	// Usar la URL de la aplicación de la configuración
	appURL := config.Config.App.URL

	resetURL := fmt.Sprintf("%s/restablecer-password?token=%s", appURL, token)
	subject := "Restablecimiento de contraseña - Comparador de Precios"

	htmlBody := fmt.Sprintf(`
		<html>
		<head>
			<style>
				%s
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header header-purple">
					<h2>Restablecimiento de Contraseña</h2>
				</div>
				<div class="content">
					<div class="icon">🔐</div>
					<h3>¡Hola <span class="highlight highlight-purple">%s</span>!</h3>
					<p>Hemos recibido una solicitud para restablecer tu contraseña.</p>
					<p>Si fuiste tú, haz clic en el siguiente botón para crear una nueva contraseña:</p>
					
					<a href="%s" class="button button-purple">Restablecer mi contraseña</a>
					
					<p><small>¿El botón no funciona? Copia y pega este enlace en tu navegador:</small></p>
					<p class="link">%s</p>
					<p><small>Este enlace expirará en 24 horas por seguridad.</small></p>
					<p><small>Si no has solicitado el restablecimiento de contraseña, puedes ignorar este mensaje.</small></p>
				</div>
				<div class="footer">
					<p>© Comparador de Precios - Ahorra en tus compras online</p>
					<p>Este correo es automático, por favor no lo respondas.</p>
				</div>
			</div>
		</body>
		</html>
	`, m.cssStyle, username, resetURL, resetURL)

	return m.sendMail(to, subject, htmlBody)
}

// SendPriceAlertEmail envía un correo cuando un producto alcanza el precio objetivo
func (m *Mailer) SendPriceAlertEmail(to string, username string, productName string, productID uint,
	targetPrice float64, currentPrice float64, store string, productURL string) error {

	// Construir el asunto del correo
	subject := fmt.Sprintf("¡Alerta de precio para %s! - Comparador de Precios", productName)

	// URL del producto en nuestra plataforma
	ourProductURL := fmt.Sprintf("%s/producto/%d", config.Config.App.URL, productID)

	// Calcular el porcentaje de ahorro
	savingsPercent := ((targetPrice - currentPrice) / targetPrice) * 100

	// Construir el cuerpo del correo en formato HTML
	body := fmt.Sprintf(`
		<html>
		<head>
			<style>
				%s
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header header-success">
					<h2>¡Alerta de Precio!</h2>
				</div>
				<div class="content">
					<div class="icon">🎉</div>
					<h3>¡Buenas noticias, <span class="highlight highlight-success">%s</span>!</h3>
					<p>El producto que estabas siguiendo ha alcanzado tu precio objetivo.</p>
					
					<div class="product-card">
						<h3>%s</h3>
						<p>
							<span class="price-tag">%.2f€</span>
							<span class="old-price">%.2f€</span>
							<span class="store-badge">%s</span>
						</p>
						<div class="savings">¡Ahorras un %.1f%% (%.2f€)!</div>
					</div>
					
					<div class="button-container">
						<a href="%s" class="button button-success" style="width: 45%%;">
							Ver oferta en %s
						</a>
						<a href="%s" class="button button-primary" style="width: 45%%;">
							Ver en Comparador
						</a>
					</div>
					
					<p><small>Esta alerta de precio se ha activado porque el precio actual del producto es igual o inferior al precio objetivo que configuraste.</small></p>
					<p><small>Puedes gestionar tus alertas de precio en tu perfil dentro de nuestra plataforma.</small></p>
				</div>
				<div class="footer">
					<p>© Comparador de Precios - Ahorra en tus compras online</p>
					<p>Este correo es automático, por favor no lo respondas.</p>
				</div>
			</div>
		</body>
		</html>
	`, m.cssStyle, username, productName, currentPrice, targetPrice, store, savingsPercent, targetPrice-currentPrice, productURL, store, ourProductURL)

	// Enviar el correo
	return m.sendMail(to, subject, body)
}

// sendMail envía un correo electrónico
func (m *Mailer) sendMail(to string, subject string, body string) error {
	// Verificar que la configuración SMTP está completa
	if m.smtpHost == "" || m.smtpPort == "" || m.smtpUser == "" || m.smtpPass == "" {
		log.Printf("[ERROR] Configuración SMTP incompleta - Host: %s, Puerto: %s, Usuario: %s",
			m.smtpHost, m.smtpPort, m.smtpUser)
		return fmt.Errorf("configuración SMTP incompleta - Host: %s, Puerto: %s, Usuario: %s",
			m.smtpHost, m.smtpPort, m.smtpUser)
	}

	// Si no se especificó el remitente, usar el usuario SMTP
	from := m.from
	if from == "" {
		from = m.smtpUser
	}

	// Construir los encabezados del correo
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Construir el mensaje completo
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	log.Printf("[INFO] Preparando envío de correo a %s con asunto '%s' a través de SMTP %s:%s",
		to, subject, m.smtpHost, m.smtpPort)

	// Configurar la autenticación SMTP
	auth := smtp.PlainAuth("", m.smtpUser, m.smtpPass, m.smtpHost)

	// Enviar el correo
	err := smtp.SendMail(
		m.smtpHost+":"+m.smtpPort,
		auth,
		from,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		// Proporcionar más información en caso de error
		log.Printf("[ERROR] Error al enviar correo: %v (Host: %s, Puerto: %s, Usuario: %s)",
			err, m.smtpHost, m.smtpPort, m.smtpUser)
		return fmt.Errorf("error al enviar correo: %v (Host: %s, Puerto: %s, Usuario: %s)",
			err, m.smtpHost, m.smtpPort, m.smtpUser)
	}

	log.Printf("[INFO] Correo enviado exitosamente a %s con asunto '%s'", to, subject)
	return nil
}
