package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IamNewInThis/mercado_libre_service/internal/odoo"
)

// Server representa el servidor HTTP
type Server struct {
	port       string
	odooClient *odoo.Client
	httpServer *http.Server
}

// NewServer crea una nueva instancia del servidor
func NewServer(port string, odooClient *odoo.Client) *Server {
	return &Server{
		port:       port,
		odooClient: odooClient,
	}
}

// Start inicia el servidor HTTP
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Rutas
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/odoo/status", s.handleOdooStatus)
	mux.HandleFunc("/odoo/product/", s.handleGetProduct)

	s.httpServer = &http.Server{
		Addr:         ":" + s.port,
		Handler:      s.loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Servidor iniciado en http://localhost:%s\n", s.port)
	return s.httpServer.ListenAndServe()
}

// Middleware de logging
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("📥 %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("📤 %s %s - %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// handleHome maneja la ruta principal
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service": "Mercado Libre Service",
		"version": "1.0.1",
		"status":  "running",
		"message": "Middleware para integración Odoo - Mercado Libre",
	}
	s.sendJSON(w, http.StatusOK, response)
}

// handleHealth verifica el estado del servicio
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}
	s.sendJSON(w, http.StatusOK, response)
}

// handleOdooStatus verifica la conexión con Odoo
func (s *Server) handleOdooStatus(w http.ResponseWriter, r *http.Request) {
	if s.odooClient == nil {
		s.sendJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status":  "error",
			"message": "Cliente Odoo no configurado",
		})
		return
	}

	// Verificar si ya está autenticado
	if s.odooClient.UID == 0 {
		err := s.odooClient.Authenticate()
		if err != nil {
			s.sendJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
				"status":  "error",
				"message": fmt.Sprintf("Error autenticando con Odoo: %v", err),
			})
			return
		}
	}

	response := map[string]interface{}{
		"status":      "connected",
		"client_name": s.odooClient.ClientName,
		"uid":         s.odooClient.UID,
		"database":    s.odooClient.Database,
	}
	s.sendJSON(w, http.StatusOK, response)
}

// handleGetProduct obtiene información de un producto por ID
func (s *Server) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	if s.odooClient == nil {
		s.sendJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"error": "Cliente Odoo no configurado",
		})
		return
	}

	// Verificar autenticación
	if s.odooClient.UID == 0 {
		if err := s.odooClient.Authenticate(); err != nil {
			s.sendJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
				"error": fmt.Sprintf("Error autenticando con Odoo: %v", err),
			})
			return
		}
	}

	// Extraer ID del path: /odoo/product/{id}
	path := r.URL.Path
	var productID int
	_, err := fmt.Sscanf(path, "/odoo/product/%d", &productID)
	if err != nil || productID <= 0 {
		s.sendJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "ID de producto inválido. Usa: /odoo/product/{id}",
		})
		return
	}

	// Obtener producto
	productService := odoo.NewProductService(s.odooClient)
	product, err := productService.GetProductByID(productID)
	if err != nil {
		s.sendJSON(w, http.StatusNotFound, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	s.sendJSON(w, http.StatusOK, product)
}

// sendJSON envía una respuesta JSON
func (s *Server) sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("❌ Error encoding JSON: %v", err)
	}
}
