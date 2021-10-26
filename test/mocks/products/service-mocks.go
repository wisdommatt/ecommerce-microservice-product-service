package products

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
)

type ServiceMock struct {
	AddProductFunc func(ctx context.Context, jwtToken string, newProduct *products.Product) (*products.Product, error)
	GetProductFunc func(ctx context.Context, sku string) (*products.Product, error)
}

func (s *ServiceMock) AddProduct(ctx context.Context, jwtToken string, newProduct *products.Product) (*products.Product, error) {
	return s.AddProductFunc(ctx, jwtToken, newProduct)
}

func (s *ServiceMock) GetProduct(ctx context.Context, sku string) (*products.Product, error) {
	return s.GetProductFunc(ctx, sku)
}
