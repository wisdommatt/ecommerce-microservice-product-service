package services

import (
	"context"
	"errors"

	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
)

// ProductService is the interface that describes a product service.
type ProductService interface {
	AddProduct(ctx context.Context, newProduct *products.Product) (*products.Product, error)
	GetProduct(ctx context.Context, sku string) (*products.Product, error)
}

// ProductServiceImpl is the default implementation for ProductService
// interface.
type ProductServiceImpl struct {
	productRepo products.Repository
}

// NewProductService returns a new product service object.
func NewProductService(productRepo products.Repository) *ProductServiceImpl {
	return &ProductServiceImpl{
		productRepo: productRepo,
	}
}

func (s *ProductServiceImpl) AddProduct(ctx context.Context, newProduct *products.Product) (*products.Product, error) {
	err := s.productRepo.SaveProduct(ctx, newProduct)
	if err != nil {
		return nil, errors.New("an error occured while adding product, please try again later")
	}
	return newProduct, nil
}

func (s *ProductServiceImpl) GetProduct(ctx context.Context, sku string) (*products.Product, error) {
	if sku == "" {
		return nil, errors.New("sku must be provided")
	}
	product, err := s.productRepo.GetProductBySKU(ctx, sku)
	if err != nil {
		return nil, errors.New("product does not exist")
	}
	return product, nil
}
