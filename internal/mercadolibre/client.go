package mercadolibre

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client representa el cliente para Mercado Libre API
type Client struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
	httpClient   *http.Client
}

// Config contiene la configuración de Mercado Libre
type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// NewClient crea un nuevo cliente de Mercado Libre
func NewClient(config *Config) *Client {
	return &Client{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURI:  config.RedirectURI,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

// TokenResponse representa la respuesta del endpoint de token
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	UserID       int64  `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

// GetAuthURL genera la URL para que el usuario autorice la aplicación
func (c *Client) GetAuthURL(state string) string {
	baseURL := "https://auth.mercadolibre.com.co/authorization"
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", c.ClientID)
	params.Add("redirect_uri", c.RedirectURI)
	params.Add("state", state)

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

// ExchangeCodeForToken intercambia el código de autorización por un access token
func (c *Client) ExchangeCodeForToken(code string) error {
	fmt.Printf("🔑 Intercambiando código por access token...\n")

	tokenURL := "https://api.mercadolibre.com/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", c.RedirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error en la petición: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error HTTP %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("error parseando respuesta: %w", err)
	}

	c.AccessToken = tokenResp.AccessToken
	c.RefreshToken = tokenResp.RefreshToken
	c.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	fmt.Printf("✅ Access token obtenido. Expira en: %v\n", c.ExpiresAt)
	return nil
}

// RefreshAccessToken renueva el access token usando el refresh token
func (c *Client) RefreshAccessToken() error {
	if c.RefreshToken == "" {
		return fmt.Errorf("no hay refresh token disponible")
	}

	fmt.Printf("🔄 Renovando access token...\n")

	tokenURL := "https://api.mercadolibre.com/oauth/token"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("refresh_token", c.RefreshToken)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("error creando request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error en la petición: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error leyendo respuesta: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error HTTP %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("error parseando respuesta: %w", err)
	}

	c.AccessToken = tokenResp.AccessToken
	c.RefreshToken = tokenResp.RefreshToken
	c.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	fmt.Printf("✅ Access token renovado. Expira en: %v\n", c.ExpiresAt)
	return nil
}

// IsTokenValid verifica si el token actual es válido
func (c *Client) IsTokenValid() bool {
	if c.AccessToken == "" {
		return false
	}
	return time.Now().Before(c.ExpiresAt.Add(-5 * time.Minute)) // 5 min de margen
}

// EnsureValidToken asegura que hay un token válido, renovándolo si es necesario
func (c *Client) EnsureValidToken() error {
	if c.IsTokenValid() {
		return nil
	}

	if c.RefreshToken != "" {
		return c.RefreshAccessToken()
	}

	return fmt.Errorf("no hay token válido y no hay refresh token para renovar")
}
