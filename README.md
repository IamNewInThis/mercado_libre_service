# Mercado Libre Service

Middleware desarrollado en Go para la integración entre Odoo y Mercado Libre.

## Descripción

Este servicio actúa como un middleware que facilita la comunicación y sincronización de datos entre sistemas Odoo y la API de Mercado Libre. Proporciona una capa de abstracción que permite gestionar productos, órdenes, inventario y otras operaciones comerciales de manera eficiente.

## Características

- 🔄 Sincronización bidireccional entre Odoo y Mercado Libre
- 📦 Gestión de productos y categorías
- 🛒 Procesamiento de órdenes y ventas
- 📊 Actualización automática de inventario
- 🔐 Autenticación OAuth con Mercado Libre
- 🚀 Alta performance gracias a Go
- 📝 Logging detallado de operaciones
- ⚡ Procesamiento asíncrono de tareas

## Tecnologías

- **Lenguaje:** Go (Golang)
- **API:** Mercado Libre API
- **ERP:** Odoo
- **Arquitectura:** Microservicios

## Requisitos Previos

- Go 1.21 o superior
- Acceso a instancia de Odoo
- Credenciales de aplicación de Mercado Libre (App ID y Secret Key)
- Git

## Instalación

```bash
# Clonar el repositorio
git clone https://github.com/IamNewInThis/mercado_libre_service.git

# Navegar al directorio del proyecto
cd mercado_libre_service

# Descargar dependencias
go mod download

# Compilar el proyecto
go build -o mercado_libre_service
```

## Configuración

Crear un archivo `.env` en la raíz del proyecto con las siguientes variables:

```env
# Mercado Libre Configuration
ML_APP_ID=your_app_id
ML_SECRET_KEY=your_secret_key
ML_REDIRECT_URI=your_redirect_uri
ML_COUNTRY=CL

# Odoo Configuration
ODOO_URL=https://your-odoo-instance.com
ODOO_DB=your_database
ODOO_USERNAME=your_username
ODOO_PASSWORD=your_password

# Service Configuration
PORT=8080
LOG_LEVEL=info
ENVIRONMENT=production
```

## Uso

```bash
# Ejecutar el servicio
./mercado_libre_service

# O usando go run
go run main.go
```

## Estructura del Proyecto

```
mercado_libre_service/
├── cmd/                    # Punto de entrada de la aplicación
├── internal/              # Código privado de la aplicación
│   ├── api/              # Handlers y rutas HTTP
│   ├── config/           # Configuración de la aplicación
│   ├── mercadolibre/     # Cliente y lógica de Mercado Libre
│   ├── odoo/             # Cliente y lógica de Odoo
│   ├── models/           # Modelos de datos
│   └── services/         # Lógica de negocio
├── pkg/                   # Código reutilizable
├── tests/                 # Tests
├── go.mod                # Dependencias del proyecto
├── go.sum                # Checksums de dependencias
└── README.md             # Este archivo
```

## API Endpoints

### Autenticación
- `POST /auth/mercadolibre` - Iniciar OAuth con Mercado Libre
- `GET /auth/callback` - Callback de OAuth

### Productos
- `POST /products/sync` - Sincronizar productos desde Odoo a ML
- `GET /products/:id` - Obtener producto
- `PUT /products/:id` - Actualizar producto

### Órdenes
- `GET /orders/sync` - Sincronizar órdenes desde ML a Odoo
- `GET /orders/:id` - Obtener detalle de orden

### Inventario
- `POST /inventory/update` - Actualizar inventario

## Testing

```bash
# Ejecutar todos los tests
go test ./...

# Ejecutar tests con cobertura
go test -cover ./...

# Ejecutar tests de un paquete específico
go test ./internal/services/...
```

## Deployment

### Docker

```bash
# Construir imagen
docker build -t mercado_libre_service .

# Ejecutar contenedor
docker run -p 8080:8080 --env-file .env mercado_libre_service
```

### Docker Compose

```bash
docker-compose up -d
```

## Contribución

1. Fork el proyecto
2. Crear una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abrir un Pull Request

## Roadmap

- [ ] Implementar autenticación OAuth con Mercado Libre
- [ ] Desarrollar cliente para API de Odoo
- [ ] Desarrollar cliente para API de Mercado Libre
- [ ] Sincronización de productos
- [ ] Sincronización de órdenes
- [ ] Actualización de inventario
- [ ] Sistema de webhooks
- [ ] Panel de administración
- [ ] Documentación API (Swagger)
- [ ] Métricas y monitoreo

## Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## Soporte

Para reportar bugs o solicitar features, por favor abrir un issue en el repositorio de GitHub.

## Autores

- **Simple Digital** - [simpledigitalCL](https://github.com/simpledigitalCL)

## Agradecimientos

- Mercado Libre por su API
- Comunidad de Odoo
- Comunidad de Go
