package main

import (
	"log"
	"os"

	"github.com/IamNewInThis/mercado_libre_service/internal/mercadolibre"
	"github.com/IamNewInThis/mercado_libre_service/internal/odoo"
	"github.com/IamNewInThis/mercado_libre_service/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env (si existe)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Configurar cliente Odoo
	odooConfig, err := odoo.NewConfigFromEnv()
	if err != nil {
		log.Printf("⚠️ Error configurando Odoo: %v", err)
		log.Println("ℹ️ El servidor iniciará sin conexión a Odoo")
	}

	var odooClient *odoo.Client
	if odooConfig != nil {
		odooClient = odoo.NewClient(odooConfig)

		// Intentar autenticar al inicio
		if err := odooClient.Authenticate(); err != nil {
			log.Printf("⚠️ Error autenticando con Odoo: %v", err)
			log.Println("ℹ️ El servidor iniciará, pero la conexión a Odoo no está disponible")
		}
	}

	// Configurar cliente Mercado Libre
	mlConfig, err := mercadolibre.NewConfigFromEnv()
	if err != nil {
		log.Printf("⚠️ Error configurando Mercado Libre: %v", err)
		log.Println("ℹ️ El servidor iniciará sin integración con Mercado Libre")
	}

	var mlClient *mercadolibre.Client
	if mlConfig != nil {
		mlClient = mercadolibre.NewClient(mlConfig)
		log.Println("✅ Cliente Mercado Libre configurado")
	}

	// Obtener puerto del entorno o usar 8081 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Crear e iniciar servidor
	srv := server.NewServer(port, odooClient, mlClient)

	log.Printf("🎯 Mercado Libre Service - Middleware Odoo/MercadoLibre")
	log.Printf("🌐 Escuchando en puerto %s", port)

	if err := srv.Start(); err != nil {
		log.Fatalf("❌ Error iniciando servidor: %v", err)
	}
}
