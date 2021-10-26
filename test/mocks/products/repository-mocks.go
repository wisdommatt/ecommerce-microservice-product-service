package products

import (
	"context"

	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
)

type RepositoryMock struct {
	SaveProductFunc func(ctx context.Context, product *products.Product) error
}

func (r *RepositoryMock) SaveProduct(ctx context.Context, product *products.Product) error {
	return r.SaveProductFunc(ctx, product)
}
