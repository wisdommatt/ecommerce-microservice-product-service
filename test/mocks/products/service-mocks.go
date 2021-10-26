package products

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
)

type ServiceMock struct {
	AddProductFunc func(ctx context.Context, newProduct *products.Product) (*products.Product, error)
}

func (s *ServiceMock) AddProduct(ctx context.Context, newProduct *products.Product) (*products.Product, error) {
	return s.AddProductFunc(ctx, newProduct)
}
