package odoo

import "fmt"

// ProductTemplate representa la información de un producto template
type ProductTemplate struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ListPrice float64 `json:"list_price"`
}

// ProductService proporciona operaciones para productos
type ProductService struct {
	client *Client
}

// NewProductService crea un nuevo servicio de productos
func NewProductService(client *Client) *ProductService {
	return &ProductService{
		client: client,
	}
}

// GetProductByID obtiene un producto por su ID de product.template
// 📦 Obtiene nombre y precio de venta de un producto
func (s *ProductService) GetProductByID(productID int) (*ProductTemplate, error) {
	if s.client.UID == 0 {
		return nil, fmt.Errorf("cliente no autenticado")
	}

	fmt.Printf("🔍 Buscando producto ID: %d\n", productID)

	// Ejecutar método read en product.template
	payload := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  "call",
		Params: map[string]interface{}{
			"service": "object",
			"method":  "execute_kw",
			"args": []interface{}{
				s.client.Database,
				s.client.UID,
				s.client.Password,
				"product.template",
				"read",
				[]interface{}{[]int{productID}},
				map[string]interface{}{
					"fields": []string{"id", "name", "list_price"},
				},
			},
		},
		ID: 1,
	}

	response, err := s.client.doRequest(payload)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo producto: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("error de Odoo: %s", response.Error.Message)
	}

	// Procesar resultado
	resultSlice, ok := response.Result.([]interface{})
	if !ok || len(resultSlice) == 0 {
		return nil, fmt.Errorf("producto no encontrado con ID: %d", productID)
	}

	productData, ok := resultSlice[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("formato de respuesta inválido")
	}

	product := &ProductTemplate{
		ID:        int(productData["id"].(float64)),
		Name:      productData["name"].(string),
		ListPrice: productData["list_price"].(float64),
	}

	fmt.Printf("✅ Producto encontrado: %s (Precio: %.2f)\n", product.Name, product.ListPrice)
	return product, nil
}
