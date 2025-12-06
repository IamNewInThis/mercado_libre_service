package mercadolibre

import (
	"fmt"
	"os"
)

// NewConfigFromEnv crea una configuración desde variables de entorno
func NewConfigFromEnv() (*Config, error) {
	clientID := os.Getenv("ML_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("ML_CLIENT_ID no está configurado")
	}

	clientSecret := os.Getenv("ML_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("ML_CLIENT_SECRET no está configurado")
	}

	redirectURI := os.Getenv("ML_REDIRECT_URI")
	if redirectURI == "" {
		return nil, fmt.Errorf("ML_REDIRECT_URI no está configurado")
	}

	return &Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}, nil
}
