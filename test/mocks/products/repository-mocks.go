package products

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
)

type RepositoryMock struct {
	SaveProductFunc     func(ctx context.Context, product *products.Product) error
	GetProductBySKUFunc func(ctx context.Context, sku string) (*products.Product, error)
}

func (r *RepositoryMock) SaveProduct(ctx context.Context, product *products.Product) error {
	return r.SaveProductFunc(ctx, product)
}

func (r *RepositoryMock) GetProductBySKU(ctx context.Context, sku string) (*products.Product, error) {
	return r.GetProductBySKUFunc(ctx, sku)
}
