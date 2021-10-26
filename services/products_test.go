package services

import (
	"context"
	"errors"
	"testing"

	"github.com/wisdommatt/ecommerce-microservice-product-service/internal/products"
	productmock "github.com/wisdommatt/ecommerce-microservice-product-service/test/mocks/products"
)

func TestProductServiceImpl_AddProduct(t *testing.T) {
	type args struct {
		newProduct *products.Product
	}
	tests := []struct {
		name                string
		args                args
		repoSaveProductFunc func(ctx context.Context, product *products.Product) error
		wantErr             bool
	}{
		{
			name: "SaveProduct repository implementation with error",
			args: args{newProduct: &products.Product{
				Name:        "Yellow Speaker",
				Description: "Speaker with yellow sounds",
			}},
			repoSaveProductFunc: func(ctx context.Context, product *products.Product) error {
				return errors.New("an error occured while saving product")
			},
			wantErr: true,
		},
		{
			name: "SaveProduct repository implementation without error",
			args: args{newProduct: &products.Product{
				Name:        "Green Speaker",
				Description: "Speaker with green charger",
				Price:       10009999,
			}},
			repoSaveProductFunc: func(ctx context.Context, product *products.Product) error {
				return nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductService(&productmock.RepositoryMock{
				SaveProductFunc: tt.repoSaveProductFunc,
			})
			_, err := s.AddProduct(context.TODO(), tt.args.newProduct)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductServiceImpl.AddProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
