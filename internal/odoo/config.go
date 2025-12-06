package odoo

import (
	"fmt"
	"os"
)

// Config contiene la configuración para conectarse a Odoo
type Config struct {
	// Configuración de conexión Odoo
	URL      string
	Database string
	Username string
	Password string

	// Información del cliente (multi-tenant)
	ClientID   string
	ClientName string
}

// NewConfigFromEnv crea una configuración desde variables de entorno
// ⚠️ Modo legacy (retrocompatibilidad)
func NewConfigFromEnv() (*Config, error) {
	url := os.Getenv("ODOO_URL")
	if url == "" {
		return nil, fmt.Errorf("ODOO_URL no está configurado")
	}

	database := os.Getenv("ODOO_DATABASE")
	if database == "" {
		return nil, fmt.Errorf("ODOO_DATABASE no está configurado")
	}

	username := os.Getenv("ODOO_USERNAME")
	if username == "" {
		return nil, fmt.Errorf("ODOO_USERNAME no está configurado")
	}

	password := os.Getenv("ODOO_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("ODOO_PASSWORD no está configurado")
	}

	return &Config{
		URL:        url,
		Database:   database,
		Username:   username,
		Password:   password,
		ClientID:   "default",
		ClientName: "Default Client",
	}, nil
}
