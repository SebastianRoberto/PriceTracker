// Package config proporciona funcionalidades para la configuración de la aplicación
package config

import (
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

var (
	// Config es la instancia global de configuración
	Config *Configuration
)

// Configuration representa la estructura de toda la configuración de la aplicación
type Configuration struct {
	App      AppConfig
	Database DatabaseConfig
	Scraper  ScraperConfig
	Email    EmailConfig
}

// AppConfig contiene la configuración general de la aplicación
type AppConfig struct {
	Name        string
	Environment string
	Port        int
	URL         string
	SessionTTL  int
}

// DatabaseConfig contiene la configuración de la base de datos
type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            int
	Username        string
	Password        string
	Name            string
	Charset         string
	Loc             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// ScraperConfig contiene la configuración para los scrapers
type ScraperConfig struct {
	UpdateInterval time.Duration
	UserAgent      string
	MaxRetries     int
	RetryDelay     time.Duration
}

// EmailConfig contiene la configuración para el servicio de correo electrónico
type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string
}

// InitConfig inicializa la configuración global de la aplicación
func InitConfig() {
	// Establecer las configuraciones por defecto
	viper.SetDefault("app.name", "Comparador de Precios")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.port", 8080)
	viper.SetDefault("app.url", "http://localhost:8080")
	viper.SetDefault("app.session_ttl", 86400) // 24 horas

	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.loc", "Local")
	viper.SetDefault("database.max_idle_conns", 10)
	viper.SetDefault("database.max_open_conns", 100)
	viper.SetDefault("database.conn_max_lifetime", "1h")

	viper.SetDefault("scraper.update_interval", "48h")
	viper.SetDefault("scraper.user_agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	viper.SetDefault("scraper.max_retries", 3)
	viper.SetDefault("scraper.retry_delay", "5s")

	viper.SetDefault("email.smtp_host", "smtp.gmail.com")
	viper.SetDefault("email.smtp_port", 587)
	viper.SetDefault("email.smtp_user", "")
	viper.SetDefault("email.smtp_pass", "")
	viper.SetDefault("email.smtp_from", "")

	// Configurar Viper para leer del archivo
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// Leer el archivo de configuración
	if err := viper.ReadInConfig(); err != nil {
		// Si el archivo no existe, usar valores predeterminados
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error al leer el archivo de configuración: %s", err)
		}
	}

	// Sobrescribir con variables de entorno si existen
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	// Obtener valores de la base de datos, prefiriendo variables de entorno, pero con respaldo del archivo config.yaml
	dbUsername := os.Getenv("DB_USER")
	if dbUsername == "" {
		dbUsername = viper.GetString("database.username")
	}

	dbPassword := os.Getenv("DB_PASS")
	if dbPassword == "" {
		dbPassword = viper.GetString("database.password")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = viper.GetString("database.name")
	}

	// Obtener valores de email, prefiriendo variables de entorno, pero con respaldo del archivo config.yaml
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = viper.GetString("email.smtp_host")
	}

	smtpPort := 0
	if os.Getenv("SMTP_PORT") != "" {
		// No hacer la conversión aquí para evitar errores, simplemente usar GetInt más adelante
	} else {
		smtpPort = viper.GetInt("email.smtp_port")
	}

	smtpUser := os.Getenv("SMTP_USER")
	if smtpUser == "" {
		smtpUser = viper.GetString("email.smtp_user")
	}

	smtpPass := os.Getenv("SMTP_PASS")
	if smtpPass == "" {
		smtpPass = viper.GetString("email.smtp_pass")
	}

	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = viper.GetString("email.smtp_from")
	}

	// Parsear la configuración
	Config = &Configuration{
		App: AppConfig{
			Name:        viper.GetString("app.name"),
			Environment: viper.GetString("app.environment"),
			Port:        viper.GetInt("app.port"),
			URL:         viper.GetString("app.url"),
			SessionTTL:  viper.GetInt("app.session_ttl"),
		},
		Database: DatabaseConfig{
			Driver:          viper.GetString("database.driver"),
			Host:            viper.GetString("database.host"),
			Port:            viper.GetInt("database.port"),
			Username:        dbUsername,
			Password:        dbPassword,
			Name:            dbName,
			Charset:         viper.GetString("database.charset"),
			Loc:             viper.GetString("database.loc"),
			MaxIdleConns:    viper.GetInt("database.max_idle_conns"),
			MaxOpenConns:    viper.GetInt("database.max_open_conns"),
			ConnMaxLifetime: viper.GetDuration("database.conn_max_lifetime"),
		},
		Scraper: ScraperConfig{
			UpdateInterval: viper.GetDuration("scraper.update_interval"),
			UserAgent:      viper.GetString("scraper.user_agent"),
			MaxRetries:     viper.GetInt("scraper.max_retries"),
			RetryDelay:     viper.GetDuration("scraper.retry_delay"),
		},
		Email: EmailConfig{
			SMTPHost: smtpHost,
			SMTPPort: smtpPort,
			SMTPUser: smtpUser,
			SMTPPass: smtpPass,
			SMTPFrom: smtpFrom,
		},
	}

	log.Printf("Configuración cargada correctamente. Ambiente: %s", Config.App.Environment)
}
